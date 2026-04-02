package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/chishkin-afk/todo/docs"
	"github.com/chishkin-afk/todo/internal/app"
)

// @title Todo list API
// @version 1.0
// @description API for creating task & groups to do something useful
// @termsOfService http://swagger.io/terms/

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:9000
// @BasePath /api/v1
// @schemes http https

// @securityDefinitions.apikey jwt
// @in header
// @name Authorization
// @description Type JWT token.
func main() {
	docs.SwaggerInfo.Title = "Todo list API"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Description = "API for creating task & groups to do something useful"
	docs.SwaggerInfo.Host = "localhost:9000"
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	app, cleanup, err := app.New()
	if err != nil {
		slog.Error("failed to setup app",
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}
	defer cleanup()

	go func() {
		if err := app.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("failed to start server",
				slog.String("error", err.Error()),
			)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := app.Shutdown(ctx); err != nil {
		slog.Error("failed to shutdown server",
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}
}
