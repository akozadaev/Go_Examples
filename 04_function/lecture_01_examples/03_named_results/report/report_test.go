package report

import (
	"errors"
	"testing"
)

func TestBuildSummary(t *testing.T) {
	got, err := BuildSummary(10, 20, 30)
	if err != nil {
		t.Fatalf("BuildSummary() unexpected error: %v", err)
	}

	want := Summary{Min: 10, Max: 30, Avg: 20}
	if got != want {
		t.Fatalf("BuildSummary() = %+v, want %+v", got, want)
	}
}

func TestBuildSummaryEmpty(t *testing.T) {
	got, err := BuildSummary()
	if !errors.Is(err, ErrNoScores) {
		t.Fatalf("BuildSummary() error = %v, want ErrNoScores", err)
	}
	if got != (Summary{}) {
		t.Fatalf("BuildSummary() summary = %+v, want zero value", got)
	}
}

func TestPassFail(t *testing.T) {
	passed, failed := PassFail(60, 100, 40, 60, 59)
	if passed != 2 || failed != 2 {
		t.Fatalf("PassFail() = %d, %d; want 2, 2", passed, failed)
	}
}

func TestPassFailEmpty(t *testing.T) {
	passed, failed := PassFail(60)
	if passed != 0 || failed != 0 {
		t.Fatalf("PassFail() = %d, %d; want 0, 0", passed, failed)
	}
}
