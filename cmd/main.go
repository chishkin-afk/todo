package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/chishkin-afk/todo/internal/common/config"
	"github.com/chishkin-afk/todo/internal/infrastructure/persistence/postgres"
	"github.com/chishkin-afk/todo/internal/infrastructure/session/jwt"
	authservices "github.com/chishkin-afk/todo/internal/modules/auth/application/services"
	userpg "github.com/chishkin-afk/todo/internal/modules/auth/infrastructure/persistence/postgres/user"
	"github.com/chishkin-afk/todo/pkg/consts"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Println(godotenv.Load(".env"))
	l := config.New()
	cfg, err := l.Init(os.Getenv("APP_CONFIG_PATH"))
	fmt.Println(err)
	db, err := postgres.Connect(cfg)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	fmt.Println(postgres.MigrateUP(context.Background(), db))

	ur := userpg.New(db)
	jm := jwt.New(cfg)

	service := authservices.New(
		cfg,
		slog.Default(),
		ur,
		jm,
	)

	ctx := context.WithValue(context.Background(), consts.UserID, uuid.MustParse("36e3532e-ef2c-4498-bc96-3be5b1398391"))
	fmt.Println(service.Delete(ctx))

}
