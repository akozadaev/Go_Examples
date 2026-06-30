package main

import (
	"fmt"
	"log"

	"example.com/functions/part2/scopeshadowing/scopecheck"
)

func main() {
	values, err := scopecheck.ParseCSVInts("-1, 2, 3, 0, 4")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("buggy count:", scopecheck.CountPositiveWithBug(values))
	fmt.Println("correct count:", scopecheck.CountPositive(values))
	fmt.Println("label:", scopecheck.FindLabel(map[string]string{"go": "Golang"}, "go"))
}
