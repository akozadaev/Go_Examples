package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
)

func main() {
	info, err := os.Stat("mydirectory")
	switch {
	case err == nil:
		fmt.Printf("'mydirectory' существует; каталог: %t\n", info.IsDir())
	case errors.Is(err, fs.ErrNotExist):
		fmt.Println("Каталог 'mydirectory' не существует.")
	default:
		fmt.Fprintln(os.Stderr, "Не удалось проверить каталог:", err)
	}
}
