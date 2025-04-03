package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config func to get env value
func Config(key string) string {
	// load .env file
	env := os.Getenv("GO_ENV")
	if env == "" {
		env = "development"
	}

	envFile := ".env." + env

	err := godotenv.Load(envFile)
	if err != nil {
		fmt.Print("Error loading .env file")
	}
	return os.Getenv(key)
}
