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

	pos, err := file.Seek(0, io.SeekStart)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	fmt.Println("Position:", pos)

	data := make([]byte, 10)
	n, err := file.Read(data)
	if err != nil && !errors.Is(err, io.EOF) {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	fmt.Printf("Read %d bytes: %s\n", n, data[:n])
}
