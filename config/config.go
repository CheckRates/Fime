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

// S3Bucket contains all the configuration params for S3 connection
type S3Bucket struct {
	Region string
	Bucket string
	Secret string
	Access string
}

// Config contains all the configuration params for Fime
type Config struct {
	Database DBConfig
	OAuth    OAuthConfig
	S3       S3Bucket
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
		S3: S3Bucket{
			Region: getEnv("AWS_S3_REGION", ""),
			Access: getEnv("AWS_ACCESS_KEY", ""),
			Secret: getEnv("AWS_SECRET", ""),
			Bucket: getEnv("AWS_S3_BUCKET", ""),
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
