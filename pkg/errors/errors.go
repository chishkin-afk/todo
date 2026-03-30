package errs

import "errors"

var (
	// Domain's
	ErrInvalidID       = errors.New("invalid id of user")
	ErrInvalidEmail    = errors.New("invalid email of user")
	ErrInvalidPassword = errors.New("len of password must be more than 6 and less than 32")
	ErrInvalidUsername = errors.New("len of username must be more than 3 and less than 64")
)
