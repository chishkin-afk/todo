package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/chishkin-afk/todo/internal/common/config"
	"github.com/chishkin-afk/todo/internal/infrastructure/persistence/postgres"
	userpg "github.com/chishkin-afk/todo/internal/modules/auth/infrastructure/persistence/postgres/user"
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
	fmt.Println(ur.Delete(context.Background(), uuid.MustParse("38a17674-6d33-48c1-b4c4-7f7673d103c0")))
}
