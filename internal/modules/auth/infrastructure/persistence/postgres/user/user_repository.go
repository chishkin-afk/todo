package userpg

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/chishkin-afk/todo/internal/modules/auth/domain/user"
	errs "github.com/chishkin-afk/todo/pkg/errors"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type userPersistenceRepository struct {
	db *sql.DB
}

func New(db *sql.DB) *userPersistenceRepository {
	return &userPersistenceRepository{
		db: db,
	}
}

func (ur *userPersistenceRepository) Save(ctx context.Context, user *user.User) (*user.User, error) {
	_, err := ur.db.ExecContext(ctx, `insert into users (
		id,
		email,
		password_hash,
		username,
		created_at,
		updated_at) values (
		$1,
		$2,
		$3,
		$4,
		$5,
		$6)`, user.ID(), user.Email().String(), user.PasswordHash().String(), user.Username(), user.CreatedAt(), user.UpdatedAt())
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, err
		}

		if err, ok := ur.knownError(err); ok {
			return nil, err
		}

		return nil, fmt.Errorf("failed to save user: %w", err)
	}

	return user, nil
}

func (ur *userPersistenceRepository) Update(ctx context.Context, user *user.User) (*user.User, error) {
	_, err := ur.db.ExecContext(ctx, `update users set 
		password_hash = $1,
		username = $2,
		email = $3,
		updated_at = $4 where id = $5
	`, user.PasswordHash(), user.Username(), user.Email(), user.UpdatedAt(), user.ID())
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, err
		}

		if err, ok := ur.knownError(err); ok {
			return nil, err
		}

		return nil, fmt.Errorf("failed to update fields of user: %w", err)
	}

	return user, nil
}

func (ur *userPersistenceRepository) GetByEmail(ctx context.Context, email user.Email) (*user.User, error) {
	result := ur.db.QueryRowContext(ctx, "select id, email, password_hash, username, created_at, updated_at from users where email = $1", email)
	if err := result.Err(); err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, err
		}

		if err, ok := ur.knownError(err); ok {
			return nil, err
		}

		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return ToDomain(result)
}

func (ur *userPersistenceRepository) GetByID(ctx context.Context, id uuid.UUID) (*user.User, error) {
	result := ur.db.QueryRowContext(ctx, "select id, email, password_hash, username, created_at, updated_at from users where id = $1", id)
	if err := result.Err(); err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, err
		}

		if err, ok := ur.knownError(err); ok {
			return nil, err
		}

		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return ToDomain(result)
}

func (ur *userPersistenceRepository) Delete(ctx context.Context, id uuid.UUID) error {
	res, err := ur.db.ExecContext(ctx, "delete from users where id = $1", id)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return err
		}

		return fmt.Errorf("failed to delete user: %w", err)
	}

	if affected, _ := res.RowsAffected(); affected == 0 {
		return errs.ErrUserNotFound
	}

	return nil
}

func (ur *userPersistenceRepository) knownError(err error) (error, bool) {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		switch pqErr.Code {
		case "23505":
			return errs.ErrUserAlreadyExists, true
			// another errors
		}
	}

	return err, false
}
