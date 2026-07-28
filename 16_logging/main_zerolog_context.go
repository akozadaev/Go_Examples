package main

import (
	"net/http"
	"os"

	"github.com/rs/zerolog"
)

func someMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		l := zerolog.New(os.Stdout).With().
			Str("user_id", "usr-1234").
			Logger()

		// Сохраняем логгер в контексте
		ctx = l.WithContext(ctx)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func handler(w http.ResponseWriter, r *http.Request) {
	l := zerolog.Ctx(r.Context())
	l.Info().Str("doc_id", "doc-xyz").Msg("Документ doc-xyz удален")
	w.WriteHeader(http.StatusOK)
}

func main() {
	http.Handle("/delete", someMiddleware(http.HandlerFunc(handler)))
	println("Сервер запущен на :8080")
	http.ListenAndServe(":8080", nil)
}
