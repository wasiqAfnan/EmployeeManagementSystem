package main

import (
	"log"
	"net/http"
	"os"

	// middlewares "EMS/middlewares"
	// routes "EMS/routes"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("PORT")

	// Create a new ServeMux
	mux := http.NewServeMux()

	// Register the routes
	// routes.EmployeeRoutes(mux)

	// Add middleware
	// handler := middlewares.Logging(mux)

	// Set up the server configuration
	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Println("Server starting on http://localhost:5050")

	// Start the server
	log.Fatal(server.ListenAndServe())
}
