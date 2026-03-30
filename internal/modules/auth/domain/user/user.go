package user

import (
	"strings"
	"time"

	errs "github.com/chishkin-afk/todo/pkg/errors"
	"github.com/google/uuid"
)

type User struct {
	id           uuid.UUID
	email        Email
	passwordHash PasswordHash
	username     string

	createdAt time.Time
	updatedAt time.Time
}

func New(
	email Email,
	password string,
	username string,
) (*User, error) {
	passwordHash, err := NewPasswordHash(password)
	if err != nil {
		return nil, err
	}

	if !email.IsValid() {
		return nil, errs.ErrInvalidEmail
	}

	username = strings.TrimSpace(username)
	if len([]rune(username)) < 3 || len([]rune(username)) > 64 {
		return nil, errs.ErrInvalidUsername
	}

	return &User{
		id:           uuid.New(),
		email:        email,
		passwordHash: passwordHash,
		username:     username,
		createdAt:    time.Now().UTC(),
		updatedAt:    time.Now().UTC(),
	}, nil
}

func From(
	id uuid.UUID,
	email Email,
	passwordHash PasswordHash,
	username string,
	createdAt, updatedAt time.Time,
) (*User, error) {
	if id == uuid.Nil {
		return nil, errs.ErrInvalidID
	}

	if !email.IsValid() {
		return nil, errs.ErrInvalidEmail
	}

	username = strings.TrimSpace(username)
	if len([]rune(username)) < 3 || len([]rune(username)) > 64 {
		return nil, errs.ErrInvalidUsername
	}

	return &User{
		id:           id,
		email:        email,
		passwordHash: passwordHash,
		username:     username,
		createdAt:    createdAt,
		updatedAt:    updatedAt,
	}, nil
}

func (u *User) ChangePassword(newPassword string) error {
	newHash, err := NewPasswordHash(newPassword)
	if err != nil {
		return err
	}

	u.passwordHash = newHash
	u.updatedAt = time.Now().UTC()
	return nil
}

func (u *User) ChangeUsername(newUsername string) error {
	newUsername = strings.TrimSpace(newUsername)
	if len([]rune(newUsername)) < 3 || len([]rune(newUsername)) > 64 {
		return errs.ErrInvalidUsername
	}

	u.username = newUsername
	u.updatedAt = time.Now().UTC()
	return nil
}

func (u *User) ChangeEmail(newEmail Email) error {
	if !newEmail.IsValid() {
		return errs.ErrInvalidEmail
	}

	u.email = newEmail
	u.updatedAt = time.Now().UTC()
	return nil
}

func (u *User) ID() uuid.UUID {
	return u.id
}

func (u *User) Email() Email {
	return u.email
}

func (u *User) PasswordHash() PasswordHash {
	return u.passwordHash
}

func (u *User) Username() string {
	return u.username
}

func (u *User) CreatedAt() time.Time {
	return u.createdAt
}

func (u *User) UpdatedAt() time.Time {
	return u.updatedAt
}
