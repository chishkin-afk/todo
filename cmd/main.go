package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/chishkin-afk/todo/internal/common/config"
	"github.com/chishkin-afk/todo/internal/infrastructure/persistence/postgres"
	grouppg "github.com/chishkin-afk/todo/internal/modules/task/infrastructure/persistence/postgres/group"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")
	cfg, err := config.New().Init(os.Getenv("APP_CONFIG_PATH"))
	fmt.Println(err)
	db, err := postgres.Connect(cfg)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	fmt.Println(postgres.MigrateUP(context.TODO(), db))

	gr := grouppg.New(db)
	fmt.Println(gr.GetByID(context.Background(), uuid.MustParse("25b7d21a-b58d-46f6-9fc1-9e6e9054964e")))
}
