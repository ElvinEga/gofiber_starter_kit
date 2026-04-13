package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RefreshToken struct {
	ID        uuid.UUID      `gorm:"type:text;primaryKey" json:"id"`
	UserID    uuid.UUID      `gorm:"type:text;index" json:"user_id"`
	Token     string         `gorm:"uniqueIndex" json:"token"`
	ExpiresAt time.Time      `json:"expires_at"`
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}
