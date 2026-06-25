package main

import "fmt"

func main() {
Loop:
	for i := 1; i <= 2; i++ {
		for j := 1; j <= 3; j++ {
			if j == 2 {
				continue Loop
			}
			fmt.Printf("%d,%d ", i, j)
		}
	}

	fmt.Println()
}
