package main

import "fmt"

func main() {
	var p *int
	//p := new(3)

	if p != nil && *p > 0 {
		fmt.Println("positive")
	} else {
		fmt.Println("nil or non-positive")
	}
}
