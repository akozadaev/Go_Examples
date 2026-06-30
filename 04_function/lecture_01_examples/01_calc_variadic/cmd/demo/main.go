package main

import (
	"fmt"
	"log"

	"example.com/functions/part1/calcvariadic/calc"
)

func main() {
	fmt.Println("sum:", calc.Sum(1, 2, 3, 4))

	avg, err := calc.Avg(2, 4, 6)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("avg:", avg)

	increment := func(x int) int {
		return x + 1
	}
	fmt.Println("apply twice:", calc.ApplyTwice(increment, 10))
}
