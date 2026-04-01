package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/chishkin-afk/todo/internal/common/config"
	"github.com/chishkin-afk/todo/internal/infrastructure/persistence/postgres"
	grouppg "github.com/chishkin-afk/todo/internal/modules/task/infrastructure/persistence/postgres/group"
	taskpg "github.com/chishkin-afk/todo/internal/modules/task/infrastructure/persistence/postgres/task"
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
	fmt.Println(gr.GetByID(context.Background(), uuid.MustParse("873b6077-ccca-473e-b4a4-61abcb3c7ee2")))

	tr := taskpg.New(db)

	fmt.Println(tr.Delete(context.TODO(), uuid.MustParse("b4e1aec2-a681-4f11-a786-0a8def83af2a")))
}
