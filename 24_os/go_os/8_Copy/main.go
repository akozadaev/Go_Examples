package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() (err error) {
	sourceFile, err := os.Open("source.txt")
	if err != nil {
		return fmt.Errorf("open source: %w", err)
	}
	defer sourceFile.Close()

	destFile, err := os.Create("destination.txt")
	if err != nil {
		return fmt.Errorf("create destination: %w", err)
	}
	defer func() {
		if closeErr := destFile.Close(); err == nil {
			err = closeErr
		}
	}()

	n, err := io.Copy(destFile, sourceFile)
	if err != nil {
		return fmt.Errorf("copy: %w", err)
	}
	fmt.Printf("Copied %d bytes\n", n)
	return nil
}
