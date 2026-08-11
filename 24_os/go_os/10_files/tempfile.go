package main

import (
	"fmt"
	"os"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	tempFile, err := os.CreateTemp("", "log-*.txt")
	if err != nil {
		return fmt.Errorf("создание временного файла: %w", err)
	}
	name := tempFile.Name()
	defer os.Remove(name)

	fmt.Printf("Создан временный файл: %s\n", name)
	if _, err := tempFile.WriteString("Это временный лог-файл."); err != nil {
		tempFile.Close()
		return fmt.Errorf("запись во временный файл: %w", err)
	}

	// Close выполняется сейчас, а не через defer: только после него файл читается заново.
	if err := tempFile.Close(); err != nil {
		return fmt.Errorf("закрытие временного файла: %w", err)
	}

	content, err := os.ReadFile(name)
	if err != nil {
		return fmt.Errorf("чтение временного файла: %w", err)
	}
	fmt.Printf("Содержимое файла: %s\n", content)
	return nil
}
