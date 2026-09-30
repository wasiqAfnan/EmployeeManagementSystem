package main

import (
	"log"
	"net/http"

	"EMS/internal/cognito"
	"EMS/internal/config"
	"EMS/internal/database"
	"EMS/internal/handler"
	"EMS/internal/middleware"
	"EMS/internal/repository"
	"EMS/internal/routes"
	"EMS/internal/service"
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

	// Get the users collection & repo
	userCollection := client.Database(cfg.DatabaseName).Collection("users")
	userRepo := repository.NewUserRepository(userCollection)
	userService := service.NewUserService(userRepo)

	// Create repository, service, and handler for Employee
	employeeRepo := repository.NewEmployeeRepository(collection)
	employeeService := service.NewEmployeeService(employeeRepo, userRepo)
	employeeHandler := handler.NewEmployeeHandler(employeeService)

	// Create a new ServeMux
	mux := http.NewServeMux()

	// Register the routes
	routes.EmployeeRoutes(mux, employeeHandler)

	// Initialize Cognito JWT Verifier
	var jwtVerifier *cognito.JWTVerifier
	if cfg.CognitoOpenIDConfigURL != "" && cfg.CognitoClientID != "" {
		v, err := cognito.NewJWTVerifier(cfg)
		if err != nil {
			log.Fatalf("Failed to initialize JWT verifier: %v", err)
		}
		jwtVerifier = v
	}

	userHandler := handler.NewUserHandler(userService, jwtVerifier)

	// Register User routes
	routes.UserRoutes(mux, userHandler)

	// Add middlewares (Execution flow: Logging -> Auth -> CORS -> Mux)
	// To achieve this flow, we wrap inside-out: CORS, then Auth, then Logging
	var finalHandler http.Handler = mux
	finalHandler = middleware.CORS(finalHandler, cfg.FrontendURL)
	if jwtVerifier != nil {
		finalHandler = middleware.Auth(jwtVerifier, userRepo)(finalHandler)
	}
	finalHandler = middleware.Logging(finalHandler)

	// Set up the server
	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: finalHandler,
	}

	log.Printf("Server starting on http://localhost:%s", cfg.Port)

	// Start the server
	log.Fatal(server.ListenAndServe())
}
