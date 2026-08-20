package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() (err error) {
	file, err := os.Create("output.txt")
	if err != nil {
		return fmt.Errorf("create output.txt: %w", err)
	}
	defer func() {
		if closeErr := file.Close(); err == nil {
			err = closeErr
		}
	}()

	writer := bufio.NewWriter(file)
	n, err := writer.WriteString("Hello, World!")
	if err != nil {
		return fmt.Errorf("buffer write: %w", err)
	}
	if err := writer.Flush(); err != nil {
		return fmt.Errorf("flush to file: %w", err)
	}

	fmt.Printf("Written %d bytes\n", n)

	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 0, 4096), 4096)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
	return nil

}
