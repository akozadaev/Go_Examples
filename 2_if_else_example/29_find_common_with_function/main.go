package main

import "fmt"

func findCommon(a, b []int) (int, bool) {
	for _, x := range a {
		for _, y := range b {
			if x == y {
				return x, true
			}
		}
	}

	return 0, false
}

func main() {
	a := []int{1, 2, 3}
	b := []int{4, 5, 3}

	if value, ok := findCommon(a, b); ok {
		fmt.Println("found:", value)
	} else {
		fmt.Println("not found")
	}
}
