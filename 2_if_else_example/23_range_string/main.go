package main

import "fmt"

func main() {
	for index, r := range "Go 世界" {
		fmt.Printf("%d %c\n", index, r)
	}
}
