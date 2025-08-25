package utils

import (
	"github.com/garrettladley/fiberpaginate/v2"
	"github.com/gofiber/fiber/v2"
)

// PaginationResponse formats a response with pagination metadata
func PaginationResponse(c *fiber.Ctx, data interface{}, totalCount int64) fiber.Map {
	// Get pagination info from context
	pageInfo, ok := fiberpaginate.FromContext(c)
	if !ok {
		// If pagination info is not available, use default values
		pageInfo = &fiberpaginate.PageInfo{
			Page:  1,
			Limit: 10,
		}
	}

	// Calculate total pages
	totalPages := (totalCount + int64(pageInfo.Limit) - 1) / int64(pageInfo.Limit)

	// Return formatted response
	return fiber.Map{
		"data": data,
		"meta": fiber.Map{
			"page":        pageInfo.Page,
			"limit":       pageInfo.Limit,
			"total":       totalCount,
			"total_pages": totalPages,
		},
	}
}
