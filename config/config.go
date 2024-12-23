package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Config struct {
	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string
	Env                string
	Host               string
	Port               string
}

var AppConfig *Config
var GoogleOauthConfig *oauth2.Config

func LoadConfig() {
	if os.Getenv("ENV") != "production" {
		if err := godotenv.Load(".env"); err != nil {
			log.Fatal("Error loading .env file:", err)
		}
	}

	AppConfig = &Config{
		GoogleClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		GoogleClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		GoogleRedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
		Env:                os.Getenv("ENV"),
		Host:               os.Getenv("HOST"),
		Port:               os.Getenv("PORT"),
	}

	GoogleOauthConfig = &oauth2.Config{
		ClientID:     AppConfig.GoogleClientID,
		ClientSecret: AppConfig.GoogleClientSecret,
		RedirectURL:  AppConfig.GoogleRedirectURL,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.profile", "https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
}
