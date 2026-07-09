package testingdemo

import (
	"fmt"
	"sort"
	"strings"
	"unicode"
)

func WordCount(text string) map[string]int {
	result := make(map[string]int)
	for _, field := range strings.FieldsFunc(text, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r)
	}) {
		word := strings.ToLower(field)
		result[word]++
	}
	return result
}

func RenderReport(counts map[string]int) string {
	keys := make([]string, 0, len(counts))
	for key := range counts {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var b strings.Builder
	for _, key := range keys {
		fmt.Fprintf(&b, "%s: %d\n", key, counts[key])
	}
	return b.String()
}

func NormalizeSpaces(s string) string {
	return strings.Join(strings.Fields(s), " ")
}
