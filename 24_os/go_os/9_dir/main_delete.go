package main

import (
	"fmt"
	"os"
)

func main() {
	err := os.Remove("mydirectory")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Каталог 'mydirectory' удален успешно.")
}
