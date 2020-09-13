package config

// Reference: https://dev.to/craicoverflow/a-no-nonsense-guide-to-environment-variables-in-go-a2f

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// DBConfig provides all credentials information about the
// db connection
type DBConfig struct {
	ConnString string
}

// Config contains all the configuration params for Fime
type Config struct {
	Database DBConfig
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
	}
}

// Simple helper function to read an environment or return a default value
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
