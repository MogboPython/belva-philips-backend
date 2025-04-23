package utils

import "strconv"

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
