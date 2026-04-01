package userpg

import (
	"database/sql"
	"errors"

	"github.com/chishkin-afk/todo/internal/modules/auth/domain/user"
	errs "github.com/chishkin-afk/todo/pkg/errors"
)

func ToDomain(result *sql.Row) (*user.User, error) {
	var model UserModel
	if err := result.Scan(
		&model.ID,
		&model.Email,
		&model.PasswordHash,
		&model.Username,
		&model.CreatedAt,
		&model.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrUserNotFound
		}

		return nil, err
	}

	return user.From(
		model.ID,
		user.Email(model.Email),
		user.PasswordHash(model.PasswordHash),
		model.Username,
		model.CreatedAt,
		model.UpdatedAt,
	)
}
