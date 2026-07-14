package main

import (
	"context"
	"testing"
	"time"
)

func TestStreamNumbersClosesChannel(t *testing.T) {
	got := make([]int, 0, 3)
	for value := range streamNumbers(3, time.Millisecond) {
		got = append(got, value)
	}

	if len(got) != 3 {
		t.Fatalf("получено значений: %d, ожидается 3", len(got))
	}
}

func TestSlowOperationRespectsContext(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()

	_, err := slowOperation(ctx)
	if err == nil {
		t.Fatal("ожидалась ошибка отмены context")
	}
}
