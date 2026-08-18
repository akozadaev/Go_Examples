package main

import (
	"runtime"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	workers := runtime.GOMAXPROCS(0) * 4
	wg.Add(workers)
	for range workers {
		go func() {
			defer wg.Done()
			deadline := time.Now().Add(1500 * time.Millisecond)
			var n uint64
			for time.Now().Before(deadline) {
				n++
				if n%1_000_000 == 0 {
					runtime.Gosched()
				}
			}
		}()
	}
	wg.Wait()
}
