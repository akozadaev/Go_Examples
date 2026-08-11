package main

import (
	"fmt"
	"os"
)

func main() {
	files, err := os.ReadDir("mydirectory")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Файлы в каталоге 'mydirectory':")
	for _, file := range files {
		fmt.Println(file.Name())
	}
}
