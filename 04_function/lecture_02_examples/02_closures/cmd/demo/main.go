package main

import (
	"fmt"

	"example.com/functions/part2/closures/closurekit"
)

func main() {
	counter := closurekit.NewCounter()
	fmt.Println("counter:", counter(), counter(), counter())

	double := closurekit.MakeMultiplier(2)
	fmt.Println("mapped:", closurekit.Map([]int{1, 2, 3}, double))
}
