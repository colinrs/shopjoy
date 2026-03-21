package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func testPassword() {
	password := "admin123"
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Bcrypt hash for 'admin123':")
	fmt.Println(string(hashed))
}
