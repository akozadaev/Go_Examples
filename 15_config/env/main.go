package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found, using system environment")
	}

	port := getEnv("MY_APP_PORT", "8080")
	port1 := getEnv("MY_APP_PORT1", "8080")
	dbHost := getEnv("MY_APP_DB_HOST", "localhost")
	dbPort := getEnv("MY_APP_DB_PORT", "5432")
	dbUser := getEnv("MY_APP_DB_USER", "user")
	dbPass := getEnv("MY_APP_DB_PASSWORD", "password")

	fmt.Printf("Server Port: %s  %s \n", port, port1)
	fmt.Printf("DB: %s@%s:%s\n", dbUser, dbHost, dbPort)
	fmt.Printf("DB password:%s\n", dbPass)
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
