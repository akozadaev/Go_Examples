package main

import (
	"fmt"
	"os"
)

func main() {
	err := os.Remove("sourcefile.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Файл удален успешно.")
}
