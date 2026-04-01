package errs

import "errors"

var (
	// Domain's
	ErrInvalidID       = errors.New("invalid id")
	ErrInvalidEmail    = errors.New("invalid email of user")
	ErrInvalidPassword = errors.New("len of password must be more than 6 and less than 32")
	ErrInvalidUsername = errors.New("len of username must be more than 3 and less than 64")

	ErrInvalidTitle        = errors.New("len of title must be more than 3 and less than 64")
	ErrInvalidTaskDesc     = errors.New("len of desc must be more than 3 and less than 512")
	ErrInvalidTaskPriority = errors.New("invalid priority")
	ErrTaskAlreadyDone     = errors.New("task is already done")
	ErrTaskNotDone         = errors.New("task isn't done")

	// Repository's
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrGroupNotFound     = errors.New("group not found")
	ErrTaskNotFound      = errors.New("task not found")
	ErrDepsNotFound      = errors.New("some deps are not found")

	// General's
	ErrInvalidToken      = errors.New("invalid token")
	ErrTooManyRequests   = errors.New("too many requests")

	// Service's
	ErrInternalServer     = errors.New("internal server error")
	ErrNotEnoughRights    = errors.New("not enough rights")
	ErrInvalidCredentials = errors.New("invalid credentials")
)