package middlewares

import (
	"github.com/ElvinEga/gofiber_starter/utils"
	"github.com/gofiber/fiber/v2"
)

func JWTProtected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID, role, err := utils.VerifyJWT(c)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "error",
				"message": "Unauthorized",
			})
		}
		c.Locals("userID", userID)
		c.Locals("userRole", role)
		return c.Next()
	}
}
