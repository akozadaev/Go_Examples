package closurekit

import (
	"reflect"
	"testing"
)

func TestNewCounter(t *testing.T) {
	counter := NewCounter()

	for _, want := range []int{1, 2, 3} {
		if got := counter(); got != want {
			t.Fatalf("counter() = %d, want %d", got, want)
		}
	}
}

func TestCountersHaveIndependentState(t *testing.T) {
	first := NewCounter()
	second := NewCounter()

	if first() != 1 || first() != 2 {
		t.Fatal("first counter did not advance as expected")
	}
	if got := second(); got != 1 {
		t.Fatalf("second counter = %d, want 1", got)
	}
}

func TestMakeMultiplier(t *testing.T) {
	triple := MakeMultiplier(3)
	if got := triple(4); got != 12 {
		t.Fatalf("triple(4) = %d, want 12", got)
	}

	zero := MakeMultiplier(0)
	if got := zero(99); got != 0 {
		t.Fatalf("zero(99) = %d, want 0", got)
	}
}

func TestMap(t *testing.T) {
	double := MakeMultiplier(2)
	got := Map([]int{1, 2, 3}, double)
	want := []int{2, 4, 6}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Map() = %v, want %v", got, want)
	}
}

func TestMapNilFunction(t *testing.T) {
	input := []int{1, 2, 3}
	got := Map(input, nil)
	if !reflect.DeepEqual(got, input) {
		t.Fatalf("Map(nil) = %v, want %v", got, input)
	}
	if &got[0] == &input[0] {
		t.Fatal("Map() reused backing array, want independent result")
	}
}
