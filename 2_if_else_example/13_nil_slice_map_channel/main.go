package main

import "fmt"

func main() {
	var nums []int
	var dict map[string]int

	fmt.Println("len nums:", len(nums))
	fmt.Println("len dict:", len(dict))

	nums = append(nums, 10)
	fmt.Println("nums:", nums)

	if dict != nil {
		dict["x"] = 1
	} else {
		fmt.Println("cannot write to nil map")
	}
}
