package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func Env(key string) string {
	appMode := os.Getenv("APP_MODE")

	if appMode != "production" {
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatalf("Error loading .env file")
		}
	}

	result := os.Getenv(key)

	if result == "" {
		message := fmt.Sprintf("Environment variable %s is missing", key)
		panic(message)
	}

	return result
}
