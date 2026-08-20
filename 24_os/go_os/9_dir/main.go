package main

import (
	"fmt"
	"os"
)

func main() {
	err := os.Mkdir("mydirectory", 0755)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Каталог 'mydirectory' создан успешно.")

	file, err := os.OpenInRoot("/mnt", "wsl")
	defer file.Close()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(file.Name())
}
