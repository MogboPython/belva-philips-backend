package utils

import (
	"errors"
	"strconv"
	"strings"

	"log"

	"github.com/MogboPython/belvaphilips_backend/internal/config"
	"gorm.io/gorm"
)

const (
	PAGE  = 1
	LIMIT = 10
)

func GetPageAndLimitInt(pageStr, limitStr string) (offset, limit int) {
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = PAGE
	}

	limit, err = strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = LIMIT
	}

	// Calculate offset
	offset = (page - 1) * limit

	return offset, limit
}

func ToSnakeCase(input string) string {
	lower := strings.ToLower(input)
	snake := strings.ReplaceAll(lower, " ", "_")

	return snake
}

// ExistsByID checks if a record exists for the given model and ID.
func ExistsByID(db *gorm.DB, model any, id string) (bool, error) {
	err := db.Where("id = ?", id).First(model).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}

func RemoveFile(file string) error {
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
		log.Println("Error deleting image:", err)
		return errors.New("error deleting image")
	}

	return nil
}
