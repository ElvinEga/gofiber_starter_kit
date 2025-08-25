package services

import (
	"github.com/ElvinEga/gofiber_starter/database"
	"github.com/ElvinEga/gofiber_starter/models"
	"github.com/ElvinEga/gofiber_starter/responses"
	"github.com/ElvinEga/gofiber_starter/utils"
	"github.com/gofiber/fiber/v2"
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

func UpdateUser(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	var updateData struct {
		Name     string `json:"name"`
		Username string `json:"username"`
	}

	if err := c.BodyParser(&updateData); err != nil {
		return utils.HandleError(c, fiber.StatusBadRequest, "Invalid input")
	}

	var user models.User
	if err := database.DB.First(&user, "id = ?", userID).Error; err != nil {
		return utils.HandleError(c, fiber.StatusNotFound, "User not found")
	}

	// Update allowed fields only
	if updateData.Name != "" {
		user.Name = updateData.Name
	}
	if updateData.Username != "" {
		// Check if username is already taken
		var existingUser models.User
		if err := database.DB.Where("username = ? AND id != ?", updateData.Username, userID).First(&existingUser).Error; err == nil {
			return utils.HandleError(c, fiber.StatusConflict, "Username already taken")
		}
		user.Username = updateData.Username
	}

	database.DB.Save(&user)
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "User updated successfully",
		"data":    responses.ToUserResponse(user),
	})
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

	if req.CurrentPassword == "" || req.NewPassword == "" {
		return utils.HandleError(c, fiber.StatusBadRequest, "Current password and new password are required")
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

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Password updated successfully",
	})
}

func GetUserByID(id string) (*models.User, error) {
	var user models.User
	err := database.DB.First(&user, "id = ?", id).Error
	return &user, err
}
