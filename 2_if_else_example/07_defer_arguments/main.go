package main

import "fmt"

func main() {
	x := 1
	//defer fmt.Println("defer:", x)

	x = 2
	defer func(x int) {
		fmt.Println("defer:", x)
	}(x)

	fmt.Println("main:", x)
}
