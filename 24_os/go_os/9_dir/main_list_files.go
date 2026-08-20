package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

func main() {
	files, err := os.ReadDir(".")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Файлы в каталоге 'mydirectory':")
	for _, file := range files {
		fmt.Println(file.Name(), file.IsDir())
	}

	err = filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		fmt.Println(".", d.IsDir(), d.Name())
		return nil
	})
}
