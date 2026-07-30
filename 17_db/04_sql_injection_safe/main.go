package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "./test.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		age INTEGER NOT NULL
	)`)
	if err != nil {
		log.Fatal(err)
	}

	name := "David"
	age := 22

	// Подготовленный запрос
	stmt, err := db.Prepare("INSERT INTO users (name, age) VALUES (?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(name, age)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Данные безопасно вставлены")
}
