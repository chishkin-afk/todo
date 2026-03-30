package jwt

import (
	"fmt"
	"time"

	"github.com/chishkin-afk/todo/internal/common/config"
	errs "github.com/chishkin-afk/todo/pkg/errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type customClaims struct {
	jwt.RegisteredClaims
	UserID uuid.UUID `json:"user_id"`
}

type jwtManager struct {
	cfg *config.Config
}

func New(cfg *config.Config) *jwtManager {
	return &jwtManager{
		cfg: cfg,
	}
}

func (jm *jwtManager) Generate(userID uuid.UUID) (string, error) {
	secretKey := []byte(jm.cfg.Session.SecretKey)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &customClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "todo-list",
			Subject:   "todo-user",
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(jm.cfg.Session.TokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			ID:        uuid.NewString(),
		},
	})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

func (jm *jwtManager) Validate(tokenString string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &customClaims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errs.ErrInvalidToken
		}

		return []byte(jm.cfg.Session.SecretKey), nil
	})
	if err != nil {
		return uuid.Nil, errs.ErrInvalidToken
	}

	if claims, ok := token.Claims.(*customClaims); ok {
		return claims.UserID, nil
	}

	return uuid.Nil, errs.ErrInvalidToken
}
