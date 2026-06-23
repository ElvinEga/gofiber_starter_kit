package database

import (
	"github.com/ElvinEga/gofiber_starter/models"
)

func MigrateDB() {
	err := DB.AutoMigrate(
		&models.User{},
		&models.RefreshToken{},
		// Add other models here
	)
	if err != nil {
		panic("Failed to migrate database")
	}
}
