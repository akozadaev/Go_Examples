package main

import (
	"fmt"
	"time"

	"github.com/patrickmn/go-cache"
)

func main() {
	const ttl = 2 * time.Second
	c := cache.New(ttl, time.Second)

	c.Set("key", "value", cache.DefaultExpiration)
	if val, found := c.Get("key"); found {
		fmt.Println("cache hit:", val)
	}

	fmt.Println("waiting for TTL expiration...")
	time.Sleep(ttl + time.Second)
	if _, found := c.Get("key"); !found {
		fmt.Println("cache miss: value expired")
	}
}
