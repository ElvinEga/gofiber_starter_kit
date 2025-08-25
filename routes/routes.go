package routes

import (
	"github.com/ElvinEga/gofiber_starter/controllers"
	"github.com/ElvinEga/gofiber_starter/middlewares"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	auth := api.Group("/auth")
	auth.Post("/register", controllers.Register)
	auth.Post("/login", controllers.Login)
	auth.Post("/google", controllers.GoogleSSO)
	auth.Post("/logout", controllers.Logout)

	user := api.Group("/user", middlewares.JWTProtected())
	user.Get("/profile", controllers.GetUserProfile)

	// Tour Routes
	api.Get("/tours", controllers.GetAllTours)
	api.Get("/tours/:id", controllers.GetTourByID)
	api.Get("/tours/featured", controllers.GetFeaturedTours)
	api.Get("/tours/filter", controllers.GetFilteredTours)
	api.Get("/tours/:id/reviews", controllers.GetTourReviews)
	api.Post("/tours/:id/reviews", middlewares.JWTProtected(), controllers.CreateTourReview)

	api.Get("/destinations", controllers.GetAllDestinations)
	api.Get("/destinations/:id", controllers.GetDestinationByID)

	api.Get("/categories", controllers.GetAllCategories)
	api.Get("/categories/:id", controllers.GetCategoryByID)

}
