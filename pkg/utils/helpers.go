package utils

import (
	"errors"
	"strconv"

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

// ExistsByID checks if a record exists for the given model and ID.
func ExistsByID(db *gorm.DB, model any, id string) (bool, error) {
	err := db.First(model, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}
