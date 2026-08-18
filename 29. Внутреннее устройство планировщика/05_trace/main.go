package main

import (
	"context"
	"fmt"
	"os"
	"runtime/trace"
	"sync"
	"time"
)

func main() {
	file, err := os.Create("trace.out")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	if err := trace.Start(file); err != nil {
		panic(err)
	}

	ctx, task := trace.NewTask(context.Background(), "worker-pool")
	region := trace.StartRegion(ctx, "parallel-work")
	var wg sync.WaitGroup
	for id := range 8 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			trace.WithRegion(ctx, "worker", func() {
				time.Sleep(time.Duration(id%3+1) * 10 * time.Millisecond)
			})
		}()
	}
	wg.Wait()
	region.End()
	task.End()
	trace.Stop()

	fmt.Println("Создан trace.out; откройте: go tool trace trace.out")
}
