package event

import (
	"strings"
	"testing"
	"time"
)

func validEnvelope() Envelope {
	return Envelope{
		ID:            "event-1",
		Type:          OrderCreatedType,
		SchemaVersion: CurrentSchemaVersion,
		OccurredAt:    time.Now(),
		CorrelationID: "checkout-1",
		Data: OrderCreated{
			OrderID:     "order-1",
			CustomerID:  "customer-1",
			AmountMinor: 100,
			Currency:    "RUB",
		},
	}
}

func TestEnvelopeValidate(t *testing.T) {
	tests := []struct {
		name   string
		change func(*Envelope)
		want   string
	}{
		{name: "valid"},
		{name: "missing ID", change: func(e *Envelope) { e.ID = "" }, want: "event ID"},
		{name: "wrong type", change: func(e *Envelope) { e.Type = "other" }, want: "event type"},
		{name: "wrong version", change: func(e *Envelope) { e.SchemaVersion = 2 }, want: "schema version"},
		{name: "missing order", change: func(e *Envelope) { e.Data.OrderID = "" }, want: "order_id"},
		{name: "invalid amount", change: func(e *Envelope) { e.Data.AmountMinor = 0 }, want: "amount_minor"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := validEnvelope()
			if tt.change != nil {
				tt.change(&event)
			}
			err := event.Validate()
			if tt.want == "" && err != nil {
				t.Fatalf("Validate() error = %v", err)
			}
			if tt.want != "" && (err == nil || !strings.Contains(err.Error(), tt.want)) {
				t.Fatalf("Validate() error = %v, want containing %q", err, tt.want)
			}
		})
	}
}

func TestNewOrderCreatedGeneratesUUID(t *testing.T) {
	event, err := NewOrderCreated(OrderCreated{
		OrderID:     "order-1",
		CustomerID:  "customer-1",
		AmountMinor: 100,
		Currency:    "RUB",
	}, "checkout-1")
	if err != nil {
		t.Fatalf("NewOrderCreated() error = %v", err)
	}
	if len(event.ID) != 36 || event.ID[14] != '4' {
		t.Fatalf("generated ID %q is not UUIDv4-shaped", event.ID)
	}
}
