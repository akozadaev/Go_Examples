package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
)

func main() {
	fileInfo, err := os.Stat("example.txt")
	switch {
	case err == nil:
		fmt.Println(fileInfo.Name())
		fmt.Println(fileInfo.Size())
		fmt.Println("Файл существует.")
	case errors.Is(err, fs.ErrNotExist):
		fmt.Println("Файл не существует.")
	default:
		fmt.Fprintln(os.Stderr, "Не удалось проверить файл:", err)
	}
}
