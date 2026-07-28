package main

import (
	"log/slog"
	"os"
)

func main() {
	// Человекочитаемый формат (по умолчанию)
	slog.Info("Привет от slog", "user_id", "usr-1234", "doc_id", "doc-xyz")

	// JSON-формат
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logger.Info("Привет от slog в JSON", "user_id", "usr-1234")
}
