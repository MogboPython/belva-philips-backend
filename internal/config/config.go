package config

import (
	"os"

	"github.com/gofiber/fiber/v2/log"
	"github.com/joho/godotenv"
	storage "github.com/supabase-community/storage-go"
)

func Config(key string) string {
	if os.Getenv("FLY_APP_NAME") == "" {
		err := godotenv.Load(".env.prod")
		if err != nil {
			log.Errorf("Error loading .env file: %v", err)
		}
	}

	return os.Getenv(key)
}

func CreateStorageClient() *storage.Client {
	return storage.NewClient(
		Config("SUPABASE_URL"),
		Config("SUPABASE_API_KEY"),
		nil,
	)
}
