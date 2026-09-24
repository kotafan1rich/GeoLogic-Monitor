package app

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/config"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/middleware"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/scheduler"
	"github.com/pressly/goose/v3"
)

type App struct {
	cfg             *config.Config
	diContainer     *diContainer
	httpServer      *http.Server
	ratingScheduler *scheduler.Rating
}

func New(ctx context.Context, cfg *config.Config) (*App, error) {
	a := &App{
		cfg:         cfg,
		diContainer: NewDIContainer(cfg),
	}
	err := a.initDeps(ctx)
	if err != nil {
		if a.diContainer.db != nil {
			a.diContainer.db.Close()
		}
		return nil, err
	}

	return a, nil
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(ctx context.Context) error{
		a.migrate,
		a.initHTTPServer,
		a.initRatingScheduler,
	}

	for _, fn := range inits {
		err := fn(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *App) migrate(_ context.Context) error {
	log := a.diContainer.Log()

	log.Info("apply migrations")
	db, err := sql.Open("pgx", a.cfg.Database.DSN())
	if err != nil {
		log.Error(
			"error to connect to database",
			"err", err,
		)
		return err
	}
	defer db.Close()

	if err := goose.SetDialect("postgres"); err != nil {
		log.Error("error to setup dialect", "err", err)
		return err
	}

	if err := goose.Up(db, "./migrations"); err != nil {
		log.Error(
			"error migrate",
			"err", err,
		)
		return err
	}

	log.Info("migrations completed")
	return nil
}

func (a *App) initHTTPServer(ctx context.Context) error {
	a.httpServer = &http.Server{
		Addr:         fmt.Sprintf(":%d", a.cfg.HttpServer.ServerPort),
		ReadTimeout:  a.cfg.HttpServer.ReadTimeout,
		WriteTimeout: a.cfg.HttpServer.WriteTimeout,
		IdleTimeout:  a.cfg.HttpServer.IdleTimeout,

		Handler: middleware.LoggerMiddleware(
			a.diContainer.Log(),
			a.diContainer.Handler(ctx).Routes(),
		),
	}
	return nil
}

func (a *App) initRatingScheduler(ctx context.Context) error {
	var err error
	a.ratingScheduler, err = a.diContainer.RatingScheduler(ctx)
	return err
}

func (a *App) gracefullShutdown() error {
	log := a.diContainer.Log()
	schedulerErr := a.ratingScheduler.Shutdown()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	httpErr := a.httpServer.Shutdown(ctx)
	if httpErr != nil {
		log.Error("HTTP server graceful shutdown failed", "err", httpErr)
		return errors.Join(schedulerErr, fmt.Errorf("server shutdown failed: %w", httpErr))
	}

	db := a.diContainer.DB(ctx)
	db.Close()
	return schedulerErr
}

func (a *App) Run() (runErr error) {
	defer func() {
		runErr = errors.Join(runErr, a.gracefullShutdown())
		if runErr == nil {
			a.diContainer.Log().Info("server stopped cleanly")
		}
	}()
	a.ratingScheduler.Start()
	log := a.diContainer.Log()
	log.Info(
		"server started",
		"addr", a.httpServer.Addr,
	)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(quit)

	errChan := make(chan error, 1)

	go func() {
		err := a.httpServer.ListenAndServe()
		if err != nil {
			log.Error("HTTP server failed to listen", "err", err)
			errChan <- err
		}
		close(errChan)
	}()

	select {
	case sig := <-quit:
		log.Info("shutdown signal received, starting graceful shutdown...", "signal", sig.String())
	case err := <-errChan:
		return err
	}
	return nil
}
