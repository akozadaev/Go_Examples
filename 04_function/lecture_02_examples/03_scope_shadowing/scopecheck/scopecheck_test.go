package scopecheck

import (
	"errors"
	"reflect"
	"testing"
)

func TestCountPositiveWithBugDocumentsShadowing(t *testing.T) {
	got := CountPositiveWithBug([]int{-1, 2, 3})
	if got != 0 {
		t.Fatalf("CountPositiveWithBug() = %d, want 0", got)
	}
}

func TestCountPositive(t *testing.T) {
	got := CountPositive([]int{-1, 2, 3, 0, 4})
	if got != 3 {
		t.Fatalf("CountPositive() = %d, want 3", got)
	}
}

func TestCountPositiveEmpty(t *testing.T) {
	if got := CountPositive(nil); got != 0 {
		t.Fatalf("CountPositive(nil) = %d, want 0", got)
	}
}

func TestParseCSVInts(t *testing.T) {
	got, err := ParseCSVInts("1, 2, -3")
	if err != nil {
		t.Fatalf("ParseCSVInts() unexpected error: %v", err)
	}

	want := []int{1, 2, -3}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("ParseCSVInts() = %v, want %v", got, want)
	}
}

func TestParseCSVIntsEmptyInput(t *testing.T) {
	got, err := ParseCSVInts(" ")
	if err != nil {
		t.Fatalf("ParseCSVInts() unexpected error: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("ParseCSVInts() = %v, want empty slice", got)
	}
}

func TestParseCSVIntsEmptyToken(t *testing.T) {
	_, err := ParseCSVInts("1,,2")
	if !errors.Is(err, ErrEmptyToken) {
		t.Fatalf("ParseCSVInts() error = %v, want ErrEmptyToken", err)
	}
}

func TestFindLabel(t *testing.T) {
	labels := map[string]string{"go": "Golang"}
	if got := FindLabel(labels, "go"); got != "Golang" {
		t.Fatalf("FindLabel() = %q, want Golang", got)
	}
	if got := FindLabel(labels, "rust"); got != "unknown" {
		t.Fatalf("FindLabel() = %q, want unknown", got)
	}
}
