package main

import "testing"

func TestLazyFillWarmupAndInvalidation(t *testing.T) {
	repo := NewRepository(map[string]string{"a": "old", "b": "warm"})
	service := NewService(repo, NewCache())

	if err := service.Warm("b"); err != nil {
		t.Fatal(err)
	}
	readsAfterWarmup := repo.reads
	if value, err := service.Get("b"); err != nil || value != "warm" {
		t.Fatalf("warm cache hit: value=%q err=%v", value, err)
	}
	if repo.reads != readsAfterWarmup {
		t.Fatal("cache hit unexpectedly read repository")
	}

	if value, err := service.Get("a"); err != nil || value != "old" {
		t.Fatalf("lazy cache miss: value=%q err=%v", value, err)
	}
	readsAfterMiss := repo.reads
	if _, err := service.Get("a"); err != nil {
		t.Fatal(err)
	}
	if repo.reads != readsAfterMiss {
		t.Fatal("second read must be a cache hit")
	}

	service.Update("a", "new")
	if value, err := service.Get("a"); err != nil || value != "new" {
		t.Fatalf("read after invalidation: value=%q err=%v", value, err)
	}
	if repo.reads != readsAfterMiss+1 {
		t.Fatal("read after invalidation must reach repository")
	}
}
