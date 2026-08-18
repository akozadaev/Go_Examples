package main

import (
	"context"
	"testing"
)

func TestPipeline(t *testing.T) {
	ctx := context.Background()
	var got []int
	for value := range square(ctx, square(ctx, generate(ctx, 2, 3))) {
		got = append(got, value)
	}
	want := []int{16, 81}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got %v, want %v", got, want)
		}
	}
}
