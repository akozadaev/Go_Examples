package main

import (
	"fmt"
	"runtime/metrics"
	"sync"
)

var schedulerMetrics = []string{
	"/sched/gomaxprocs:threads",
	"/sched/goroutines:goroutines",
	"/sched/goroutines/runnable:goroutines",
	"/sched/goroutines/running:goroutines",
	"/sched/goroutines/waiting:goroutines",
	"/sched/goroutines/not-in-go:goroutines",
	"/sched/goroutines-created:goroutines",
	"/sched/threads/total:threads",
}

func readMetrics() {
	samples := make([]metrics.Sample, len(schedulerMetrics))
	for i, name := range schedulerMetrics {
		samples[i].Name = name
	}
	metrics.Read(samples)

	for _, sample := range samples {
		if sample.Value.Kind() != metrics.KindUint64 {
			fmt.Printf("%-46s unavailable\n", sample.Name)
			continue
		}
		fmt.Printf("%-46s %d\n", sample.Name, sample.Value.Uint64())
	}
}

func main() {
	release := make(chan struct{})
	var started sync.WaitGroup
	started.Add(100)
	for range 100 {
		go func() {
			started.Done()
			<-release
		}()
	}
	started.Wait()

	readMetrics()
	close(release)
}
