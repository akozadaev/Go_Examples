package main

import (
	"fmt"
	"time"

	"github.com/akozadaev/go-code-generation/impl_example"
)

func main() {
	// Демонстрация использования impl
	fmt.Println("=== Демонстрация impl ===")

	// Используем реализацию, созданную с помощью impl
	logger := &impl_example.FileLogger{
		Filename: "app.log",
	}

	// Тестируем методы интерфейса Logger
	logger.Log("Приложение запущено")
	logger.Logf("Текущее время: %s", time.Now().Format(time.RFC3339))
	logger.Error("Это тестовая ошибка")
	logger.Errorf("Ошибка в функции %s с кодом %d", "processData", 500)

	fmt.Println("Логи записаны в app.log")
}

