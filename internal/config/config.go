package config

import (
	"os"

	"github.com/gofiber/fiber/v2/log"
	"github.com/joho/godotenv"
)

// Config func to get env value
func Config(key string) string {
	err := godotenv.Load()
	if err != nil {
		log.Errorf("Error loading .env file: ", err)
	}

	return os.Getenv(key)
}

// load .env file
// env := os.Getenv("ENV")
// if env == "" {
// 	env = "development"
// }
// envFile := ".env." + env
// envFile
