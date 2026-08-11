package main

import (
	"fmt"
	"os"
)

func main() {
	data := []byte("Hello")
	n, err := os.Stdout.Write(data)
	if err != nil {
		fmt.Fprintln(os.Stderr, "write stdout:", err)
		return
	}
	fmt.Printf("\nЗаписали %d байт\n", n)
}
