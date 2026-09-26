package app

import (
	"context"
	"fmt"
	"log/slog"
	"os/signal"
	"syscall"
	"time"

	"github.com/kotafan1rich/GeoLogic-Monitor/ingestion/internal/closer"
)

const closeDuration = 10 * time.Second

type App struct {
	di *diContainer
}

func New() *App {
	return &App{di: newDIContainer()}
}

func (a *App) Run(ctx context.Context) error {
	const op = "ingestion.app.App.Run"

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
