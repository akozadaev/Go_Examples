package main

import (
	"fmt"

	"github.com/akozadaev/go_cache/4_lfu/custom_thread_safe"
)

func main() {
	cache := custom_thread_safe.NewLFUCache[string, int](2)

	cache.Add("a", 1)
	cache.Add("b", 2)

	cache.Get("a") // увеличиваем частоту "a"

	cache.Add("c", 3) // вытеснит "b"

	fmt.Println(cache.Get("a")) // 1 true
	fmt.Println(cache.Get("b")) // 0 false
	fmt.Println(cache.Get("c")) // 3 true
	fmt.Println(cache.Size())   // 2
}
