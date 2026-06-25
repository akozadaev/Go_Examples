package main

import "fmt"

func main() {
	n := 1

	switch n {
	case 1:
		fmt.Println("one")
		fallthrough
	case 2:
		fmt.Println("two")
	default:
		fmt.Println("other")
	}
}
