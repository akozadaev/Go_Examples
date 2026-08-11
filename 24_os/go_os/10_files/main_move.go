package main

import (
	"fmt"
	"os"
)

func main() {
	err := os.Rename("example.txt", "destination/folder/sourcefile.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Файл перемещен.")
}
