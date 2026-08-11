package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"example.com/rabbit-example/internal/config"
	"example.com/rabbit-example/internal/messaging"
)

func main() {
	if err := run(); err != nil {
		log.Printf("level=error component=dlq-reader error=%q", err)
		os.Exit(1)
	}
}

func run() error {
	ack := flag.Bool("ack", false, "remove the inspected message from the DLQ")
	flag.Parse()

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
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

	delivery, ok, err := ch.Get(messaging.DeadQueue, false)
	if err != nil {
		return fmt.Errorf("get message from DLQ: %w", err)
	}
	if !ok {
		log.Printf("level=info event=dlq_empty queue=%q", messaging.DeadQueue)
		return nil
	}

	headers, _ := json.Marshal(delivery.Headers)
	log.Printf(
		"level=info event=dlq_message message_id=%q type=%q headers=%s body=%s",
		delivery.MessageId,
		delivery.Type,
		headers,
		delivery.Body,
	)
	if *ack {
		if err := delivery.Ack(false); err != nil {
			return fmt.Errorf("ack DLQ message: %w", err)
		}
		log.Printf("level=info event=dlq_message_removed message_id=%q", delivery.MessageId)
		return nil
	}

	if err := delivery.Nack(false, true); err != nil {
		return fmt.Errorf("return message to DLQ: %w", err)
	}
	log.Printf("level=info event=dlq_message_returned message_id=%q", delivery.MessageId)
	return nil
}
