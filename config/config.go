package config

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
)

// InitConfig loads the environment variables from the .env file
func InitConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// You can add more configuration setups here
	fmt.Println("Configuration loaded successfully!")
}
