//go:build ignore

package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte("my-secret-key")

func main() {
	http.HandleFunc("/protected", func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			http.Error(w, "Expected Authorization: Bearer <token>", http.StatusUnauthorized)
			return
		}
		tokenString := strings.TrimSpace(strings.TrimPrefix(auth, "Bearer "))
		if tokenString == "" {
			http.Error(w, "No token provided", http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return secretKey, nil
		}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
		if err != nil || !token.Valid {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		username, _ := claims["username"].(string)
		fmt.Fprintf(w, "Welcome, %s!", username)
	})

	fmt.Println("JWT check server on http://localhost:8080/protected")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
