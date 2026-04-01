package groupredis

import (
	"context"
	"errors"
	"fmt"

	"github.com/chishkin-afk/todo/internal/common/config"
	"github.com/chishkin-afk/todo/internal/modules/task/domain/group"
	errs "github.com/chishkin-afk/todo/pkg/errors"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type groupCacheRepository struct {
	cfg    *config.Config
	client *redis.Client
}

func New(cfg *config.Config, client *redis.Client) *groupCacheRepository {
	return &groupCacheRepository{
		cfg:    cfg,
		client: client,
	}
}

func (gcr *groupCacheRepository) Save(ctx context.Context, group *group.Group) error {
	bytes, err := ToBytes(group)
	if err != nil {
		return fmt.Errorf("failed to convert domain into bytes: %w", err)
	}

	if err := gcr.client.Set(ctx, gcr.getKey(group.ID()), bytes, gcr.cfg.Redis.GroupTTL).Err(); err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return err
		}

		return fmt.Errorf("failed to save group into cache: %w", err)
	}

	return nil
}

func (gcr *groupCacheRepository) Get(ctx context.Context, id uuid.UUID) (*group.Group, error) {
	bytes, err := gcr.client.Get(ctx, gcr.getKey(id)).Bytes()
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, err
		}

		if errors.Is(err, redis.Nil) {
			return nil, errs.ErrGroupNotFound
		}

		return nil, fmt.Errorf("failed to get group from cache: %w", err)
	}

	return ToDomain(bytes)
}

func (gcr *groupCacheRepository) Del(ctx context.Context, id uuid.UUID) error {
	if err := gcr.client.Del(ctx, gcr.getKey(id)).Err(); err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return err
		}

		return fmt.Errorf("failed to delete group from cache: %w", err)
	}

	return nil
}

func (gcr *groupCacheRepository) getKey(id uuid.UUID) string {
	return fmt.Sprintf("group:%s", id.String())
}
