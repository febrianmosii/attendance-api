package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strconv"

	db "attendance-app/internal/database"
	"attendance-app/routes" // Import the routes package

	"github.com/joho/godotenv"
	"golang.ngrok.com/ngrok"
	"golang.ngrok.com/ngrok/config"
)

func main() {
	// Initialize the router with routes
	router := routes.SetupRoutes()

	// Run ngrok and start the server
	if err := run(context.Background(), router); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context, router http.Handler) error {
	// Load .env file
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	} else {
		log.Println("Successfully loaded .env file")
	}

	isDeploy, err := strconv.ParseBool(os.Getenv("NGROK_DEPLOY"))

	if err != nil {
		log.Println("Error parsing NGROK_DEPLOY")
	}

	// Initialize the database
	db.InitializeDB()

	if isDeploy {
		// Start ngrok tunnel
		listener, err := ngrok.Listen(ctx,
			config.HTTPEndpoint(),
			ngrok.WithAuthtokenFromEnv(), // Uses the authtoken from the environment variables
		)
		if err != nil {
			return err
		}

		log.Println("App URL", listener.URL())

		return http.Serve(listener, router)
	} else {
		// Start the server
		log.Println("Server running on port 8080...")

		return http.ListenAndServe(":8080", router)
	}

	// Start the server and use the router
	// return http.Serve(listener, router)
}
