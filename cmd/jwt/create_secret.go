package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
)

func generateSecretKey(length int) string {
	// Create a byte slice of the desired length
	bytes := make([]byte, length)

	// Fill it with random bytes
	_, err := rand.Read(bytes)
	if err != nil {
		log.Fatalf("Failed to generate secret key: %v", err)
	}

	// Return the key as a hex string
	return hex.EncodeToString(bytes)
}

func main() {
	// Generate a 32-byte secret key (256-bit key, suitable for JWT)
	secretKey := generateSecretKey(32)
	fmt.Println("Generated Secret Key:", secretKey)
}
