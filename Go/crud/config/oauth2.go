package config

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"log"
	"os"

	"github.com/joho/godotenv"
)

var OAuth2Config *oauth2.Config

func init() {
	err := godotenv.Load() // This will load the .env file
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	OAuth2Config = &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  "http://localhost:8086/callback",
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
}
