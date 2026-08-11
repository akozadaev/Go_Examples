package main

import (
	"errors"
	"fmt"
	"io"
	"os"
)

func main() {
	file, err := os.Open("example.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	defer file.Close()

	data := make([]byte, 1024)
	n, err := file.Read(data)
	if err != nil && !errors.Is(err, io.EOF) {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	fmt.Printf("Прочитано %d байт: %s\n", n, data[:n])
}
