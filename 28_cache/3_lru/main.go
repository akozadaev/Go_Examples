package main

import (
	"fmt"

	"github.com/hashicorp/golang-lru/v2"
)

func main() {
	cache, _ := lru.New[string, int](3)

	cache.Add("a", 1)
	cache.Add("b", 2)
	cache.Add("c", 3)

	fmt.Println(cache.Get("a")) // (1, true)
	cache.Add("d", 4)           // "b" будет удален

	fmt.Println(cache.Get("b")) // (0, false)
}
