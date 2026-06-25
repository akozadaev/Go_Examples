package main

import "fmt"

func main() {
	nums := []int{10, 20, 30}

	for index, value := range nums {
		fmt.Println(index, value)
	}

	for _, value := range nums {
		fmt.Println("value:", value)
	}

	for index, _ := range nums {
		fmt.Println("index:", index)
	}
}
