package errs

import "errors"

var (
	// Domain's
	ErrInvalidID       = errors.New("invalid id of user")
	ErrInvalidEmail    = errors.New("invalid email of user")
	ErrInvalidPassword = errors.New("len of password must be more than 6 and less than 32")
	ErrInvalidUsername = errors.New("len of username must be more than 3 and less than 64")

	// Repository's
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")

	// General's
	ErrInvalidToken = errors.New("invalid token")

	// Service's
	ErrInternalServer     = errors.New("internal server error")
	ErrInvalidCredentials = errors.New("invalid credentials")
)
