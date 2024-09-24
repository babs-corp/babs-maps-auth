package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	PassHash  []byte    `json:"password_hash" db:"pass_hash"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
