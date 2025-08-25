package routes

import (
	"github.com/ElvinEga/gofiber_starter/controllers"
	"github.com/ElvinEga/gofiber_starter/middlewares"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	// Apply security headers and rate limiting globally
	app.Use(middlewares.SecurityHeaders())
	app.Use(middlewares.RateLimit())

	// API group
	api := app.Group("/api")

	// Auth routes (public)
	auth := api.Group("/auth")
	auth.Post("/register", controllers.Register)
	auth.Post("/login", controllers.Login)
	auth.Post("/google", controllers.GoogleSSO)
	auth.Get("/google/callback", controllers.GoogleCallback)
	auth.Post("/refresh", controllers.RefreshToken)
	auth.Get("/verify", controllers.VerifyEmail)
	auth.Post("/forgot-password", controllers.RequestPasswordReset)
	auth.Post("/reset-password", controllers.ResetPassword)

	// Protected routes
	protected := api.Group("/", middlewares.JWTProtected())

	// User routes
	user := protected.Group("/user")
	user.Get("/profile", controllers.GetUserProfile)
	user.Put("/profile", controllers.UpdateUser)
	user.Put("/password", controllers.ChangePassword)

	// Logout route (protected)
	protected.Post("/logout", controllers.Logout)
}
