package messaging

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"example.com/rabbit-example/internal/event"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	HeaderRetryCount       = "x-retry-count"
	HeaderDemoFailure      = "x-demo-failure"
	HeaderDemoFailAttempts = "x-demo-fail-attempts"
)

type Publisher struct {
	ch      *amqp.Channel
	returns <-chan amqp.Return
	mu      sync.Mutex
}

func NewPublisher(ch *amqp.Channel) (*Publisher, error) {
	if err := ch.Confirm(false); err != nil {
		return nil, fmt.Errorf("enable publisher confirms: %w", err)
	}
	returns := ch.NotifyReturn(make(chan amqp.Return, 1))
	return &Publisher{ch: ch, returns: returns}, nil
}

func (p *Publisher) PublishOrderCreated(
	ctx context.Context,
	eventEnvelope event.Envelope,
	headers amqp.Table,
) error {
	if err := eventEnvelope.Validate(); err != nil {
		return fmt.Errorf("validate event: %w", err)
	}

	body, err := json.Marshal(eventEnvelope)
	if err != nil {
		return fmt.Errorf("marshal event: %w", err)
	}

	return p.Publish(ctx, EventsExchange, OrdersRoutingKey, amqp.Publishing{
		Headers:       cloneTable(headers),
		ContentType:   "application/json",
		DeliveryMode:  amqp.Persistent,
		MessageId:     eventEnvelope.ID,
		CorrelationId: eventEnvelope.CorrelationID,
		Type:          eventEnvelope.Type,
		Timestamp:     eventEnvelope.OccurredAt,
		Body:          body,
	})
}

func (p *Publisher) Publish(
	ctx context.Context,
	exchange string,
	routingKey string,
	message amqp.Publishing,
) error {
	// Нельзя использовать один Channel из нескольких publisher без координации.
	// Блокировка сохраняется до получения confirm, поэтому одновременно в
	// обработке находится только одно сообщение. Для учебного проекта такое
	// поведение остаётся простым и предсказуемым.
	p.mu.Lock()
	defer p.mu.Unlock()
	if err := p.ensureNoStaleReturn(); err != nil {
		return err
	}

	confirm, err := p.ch.PublishWithDeferredConfirmWithContext(
		ctx,
		exchange,
		routingKey,
		true, // возвращать немаршрутизируемые сообщения, а не терять их без ошибки
		false,
		message,
	)
	if err != nil {
		return fmt.Errorf("publish message %q: %w", message.MessageId, err)
	}
	if confirm == nil {
		return fmt.Errorf("publisher confirm unavailable for message %q", message.MessageId)
	}

	acknowledged, err := confirm.WaitContext(ctx)
	if err != nil {
		return fmt.Errorf("wait for confirm of message %q: %w", message.MessageId, err)
	}
	if !acknowledged {
		return fmt.Errorf("RabbitMQ negatively acknowledged message %q", message.MessageId)
	}

	// amqp091-go гарантирует, что немаршрутизируемая mandatory-публикация получает
	// confirm только после уведомления подписчиков NotifyReturn. Publisher
	// выполняет публикации последовательно, а канал возвратов имеет буфер, поэтому
	// неблокирующее чтение после confirm однозначно относится к единственному
	// сообщению, находящемуся в обработке.
	select {
	case returned, ok := <-p.returns:
		if !ok {
			return fmt.Errorf("return notification channel closed for message %q", message.MessageId)
		}
		return fmt.Errorf(
			"message %q was returned as unroutable: code=%d reason=%q exchange=%q routing_key=%q",
			returned.MessageId,
			returned.ReplyCode,
			returned.ReplyText,
			returned.Exchange,
			returned.RoutingKey,
		)
	default:
	}
	return nil
}

func (p *Publisher) ensureNoStaleReturn() error {
	select {
	case returned, ok := <-p.returns:
		if !ok {
			return errors.New("return notification channel is closed")
		}
		return fmt.Errorf("unexpected stale return for message %q", returned.MessageId)
	default:
		return nil
	}
}

func retryPublishing(delivery amqp.Delivery, retryCount int) amqp.Publishing {
	headers := cloneTable(delivery.Headers)
	headers[HeaderRetryCount] = int64(retryCount)

	return amqp.Publishing{
		Headers:         headers,
		ContentType:     delivery.ContentType,
		ContentEncoding: delivery.ContentEncoding,
		DeliveryMode:    amqp.Persistent,
		Priority:        delivery.Priority,
		CorrelationId:   delivery.CorrelationId,
		ReplyTo:         delivery.ReplyTo,
		// У RetryQueue есть собственный фиксированный TTL. Если перенести время
		// жизни исходного сообщения, повтор может произойти раньше заданной задержки.
		Expiration: "",
		MessageId:  delivery.MessageId,
		Timestamp:  delivery.Timestamp,
		Type:       delivery.Type,
		UserId:     delivery.UserId,
		AppId:      delivery.AppId,
		Body:       append([]byte(nil), delivery.Body...),
	}
}

func cloneTable(source amqp.Table) amqp.Table {
	if len(source) == 0 {
		return amqp.Table{}
	}
	result := make(amqp.Table, len(source))
	for key, value := range source {
		result[key] = value
	}
	return result
}
