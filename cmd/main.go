package main

import (
	"log"

	"github.com/ElvinEga/gofiber_starter/config"
	"github.com/ElvinEga/gofiber_starter/database"
	"github.com/ElvinEga/gofiber_starter/routes"
	"github.com/ElvinEga/gofiber_starter/utils"
	"github.com/garrettladley/fiberpaginate/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	config.InitConfig()
	database.ConnectDB()
	database.SeedSuperAdmin()
	database.MigrateDB()

	err := utils.InitCloudinary()
	if err != nil {
		log.Fatalf("Failed to initialize Cloudinary: %v", err)
	}

	app := fiber.New(fiber.Config{
		BodyLimit: 10 * 1024 * 1024, // 10MB limit
	})
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}))
	app.Static("/docs", "./docs")
	app.Use(fiberpaginate.New())
	routes.SetupRoutes(app)
	routes.RegisterAdminRoutes(app)

	log.Fatal(app.Listen(":8000"))
}
