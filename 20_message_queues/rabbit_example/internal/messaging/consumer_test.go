package messaging

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	"example.com/rabbit-example/internal/event"
	amqp "github.com/rabbitmq/amqp091-go"
)

func testEnvelope() event.Envelope {
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

func TestDecodeDelivery(t *testing.T) {
	eventEnvelope := testEnvelope()
	body, err := json.Marshal(eventEnvelope)
	if err != nil {
		t.Fatal(err)
	}

	got, attempt, err := decodeDelivery(amqp.Delivery{
		Headers:     amqp.Table{HeaderRetryCount: int64(2)},
		ContentType: "application/json",
		MessageId:   eventEnvelope.ID,
		Body:        body,
	})
	if err != nil {
		t.Fatalf("decodeDelivery() error = %v", err)
	}
	if got.ID != eventEnvelope.ID || attempt != 2 {
		t.Fatalf("decodeDelivery() = (%+v, %d), want event-1 and 2", got, attempt)
	}
}

func TestDecodeDeliveryPermanentErrors(t *testing.T) {
	eventEnvelope := testEnvelope()
	body, _ := json.Marshal(eventEnvelope)
	tests := []amqp.Delivery{
		{ContentType: "text/plain", MessageId: eventEnvelope.ID, Body: body},
		{ContentType: "application/json", MessageId: eventEnvelope.ID, Body: []byte("{")},
		{ContentType: "application/json", MessageId: "different", Body: body},
		{Headers: amqp.Table{HeaderRetryCount: -1}, ContentType: "application/json", MessageId: eventEnvelope.ID, Body: body},
	}

	for i, delivery := range tests {
		_, _, err := decodeDelivery(delivery)
		if !errors.Is(err, ErrPermanent) {
			t.Errorf("case %d: error = %v, want ErrPermanent", i, err)
		}
	}
}

func TestRetryCountSupportedTypes(t *testing.T) {
	for _, value := range []any{int8(2), int16(2), int32(2), int64(2), int(2), "2"} {
		got, err := retryCount(amqp.Table{HeaderRetryCount: value})
		if err != nil || got != 2 {
			t.Errorf("retryCount(%T(%v)) = (%d, %v), want (2, nil)", value, value, got, err)
		}
	}
}

func TestRetryPublishingCopiesAndIncrementsHeaders(t *testing.T) {
	original := amqp.Delivery{
		Headers:     amqp.Table{"original": "value"},
		MessageId:   "event-1",
		ContentType: "application/json",
		Expiration:  "1000",
		Body:        []byte("body"),
	}
	published := retryPublishing(original, 3)
	published.Headers["original"] = "changed"
	published.Body[0] = 'B'

	if original.Headers["original"] != "value" || string(original.Body) != "body" {
		t.Fatal("retryPublishing mutated delivery data")
	}
	if published.Headers[HeaderRetryCount] != int64(3) {
		t.Fatalf("retry count = %v, want 3", published.Headers[HeaderRetryCount])
	}
	if published.Expiration != "" {
		t.Fatalf("retry expiration = %q, want queue TTL only", published.Expiration)
	}
}
