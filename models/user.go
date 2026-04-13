package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID                uuid.UUID      `gorm:"type:text;primaryKey" json:"id"`
	Name              string         `json:"name"`
	Username          string         `gorm:"uniqueIndex" json:"username"`
	Email             string         `gorm:"uniqueIndex" json:"email"`
	Password          string         `json:"-"`
	Role              string         `json:"role"`
	IsVerified        bool           `json:"is_verified"`
	EmailVerifiedAt   time.Time      `json:"email_verified_at"`
	VerificationToken string         `json:"-"`
	ResetToken        string         `json:"-"`
	ResetExpiresAt    time.Time      `json:"reset_expires_at"`
	CreatedAt         time.Time      `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}
