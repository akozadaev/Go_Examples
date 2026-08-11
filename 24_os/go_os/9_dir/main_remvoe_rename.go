package main

import (
	"fmt"
	"os"
)

func main() {
	err := os.Rename("mydirectory", "newdirectory")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Каталог 'mydirectory' переименован в 'newdirectory'.")
}
