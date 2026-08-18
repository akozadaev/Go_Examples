package main

import (
	"context"
	"fmt"
	"time"
)

func operation(ctx context.Context, duration time.Duration) error {
	timer := time.NewTimer(duration)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	if err := operation(ctx, 2*time.Second); err != nil {
		fmt.Println("Операция завершена:", err)
		return
	}
	fmt.Println("Операция выполнена")
}
