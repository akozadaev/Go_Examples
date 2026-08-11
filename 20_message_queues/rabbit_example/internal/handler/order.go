package handler

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"example.com/rabbit-example/internal/event"
	"example.com/rabbit-example/internal/messaging"
)

const processingDuration = 200 * time.Millisecond

type OrderHandler struct {
	processed ProcessedStore
}

func NewOrderHandler(processed ProcessedStore) *OrderHandler {
	return &OrderHandler{processed: processed}
}

func (h *OrderHandler) HandleOrderCreated(
	ctx context.Context,
	eventEnvelope event.Envelope,
	metadata messaging.DeliveryMetadata,
) error {
	if h.processed.IsProcessed(eventEnvelope.ID) {
		log.Printf("event=duplicate_skipped message_id=%q", eventEnvelope.ID)
		return nil
	}

	failureMode, _ := metadata.Headers[messaging.HeaderDemoFailure].(string)
	switch failureMode {
	case "":
	case "permanent":
		return fmt.Errorf("%w: requested demo permanent failure", messaging.ErrPermanent)
	case "transient":
		failAttempts, err := demoFailAttempts(metadata.Headers[messaging.HeaderDemoFailAttempts])
		if err != nil {
			return fmt.Errorf("%w: %v", messaging.ErrPermanent, err)
		}
		if metadata.Attempt < failAttempts {
			return fmt.Errorf(
				"requested demo transient failure %d/%d",
				metadata.Attempt+1,
				failAttempts,
			)
		}
	default:
		return fmt.Errorf("%w: unknown demo failure mode %q", messaging.ErrPermanent, failureMode)
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(processingDuration):
	}

	// Здесь настоящий обработчик фиксировал бы бизнес-транзакцию.
	h.processed.MarkProcessed(eventEnvelope.ID)
	log.Printf(
		"event=order_processed message_id=%q order_id=%q amount_minor=%d currency=%q",
		eventEnvelope.ID,
		eventEnvelope.Data.OrderID,
		eventEnvelope.Data.AmountMinor,
		eventEnvelope.Data.Currency,
	)
	return nil
}

func demoFailAttempts(value any) (int, error) {
	if value == nil {
		return 1, nil
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
			return 0, fmt.Errorf("invalid demo fail-attempts value %q", typed)
		}
		result = parsed
	default:
		return 0, fmt.Errorf("invalid demo fail-attempts type %T", value)
	}

	if result < 1 || result > int64(^uint(0)>>1) {
		return 0, fmt.Errorf("demo fail-attempts must be a positive integer")
	}
	return int(result), nil
}
