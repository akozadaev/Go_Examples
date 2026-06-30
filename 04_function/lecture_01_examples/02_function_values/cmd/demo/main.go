package main

import (
	"fmt"

	"example.com/functions/part1/functionvalues/transform"
)

func main() {
	result := transform.Pipeline(
		"  Go Functions  ",
		transform.Trim,
		transform.Lower,
		transform.Prefix("topic: "),
		func(s string) string {
			return s + "."
		},
	)

	fmt.Println(result)
}
