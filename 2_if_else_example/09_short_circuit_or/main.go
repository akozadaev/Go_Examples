package main

import "fmt"

func fallback() bool {
	fmt.Println("fallback called")
	return true
}

func main() {
	if true || fallback() {
		fmt.Println("allowed")
	}
}
