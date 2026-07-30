package main

import (
	"fmt"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type User struct {
	ID   uint
	Name string
	Age  int
}

func main() {
	db, err := gorm.Open(sqlite.Open("gorm_test.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// Автомиграция
	if err := db.AutoMigrate(&User{}); err != nil {
		log.Fatal(err)
	}

	// Вставка
	user := User{Name: "Alexey", Age: 45}
	if err := db.Debug().Create(&user).Error; err != nil {
		log.Fatal(err)
	}

	// Чтение
	var users []User
	if err := db.Debug().Find(&users).Error; err != nil {
		log.Fatal(err)
	}

	for _, u := range users {
		fmt.Printf("ID: %d, Имя: %s, Возраст: %d\n", u.ID, u.Name, u.Age)
	}
}
