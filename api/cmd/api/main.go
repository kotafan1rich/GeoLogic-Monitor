package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/app"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/config"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	app, err := app.New(ctx, cfg)
	if err != nil {
		slog.Error("creating app error", "err", err)
		os.Exit(1)
	}

	if err := app.Run(); err != nil {
		slog.Error("app error", "err", err)
		os.Exit(1)
	}
}
