package session

import "github.com/google/uuid"

type JWTManager interface {
	Generate(userID uuid.UUID) (string, error)
	Validate(tokenString string) (uuid.UUID, error)
}
