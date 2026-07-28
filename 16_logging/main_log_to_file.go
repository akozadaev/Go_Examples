package main

import (
	"log"
	"os"
)

func main() {
	logFile, err := os.Create("app.log")
	if err != nil {
		log.Fatal("Не удалось создать файл лога:", err)
	}
	defer logFile.Close()

	log.SetOutput(logFile)
	log.Println("Это сообщение будет записано в файл app.log")
}
