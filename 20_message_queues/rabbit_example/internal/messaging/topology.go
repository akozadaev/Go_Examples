package messaging

import (
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	EventsExchange = "orders.events"
	RetryExchange  = "orders.retry"
	DeadExchange   = "orders.dlx"

	OrdersQueue = "orders.created.worker"
	RetryQueue  = "orders.created.retry"
	DeadQueue   = "orders.created.dead"

	OrdersRoutingKey = "orders.created"
	RetryRoutingKey  = "orders.created.retry"
	DeadRoutingKey   = "orders.created.failed"
)

func DeclareTopology(ch *amqp.Channel, retryDelay time.Duration) error {
	if retryDelay <= 0 {
		return fmt.Errorf("retry delay must be positive")
	}

	if err := declareExchange(ch, EventsExchange, "topic"); err != nil {
		return err
	}
	if err := declareExchange(ch, RetryExchange, "direct"); err != nil {
		return err
	}
	if err := declareExchange(ch, DeadExchange, "direct"); err != nil {
		return err
	}

	if _, err := ch.QueueDeclare(DeadQueue, true, false, false, false, amqp.Table{
		amqp.QueueTypeArg:  amqp.QueueTypeQuorum,
		"x-delivery-limit": int64(-1),
	}); err != nil {
		return fmt.Errorf("declare dead-letter queue: %w", err)
	}
	if err := ch.QueueBind(DeadQueue, DeadRoutingKey, DeadExchange, false, nil); err != nil {
		return fmt.Errorf("bind dead-letter queue: %w", err)
	}

	if _, err := ch.QueueDeclare(OrdersQueue, true, false, false, false, amqp.Table{
		amqp.QueueTypeArg:           amqp.QueueTypeQuorum,
		"x-dead-letter-exchange":    DeadExchange,
		"x-dead-letter-routing-key": DeadRoutingKey,
		"x-dead-letter-strategy":    "at-least-once",
		"x-overflow":                "reject-publish",
	}); err != nil {
		return fmt.Errorf("declare orders queue: %w", err)
	}
	if err := ch.QueueBind(OrdersQueue, OrdersRoutingKey, EventsExchange, false, nil); err != nil {
		return fmt.Errorf("bind orders queue: %w", err)
	}

	retryDelayMS := retryDelay.Milliseconds()
	if retryDelayMS < 1 {
		return fmt.Errorf("retry delay must be at least 1ms")
	}
	if _, err := ch.QueueDeclare(RetryQueue, true, false, false, false, amqp.Table{
		amqp.QueueTypeArg:           amqp.QueueTypeQuorum,
		"x-message-ttl":             retryDelayMS,
		"x-dead-letter-exchange":    EventsExchange,
		"x-dead-letter-routing-key": OrdersRoutingKey,
		"x-dead-letter-strategy":    "at-least-once",
		"x-overflow":                "reject-publish",
	}); err != nil {
		return fmt.Errorf("declare retry queue: %w", err)
	}
	if err := ch.QueueBind(RetryQueue, RetryRoutingKey, RetryExchange, false, nil); err != nil {
		return fmt.Errorf("bind retry queue: %w", err)
	}

	return nil
}

func declareExchange(ch *amqp.Channel, name, kind string) error {
	if err := ch.ExchangeDeclare(name, kind, true, false, false, false, nil); err != nil {
		return fmt.Errorf("declare exchange %q: %w", name, err)
	}
	return nil
}
