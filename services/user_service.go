package services

import (
	"fmt"

	"github.com/ElvinEga/gofiber_starter/database"
	"github.com/ElvinEga/gofiber_starter/models"
	"github.com/ElvinEga/gofiber_starter/responses"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// Profile godoc
// @Summary Get user profile
// @Description Retrieve the user's profile information
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} responses.UserResponse
// @Failure 401 {object} responses.UserResponse
// @Router /api/profile [get]
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

// services/auth_service.go
func SendVerificationEmail(userID uuid.UUID) error {
	user, err := GetUserByID(userID.String())
	if err != nil {
		return err
	}

	token := utils.GenerateSecureToken(32)
	user.VerificationToken = token
	database.DB.Save(user)

	// Send email with verification link
	verificationLink := fmt.Sprintf("https://yourdomain.com/verify?token=%s", token)
	// Use email service to send verificationLink
	return nil
}
