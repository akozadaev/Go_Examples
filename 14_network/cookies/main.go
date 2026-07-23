package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/securecookie"
)

// Фиксированные ключи только для учебного примера.
// В реальном сервисе ключи берут из конфигурации/секретницы и не хранят в git.
// GenerateRandomKey при каждом запуске сделал бы невозможным чтение cookie после рестарта.
var cookieHandler = securecookie.New(
	[]byte("0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"), // hash key, 64 bytes
	[]byte("0123456789abcdef0123456789abcdef"),                                 // block key, 32 bytes
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Сначала пытаемся прочитать уже установленную cookie.
		if cookie, err := r.Cookie("cookie-name"); err == nil {
			value := make(map[string]interface{})
			if err := cookieHandler.Decode("cookie-name", cookie.Value, &value); err == nil {
				fmt.Fprintf(w, "Cookie Value: %v\n", value["key"])
				return
			}
			fmt.Fprintln(w, "Cookie найдена, но не удалось декодировать (возможно, другой ключ).")
			return
		}

		// Cookie ещё нет - устанавливаем и просим повторить запрос.
		value := map[string]interface{}{"key": "value"}
		encoded, err := cookieHandler.Encode("cookie-name", value)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:     "cookie-name",
			Value:    encoded,
			Path:     "/",
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		})
		fmt.Fprintln(w, "Cookie установлена. Обновите страницу или повторите запрос с сохранёнными cookies.")
	})

	fmt.Println("Cookies example on http://localhost:8080/")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
