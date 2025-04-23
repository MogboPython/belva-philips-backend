package config

import (
	"os"

	"github.com/gofiber/fiber/v2/log"
	"github.com/joho/godotenv"
	storage "github.com/supabase-community/storage-go"
)

// Config func to get env value
func Config(key string) string {
	err := godotenv.Load()
	if err != nil {
		log.Errorf("Error loading .env file: ", err)
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
