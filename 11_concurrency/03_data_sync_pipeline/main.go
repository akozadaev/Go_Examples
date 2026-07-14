package main

import (
	"fmt"
	"sync"
	"time"
)

type Job struct {
	ID    int
	Value int
}

type ProcessedJob struct {
	ID     int
	Result int
}

func main() {
	fmt.Println("Pipeline: генератор -> workers -> сборщик")

	jobs := generateJobs([]int{2, 3, 4, 5, 6})
	results := runWorkers(3, jobs)

	total := 0
	for result := range results {
		fmt.Printf("результат задачи %d: %d\n", result.ID, result.Result)
		total += result.Result
	}

	fmt.Println("сумма результатов:", total)
}

func generateJobs(values []int) <-chan Job {
	out := make(chan Job)

	go func() {
		defer close(out)
		for i, value := range values {
			out <- Job{ID: i + 1, Value: value}
		}
	}()

	return out
}

func runWorkers(count int, jobs <-chan Job) <-chan ProcessedJob {
	results := make(chan ProcessedJob)
	var wg sync.WaitGroup

	for workerID := 1; workerID <= count; workerID++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for job := range jobs {
				fmt.Printf("worker %d обрабатывает задачу %d\n", id, job.ID)
				time.Sleep(50 * time.Millisecond)
				results <- ProcessedJob{ID: job.ID, Result: job.Value * job.Value}
			}
		}(workerID)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	return results
}
