package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	password := "admin123"

	// Generate hash
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Password: %s\n", password)
	fmt.Printf("Hash: %s\n", string(hash))

	// Test compare
	err = bcrypt.CompareHashAndPassword(hash, []byte(password))
	if err != nil {
		fmt.Println("❌ Password does NOT match")
	} else {
		fmt.Println("✓ Password matches!")
	}

	// Test with existing hash from DB
	existingHash := "$2a$10$rZCimQKup8dZInPf92d8l.sd6ZKtHEH1xm0cqj6HUkW6YqbVqM1hy"
	err = bcrypt.CompareHashAndPassword([]byte(existingHash), []byte(password))
	if err != nil {
		fmt.Printf("❌ Existing hash does NOT match password '%s'\n", password)
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("✓ Existing hash matches password '%s'\n", password)
	}
}
