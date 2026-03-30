package userpg

import (
	"time"

	"github.com/google/uuid"
)

type UserModel struct {
	ID           uuid.UUID
	Email        string
	PasswordHash string
	Username     string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
