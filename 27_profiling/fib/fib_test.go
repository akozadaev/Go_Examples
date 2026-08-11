package fib

import "testing"

var result uint64

func TestImplementations(t *testing.T) {
	for n := range uint(20) {
		if slow, fast := Slow(n), Fast(n); slow != fast {
			t.Fatalf("n=%d: Slow=%d Fast=%d", n, slow, fast)
		}
	}
}

func BenchmarkSlow20(b *testing.B) {
	for b.Loop() {
		result = Slow(20)
	}
}

func BenchmarkFast20(b *testing.B) {
	for b.Loop() {
		result = Fast(20)
	}
}
