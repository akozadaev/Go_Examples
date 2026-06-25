package main

import "fmt"

func check() bool {
	fmt.Println("check called")
	return true
}

func main() {
	if false && check() {
		fmt.Println("inside if")
	}

	fmt.Println("after if")
}
