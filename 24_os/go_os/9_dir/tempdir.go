package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	tempDir, err := os.MkdirTemp("", "app-cache-*")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка создания временного каталога: %v\n", err)
		return
	}
	defer os.RemoveAll(tempDir)

	fmt.Printf("Создан временный каталог: %s\n", tempDir)
	jsonPath := filepath.Join(tempDir, "data.json")

	jsonData := `{"status": "ok"}`
	if err := os.WriteFile(jsonPath, []byte(jsonData), 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка записи файла: %v\n", err)
		return
	}

	content, err := os.ReadFile(jsonPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка чтения файла: %v\n", err)
		return
	}
	fmt.Printf("Содержимое data.json: %s\n", content)
}
