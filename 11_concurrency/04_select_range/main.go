package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	fmt.Println("1. range читает канал до close")
	for value := range streamNumbers(3, 40*time.Millisecond) {
		fmt.Println("   range получил:", value)
	}

	fmt.Println()
	fmt.Println("2. select выбирает готовый канал и умеет ждать timeout")
	selectWithTimeout()

	fmt.Println()
	fmt.Println("3. default делает select неблокирующим")
	nonBlockingReceive()

	fmt.Println()
	fmt.Println("4. context отменяет долгую работу")
	cancelWithContext()
}

func streamNumbers(count int, delay time.Duration) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)
		for i := 1; i <= count; i++ {
			time.Sleep(delay)
			out <- i
		}
	}()

	return out
}

func selectWithTimeout() {
	fast := streamNumbers(2, 30*time.Millisecond)
	slow := streamNumbers(3, 90*time.Millisecond)
	timeout := time.After(220 * time.Millisecond)

	for fast != nil || slow != nil {
		select {
		case value, ok := <-fast:
			if !ok {
				fast = nil
				continue
			}
			fmt.Println("   fast:", value)
		case value, ok := <-slow:
			if !ok {
				slow = nil
				continue
			}
			fmt.Println("   slow:", value)
		case <-timeout:
			fmt.Println("   timeout: прекращаю ожидание")
			return
		}
	}
}

func nonBlockingReceive() {
	ch := make(chan string)

	select {
	case message := <-ch:
		fmt.Println("   получили:", message)
	default:
		fmt.Println("   данных сейчас нет, поэтому не блокируемся")
	}
}

func cancelWithContext() {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Millisecond)
	defer cancel()

	result, err := slowOperation(ctx)
	if err != nil {
		fmt.Println("   операция остановлена:", err)
		return
	}
	fmt.Println("   результат:", result)
}

func slowOperation(ctx context.Context) (string, error) {
	select {
	case <-time.After(300 * time.Millisecond):
		return "готово", nil
	case <-ctx.Done():
		return "", ctx.Err()
	}
}
