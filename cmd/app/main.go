package main

import (
	db "attendance-app/internal/database"
	"attendance-app/routes" // Import the routes package
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	} else {
		log.Println("Successfully loaded .env file")
	}

	// Initialize the router with routes
	router := routes.SetupRoutes()

	db.InitializeDB()

	// Start the server
	log.Println("Server running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", router))
}
