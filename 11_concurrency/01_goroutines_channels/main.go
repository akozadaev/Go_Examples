package main

import (
	"fmt"
	"time"
)

type Result struct {
	Worker string
	Value  int
}

func main() {
	fmt.Println("1. Горутина запускается через go, а канал передает результат")

	results := make(chan Result)

	go calculateSquare("worker-1", 7, results)

	fmt.Println("main: жду результат из канала")
	result := <-results
	fmt.Printf("main: получил %s -> %d\n", result.Worker, result.Value)

	fmt.Println()
	fmt.Println("2. Канал можно использовать как сигнал завершения")

	done := make(chan struct{})
	go backgroundJob(done)

	<-done
	fmt.Println("main: фоновая работа завершена")

	fmt.Println()
	fmt.Println("3. Закрытый канал можно читать: значение + ok")
	closed := make(chan int, 1)
	closed <- 10
	close(closed)

	value, ok := <-closed
	fmt.Printf("первое чтение: value=%d ok=%t\n", value, ok)
	value, ok = <-closed
	fmt.Printf("второе чтение: value=%d ok=%t\n", value, ok)
}

func calculateSquare(worker string, n int, out chan<- Result) {
	fmt.Println(worker + ": начал вычисление")
	time.Sleep(100 * time.Millisecond)
	out <- Result{Worker: worker, Value: n * n}
}

func backgroundJob(done chan<- struct{}) {
	defer close(done)

	for step := 1; step <= 3; step++ {
		fmt.Printf("background: шаг %d\n", step)
		time.Sleep(50 * time.Millisecond)
	}
}
