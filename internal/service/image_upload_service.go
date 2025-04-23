package service

import (
	"errors"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/MogboPython/belvaphilips_backend/internal/config"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"

	storage "github.com/supabase-community/storage-go"
)

func uploadImage(imageFile *multipart.FileHeader, bucketID string) (string, error) {
	if imageFile == nil {
		return "", nil
	}

	storageClient := config.CreateStorageClient()

	// Extract file extension
	fileExt := filepath.Ext(imageFile.Filename)
	if fileExt == "" {
		return "", errors.New("invalid file extension")
	}

	fileExt = strings.TrimPrefix(fileExt, ".")

	// Generate unique filename
	uniqueID := uuid.New()
	filename := strings.ReplaceAll(uniqueID.String(), "-", "")
	imagePath := fmt.Sprintf("%s.%s", filename, fileExt)

	// Check file size (e.g., 5MB limit)
	const maxFileSize = 5 * 1024 * 1024 // 5MB
	if imageFile.Size > maxFileSize {
		return "", errors.New("file size exceeds 5MB limit")
	}

	// Check content type
	contentType := imageFile.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		return "", errors.New("file is not an image")
	}

	// Open the file
	file, err := imageFile.Open()
	if err != nil {
		log.Error("Error opening file:", err)
		return "", fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	result, err := storageClient.UploadFile(bucketID, imagePath, file, storage.FileOptions{
		ContentType: &contentType,
	})
	if err != nil {
		log.Error("Error saving image:", err)
		return "", errors.New("error saving image")
	}

	return result.Key, nil
}

func removeImage(file string) error {
	storageClient := config.CreateStorageClient()

	bucketName := strings.Split(file, "/")[0]
	fileName := strings.Split(file, "/")[1]
	_, err := storageClient.RemoveFile(bucketName, []string{fileName})

	if err != nil {
		log.Error("Error deleting image:", err)
		return errors.New("error deleting image")
	}

	return nil
}

// Constructs and returns the public URL of a image
func PublicImageURL(imageName string) string {
	if imageName == "" {
		return ""
	}

	return fmt.Sprintf("%s/object/public/%s",
		config.Config("SUPABASE_URL"),
		imageName,
	)
}
