package transform

import "testing"

func TestPipeline(t *testing.T) {
	got := Pipeline("  GO  ", Trim, Lower, Prefix("lang:"))
	want := "lang:go"
	if got != want {
		t.Fatalf("Pipeline() = %q, want %q", got, want)
	}
}

func TestPipelineSkipsNilStep(t *testing.T) {
	got := Pipeline("  Go  ", nil, Trim)
	if got != "Go" {
		t.Fatalf("Pipeline() = %q, want %q", got, "Go")
	}
}

func TestPipelineWithoutSteps(t *testing.T) {
	got := Pipeline("same")
	if got != "same" {
		t.Fatalf("Pipeline() = %q, want %q", got, "same")
	}
}

func TestAnonymousStep(t *testing.T) {
	exclaim := func(s string) string {
		return s + "!"
	}

	got := Pipeline("go", exclaim)
	if got != "go!" {
		t.Fatalf("Pipeline() = %q, want %q", got, "go!")
	}
}
