package main

import "fmt"

func main() {
	fmt.Println("break example")
	for i := 0; i < 10; i++ {
		if i == 3 {
			break
		}
		fmt.Println(i)
	}

	fmt.Println("continue example")
	for i := 1; i <= 10; i++ {
		if i%2 == 0 {
			continue
		}
		fmt.Println(i)
	}
}
