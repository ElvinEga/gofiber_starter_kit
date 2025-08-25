package models

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID                uuid.UUID `gorm:"type:text;primaryKey" json:"id"`
	Name              string    `json:"name"`
	Username          string    `gorm:"uniqueIndex" json:"username"`
	Email             string    `gorm:"uniqueIndex" json:"email"`
	Password          string    `json:"-"`
	Role              string    `json:"role"`
	IsVerified        bool      `json:"is_verified"`
	EmailVerifiedAt   time.Time `json:"email_verified_at"`
	VerificationToken string    `json:"-"`
	ResetToken        string    `json:"-"`
	ResetExpiresAt    time.Time `json:"reset_expires_at"`
	Base
}
