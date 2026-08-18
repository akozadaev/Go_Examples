package main

import (
	"fmt"
	"runtime"
	"sync/atomic"
	"time"
)

func main() {
	previous := runtime.GOMAXPROCS(1)
	defer runtime.GOMAXPROCS(previous)

	var stop atomic.Bool
	done := make(chan struct{})
	go func() {
		defer close(done)
		var n uint64
		for !stop.Load() {
			n++
		}
		fmt.Println("CPU-bound goroutine stopped; iterations:", n)
	}()

	// Даже при одном P планировщик должен дать main возможность продолжить
	// работу и обработать таймер. Конкретный момент переключения не гарантирован.
	time.Sleep(100 * time.Millisecond)
	fmt.Println("main продолжил выполнение при GOMAXPROCS=1")
	stop.Store(true)
	<-done
}
