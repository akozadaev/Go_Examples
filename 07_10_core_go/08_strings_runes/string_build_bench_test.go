package main

import (
	"fmt"
	"strings"
	"testing"
)

var benchParts = []string{"go", "строки", "руны", "ошибки", "тесты"}

func BenchmarkConcatPlus(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var result string
		for _, part := range benchParts {
			result += part
		}
		_ = result
	}
}

func BenchmarkFmtSprintf(b *testing.B) {
	for i := 0; i < b.N; i++ {
		result := fmt.Sprintf("%s%s%s%s%s", benchParts[0], benchParts[1], benchParts[2], benchParts[3], benchParts[4])
		_ = result
	}
}

func BenchmarkStringsBuilder(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var builder strings.Builder
		builder.Grow(64)
		for _, part := range benchParts {
			builder.WriteString(part)
		}
		_ = builder.String()
	}
}
