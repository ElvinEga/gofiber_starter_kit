package controllers

import (
	"github.com/ElvinEga/gofiber_starter/services"
	"github.com/gofiber/fiber/v2"
)

func Register(c *fiber.Ctx) error {
	return services.Register(c)
}

func Login(c *fiber.Ctx) error {
	return services.Login(c)
}

func GoogleSSO(c *fiber.Ctx) error {
	return services.GoogleLogin(c)
}

func GoogleCallback(c *fiber.Ctx) error {
	return services.GoogleCallback(c)
}

func Logout(c *fiber.Ctx) error {
	return services.Logout(c)
}
