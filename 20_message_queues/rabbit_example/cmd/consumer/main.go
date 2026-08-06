package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"example.com/rabbit-example/internal/config"
	"example.com/rabbit-example/internal/handler"
	"example.com/rabbit-example/internal/messaging"
)

func main() {
	if err := run(); err != nil {
		log.Printf("level=error component=consumer error=%q", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	conn, err := messaging.Dial(cfg.AMQPURL)
	if err != nil {
		return err
	}
	defer conn.Close()

	consumeChannel, err := messaging.OpenChannel(conn)
	if err != nil {
		return err
	}
	defer consumeChannel.Close()
	if err := messaging.DeclareTopology(consumeChannel, cfg.RetryDelay); err != nil {
		return err
	}

	retryChannel, err := messaging.OpenChannel(conn)
	if err != nil {
		return err
	}
	defer retryChannel.Close()
	retryPublisher, err := messaging.NewPublisher(retryChannel)
	if err != nil {
		return err
	}

	orderHandler := handler.NewOrderHandler(handler.NewMemoryProcessedStore())
	consumer, err := messaging.NewConsumer(
		consumeChannel,
		retryPublisher,
		orderHandler,
		cfg.Prefetch,
		cfg.MaxRetries,
		cfg.PublishTimeout,
	)
	if err != nil {
		return err
	}

	runCtx, cancelRun := context.WithCancel(context.Background())
	defer cancelRun()
	result := make(chan error, 1)
	go func() {
		result <- consumer.Run(runCtx)
	}()

	signalCtx, stopSignals := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stopSignals()
	log.Printf(
		"level=info event=consumer_started queue=%q prefetch=%d max_retries=%d",
		messaging.OrdersQueue,
		cfg.Prefetch,
		cfg.MaxRetries,
	)

	select {
	case err := <-result:
		return err
	case <-signalCtx.Done():
		log.Printf("level=info event=shutdown_started")
		if err := consumer.StopReceiving(); err != nil {
			cancelRun()
			return err
		}
	}

	timer := time.NewTimer(cfg.ShutdownTimeout)
	defer timer.Stop()
	select {
	case err := <-result:
		if err != nil && !errors.Is(err, context.Canceled) {
			return err
		}
		log.Printf("level=info event=shutdown_completed")
		return nil
	case <-timer.C:
		cancelRun()
		return fmt.Errorf("consumer shutdown timed out after %s", cfg.ShutdownTimeout)
	}
}
