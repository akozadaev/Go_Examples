package handler

import (
	"context"
	"errors"
	"testing"
	"time"

	"example.com/rabbit-example/internal/event"
	"example.com/rabbit-example/internal/messaging"
	amqp "github.com/rabbitmq/amqp091-go"
)

func handlerEvent() event.Envelope {
	return event.Envelope{
		ID:            "event-1",
		Type:          event.OrderCreatedType,
		SchemaVersion: event.CurrentSchemaVersion,
		OccurredAt:    time.Now(),
		Data: event.OrderCreated{
			OrderID:     "order-1",
			CustomerID:  "customer-1",
			AmountMinor: 100,
			Currency:    "RUB",
		},
	}
}

func TestOrderHandlerTransientThenSuccess(t *testing.T) {
	store := NewMemoryProcessedStore()
	handler := NewOrderHandler(store)
	metadata := messaging.DeliveryMetadata{Headers: amqp.Table{
		messaging.HeaderDemoFailure:      "transient",
		messaging.HeaderDemoFailAttempts: int64(2),
	}}

	metadata.Attempt = 0
	if err := handler.HandleOrderCreated(context.Background(), handlerEvent(), metadata); err == nil {
		t.Fatal("attempt 0 error = nil, want transient failure")
	}
	metadata.Attempt = 1
	if err := handler.HandleOrderCreated(context.Background(), handlerEvent(), metadata); err == nil {
		t.Fatal("attempt 1 error = nil, want transient failure")
	}
	metadata.Attempt = 2
	if err := handler.HandleOrderCreated(context.Background(), handlerEvent(), metadata); err != nil {
		t.Fatalf("attempt 2 error = %v, want success", err)
	}
	if !store.IsProcessed("event-1") {
		t.Fatal("successful event was not marked processed")
	}
}

func TestOrderHandlerPermanentFailure(t *testing.T) {
	handler := NewOrderHandler(NewMemoryProcessedStore())
	err := handler.HandleOrderCreated(context.Background(), handlerEvent(), messaging.DeliveryMetadata{
		Headers: amqp.Table{messaging.HeaderDemoFailure: "permanent"},
	})
	if !errors.Is(err, messaging.ErrPermanent) {
		t.Fatalf("error = %v, want ErrPermanent", err)
	}
}

func TestOrderHandlerSkipsDuplicate(t *testing.T) {
	store := NewMemoryProcessedStore()
	store.MarkProcessed("event-1")
	handler := NewOrderHandler(store)
	err := handler.HandleOrderCreated(context.Background(), handlerEvent(), messaging.DeliveryMetadata{
		Headers: amqp.Table{messaging.HeaderDemoFailure: "permanent"},
	})
	if err != nil {
		t.Fatalf("duplicate error = %v, want nil", err)
	}
}
