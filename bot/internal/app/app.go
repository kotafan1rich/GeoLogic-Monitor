package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kotafan1rich/GeoLogic-Monitor/bot/internal/config"
	"github.com/kotafan1rich/GeoLogic-Monitor/bot/internal/middleware"
)

type App struct {
	cfg         *config.Config
	diContainer *diContainer
	httpServer  *http.Server
}

func New(ctx context.Context, cfg *config.Config) (*App, error) {
	a := &App{
		cfg:         cfg,
		diContainer: NewDIContainer(cfg),
	}
	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(ctx context.Context) error{
		a.initHTTPServer,
	}

	for _, fn := range inits {
		err := fn(ctx)
		if err != nil {
			return err
		}
	}
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

func (a *App) gracefullShutdown() error {
	log := a.diContainer.Log()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := a.httpServer.Shutdown(ctx); err != nil {
		log.Error("HTTP server graceful shutdown failed", "err", err)
		return fmt.Errorf("server shutdown failed: %w", err)
	}

	return nil
}

func (a *App) Run(ctx context.Context) error {
	appCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	log := a.diContainer.Log()
	log.Info(
		"server started",
		"addr", a.httpServer.Addr,
	)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	errChan := make(chan error, 2)

	go func() {
		err := a.httpServer.ListenAndServe()
		if err != nil {
			log.Error("HTTP server failed to listen", "err", err)
			errChan <- err
		}
	}()

	go func() {
		log.Info("registering MAX webhook")

		result, err := a.diContainer.MAXClient().Subscriptions.Subscribe(
			appCtx,
			a.cfg.Security.WebhookUrl,
			a.cfg.Security.WebhookSecret,
			[]string{"message_created", "bot_started"},
			"",
		)
		if err != nil {
			log.Error("failed to register MAX webhook", "err", err)
			errChan <- fmt.Errorf("subscribe MAX webhook: %w", err)
			return
		}
		if !result.Success {
			log.Error("MAX rejected webhook", "message", result.Message)
			errChan <- fmt.Errorf("MAX rejected webhook: %s", result.Message)
			return
		}

		log.Info("MAX webhook registered")
	}()

	consumer, err := a.diContainer.NewBrokerConsumer()
	if err != nil {
		log.Error("failed to create Kafka consumer", "err", err)
	} else {
		go func() {
			err := consumer.Consume(
				appCtx,
				a.diContainer.NotificationHandler().Handle,
			)
			if err != nil {
				log.Error("Kafka consumer stopped", "err", err)
			}
		}()
	}

	select {
	case sig := <-quit:
		log.Info(
			"shutdown signal received",
			"signal",
			sig.String(),
		)

		cancel()

		if consumer != nil {
			if err := consumer.Close(); err != nil {
				log.Error("failed to close Kafka consumer", "err", err)
			}
		}

		if err := a.gracefullShutdown(); err != nil {
			return err
		}
	case err := <-errChan:
		return err
	}
	return nil
}
