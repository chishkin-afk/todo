package app

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/chishkin-afk/todo/internal/common/config"
	redissdk "github.com/chishkin-afk/todo/internal/infrastructure/cache/redis"
	httpserver "github.com/chishkin-afk/todo/internal/infrastructure/http"
	"github.com/chishkin-afk/todo/internal/infrastructure/http/handlers"
	"github.com/chishkin-afk/todo/internal/infrastructure/http/middlewares"
	"github.com/chishkin-afk/todo/internal/infrastructure/persistence/postgres"
	"github.com/chishkin-afk/todo/internal/infrastructure/session/jwt"
	authservices "github.com/chishkin-afk/todo/internal/modules/auth/application/services"
	userredis "github.com/chishkin-afk/todo/internal/modules/auth/infrastructure/cache/redis/user"
	userpg "github.com/chishkin-afk/todo/internal/modules/auth/infrastructure/persistence/postgres/user"
	taskservices "github.com/chishkin-afk/todo/internal/modules/task/application/services"
	groupredis "github.com/chishkin-afk/todo/internal/modules/task/infrastructure/cache/redis/group"
	grouppg "github.com/chishkin-afk/todo/internal/modules/task/infrastructure/persistence/postgres/group"
	taskpg "github.com/chishkin-afk/todo/internal/modules/task/infrastructure/persistence/postgres/task"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

type App struct {
	log *slog.Logger
	srv *httpserver.Server
}

func (a *App) Start() error {
	a.log.Info("server is running")
	return a.srv.Start()
}

func (a *App) Shutdown(ctx context.Context) error {
	a.log.Info("server shutdown")
	return a.srv.Shutdown(ctx)
}

func New() (*App, func(), error) {
	if err := godotenv.Load(".env"); err != nil {
		slog.Warn(".env file doesn't exist",
			slog.String("error", err.Error()),
		)
	}

	cfg, err := loadConfig(os.Getenv("APP_CONFIG_PATH"))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load config: %w", err)
	}

	slog.Info("config was loaded", slog.Any("server", cfg.Server))

	db, err := providePersistence(cfg)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to provide persistence: %w", err)
	}

	client, err := provideCache(cfg)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to provide cache: %w", err)
	}

	server, err := provideServer(cfg, slog.Default(), db, client)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to provide server: %w", err)
	}

	return &App{
			log: slog.Default(),
			srv: server,
		}, func() {
			if err := postgres.Close(db); err != nil {
				slog.Error("failed to close connection with postgres", slog.String("error", err.Error()))
			} else {
				slog.Info("connection with postgres was closed")
			}

			if err := redissdk.Close(client); err != nil {
				slog.Error("failed to close connection with redis", slog.String("error", err.Error()))
			} else {
				slog.Info("connection with redis was closed")
			}
		}, nil
}

func loadConfig(path string) (*config.Config, error) {
	loader := config.New()
	return loader.Init(path)
}

func providePersistence(cfg *config.Config) (*sql.DB, error) {
	db, err := postgres.Connect(cfg)
	if err != nil {
		return nil, err
	}

	if err := postgres.MigrateUP(context.Background(), db); err != nil {
		return nil, err
	}

	return db, nil
}

func provideCache(cfg *config.Config) (*redis.Client, error) {
	return redissdk.Connect(cfg)
}

func provideServer(cfg *config.Config, log *slog.Logger, db *sql.DB, client *redis.Client) (*httpserver.Server, error) {
	jwtManager := jwt.New(cfg)

	authService := authservices.New(
		cfg,
		log,
		userpg.New(db),
		userredis.New(cfg, client),
		jwtManager,
	)

	taskService := taskservices.New(
		cfg,
		log,
		taskpg.New(db),
		grouppg.New(db),
		groupredis.New(cfg, client),
	)

	handler, err := handlers.New(cfg, authService, taskService, []gin.HandlerFunc{
		middlewares.AuthMiddleware(jwtManager, map[string]bool{
			"/api/v1/register": true,
			"/api/v1/login":    true,
			"/swagger/*any":    true,
		}),
		middlewares.RateLimitMiddleware(10, time.Second),
	})
	if err != nil {
		return nil, err
	}

	return httpserver.New(cfg, handler), nil
}
