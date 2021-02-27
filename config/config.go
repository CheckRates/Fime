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

// JWT provides all info for token signing
type JWT struct {
	AccessSecret     string
	RefreshSecret    string
	RefreshPublicKey string
	AccessPublicKey  string
	Expiration       int64
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
	Token    JWT
	S3       S3Bucket
	Address  string
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
		Token: JWT{
			AccessSecret:     getEnv("ACCESS_TOKEN_SECRET", ""),
			RefreshSecret:    getEnv("REFRESH_TOKEN_SECRET", ""),
			RefreshPublicKey: getEnv("REFRESH_TOKEN_PUBLIC_PATH", ""),
			Expiration:       5,
		},
		S3: S3Bucket{
			Region: getEnv("AWS_S3_REGION", ""),
			Access: getEnv("AWS_ACCESS_KEY", ""),
			Secret: getEnv("AWS_SECRET", ""),
			Bucket: getEnv("AWS_S3_BUCKET", ""),
		},
		Address: getEnv("PORT", ""),
	}
}

// Simple helper function to read an environment or return a default value
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
