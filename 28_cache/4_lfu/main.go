package main

import "fmt"

func main() {
	cache := NewLFUCache(2)

	cache.Add("a", 1)
	cache.Add("b", 2)

	cache.Get("a") // повышаем частоту ключа "a"

	cache.Add("c", 3) // вытеснит "b", так как у него самая низкая частота

	fmt.Println(cache.Get("a")) // 1 true
	fmt.Println(cache.Get("b")) // <nil> false
	fmt.Println(cache.Get("c")) // 3 true

	cache.Remove("a")
	fmt.Println(cache.Size()) // 1
}
