package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ✅ User Model (Main Table)
type User struct {
	ID        uuid.UUID      `gorm:"primaryKey"`
	Email     string         `gorm:"unique;not null"`
	Username  string         `gorm:"unique;not null"`
	Password  string         `gorm:"not null"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// ✅ Email Change Request Model
type UserEmailChange struct {
	ID        uuid.UUID `gorm:"primaryKey"`
	UserID    uuid.UUID `gorm:"index;not null;constraint:OnDelete:CASCADE"` // Foreign key reference to User
	NewEmail  string    `gorm:"unique;not null"`
	Token     string    `gorm:"not null"`
	CreatedAt time.Time
}

// ✅ User Sessions Model
type UserSession struct {
	ID        uuid.UUID `gorm:"primaryKey"`
	UserID    uuid.UUID `gorm:"index;not null;constraint:OnDelete:CASCADE"` // Foreign key reference to User
	TokenHash string    `gorm:"not null"`
	IPAddress string
	UserAgent string
	CreatedAt time.Time
}
