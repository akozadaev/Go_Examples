package main

import (
	"fmt"
	"os"
)

func main() {
	err := os.Rename("tmp4247572260", "example.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Файл переименован.")
}
