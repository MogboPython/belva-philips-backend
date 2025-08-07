package utils

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/MogboPython/belvaphilips_backend/internal/config"
	"gorm.io/gorm"
)

const (
	PAGE  = 1
	LIMIT = 10
	REGEX = 2
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

func PublicImageURL(imageName string) string {
	if imageName == "" {
		return ""
	}

	return fmt.Sprintf("%s/object/public/%s",
		config.Config("SUPABASE_URL"),
		imageName,
	)
}

func ExtractPublicID(url string) (string, error) {
	// Regular expression to match the public ID (last segment before the file extension)
	re := regexp.MustCompile(`/v\d+/([^/]+)\.\w+$`)
	matches := re.FindStringSubmatch(url)

	if len(matches) < REGEX {
		return "", fmt.Errorf("could not extract public ID from URL: %s", url)
	}

	return matches[1], nil
}
