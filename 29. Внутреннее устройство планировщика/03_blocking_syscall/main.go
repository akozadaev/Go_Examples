package main

import (
	"fmt"
	"os"
	"runtime"
	"time"
)

func main() {
	previous := runtime.GOMAXPROCS(1)
	defer runtime.GOMAXPROCS(previous)

	reader, writer, err := os.Pipe()
	if err != nil {
		panic(err)
	}
	defer reader.Close()
	defer writer.Close()

	done := make(chan struct{})
	go func() {
		defer close(done)
		buffer := make([]byte, 1)
		fmt.Println("reader: blocking read")
		if _, err := reader.Read(buffer); err != nil {
			fmt.Println("read error:", err)
			return
		}
		fmt.Println("reader: received", string(buffer))
	}()

	time.Sleep(100 * time.Millisecond)
	fmt.Println("main работает, хотя другая горутина ждёт системный вызов")
	if _, err := writer.Write([]byte("x")); err != nil {
		panic(err)
	}
	<-done
}
