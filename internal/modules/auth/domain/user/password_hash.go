package user

import (
	"strings"

	errs "github.com/chishkin-afk/todo/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type PasswordHash string

func (ph PasswordHash) String() string {
	return string(ph)
}

func (ph PasswordHash) Compare(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(ph), []byte(password)) == nil
}

func NewPasswordHash(raw string) (PasswordHash, error) {
	raw = strings.TrimSpace(raw)
	if len([]rune(raw)) < 6 || len([]rune(raw)) > 32 {
		return "", errs.ErrInvalidPassword
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(raw), bcrypt.DefaultCost)
	if err != nil {
		return "", errs.ErrInvalidPassword
	}

	return PasswordHash(hash), nil
}
