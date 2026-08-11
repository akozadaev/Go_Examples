//go:build integration

package messaging

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func TestPublisherDetectsUnroutableMessage(t *testing.T) {
	url := os.Getenv("AMQP_URL")
	if url == "" {
		url = "amqp://app:app@localhost:5672/app"
	}

	conn, err := Dial(url)
	if err != nil {
		t.Fatalf("Dial() error = %v", err)
	}
	defer conn.Close()
	ch, err := OpenChannel(conn)
	if err != nil {
		t.Fatalf("OpenChannel() error = %v", err)
	}
	defer ch.Close()
	if err := DeclareTopology(ch, 5*time.Second); err != nil {
		t.Fatalf("DeclareTopology() error = %v", err)
	}
	publisher, err := NewPublisher(ch)
	if err != nil {
		t.Fatalf("NewPublisher() error = %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = publisher.Publish(ctx, EventsExchange, "route.without.binding", amqp.Publishing{
		ContentType:  "application/json",
		DeliveryMode: amqp.Persistent,
		MessageId:    "integration-unroutable",
		Body:         []byte(`{}`),
	})
	if err == nil || !strings.Contains(err.Error(), "unroutable") {
		t.Fatalf("Publish() error = %v, want unroutable error", err)
	}
}
