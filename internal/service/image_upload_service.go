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

func uploadFile(imageFile *multipart.FileHeader, bucketID string, subPath ...string) (string, error) {
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

	// Determine path based on optional subPath parameter
	var imagePath string
	if len(subPath) > 0 && subPath[0] != "" {
		imagePath = fmt.Sprintf("%s/%s.%s", subPath[0], filename, fileExt)
	} else {
		imagePath = fmt.Sprintf("%s.%s", filename, fileExt)
	}

	// Check file size (e.g., 5MB limit)
	const maxFileSize = 5 * 1024 * 1024 // 5MB
	if imageFile.Size > maxFileSize {
		return "", errors.New("file size exceeds 5MB limit")
	}

	// Check content type
	contentType := imageFile.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") && contentType != "application/pdf" {
		return "", errors.New("file is neither an image nor a PDF")
	}

	// Open the file
	file, err := imageFile.Open()
	if err != nil {
		log.Error("Error opening file:", err)
		return "", errors.New("error opening file")
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

func removeFile(file string) error {
	storageClient := config.CreateStorageClient()

	var bucketName, fileName string

	const zero, one, two = 0, 1, 2

	filePath := strings.Split(file, "/")
	switch len(filePath) {
	case two:
		bucketName = filePath[0]
		fileName = filePath[1]
	case zero, one:
		return errors.New("invalid file path")
	default:
		bucketName = filePath[0]
		fileName = strings.Join(filePath[1:], "/")
	}

	_, err := storageClient.RemoveFile(bucketName, []string{fileName})

	if err != nil {
		log.Error("Error deleting image:", err)
		return errors.New("error deleting image")
	}

	return nil
}

// Constructs and returns the public URL of a image
func publicImageURL(imageName string) string {
	if imageName == "" {
		return ""
	}

	return fmt.Sprintf("%s/object/public/%s",
		config.Config("SUPABASE_URL"),
		imageName,
	)
}
