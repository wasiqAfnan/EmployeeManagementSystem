package main

import (
	"log"
	"net/http"

	"EMS/internal/config"
	"EMS/internal/database"
	"EMS/internal/handler"
	"EMS/internal/middleware"
	"EMS/internal/routes"
)

func main() {
	// Load config from .env
	cfg := config.Load()

	// Connect to MongoDB
	client, err := database.Connect(cfg.MongoURI)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	log.Println("Connected to MongoDB.")

	// Get the employees collection directly
	collection := client.Database(cfg.DatabaseName).Collection("employees")

	// Create handler: pass the collection directly, no repository needed
	employeeHandler := handler.NewEmployeeHandler(collection)

	// Create a new ServeMux
	mux := http.NewServeMux()

	// Register the routes
	routes.EmployeeRoutes(mux, employeeHandler)

	// Add logging middleware
	muxWithCORS := middleware.CORS(mux, cfg.FrontendURL)
	muxWithLogging := middleware.Logging(muxWithCORS)

	// Set up the server
	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: muxWithLogging,
	}

	log.Printf("Server starting on http://localhost:%s", cfg.Port)

	// Start the server
	log.Fatal(server.ListenAndServe())
}
