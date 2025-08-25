package utils

import (
	"github.com/gofiber/fiber/v2"
)

type ErrorResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

func HandleError(c *fiber.Ctx, status int, message string, details ...interface{}) error {
	response := ErrorResponse{
		Status:  "error",
		Message: message,
	}

	if len(details) > 0 {
		response.Details = details[0]
	}

	return c.Status(status).JSON(response)
}

func HandleSuccess(c *fiber.Ctx, message string, data ...interface{}) error {
	response := fiber.Map{
		"status":  "success",
		"message": message,
	}

	if len(data) > 0 {
		response["data"] = data[0]
	}

	return c.JSON(response)
}
