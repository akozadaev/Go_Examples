package main

import "testing"

func TestCacheHitAndMiss(t *testing.T) {
	cache := &Cache{}
	cache.Set("user:1", "Алексей")

	value, ok := cache.Get("user:1")
	if !ok || value != "Алексей" {
		t.Fatalf("cache hit: got (%v, %v)", value, ok)
	}
	if _, ok := cache.Get("missing"); ok {
		t.Fatal("cache miss: key unexpectedly found")
	}
}
