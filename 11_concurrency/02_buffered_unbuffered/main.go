package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	fmt.Println("1. Небуферизированный канал: отправитель ждет получателя")
	unbuffered()

	fmt.Println()
	fmt.Println("2. Буферизированный канал: отправка проходит, пока есть место в буфере")
	buffered()

	fmt.Println()
	fmt.Println("3. Буфер не делает обработку параллельной сам по себе")
	bufferAsQueue()
}

func unbuffered() {
	ch := make(chan string)

	go func() {
		time.Sleep(100 * time.Millisecond)
		fmt.Println("   получатель готов")
		fmt.Println("   получил:", <-ch)
	}()

	fmt.Println("   отправитель пытается отправить")
	ch <- "сообщение"
	fmt.Println("   отправитель продолжил работу только после приема")
}

func buffered() {
	ch := make(chan string, 2)

	fmt.Println("   отправляю первое значение")
	ch <- "первое"
	fmt.Println("   отправляю второе значение")
	ch <- "второе"
	fmt.Println("   оба значения помещены в буфер без отдельного получателя")

	fmt.Println("   читаю:", <-ch)
	fmt.Println("   читаю:", <-ch)

	fmt.Println("   если буфер заполнен, следующая отправка снова будет ждать получателя")
	ch <- "буфер-1"
	ch <- "буфер-2"

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		ch <- "третье"
		fmt.Println("   третье значение отправлено после освобождения места")
	}()

	time.Sleep(50 * time.Millisecond)
	fmt.Println("   освобождаю место в буфере чтением:", <-ch)
	wg.Wait()
	fmt.Println("   дочитываю:", <-ch)
	fmt.Println("   дочитываю:", <-ch)
}

func bufferAsQueue() {
	jobs := make(chan int, 3)

	for job := 1; job <= 3; job++ {
		jobs <- job
		fmt.Printf("   задача %d поставлена в очередь\n", job)
	}
	close(jobs)

	for job := range jobs {
		fmt.Printf("   задача %d обработана последовательно\n", job)
	}
}
