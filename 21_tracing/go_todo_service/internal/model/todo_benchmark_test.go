package model

import "testing"

var benchmarkTodo *Todo
var benchmarkTodoValue Todo

func BenchmarkTodoCreateRequestToTodo(b *testing.B) {
	req := &TodoCreateRequest{
		Title:       "Изучить профилирование",
		Description: "Снять CPU и heap профили",
		Done:        false,
	}

	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		benchmarkTodo = req.ToTodo()
	}
}

func BenchmarkTodoCreateRequestValue(b *testing.B) {
	req := TodoCreateRequest{
		Title:       "Изучить профилирование",
		Description: "Снять CPU и heap профили",
		Done:        false,
	}

	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		benchmarkTodoValue = Todo{
			Title:       req.Title,
			Description: req.Description,
			Done:        req.Done,
		}
	}
}
