package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

func cpuWork(id int, start <-chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()
	<-start

	deadline := time.Now().Add(300 * time.Millisecond)
	var iterations uint64
	for time.Now().Before(deadline) {
		iterations++
	}
	fmt.Printf("worker %d: %d iterations\n", id, iterations)
}

func main() {
	fmt.Println("NumCPU:", runtime.NumCPU())
	fmt.Println("GOMAXPROCS:", runtime.GOMAXPROCS(0))

	workers := runtime.GOMAXPROCS(0) * 2
	start := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(workers)
	for id := range workers {
		go cpuWork(id+1, start, &wg)
	}
	close(start)
	wg.Wait()
}
