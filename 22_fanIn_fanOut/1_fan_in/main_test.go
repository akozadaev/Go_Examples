package main

import (
	"context"
	"sort"
	"testing"
	"time"
)

func TestFanIn(t *testing.T) {
	ctx := context.Background()
	var got []int
	for value := range fanIn(ctx, produce(ctx, 1, 2), produce(ctx, 3, 4)) {
		got = append(got, value)
	}
	sort.Ints(got)
	want := []int{1, 2, 3, 4}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got %v, want %v", got, want)
		}
	}
}

func TestFanInCancellationClosesOutput(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	input := make(chan int)
	out := fanIn(ctx, input)
	cancel()

	select {
	case _, ok := <-out:
		if ok {
			t.Fatal("output must be closed")
		}
	case <-time.After(time.Second):
		t.Fatal("fanIn did not stop after cancellation")
	}
}
