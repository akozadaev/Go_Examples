package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"
	"time"
)

func run(ctx context.Context) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Рабочая горутина остановлена")
			return
		case <-ticker.C:
			fmt.Println("Приложение работает")
		}
	}
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	done := make(chan struct{})
	go func() {
		defer close(done)
		run(ctx)
	}()

	fmt.Println("Приложение запущено. Нажмите Ctrl+C")
	<-ctx.Done()
	stop() // прекращаем доставку сигналов и освобождаем связанные ресурсы
	fmt.Println("Получен сигнал, ожидаем завершения")

	shutdownTimer := time.NewTimer(3 * time.Second)
	defer shutdownTimer.Stop()
	select {
	case <-done:
		fmt.Println("Приложение корректно завершено")
	case <-shutdownTimer.C:
		fmt.Println("Истёк таймаут завершения")
	}
}
