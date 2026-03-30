package userredis

import (
	"context"
	"errors"
	"fmt"

	"github.com/chishkin-afk/todo/internal/common/config"
	"github.com/chishkin-afk/todo/internal/modules/auth/domain/user"
	errs "github.com/chishkin-afk/todo/pkg/errors"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type userCacheRepository struct {
	cfg    *config.Config
	client *redis.Client
}

func New(cfg *config.Config, client *redis.Client) *userCacheRepository {
	return &userCacheRepository{
		cfg:    cfg,
		client: client,
	}
}

func (ucr *userCacheRepository) Save(ctx context.Context, user *user.User) error {
	bytes, err := ToBytes(user)
	if err != nil {
		return fmt.Errorf("failed to convert domain into bytes: %w", err)
	}

	key := ucr.getKey(user.ID())
	if err := ucr.client.Set(ctx, key, bytes, ucr.cfg.Redis.UserTTL).Err(); err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return err
		}

		return fmt.Errorf("failed to save user into cache: %w", err)
	}

	return nil
}

func (ucr *userCacheRepository) Get(ctx context.Context, id uuid.UUID) (*user.User, error) {
	key := ucr.getKey(id)
	bytes, err := ucr.client.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, err
		}

		if errors.Is(err, redis.Nil) {
			return nil, errs.ErrUserNotFound
		}

		return nil, fmt.Errorf("failed to get user from cache: %w", err)
	}

	return ToDomain(bytes)
}

func (ucr *userCacheRepository) Del(ctx context.Context, id uuid.UUID) error {
	key := ucr.getKey(id)
	if err := ucr.client.Del(ctx, key).Err(); err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return err
		}

		return fmt.Errorf("failed to delete user from cache: %w", err)
	}

	return nil
}

func (ucr *userCacheRepository) getKey(id uuid.UUID) string {
	return fmt.Sprintf("us:%s", id.String())
}
