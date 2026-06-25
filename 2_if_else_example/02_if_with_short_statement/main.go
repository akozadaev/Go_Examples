package main

import (
	"fmt"
	"os"
)

func main() {
	if file, err := os.Open("main.go"); err != nil {
		fmt.Println("open error:", err)
	} else {
		defer file.Close()
		fmt.Println("file opened")
	}
}
