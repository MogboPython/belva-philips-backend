package storage

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/MogboPython/belvaphilips_backend/internal/config"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/admin"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
	storage "github.com/supabase-community/storage-go"
)

type StorageService interface { //nolint:revive // it works
	UploadFile(imageFile *multipart.FileHeader, bucketID string, subPath ...string) (string, error)
	RemoveFile(file string) error
	RemoveFolder(bucketName, folderPath string) error
	BulkDeleteCloudAssets(publicIDs []string) error
}

type storageService struct {
	client *storage.Client
}

var (
	supabaseURL    = config.Config("SUPABASE_URL")
	supabaseAPIKey = config.Config("SUPABASE_API_KEY")
	client         = storage.NewClient(supabaseURL, supabaseAPIKey, nil)
)

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

func (*storageService) RemoveFile(file string) error {
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

	// FIXME: very interesting, it works with client.RemoveFile but not with s.client.RemoveFile
	_, err := client.RemoveFile(bucketName, []string{fileName})

	if err != nil {
		log.Error("Error deleting image: ", err)
		return errors.New("error deleting image")
	}

	log.Infof("File %s deleted successfully", file)

	return nil
}

func (*storageService) RemoveFolder(bucketName, folderPath string) error {
	if folderPath[len(folderPath)-1] != '/' {
		folderPath += "/"
	}

	files, err := client.ListFiles(bucketName, folderPath, storage.FileSearchOptions{})
	if err != nil {
		log.Error("Error listing files in folder:", err)
		return errors.New("error listing folder contents")
	}

	filePaths := make([]string, 0, len(files))

	for i := range files {
		filePaths = append(filePaths, folderPath+files[i].Name)
	}

	if len(filePaths) > 0 {
		_, err := client.RemoveFile(bucketName, filePaths)
		if err != nil {
			log.Error("Error deleting files:", err)
			return errors.New("error deleting folder contents")
		}
	}

	log.Infof("Files in %s deleted successfully", folderPath)

	return nil
}

func getCloudinaryClient() (*cloudinary.Cloudinary, error) {
	return cloudinary.NewFromParams(
		config.Config("CLOUDINARY_CLOUD_NAME"),
		config.Config("CLOUDINARY_API_KEY"),
		config.Config("CLOUDINARY_API_SECRET"),
	)
}

// func deleteCloudImage(publicID string) error {
// 	cld, err := getCloudinaryClient()
// 	if err != nil {
// 		log.Error("Failed to initialize Cloudinary: ", err)
// 		return errors.New("Failed to initialize Cloudinary")
// 	}

// 	ctx := context.Background()

// 	resp, err := cld.Upload.Destroy(ctx, uploader.DestroyParams{
// 		PublicID:     publicID,
// 		ResourceType: "image",
// 	})
// 	if err != nil {
// 		fmt.Printf("Failed to delete image %s: %v\n", publicID, err)
// 		return errors.New("failed to delete image")
// 	}

// 	log.Infof("Image %s deleted successfully.", publicID)
// 	return nil
// }

func (*storageService) BulkDeleteCloudAssets(publicIDs []string) error {
	cld, err := getCloudinaryClient()
	if err != nil {
		log.Error("Failed to initialize Cloudinary: ", err)
		return errors.New("failed to initialize Cloudinary")
	}

	ctx := context.Background()
	resp, err := cld.Admin.DeleteAssets(
		ctx,
		admin.DeleteAssetsParams{PublicIDs: publicIDs},
	)

	if err != nil {
		log.Error("failed to bulk delete assets: ", err)
		return errors.New("failed to delete images")
	}

	log.Infof("Deleted assets count: %d", len(resp.Deleted))

	return nil
}
