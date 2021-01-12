package config

// Reference: https://dev.to/craicoverflow/a-no-nonsense-guide-to-environment-variables-in-go-a2f

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// DBConfig provides all credentials information about the db connection
type DBConfig struct {
	ConnString string
}

// OAuthConfig provides all info for Google OAuth
type OAuthConfig struct {
	ClientID     string
	ClientSecret string
}

// Config contains all the configuration params for Fime
type Config struct {
	Database DBConfig
	OAuth    OAuthConfig
}

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

// New -> Constructor for Config
func New() *Config {
	return &Config{
		Database: DBConfig{
			ConnString: getEnv("DATABASE_URL", ""),
		},
		OAuth: OAuthConfig{
			ClientID:     getEnv("AUTH0_CLIENT_ID", ""),
			ClientSecret: getEnv("AUTH0_CLIENT_SECRET", ""),
		},
	}
}

// Simple helper function to read an environment or return a default value
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
