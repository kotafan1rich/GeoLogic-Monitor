package app

import (
	"context"
	"os"

	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/config"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/database"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/database/postgres"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/logger"
)

type diContainer struct {
	cfg       *config.Config
	db        database.DBTX
	txManager database.TxManager
	log       *logger.Logger
}

func NewDIContainer(cfg *config.Config) *diContainer {
	if cfg == nil {
		panic("diContainer: config cannot be nil")
	}
	return &diContainer{
		cfg: cfg,
	}
}

func (d *diContainer) DB(ctx context.Context) database.DBTX {
	if d.db == nil {
		cfg := d.cfg.Database
		pool, err := postgres.NewPool(
			ctx,
			cfg.DSN(),
			cfg.MaxOpenConns,
			cfg.MaxConnLifetime,
		)
		if err != nil {
			d.Log().Error("error to init database pool", "err", err)
			os.Exit(1)
		}
		d.db = pool
	}
	return d.db
}

func (d *diContainer) TxManager(ctx context.Context) database.TxManager {
	if d.txManager == nil {
		d.txManager = postgres.NewManager(d.DB(ctx), d.Log())
	}
	return d.txManager
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
