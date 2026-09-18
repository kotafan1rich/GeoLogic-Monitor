package app

import (
	"log/slog"
	"os"

	"github.com/kotafan1rich/GeoLogic-Monitor/ingestion/internal/config"
	"github.com/kotafan1rich/GeoLogic-Monitor/ingestion/internal/logger"
)

type diContainer struct {
	log *slog.Logger
}

func newDIContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) Logger() *slog.Logger {
	if d.log == nil {
		cfg := config.Get()
		d.log = logger.MustNew(cfg.Logger.Level, cfg.Logger.Format, os.Stdout)
	}
	return d.log
}
