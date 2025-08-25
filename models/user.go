package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID              uuid.UUID `gorm:"type:text;primaryKey" json:"id"`
	Name            string    `json:"name"`
	Username        string    `gorm:"uniqueIndex" json:"username"`
	Email           string    `gorm:"uniqueIndex" json:"email"`
	Password        string    `json:"-"`
	Role            string    `json:"role"`
	IsVerified      bool      `json:"is_verified"`
	EmailVerifiedAt time.Time `json:"is_email_verified"`
	Base
}
