package main

import (
	"context"
	"sort"
	"testing"
)

func TestFanOutProcessesEveryJobOnce(t *testing.T) {
	ctx := context.Background()
	var got []int
	for item := range fanOut(ctx, 3, generate(ctx, 5)) {
		got = append(got, item.value)
	}
	sort.Ints(got)
	want := []int{1, 4, 9, 16, 25}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got %v, want %v", got, want)
		}
	}
}
