package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	fmt.Println("Приложение запущено. Для завершения нажмите Ctrl+C или отправьте SIGTERM.")
	<-ctx.Done()

	// Начиная с Go 1.26 NotifyContext сохраняет полученный сигнал как причину отмены.
	fmt.Printf("Получен сигнал завершения: %v\n", context.Cause(ctx))
	fmt.Println("Завершаем приложение корректно.")

	// Возврат из main выполняет defer. Вызов os.Exit здесь не нужен.
}
