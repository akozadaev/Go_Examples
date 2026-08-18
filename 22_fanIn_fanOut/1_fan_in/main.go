package main

import (
	"context"
	"fmt"
	"sync"
)

// fanIn объединяет несколько каналов в один. Результирующий канал закрывается
// только после завершения всех копирующих горутин.
func fanIn(ctx context.Context, inputs ...<-chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup
	wg.Add(len(inputs))

	for _, input := range inputs {
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case value, ok := <-input:
					if !ok {
						return
					}
					select {
					case <-ctx.Done():
						return
					case out <- value:
					}
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func produce(ctx context.Context, values ...int) <-chan int {
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

func main() {
	ctx := context.Background()

	first := produce(ctx, 1, 2, 3)
	second := produce(ctx, 10, 20, 30)

	for value := range fanIn(ctx, first, second) {
		fmt.Println(value)
	}
}
