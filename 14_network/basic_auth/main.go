package main

import (
	"fmt"
	"log"
	"net/http"
)

func protectedHandler(w http.ResponseWriter, r *http.Request) {
	user, pass, ok := r.BasicAuth()
	if !ok || user != "admin" || pass != "secret" {
		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	fmt.Fprintln(w, "Welcome, admin!")
}

func main() {
	http.HandleFunc("/admin", protectedHandler)
	fmt.Println("Basic Auth example on http://localhost:8080/admin")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
