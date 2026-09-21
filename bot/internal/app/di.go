package app

import (
	"context"
	"os"

	"github.com/kotafan1rich/GeoLogic-Monitor/bot/internal/bot"
	"github.com/kotafan1rich/GeoLogic-Monitor/bot/internal/config"
	"github.com/kotafan1rich/GeoLogic-Monitor/bot/internal/logger"
)

type diContainer struct {
	cfg     *config.Config
	log     *logger.Logger
	handler bot.Handler
}

func NewDIContainer(cfg *config.Config) *diContainer {
	if cfg == nil {
		panic("diContainer: config cannot be nil")
	}
	return &diContainer{
		cfg: cfg,
	}
}

func (d *diContainer) Log() *logger.Logger {
	if d.log == nil {
		d.log = logger.New(
			d.cfg.Logging.LogLevel,
			d.cfg.Logging.LogFormat,
			d.cfg.Logging.LogAddSource,
			os.Stdout,
		)
	}
	return d.log
}

func (d *diContainer) Handler(ctx context.Context) bot.Handler {
	if d.handler == nil {
		d.handler = bot.NewHandler(
		)
	}
	return d.handler
}
