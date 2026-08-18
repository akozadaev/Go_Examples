package main

import (
	"context"
	"fmt"
	"sync"
)

type result struct {
	worker int
	value  int
}

// fanOut запускает workers конкурирующих обработчиков одного канала заданий.
// Каждое задание получает ровно один воркер — это распределение работы,
// а не широковещательная рассылка одного значения всем получателям.
func fanOut(ctx context.Context, workers int, jobs <-chan int) <-chan result {
	if workers < 1 {
		workers = 1
	}

	out := make(chan result)
	var wg sync.WaitGroup
	wg.Add(workers)

	for id := 1; id <= workers; id++ {
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case job, ok := <-jobs:
					if !ok {
						return
					}
					select {
					case <-ctx.Done():
						return
					case out <- result{worker: id, value: job * job}:
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

func generate(ctx context.Context, count int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for value := 1; value <= count; value++ {
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
	for item := range fanOut(ctx, 3, generate(ctx, 9)) {
		fmt.Printf("worker %d: %d\n", item.worker, item.value)
	}
}
