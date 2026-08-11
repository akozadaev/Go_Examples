package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"time"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "Ошибка выполнения команды:", err)
		os.Exit(1)
	}
}

func run() error {
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("get working directory: %w", err)
	}
	fmt.Println("Текущая рабочая директория:", wd)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Команда go version переносимее учебного вызова Unix-команды ls.
	cmd := exec.CommandContext(ctx, "go", "version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return fmt.Errorf("go version timed out: %w", ctx.Err())
		}
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return fmt.Errorf("go version exited with code %d: %s", exitErr.ExitCode(), output)
		}
		return fmt.Errorf("start go version: %w", err)
	}

	fmt.Print(string(output))
	return nil
}
