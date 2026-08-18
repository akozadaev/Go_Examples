package main

import "testing"

func TestFIFOEvictsOldestInserted(t *testing.T) {
	cache := NewFIFOCache(2)
	cache.Add("a", 1)
	cache.Add("b", 2)
	cache.Get("a")
	cache.Add("c", 3)

	if _, ok := cache.Get("a"); ok {
		t.Fatal("key a must be evicted despite the read")
	}
	if value, ok := cache.Get("b"); !ok || value != 2 {
		t.Fatalf("key b must remain: got (%v, %v)", value, ok)
	}
}
