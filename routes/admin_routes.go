package routes

import (
	"github.com/ElvinEga/gofiber_starter/controllers"
	"github.com/ElvinEga/gofiber_starter/middlewares"
	"github.com/gofiber/fiber/v2"
)

func RegisterAdminRoutes(app *fiber.App) {
	admin := app.Group("/admin", middlewares.AdminOnly)
	admin.Post("/login", controllers.AdminLogin)

	// Tour Routes
	admin.Get("/tours", controllers.GetAllTours)
	admin.Get("/tours/:id", controllers.GetTourByID)
	admin.Post("/tours", controllers.CreateTour)
	admin.Put("/tours/:id", controllers.UpdateTour)
	admin.Delete("/tours/:id", controllers.DeleteTour)

	//Events Routes
	admin.Get("/events", controllers.GetAllEvents)
	admin.Get("/events/:id", controllers.GetEventByID)
	admin.Post("/events", controllers.CreateEvent)
	admin.Put("/events/:id", controllers.UpdateEvent)
	admin.Delete("/events/:id", controllers.DeleteEvent)

	// Destination Routes
	admin.Get("/destinations", controllers.GetAllDestinations)
	admin.Get("/destinations/:id", controllers.GetDestinationByID)
	admin.Post("/destinations", controllers.CreateDestination)
	admin.Put("/destinations/:id", controllers.UpdateDestination)
	admin.Delete("/destinations/:id", controllers.DeleteDestination)

	// Category Routes
	admin.Get("/categories", controllers.GetAllCategories)
	admin.Get("/categories/:id", controllers.GetCategoryByID)
	admin.Post("/categories", controllers.CreateCategory)
	admin.Put("/categories/:id", controllers.UpdateCategory)
	admin.Delete("/categories/:id", controllers.DeleteCategory)

	// Review Routes
	admin.Get("/reviews", controllers.GetAllReviews)
	admin.Get("/reviews/:id", controllers.GetReviewByID)
	admin.Post("/reviews", controllers.CreateReview)
	admin.Put("/reviews/:id", controllers.UpdateReview)
	admin.Delete("/reviews/:id", controllers.DeleteReview)

	admin.Put("/me/password", controllers.UpdateAdminPassword)
	admin.Get("/user/me", controllers.GetCurrentAdminProfile)

	userAdmin := admin.Group("/users")
	userAdmin.Get("/", controllers.GetAllUsers)
	userAdmin.Get("/:id", controllers.GetUserByID)
	userAdmin.Post("/", controllers.CreateUser)
	userAdmin.Put("/:id", controllers.UpdateUser)
	userAdmin.Delete("/:id", controllers.DeleteUser)

	adminUsers := admin.Group("/managers", middlewares.SuperAdminOnly())
	adminUsers.Get("/", controllers.ListAdmins)
	adminUsers.Post("/", controllers.CreateAdmin)
	adminUsers.Put("/:id", controllers.UpdateAdmin)
	adminUsers.Delete("/:id", controllers.DeleteAdmin)

}
