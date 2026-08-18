package custom_thread_safe

import (
	"sync"
	"testing"
)

func TestLFUConcurrentAccess(t *testing.T) {
	cache := NewLFUCache[int, int](32)
	var wg sync.WaitGroup
	for worker := range 8 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for key := range 1_000 {
				cache.Add((key+worker)%64, key)
				cache.Get((key + worker) % 64)
			}
		}()
	}
	wg.Wait()
	if cache.Size() > 32 {
		t.Fatalf("size exceeds capacity: %d", cache.Size())
	}
}
