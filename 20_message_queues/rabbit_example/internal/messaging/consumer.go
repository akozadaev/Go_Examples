package messaging

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"sync/atomic"
	"time"

	"example.com/rabbit-example/internal/event"
	amqp "github.com/rabbitmq/amqp091-go"
)

var ErrPermanent = errors.New("permanent message error")

const consumerTag = "orders-created-worker"

type DeliveryMetadata struct {
	Attempt int
	Headers amqp.Table
}

type Handler interface {
	HandleOrderCreated(context.Context, event.Envelope, DeliveryMetadata) error
}

type Consumer struct {
	consumeChannel *amqp.Channel
	retryPublisher *Publisher
	handler        Handler
	prefetch       int
	maxRetries     int
	publishTimeout time.Duration
	stopping       atomic.Bool
}

func NewConsumer(
	consumeChannel *amqp.Channel,
	retryPublisher *Publisher,
	handler Handler,
	prefetch int,
	maxRetries int,
	publishTimeout time.Duration,
) (*Consumer, error) {
	switch {
	case consumeChannel == nil:
		return nil, errors.New("consume channel is required")
	case retryPublisher == nil:
		return nil, errors.New("retry publisher is required")
	case handler == nil:
		return nil, errors.New("handler is required")
	case prefetch <= 0:
		return nil, errors.New("prefetch must be positive")
	case maxRetries < 0:
		return nil, errors.New("max retries must be non-negative")
	case publishTimeout <= 0:
		return nil, errors.New("publish timeout must be positive")
	}

	return &Consumer{
		consumeChannel: consumeChannel,
		retryPublisher: retryPublisher,
		handler:        handler,
		prefetch:       prefetch,
		maxRetries:     maxRetries,
		publishTimeout: publishTimeout,
	}, nil
}

func (c *Consumer) Run(ctx context.Context) error {
	if err := c.consumeChannel.Qos(c.prefetch, 0, false); err != nil {
		return fmt.Errorf("set consumer QoS: %w", err)
	}

	deliveries, err := c.consumeChannel.ConsumeWithContext(
		ctx,
		OrdersQueue,
		consumerTag,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("start consumer: %w", err)
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case delivery, ok := <-deliveries:
			if !ok {
				if ctx.Err() != nil || c.stopping.Load() {
					return nil
				}
				return errors.New("delivery channel closed unexpectedly")
			}
			if err := c.processDelivery(ctx, delivery); err != nil {
				return err
			}
		}
	}
}

// StopReceiving просит RabbitMQ прекратить новые доставки. Уже отправленные
// consumer сообщения остаются доступны Run и могут быть обработаны до закрытия
// каналов. Контекст обработки намеренно не отменяется в этом методе.
func (c *Consumer) StopReceiving() error {
	c.stopping.Store(true)
	if err := c.consumeChannel.Cancel(consumerTag, false); err != nil {
		return fmt.Errorf("cancel RabbitMQ consumer: %w", err)
	}
	return nil
}

func (c *Consumer) processDelivery(ctx context.Context, delivery amqp.Delivery) error {
	eventEnvelope, attempt, err := decodeDelivery(delivery)
	if err == nil {
		err = c.handler.HandleOrderCreated(ctx, eventEnvelope, DeliveryMetadata{
			Attempt: attempt,
			Headers: cloneTable(delivery.Headers),
		})
	}

	if err == nil {
		if ackErr := delivery.Ack(false); ackErr != nil {
			return fmt.Errorf("ack message %q: %w", delivery.MessageId, ackErr)
		}
		log.Printf("event=acked message_id=%q attempt=%d", delivery.MessageId, attempt)
		return nil
	}

	if ctx.Err() != nil {
		// Оставляем сообщение неподтверждённым. При закрытии канала оно вернётся
		// в RabbitMQ, поэтому незавершённая работа не потеряется при остановке.
		return nil
	}

	if errors.Is(err, ErrPermanent) || attempt >= c.maxRetries {
		reason := "permanent_error"
		if !errors.Is(err, ErrPermanent) {
			reason = "retries_exhausted"
		}
		log.Printf(
			"event=dead_letter message_id=%q attempt=%d reason=%s error=%q",
			delivery.MessageId,
			attempt,
			reason,
			err,
		)
		if nackErr := delivery.Nack(false, false); nackErr != nil {
			return fmt.Errorf("dead-letter message %q: %w", delivery.MessageId, nackErr)
		}
		return nil
	}

	nextAttempt := attempt + 1
	publishCtx, cancel := context.WithTimeout(ctx, c.publishTimeout)
	defer cancel()
	if publishErr := c.retryPublisher.Publish(
		publishCtx,
		RetryExchange,
		RetryRoutingKey,
		retryPublishing(delivery, nextAttempt),
	); publishErr != nil {
		// Исходное сообщение остаётся неподтверждённым. После возврата ошибки
		// каналы процесса закроются, и RabbitMQ сможет доставить его повторно.
		// Таймаут confirm даёт неопределённый результат, поэтому возможна вторая
		// retry-копия, а обработчик обязан быть идемпотентным.
		return fmt.Errorf("schedule retry for message %q: %w", delivery.MessageId, publishErr)
	}

	if ackErr := delivery.Ack(false); ackErr != nil {
		return fmt.Errorf("ack original message %q after retry publish: %w", delivery.MessageId, ackErr)
	}
	log.Printf(
		"event=retry_scheduled message_id=%q attempt=%d error=%q",
		delivery.MessageId,
		nextAttempt,
		err,
	)
	return nil
}

func decodeDelivery(delivery amqp.Delivery) (event.Envelope, int, error) {
	attempt, err := retryCount(delivery.Headers)
	if err != nil {
		return event.Envelope{}, 0, fmt.Errorf("%w: %v", ErrPermanent, err)
	}
	if delivery.ContentType != "application/json" {
		return event.Envelope{}, attempt, fmt.Errorf(
			"%w: unsupported content type %q",
			ErrPermanent,
			delivery.ContentType,
		)
	}

	var eventEnvelope event.Envelope
	if err := json.Unmarshal(delivery.Body, &eventEnvelope); err != nil {
		return event.Envelope{}, attempt, fmt.Errorf("%w: decode JSON: %v", ErrPermanent, err)
	}
	if err := eventEnvelope.Validate(); err != nil {
		return event.Envelope{}, attempt, fmt.Errorf("%w: validate event: %v", ErrPermanent, err)
	}
	if delivery.MessageId == "" || delivery.MessageId != eventEnvelope.ID {
		return event.Envelope{}, attempt, fmt.Errorf(
			"%w: AMQP message_id %q does not match event ID %q",
			ErrPermanent,
			delivery.MessageId,
			eventEnvelope.ID,
		)
	}
	return eventEnvelope, attempt, nil
}

func retryCount(headers amqp.Table) (int, error) {
	value, ok := headers[HeaderRetryCount]
	if !ok {
		return 0, nil
	}

	var result int64
	switch typed := value.(type) {
	case int8:
		result = int64(typed)
	case int16:
		result = int64(typed)
	case int32:
		result = int64(typed)
	case int64:
		result = typed
	case int:
		result = int64(typed)
	case string:
		parsed, err := strconv.ParseInt(typed, 10, 32)
		if err != nil {
			return 0, fmt.Errorf("invalid %s value %q", HeaderRetryCount, typed)
		}
		result = parsed
	default:
		return 0, fmt.Errorf("invalid %s type %T", HeaderRetryCount, value)
	}

	if result < 0 || result > int64(^uint(0)>>1) {
		return 0, fmt.Errorf("invalid %s value %d", HeaderRetryCount, result)
	}
	return int(result), nil
}
