package api

import (
	"fmt"
	"strconv"
)

func parseLimit(s string) (int, error) {
	limit, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("invalid limit format: %w", err)
	}
	if limit < 1 || limit > 100 {
		return 0, fmt.Errorf("limit must be between 1 and 100")
	}
	return limit, nil
}

func validateDate(dateStr string) (bool, string) {
	if dateStr == "" {
		return false, "date cannot be empty"
	}
	return true, ""
}
