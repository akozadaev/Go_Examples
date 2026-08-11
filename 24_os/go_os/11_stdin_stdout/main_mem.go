package main

import (
	"bytes"
	"fmt"
	"io"
)

func main() {
	data := []byte("Hello, World!")
	reader := bytes.NewReader(data)
	writer := &bytes.Buffer{}

	n, err := io.Copy(writer, reader)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Скопировано %d байт: %s\n", n, writer.String())
}
