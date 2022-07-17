package config

// Reference: https://dev.to/craicoverflow/a-no-nonsense-guide-to-environment-variables-in-go-a2f

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/joho/godotenv"
)

// JWT provides all info for token signing
type JWT struct {
	AccessSecret      string
	AccessExpiration  time.Duration
	RefreshExpiration time.Duration
}

// S3Bucket contains all the configuration params for AWS S3 connection
type S3Bucket struct {
	Region string
	Bucket string
	Secret string
	Access string
}

// Config contains all the configuration params for Fime
type Config struct {
	Address    string
	ConnString string
	Token      JWT
	S3         S3Bucket
}

func Load(path string) (Config, error) {
	envPath := filepath.Join(path, ".env")

	if err := godotenv.Load(envPath); err != nil {
		log.Print("No .env file found")
	}

	return Config{
		Address:    getEnv("PORT", ""),
		ConnString: getEnv("DATABASE_URL", ""),
		Token: JWT{
			AccessSecret:      getEnv("TOKEN_SYMMETRIC_KEY", ""),
			AccessExpiration:  15 * time.Minute,
			RefreshExpiration: 3 * time.Hour,
		},
		S3: S3Bucket{
			Region: getEnv("AWS_S3_REGION", ""),
			Access: getEnv("AWS_ACCESS_KEY", ""),
			Secret: getEnv("AWS_SECRET", ""),
			Bucket: getEnv("AWS_S3_BUCKET", ""),
		},
	}, nil
}

// Simple helper function to read an environment or return a default value
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
