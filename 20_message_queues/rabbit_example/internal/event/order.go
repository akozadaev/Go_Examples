package event

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"
)

const (
	OrderCreatedType     = "orders.created"
	CurrentSchemaVersion = 1
)

type Envelope struct {
	ID            string       `json:"id"`
	Type          string       `json:"type"`
	SchemaVersion int          `json:"schema_version"`
	OccurredAt    time.Time    `json:"occurred_at"`
	CorrelationID string       `json:"correlation_id"`
	Data          OrderCreated `json:"data"`
}

type OrderCreated struct {
	OrderID     string `json:"order_id"`
	CustomerID  string `json:"customer_id"`
	AmountMinor int64  `json:"amount_minor"`
	Currency    string `json:"currency"`
}

func NewOrderCreated(data OrderCreated, correlationID string) (Envelope, error) {
	id, err := newUUID()
	if err != nil {
		return Envelope{}, fmt.Errorf("generate event ID: %w", err)
	}

	event := Envelope{
		ID:            id,
		Type:          OrderCreatedType,
		SchemaVersion: CurrentSchemaVersion,
		OccurredAt:    time.Now().UTC(),
		CorrelationID: correlationID,
		Data:          data,
	}
	if err := event.Validate(); err != nil {
		return Envelope{}, err
	}
	return event, nil
}

func (e Envelope) Validate() error {
	switch {
	case e.ID == "":
		return errors.New("event ID is required")
	case e.Type != OrderCreatedType:
		return fmt.Errorf("unsupported event type %q", e.Type)
	case e.SchemaVersion != CurrentSchemaVersion:
		return fmt.Errorf("unsupported schema version %d", e.SchemaVersion)
	case e.OccurredAt.IsZero():
		return errors.New("occurred_at is required")
	case e.Data.OrderID == "":
		return errors.New("order_id is required")
	case e.Data.CustomerID == "":
		return errors.New("customer_id is required")
	case e.Data.AmountMinor <= 0:
		return errors.New("amount_minor must be positive")
	case e.Data.Currency == "":
		return errors.New("currency is required")
	default:
		return nil
	}
}

func newUUID() (string, error) {
	var value [16]byte
	if _, err := rand.Read(value[:]); err != nil {
		return "", err
	}
	value[6] = (value[6] & 0x0f) | 0x40
	value[8] = (value[8] & 0x3f) | 0x80

	buffer := make([]byte, 36)
	hex.Encode(buffer[0:8], value[0:4])
	buffer[8] = '-'
	hex.Encode(buffer[9:13], value[4:6])
	buffer[13] = '-'
	hex.Encode(buffer[14:18], value[6:8])
	buffer[18] = '-'
	hex.Encode(buffer[19:23], value[8:10])
	buffer[23] = '-'
	hex.Encode(buffer[24:36], value[10:16])
	return string(buffer), nil
}
