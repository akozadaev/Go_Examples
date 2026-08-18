package main

import (
	"fmt"
	"sync"
)

type Cache struct {
	data sync.Map // потокобезопасная map
}

func (c *Cache) Set(key string, value any) {
	c.data.Store(key, value)
}

func (c *Cache) Get(key string) (any, bool) {
	return c.data.Load(key)
}

func main() {
	cache := &Cache{}
	cache.Set("user:1", "Алексей")
	if val, ok := cache.Get("user:1"); ok {
		fmt.Println("cache hit:", val)
	}
	_, ok := cache.Get("user:2")
	fmt.Println("cache miss:", !ok)
}
