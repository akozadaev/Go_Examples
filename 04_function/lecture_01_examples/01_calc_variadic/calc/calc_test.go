package calc

import (
	"errors"
	"testing"
)

func TestSum(t *testing.T) {
	tests := []struct {
		name string
		nums []int
		want int
	}{
		{name: "positive", nums: []int{1, 2, 3}, want: 6},
		{name: "empty", nums: nil, want: 0},
		{name: "with negatives", nums: []int{-2, 5, -1}, want: 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Sum(tt.nums...); got != tt.want {
				t.Fatalf("Sum() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestAvg(t *testing.T) {
	got, err := Avg(2, 4, 6)
	if err != nil {
		t.Fatalf("Avg() unexpected error: %v", err)
	}
	if got != 4 {
		t.Fatalf("Avg() = %v, want 4", got)
	}
}

func TestAvgEmpty(t *testing.T) {
	_, err := Avg()
	if !errors.Is(err, ErrEmptyInput) {
		t.Fatalf("Avg() error = %v, want ErrEmptyInput", err)
	}
}

func TestApplyTwice(t *testing.T) {
	double := func(x int) int {
		return x * 2
	}

	if got := ApplyTwice(double, 3); got != 12 {
		t.Fatalf("ApplyTwice() = %d, want 12", got)
	}

	if got := ApplyTwice(nil, 3); got != 3 {
		t.Fatalf("ApplyTwice(nil) = %d, want 3", got)
	}
}
