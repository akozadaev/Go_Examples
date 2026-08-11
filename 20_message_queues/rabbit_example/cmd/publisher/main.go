package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"example.com/rabbit-example/internal/config"
	"example.com/rabbit-example/internal/event"
	"example.com/rabbit-example/internal/messaging"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	if err := run(); err != nil {
		log.Printf("level=error component=publisher error=%q", err)
		os.Exit(1)
	}
}

func run() error {
	var (
		orderID       = flag.String("order-id", "order-42", "order identifier")
		customerID    = flag.String("customer-id", "customer-7", "customer identifier")
		amountMinor   = flag.Int64("amount", 159900, "amount in minor currency units")
		currency      = flag.String("currency", "RUB", "ISO-like currency code")
		correlationID = flag.String("correlation-id", "", "correlation identifier")
		messageID     = flag.String("message-id", "", "override generated event ID")
		failure       = flag.String("demo-failure", "", "demo failure: transient or permanent")
		failAttempts  = flag.Int("fail-attempts", 1, "transient failures before success")
		malformed     = flag.Bool("malformed", false, "publish malformed JSON for DLQ demonstration")
	)
	flag.Parse()

	if *failure != "" && *failure != "transient" && *failure != "permanent" {
		return fmt.Errorf("demo-failure must be transient or permanent")
	}
	if *failAttempts < 1 {
		return fmt.Errorf("fail-attempts must be positive")
	}
	if *correlationID == "" {
		*correlationID = "checkout-" + *orderID
	}

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	eventEnvelope, err := event.NewOrderCreated(event.OrderCreated{
		OrderID:     *orderID,
		CustomerID:  *customerID,
		AmountMinor: *amountMinor,
		Currency:    *currency,
	}, *correlationID)
	if err != nil {
		return fmt.Errorf("create event: %w", err)
	}
	if *messageID != "" {
		eventEnvelope.ID = *messageID
	}

	conn, err := messaging.Dial(cfg.AMQPURL)
	if err != nil {
		return err
	}
	defer conn.Close()

	ch, err := messaging.OpenChannel(conn)
	if err != nil {
		return err
	}
	defer ch.Close()

	if err := messaging.DeclareTopology(ch, cfg.RetryDelay); err != nil {
		return err
	}
	publisher, err := messaging.NewPublisher(ch)
	if err != nil {
		return err
	}

	headers := amqp.Table{}
	if *failure != "" {
		headers[messaging.HeaderDemoFailure] = *failure
		headers[messaging.HeaderDemoFailAttempts] = int64(*failAttempts)
	}

	ctx, cancel := context.WithTimeout(context.Background(), cfg.PublishTimeout)
	defer cancel()

	publishedBody := []byte(nil)
	if *malformed {
		publishedBody = []byte(`{"broken":`)
		err = publisher.Publish(ctx, messaging.EventsExchange, messaging.OrdersRoutingKey, amqp.Publishing{
			Headers:       headers,
			ContentType:   "application/json",
			DeliveryMode:  amqp.Persistent,
			MessageId:     eventEnvelope.ID,
			CorrelationId: eventEnvelope.CorrelationID,
			Type:          eventEnvelope.Type,
			Timestamp:     time.Now().UTC(),
			Body:          publishedBody,
		})
	} else {
		err = publisher.PublishOrderCreated(ctx, eventEnvelope, headers)
		publishedBody, _ = json.Marshal(eventEnvelope)
	}
	if err != nil {
		return err
	}

	log.Printf(
		"level=info event=published message_id=%q routing_key=%q payload=%s",
		eventEnvelope.ID,
		messaging.OrdersRoutingKey,
		publishedBody,
	)
	return nil
}
