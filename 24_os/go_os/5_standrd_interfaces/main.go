package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	reader := strings.NewReader("Hello, World!")
	writer := &strings.Builder{}

	if _, err := io.Copy(writer, reader); err != nil {
		log.Fatal(err)
	}
	fmt.Println(writer.String())

	// strings.Reader не владеет внешним ресурсом и не реализует io.Closer.
	var r io.Reader = strings.NewReader("Some resource to close")
	if closer, ok := r.(io.Closer); ok {
		if err := closer.Close(); err != nil {
			log.Fatal(err)
		}
		fmt.Println("Resource closed")
	} else {
		fmt.Println("This reader does not need to be closed")
	}

	if err := example(); err != nil {
		log.Fatal(err)
	}
}

func example() (err error) {
	file, err := os.Open("example.txt")
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := file.Close(); err == nil {
			err = closeErr
		}
	}()
	return nil
}
