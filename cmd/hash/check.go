package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Hash from Laravel (as an example)
	storedHash := "$2y$12$RdhLjwrrMePaVipgyOl6uOtOUwWSZs6feZSLsgmQN4F8CPsvidHVy"

	// Password to verify
	password := "zzz"

	// Compare the hash with the password
	err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(password))
	if err != nil {
		fmt.Println("Password does not match!")
	} else {
		fmt.Println("Password matches!")
	}
}
