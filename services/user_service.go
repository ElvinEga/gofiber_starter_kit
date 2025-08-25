package services

import (
	"github.com/ElvinEga/gofiber_starter/database"
	"github.com/ElvinEga/gofiber_starter/models"
	"github.com/ElvinEga/gofiber_starter/responses"
	"github.com/gofiber/fiber/v2"
)

func GetUserProfile(c *fiber.Ctx) error {
	userId := c.Locals("userID").(string)
	var user models.User

	if err := database.DB.First(&user, "id = ?", userId).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}
	return c.JSON(responses.ToUserResponse(user))
}

func GetUserByID(id string) (*models.User, error) {
	var user models.User
	err := database.DB.First(&user, "id = ?", id).Error
	return &user, err
}
