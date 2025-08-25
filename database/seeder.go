package database

import (
	"fmt"
	"github.com/ElvinEga/gofiber_starter/models"
	"github.com/ElvinEga/gofiber_starter/utils"
)

func SeedSuperAdmin() {
	var count int64
	DB.Model(&models.User{}).Where("role = ?", "superadmin").Count(&count)
	if count > 0 {
		fmt.Println("✅ Superadmin already exists")
		return
	}

	superadmin := models.User{
		ID:         utils.GenerateUUID(),
		Name:       "Super Admin",
		Email:      "admin@example.com",
		Username:   "superadmin",
		Password:   utils.HashPassword("admin1234"),
		Role:       "superadmin",
		IsVerified: true,
	}

	if err := DB.Create(&superadmin).Error; err != nil {
		fmt.Println("❌ Failed to seed superadmin:", err)
		return
	}
	fmt.Println("✅ Superadmin seeded: email=admin@example.com, password=admin1234")
}
