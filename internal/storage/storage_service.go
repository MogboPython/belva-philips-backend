package storage

import (
	"errors"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
	storage "github.com/supabase-community/storage-go"
)

type StorageService interface { //nolint:revive // it works
	UploadFile(imageFile *multipart.FileHeader, bucketID string, subPath ...string) (string, error)
	RemoveFile(file string) error
	RemoveFolder(bucketName, folderPath string) error
}

type storageService struct {
	client *storage.Client
}

func NewStorageService(client *storage.Client) StorageService {
	return &storageService{
		client: client,
	}
}

func (s *storageService) UploadFile(imageFile *multipart.FileHeader, bucketID string, subPath ...string) (string, error) {
	if imageFile == nil {
		return "", nil
	}

	fileExt := filepath.Ext(imageFile.Filename)
	if fileExt == "" {
		return "", errors.New("invalid file extension")
	}

	fileExt = strings.TrimPrefix(fileExt, ".")

	uniqueID := uuid.New()
	filename := strings.ReplaceAll(uniqueID.String(), "-", "")

	var imagePath string
	if len(subPath) > 0 && subPath[0] != "" {
		imagePath = fmt.Sprintf("%s/%s.%s", subPath[0], filename, fileExt)
	} else {
		imagePath = fmt.Sprintf("%s.%s", filename, fileExt)
	}

	const maxFileSize = 5 * 1024 * 1024
	if imageFile.Size > maxFileSize {
		return "", errors.New("file size exceeds 5MB limit")
	}

	contentType := imageFile.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") && contentType != "application/pdf" {
		return "", errors.New("file is neither an image nor a PDF")
	}

	file, err := imageFile.Open()
	if err != nil {
		log.Error("Error opening file:", err)
		return "", errors.New("error opening file")
	}
	defer file.Close()

	result, err := s.client.UploadFile(bucketID, imagePath, file, storage.FileOptions{
		ContentType: &contentType,
	})
	if err != nil {
		log.Error("Error saving image:", err)
		return "", errors.New("error saving image")
	}

	return result.Key, nil
}

func (s *storageService) RemoveFile(file string) error {
	var bucketName, fileName string

	const zero, one, two, three = 0, 1, 2, 3

	filePath := strings.Split(file, "/")
	switch len(filePath) {
	case two:
		bucketName = filePath[0]
		fileName = filePath[1]
	case three:
		bucketName = strings.Join(filePath[0:2], "/")
		fileName = filePath[2]
	case zero, one:
		return errors.New("invalid file path")
	default:
		bucketName = filePath[0]
		fileName = strings.Join(filePath[1:], "/")
	}

	_, err := s.client.RemoveFile(bucketName, []string{fileName})

	if err != nil {
		log.Error("Error deleting image: ", err)
		return errors.New("error deleting image")
	}

	return nil
}

// func (s *storageService) RemoveFolder(folderPath string) error {
// 	_, err := s.client.EmptyBucket(folderPath)

// 	if err != nil {
// 		log.Error("Error emptying bucket:", err)
// 		return errors.New("error emptying bucket")
// 	}

// 	_, err = s.client.DeleteBucket(folderPath)
// 	if err != nil {
// 		log.Error("Error deleting bucket:", err)
// 		return errors.New("error deleting bucket")
// 	}

// 	return nil
// }

func (s *storageService) RemoveFolder(bucketName, folderPath string) error {
	if folderPath[len(folderPath)-1] != '/' {
		folderPath += "/"
	}

	files, err := s.client.ListFiles(bucketName, folderPath, storage.FileSearchOptions{})
	if err != nil {
		log.Error("Error listing files in folder:", err)
		return errors.New("error listing folder contents")
	}

	filePaths := make([]string, 0, len(files))

	for i := range files {
		filePaths = append(filePaths, folderPath+files[i].Name)
	}

	if len(filePaths) > 0 {
		_, err := s.client.RemoveFile(bucketName, filePaths)
		if err != nil {
			log.Error("Error deleting files:", err)
			return errors.New("error deleting folder contents")
		}
	}

	return nil
}
