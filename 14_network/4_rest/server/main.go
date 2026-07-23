package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var (
	users = map[int]User{
		1: {ID: 1, Name: "Alexey"},
		2: {ID: 2, Name: "Daria"},
	}
	nextID = 3
	mu     sync.RWMutex
)

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func parseUserID(r *http.Request) (int, error) {
	return strconv.Atoi(r.PathValue("id"))
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	mu.RLock()
	defer mu.RUnlock()

	writeJSON(w, http.StatusOK, users)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	id, err := parseUserID(r)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	mu.RLock()
	user, exists := users[id]
	mu.RUnlock()

	if !exists {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	writeJSON(w, http.StatusOK, user)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	var newUser User
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if newUser.Name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}

	mu.Lock()
	newUser.ID = nextID
	users[nextID] = newUser
	nextID++
	mu.Unlock()

	writeJSON(w, http.StatusCreated, newUser)
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	id, err := parseUserID(r)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var body User
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if body.Name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	if _, exists := users[id]; !exists {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	updated := User{ID: id, Name: body.Name}
	users[id] = updated
	writeJSON(w, http.StatusOK, updated)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := parseUserID(r)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	if _, exists := users[id]; !exists {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	delete(users, id)
	w.WriteHeader(http.StatusNoContent)
}

func main() {
	mux := http.NewServeMux()

	// Шаблоны ServeMux (Go 1.22+): метод + путь с {id}
	mux.HandleFunc("GET /users", getUsers)
	mux.HandleFunc("POST /users", createUser)
	mux.HandleFunc("GET /users/{id}", getUser)
	mux.HandleFunc("PUT /users/{id}", updateUser)
	mux.HandleFunc("DELETE /users/{id}", deleteUser)

	fmt.Println("REST API running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
