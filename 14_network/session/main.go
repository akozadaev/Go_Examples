package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
)

// CookieStore хранит Values в подписанной cookie на клиенте (не server-side store).
// Секрет только для учебного примера; в production - из переменных окружения.
var store = sessions.NewCookieStore([]byte("your-secret-key-change-me-32b!!"))

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		session, err := store.Get(r, "session-name")
		if err != nil {
			http.Error(w, "session error", http.StatusInternalServerError)
			return
		}

		value, ok := session.Values["key"].(string)
		if !ok || value == "" {
			session.Values["key"] = "value"
			session.Options.HttpOnly = true
			session.Options.SameSite = http.SameSiteLaxMode
			if err := session.Save(r, w); err != nil {
				http.Error(w, "cannot save session", http.StatusInternalServerError)
				return
			}
			fmt.Fprintln(w, "Сессия создана, значение записано. Обновите страницу.")
			return
		}

		fmt.Fprintf(w, "Session Value: %s\n", value)
	})

	fmt.Println("Session example on http://localhost:8080/")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
