package main

import (
	"log"

	"github.com/ElvinEga/gofiber_starter/config"
	"github.com/ElvinEga/gofiber_starter/database"
	"github.com/ElvinEga/gofiber_starter/routes"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

// @title Go Fiber Starter Kit API
// @version 1.0
// @description This is a sample Go Fiber API server.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email your-email@example.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8000
// @BasePath /
func main() {
	config.InitConfig()
	database.ConnectDB()
	database.SeedSuperAdmin()
	database.MigrateDB()

	app := fiber.New(fiber.Config{
		BodyLimit: 10 * 1024 * 1024, // 10MB limit
	})
	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	}))
	app.Get("/swagger/*", fiberSwagger.WrapHandler)
	// app.Static("/docs", "./docs")
	routes.SetupRoutes(app)

	log.Fatal(app.Listen(":8000"))
}
