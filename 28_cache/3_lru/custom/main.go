package main

import "fmt"

func main() {
	var cache Cache = NewLRUCache(2)

	cache.Add("a", 1)
	cache.Add("b", 2)

	value, ok := cache.Get("a")
	fmt.Println("Get a:", value, ok) // Get a: 1 true

	cache.Add("c", 3) // "b" будет удален

	value, ok = cache.Get("b")
	fmt.Println("Get b:", value, ok) // Get b: <nil> false

	fmt.Println("Size:", cache.Size()) // Size: 2 Последовательно идущие единицы
}
