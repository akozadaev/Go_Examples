package userctx

import (
	"context"
	"testing"
)

var benchmarkUserID uint

func BenchmarkContextWithUserID(b *testing.B) {
	ctx := context.Background()

	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		withUser := WithUserID(ctx, 42)
		benchmarkUserID, _ = GetUserID(withUser)
	}
}
