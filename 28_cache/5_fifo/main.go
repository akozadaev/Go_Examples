package main

import "fmt"

func main() {
	cache := NewFIFOCache(2)

	cache.Add("a", 1)
	cache.Add("b", 2)

	fmt.Println(cache.Get("a")) // 1 true

	cache.Add("c", 3) // вытеснит "a"

	fmt.Println(cache.Get("a")) // <nil> false
	fmt.Println(cache.Get("b")) // 2 true
	fmt.Println(cache.Get("c")) // 3 true

	cache.Remove("b")
	fmt.Println(cache.Size()) // 1
}
