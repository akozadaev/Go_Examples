package main

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Секрет только для учебного примера.
var secretKey = []byte("my-secret-key")

func main() {
	claims := jwt.MapClaims{
		"username": "user123",
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println("JWT:", tokenString)
	fmt.Println()
	fmt.Println("Проверка:")
	fmt.Println(`curl -H "Authorization: Bearer ` + tokenString + `" http://localhost:8080/protected`)
}
