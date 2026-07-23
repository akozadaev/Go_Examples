package main

import (
	"fmt"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Привет, мир!")
}

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("HTTP server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
