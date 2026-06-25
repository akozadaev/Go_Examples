package main

import "fmt"

func main() {
Outer:
	for i := 1; i <= 3; i++ {
		for j := 1; j <= 3; j++ {
			if i*j > 4 {
				break Outer
			}
			fmt.Printf("%d,%d ", i, j)
		}
	}

	fmt.Println()
}
