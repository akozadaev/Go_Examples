package memdemo

import "testing"

var resultString string
var resultPoint *Point

func TestExamples(t *testing.T) {
	if got := SumLocal(2, 3); got != 5 {
		t.Fatalf("SumLocal() = %d, want 5", got)
	}
	if got := NewPoint(2, 3); *got != (Point{X: 2, Y: 3}) {
		t.Fatalf("NewPoint() = %#v", got)
	}
}

func BenchmarkBuildMessageConcat(b *testing.B) {
	parts := []string{"Go", " ", "memory", " ", "demo"}
	b.ReportAllocs()
	for b.Loop() {
		resultString = BuildMessageConcat(parts)
	}
}

func BenchmarkBuildMessageBuffer(b *testing.B) {
	parts := []string{"Go", " ", "memory", " ", "demo"}
	b.ReportAllocs()
	for b.Loop() {
		resultString = BuildMessageBuffer(parts)
	}
}

func BenchmarkNewPoint(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		resultPoint = NewPoint(10, 20)
	}
}
