package main

import (
	"fmt"
	"os"
)

func main() {
	data := make([]byte, 1024)
	n, err := os.Stdin.Read(data)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Прочитано %d байт: %s\n", n, string(data[:n]))
}
