package main

import "fmt"

type User struct {
	ID           int
	Email        string
	PasswordHash string
	CreatedAt    string
}

func main() {
	newUser := User{
		ID:           888,
		Email:        "test@gmail.com",
		PasswordHash: "qwerty123",
		CreatedAt:    "28/10/2026 20:53",
	}
	fmt.Println(newUser)
}
