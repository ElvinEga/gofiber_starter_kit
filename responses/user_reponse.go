package responses

import (
	"time"

	"github.com/ElvinEga/gofiber_starter/models"
)

type UserResponse struct {
	ID         string    `json:"id"`
	Email      string    `json:"email"`
	Name       string    `json:"name"`
	Username   string    `json:"username"`
	Role       string    `json:"role"`
	IsVerified bool      `json:"is_verified"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// Converts a models.User into the public response.
func ToUserResponse(u models.User) UserResponse {
	return UserResponse{
		ID:         u.ID.String(),
		Email:      u.Email,
		Name:       u.Name,
		Username:   u.Username,
		Role:       u.Role,
		IsVerified: u.IsVerified,
		CreatedAt:  u.CreatedAt,
		UpdatedAt:  u.UpdatedAt,
	}
}

type AuthResponse struct {
	Status  string        `json:"status"`
	Message string        `json:"message"`
	Token   string        `json:"token,omitempty"` // Token is optional
	User    *UserResponse `json:"user,omitempty"`  // User is optional
}
