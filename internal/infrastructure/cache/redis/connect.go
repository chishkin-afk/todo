package redissdk

import (
	"context"
	"fmt"

	"github.com/chishkin-afk/todo/internal/common/config"
	"github.com/redis/go-redis/v9"
)

func Connect(cfg *config.Config) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Username: cfg.Redis.Auth.Username,
		Password: cfg.Redis.Auth.Password,
		DB:       cfg.Redis.Auth.DB,
	})
	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("failed to ping redis: %w", err)
	}

	return client, nil
}

func Close(client *redis.Client) error {
	if err := client.Close(); err != nil {
		return fmt.Errorf("failed to close client: %w", err)
	}

	return nil
}
