.PHONY: rabbit-up rabbit-down test integration-test race build run-consumer publish dlq

rabbit-up:
	docker compose up -d --wait

rabbit-down:
	docker compose down

test:
	go test ./...

integration-test:
	go test -tags=integration ./internal/messaging

race:
	go test -race ./...

build:
	go build ./...

run-consumer:
	go run ./cmd/consumer

publish:
	go run ./cmd/publisher

dlq:
	go run ./cmd/dlq-reader
