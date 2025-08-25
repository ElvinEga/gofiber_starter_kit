package middlewares

import (
	"github.com/gofiber/fiber/v2"
)

// RequireRole ensures the user has the given role (e.g. “admin”)
func RequireRole(role string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole := c.Locals("userRole") // set in your JWT middleware
		if userRole != role {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "forbidden"})
		}
		return c.Next()
	}
}
