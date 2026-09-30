package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                   string
	MongoURI               string
	DatabaseName           string
	FrontendURL            string
	CognitoOpenIDConfigURL string
	CognitoClientID        string
}

func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found")
	}

	return &Config{
		Port:                   os.Getenv("PORT"),
		MongoURI:               os.Getenv("MONGODB_URI"),
		DatabaseName:           os.Getenv("MONGODB_DATABASE"),
		FrontendURL:            os.Getenv("FRONTEND_URL"),
		CognitoOpenIDConfigURL: os.Getenv("COGNITO_OPENID_CONFIG_URL"),
		CognitoClientID:        os.Getenv("COGNITO_CLIENT_ID"),
	}
}
