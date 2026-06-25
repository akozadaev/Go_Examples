package main

import "fmt"

func main() {
	a := []int{1, 2, 3}
	b := []int{4, 5, 3, 2}

	found := false

Search:
	for _, x := range a {
		for _, y := range b {
			if x == y {
				fmt.Println("found:", x)
				found = true
				break Search
			}
		}
	}

	if !found {
		fmt.Println("not found")
	}
}
