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
}
