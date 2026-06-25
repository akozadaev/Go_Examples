package main

import "fmt"

type Profile struct {
	Email string
}

type User struct {
	Profile *Profile
}

func main() {
	user := &User{}

	if user != nil && user.Profile != nil && user.Profile.Email != "" {
		fmt.Println("email:", user.Profile.Email)
	} else {
		fmt.Println("email is empty or unavailable")
	}
}
