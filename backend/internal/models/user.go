package models

import (
	"time"

	"github.com/google/uuid"
)

// ✅ Corrected User Model using UUID as the primary key
type User struct {
	ID        uuid.UUID `db:"id" json:"id"`
	Email     string    `db:"email" json:"email"`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
