package controllers

import (
	"github.com/ElvinEga/gofiber_starter/services"
	"github.com/gofiber/fiber/v2"
)

func GetUserProfile(c *fiber.Ctx) error {
	return services.GetUserProfile(c)
}
