package main

import "testing"

func TestLFUEvictsLeastFrequentlyUsed(t *testing.T) {
	cache := NewLFUCache(2)
	cache.Add("a", 1)
	cache.Add("b", 2)
	cache.Get("a")
	cache.Add("c", 3)

	if _, ok := cache.Get("b"); ok {
		t.Fatal("key b must be evicted")
	}
	if value, ok := cache.Get("a"); !ok || value != 1 {
		t.Fatalf("key a must remain: got (%v, %v)", value, ok)
	}
}
