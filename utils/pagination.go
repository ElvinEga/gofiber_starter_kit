package utils

import (
	"strconv"

	"github.com/gofiber/fiber/v3"
)

// PaginationResponse formats a response with pagination metadata.
func PaginationResponse(c fiber.Ctx, data interface{}, totalCount int64) fiber.Map {
	page := parsePositiveInt(c.Query("page"), 1)
	limit := parsePositiveInt(c.Query("limit"), 10)
	totalPages := (totalCount + int64(limit) - 1) / int64(limit)

	return fiber.Map{
		"data": data,
		"meta": fiber.Map{
			"page":        page,
			"limit":       limit,
			"total":       totalCount,
			"total_pages": totalPages,
		},
	}
}

func parsePositiveInt(value string, fallback int) int {
	parsed, err := strconv.Atoi(value)
	if err != nil || parsed <= 0 {
		return fallback
	}
	return parsed
}
