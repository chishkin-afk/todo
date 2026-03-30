package authservices

import (
	"context"
	"errors"
	"log/slog"

	"github.com/chishkin-afk/todo/internal/application/dtos"
	"github.com/chishkin-afk/todo/internal/common/config"
	"github.com/chishkin-afk/todo/internal/modules/auth/domain/session"
	"github.com/chishkin-afk/todo/internal/modules/auth/domain/user"
	"github.com/chishkin-afk/todo/pkg/consts"
	errs "github.com/chishkin-afk/todo/pkg/errors"
	"github.com/google/uuid"
)

type authService struct {
	cfg                 *config.Config
	log                 *slog.Logger
	userPersistenceRepo user.UserPersistenceRepository
	jwtManager          session.JWTManager
}

func New(
	cfg *config.Config,
	log *slog.Logger,
	userPersistenceRepo user.UserPersistenceRepository,
	jwtManager session.JWTManager,
) *authService {
	return &authService{
		cfg:                 cfg,
		log:                 log,
		userPersistenceRepo: userPersistenceRepo,
		jwtManager:          jwtManager,
	}
}

func (as *authService) Register(ctx context.Context, req *dtos.RegisterRequest) (*dtos.Token, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	user, err := user.New(
		user.Email(req.Email),
		req.Password,
		req.Username,
	)
	if err != nil {
		return nil, err
	}

	savedUser, err := as.userPersistenceRepo.Save(ctx, user)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) ||
			errors.Is(err, context.Canceled) ||
			errors.Is(err, errs.ErrUserAlreadyExists) {
			return nil, err
		}

		as.log.Error("failed to save user into db",
			slog.String("error", err.Error()),
			slog.String("email", user.Email().String()),
		)
		return nil, errs.ErrInternalServer
	}

	token, err := as.jwtManager.Generate(savedUser.ID())
	if err != nil {
		as.log.Error("failed to generate new token",
			slog.String("error", err.Error()),
			slog.String("user_id", savedUser.ID().String()),
		)
		return nil, errs.ErrInternalServer
	}

	return &dtos.Token{
		Token: token,
	}, nil
}

func (as *authService) Login(ctx context.Context, req *dtos.LoginRequest) (*dtos.Token, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	email := user.Email(req.Email)
	if !email.IsValid() {
		return nil, errs.ErrInvalidEmail
	}

	user, err := as.userPersistenceRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) ||
			errors.Is(err, context.Canceled) ||
			errors.Is(err, errs.ErrUserNotFound) {
			return nil, err
		}

		as.log.Error("failed to get user by email",
			slog.String("error", err.Error()),
			slog.String("email", email.String()),
		)
		return nil, errs.ErrInternalServer
	}

	if !user.PasswordHash().Compare(req.Password) {
		return nil, errs.ErrInvalidCredentials
	}

	token, err := as.jwtManager.Generate(user.ID())
	if err != nil {
		as.log.Error("failed to generate new token",
			slog.String("error", err.Error()),
			slog.String("user_id", user.ID().String()),
		)
		return nil, errs.ErrInternalServer
	}

	return &dtos.Token{
		Token: token,
	}, nil
}

func (as *authService) GetSelf(ctx context.Context) (*dtos.User, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	userID, err := as.getUserID(ctx)
	if err != nil {
		return nil, err
	}

	user, err := as.userPersistenceRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) ||
			errors.Is(err, context.Canceled) ||
			errors.Is(err, errs.ErrUserNotFound) {
			return nil, err
		}

		as.log.Error("failed to get user by id",
			slog.String("error", err.Error()),
			slog.String("user_id", userID.String()),
		)
		return nil, errs.ErrInternalServer
	}

	return &dtos.User{
		ID:        user.ID().String(),
		Email:     user.Email().String(),
		Username:  user.Username(),
		CreatedAt: user.CreatedAt().UnixMilli(),
		UpdatedAt: user.UpdatedAt().UnixMilli(),
	}, nil
}

func (as *authService) Update(ctx context.Context, req *dtos.UpdateUserRequest) (*dtos.User, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	userID, err := as.getUserID(ctx)
	if err != nil {
		return nil, err
	}

	oldUser, err := as.userPersistenceRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) ||
			errors.Is(err, context.Canceled) ||
			errors.Is(err, errs.ErrUserNotFound) {
			return nil, err
		}

		as.log.Error("failed to get user by id",
			slog.String("error", err.Error()),
			slog.String("user_id", userID.String()),
		)
		return nil, errs.ErrInternalServer
	}

	if err := as.updateUser(oldUser, req); err != nil {
		return nil, err
	}

	updatedUser, err := as.userPersistenceRepo.Update(ctx, oldUser)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) ||
			errors.Is(err, context.Canceled) ||
			errors.Is(err, errs.ErrUserAlreadyExists) {
			return nil, err
		}

		as.log.Error("failed to update user db",
			slog.String("error", err.Error()),
			slog.String("user_id", oldUser.ID().String()),
		)
		return nil, errs.ErrInternalServer
	}

	return &dtos.User{
		ID:        updatedUser.ID().String(),
		Email:     updatedUser.Email().String(),
		Username:  updatedUser.Username(),
		CreatedAt: updatedUser.CreatedAt().UnixMilli(),
		UpdatedAt: updatedUser.UpdatedAt().UnixMilli(),
	}, nil
}

func (as *authService) Delete(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	userID, err := as.getUserID(ctx)
	if err != nil {
		return err
	}

	if err := as.userPersistenceRepo.Delete(ctx, userID); err != nil {
		if errors.Is(err, context.DeadlineExceeded) ||
			errors.Is(err, context.Canceled) ||
			errors.Is(err, errs.ErrUserNotFound) {
			return err
		}

		as.log.Error("failed to delete user",
			slog.String("error", err.Error()),
			slog.String("user_id", userID.String()),
		)
		return errs.ErrInternalServer
	}

	return nil
}

func (as *authService) updateUser(oldUser *user.User, req *dtos.UpdateUserRequest) error {
	if req.Email != nil {
		if err := oldUser.ChangeEmail(user.Email(*req.Email)); err != nil {
			return err
		}
	}
	if req.Password != nil {
		if err := oldUser.ChangePassword(*req.Password); err != nil {
			return err
		}
	}
	if req.Username != nil {
		if err := oldUser.ChangeUsername(*req.Username); err != nil {
			return err
		}
	}
	return nil
}

func (as *authService) getUserID(ctx context.Context) (uuid.UUID, error) {
	raw := ctx.Value(consts.UserID)
	if raw == nil {
		return uuid.Nil, errs.ErrInvalidToken
	}
	if id, ok := raw.(uuid.UUID); ok {
		return id, nil
	}

	return uuid.Nil, errs.ErrInvalidToken
}
