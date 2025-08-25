package services

import (
	"fmt"

	"github.com/ElvinEga/gofiber_starter/database"
	"github.com/ElvinEga/gofiber_starter/models"
	"github.com/ElvinEga/gofiber_starter/responses"
	"github.com/ElvinEga/gofiber_starter/utils"
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

// controllers/user_controller.go
func UpdateUser(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	var updateData models.User

	if err := c.BodyParser(&updateData); err != nil {
		return utils.HandleError(c, fiber.StatusBadRequest, "Invalid input")
	}

	var user models.User
	if err := database.DB.First(&user, "id = ?", userID).Error; err != nil {
		return utils.HandleError(c, fiber.StatusNotFound, "User not found")
	}

	// Update allowed fields only
	user.Name = updateData.Name
	user.Username = updateData.Username
	database.DB.Save(&user)

	return c.JSON(responses.ToUserResponse(user))
}

func ChangePassword(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	var req struct {
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password"`
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.HandleError(c, fiber.StatusBadRequest, "Invalid input")
	}

	var user models.User
	if err := database.DB.First(&user, "id = ?", userID).Error; err != nil {
		return utils.HandleError(c, fiber.StatusNotFound, "User not found")
	}

	if !utils.CheckPasswordHash(req.CurrentPassword, user.Password) {
		return utils.HandleError(c, fiber.StatusUnauthorized, "Incorrect current password")
	}

	user.Password = utils.HashPassword(req.NewPassword)
	database.DB.Save(&user)

	return c.JSON(fiber.Map{"status": "success", "message": "Password updated"})
}
