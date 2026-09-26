package app

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/tern/v2/migrate"
	"github.com/kotafan1rich/GeoLogic-Monitor/ingestion/internal/closer"
)

const (
	closeDuration = 10 * time.Second

	ternTable     = "tern_version"
	migrationsDir = "./migrations"
)

type App struct {
	di *diContainer
}

func New() *App {
	return &App{di: newDIContainer()}
}

func (a *App) Run(ctx context.Context) error {
	const op = "app.App.Run"

	log := a.di.Logger()

	stopCtx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	log.InfoContext(stopCtx, "application is starting")

	for _, j := range a.di.Jobs() {
		if err := a.di.Scheduler().Register(ctx, j); err != nil {
			log.ErrorContext(
				ctx,
				"failed to register job",
				slog.String("job", j.Name()),
				slog.Any("error", err),
			)
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	go func() { a.di.Scheduler().Start() }()

	log.InfoContext(stopCtx, "application started", slog.String("scheduler", a.di.Scheduler().Name()))

	<-stopCtx.Done()

	stop()

	log.InfoContext(ctx, "shutdown signal received, closing resources")

	closerCtx, closerCancel := context.WithTimeout(ctx, closeDuration)
	defer closerCancel()

	if err := closer.CloseAll(closerCtx); err != nil {
		log.ErrorContext(ctx, "failed to close resources", slog.Any("error", err))
		return fmt.Errorf("%s: %w", op, err)
	}

	log.InfoContext(ctx, "application stopped")

	return nil
}

func (a *App) Migrate(ctx context.Context) error {
	const op = "app.App.Migrate"

	log := a.di.Logger()

	pool := a.di.DB(ctx).Pool()

	conn, err := pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer conn.Release()

	migrator, err := migrate.NewMigrator(ctx, conn.Conn(), ternTable)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.InfoContext(ctx, "migrator initialized")

	err = migrator.LoadMigrations(os.DirFS(migrationsDir))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.InfoContext(ctx, "migrations loaded")

	err = migrator.Migrate(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.InfoContext(ctx, "migration finished")

	return nil
}
