package main

import (
	"attendance-app/config"
	"attendance-app/routes"
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Initialize configuration
	config.InitConfig()

	// Set up routes
	r := routes.InitRoutes()

	// Start server
	fmt.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
