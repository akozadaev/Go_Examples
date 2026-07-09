package testingdemo

import (
	"os"
	"reflect"
	"strings"
	"testing"
	"unicode/utf8"
)

func TestWordCount(t *testing.T) {
	tests := []struct {
		name string
		text string
		want map[string]int
	}{
		{
			name: "пустая строка",
			text: "",
			want: map[string]int{},
		},
		{
			name: "слова с Go",
			text: "Go, go, тесты!",
			want: map[string]int{"go": 2, "тесты": 1},
		},
		{
			name: "слова Юникода",
			text: "Привет мир, привет Go",
			want: map[string]int{"go": 1, "мир": 1, "привет": 2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := WordCount(tt.text)
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("WordCount() = %#v, ожидается %#v", got, tt.want)
			}
		})
	}
}

func TestRenderReport(t *testing.T) {
	counts := map[string]int{"go": 2, "тесты": 1, "ошибки": 1}
	got := RenderReport(counts)

	wantBytes, err := os.ReadFile("testdata/report.golden")
	if err != nil {
		t.Fatalf("чтение golden-файла: %v", err)
	}

	if got != string(wantBytes) {
		t.Fatalf("RenderReport() вернул неожиданный результат\n--- получено ---\n%s--- ожидается ---\n%s", got, wantBytes)
	}
}

func FuzzNormalizeSpaces(f *testing.F) {
	for _, seed := range []string{"", "привет   мир", " Привет\tмир\n", "а\u00a0б"} {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, s string) {
		got := NormalizeSpaces(s)

		if !utf8.ValidString(got) {
			t.Fatalf("NormalizeSpaces вернула некорректный UTF-8: %q", got)
		}
		if strings.Contains(got, "  ") {
			t.Fatalf("NormalizeSpaces оставила повторяющиеся пробелы: %q", got)
		}
		if got != strings.TrimSpace(got) {
			t.Fatalf("NormalizeSpaces оставила пробелы по краям: %q", got)
		}
	})
}
