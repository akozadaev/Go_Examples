package main

import (
	"context"
	"fmt"
)

func generate(ctx context.Context, values ...int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for _, value := range values {
			select {
			case <-ctx.Done():
				return
			case out <- value:
			}
		}
	}()
	return out
}

func square(ctx context.Context, in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for {
			select {
			case <-ctx.Done():
				return
			case value, ok := <-in:
				if !ok {
					return
				}
				select {
				case <-ctx.Done():
					return
				case out <- value * value:
				}
			}
		}
	}()
	return out
}

func main() {
	ctx := context.Background()

	// Каждый этап принимает канал и возвращает новый канал, поэтому этапы
	// свободно компонуются: generate -> square -> square -> main.
	for value := range square(ctx, square(ctx, generate(ctx, 2, 3, 4))) {
		fmt.Println(value) // 16, 81, 256
	}
}
