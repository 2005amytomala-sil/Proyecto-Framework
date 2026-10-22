package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	passwords := []string{"123456"}

	for _, pwd := range passwords {
		hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
		if err != nil {
			fmt.Printf("Error al hashear %s: %v\n", pwd, err)
			continue
		}
		fmt.Printf("Password: %s\nHash: %s\n", pwd, string(hash))
	}
}
