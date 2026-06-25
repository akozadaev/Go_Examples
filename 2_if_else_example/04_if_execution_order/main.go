package main

import "fmt"

func main() {
	if x := 10; x > 5 {
		fmt.Println("branch if")
	} else {
		fmt.Println("branch else")
	}

	fmt.Println("after if")
}
