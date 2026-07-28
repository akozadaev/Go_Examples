package main

import (
	"fmt"
	"log"

	"gopkg.in/ini.v1"
)

func main() {
	cfg, err := ini.Load("config.ini")
	if err != nil {
		log.Fatalf("Failed to read config.ini: %v", err)
	}

	port := cfg.Section("server").Key("port").MustInt(8080)
	host := cfg.Section("server").Key("host").MustString("localhost")
	dbName := cfg.Section("database").Key("db_name").MustString("defaultdb")

	fmt.Printf("Server: %s:%d\n", host, port)
	fmt.Printf("Database: %s\n", dbName)
}