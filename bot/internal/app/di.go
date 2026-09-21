package app

import (
	"context"
	"os"

	"github.com/kotafan1rich/GeoLogic-Monitor/bot/internal/config"
	"github.com/kotafan1rich/GeoLogic-Monitor/bot/internal/handler/bot"
	"github.com/kotafan1rich/GeoLogic-Monitor/bot/internal/handler/bot/start"
	"github.com/kotafan1rich/GeoLogic-Monitor/bot/internal/handler/http/webhook"
	"github.com/kotafan1rich/GeoLogic-Monitor/bot/internal/integrations"
	apiintegration "github.com/kotafan1rich/GeoLogic-Monitor/bot/internal/integrations/api"
	"github.com/kotafan1rich/GeoLogic-Monitor/bot/internal/logger"
	"github.com/kotafan1rich/GeoLogic-Monitor/bot/internal/repository/api"
	"github.com/kotafan1rich/GeoLogic-Monitor/bot/internal/server"
	"github.com/kotafan1rich/GeoLogic-Monitor/bot/internal/service/user"
)

type diContainer struct {
	cfg          *config.Config
	log          *logger.Logger
	startHandler bot.StartHandler
	userService  start.UserService
	handler      server.Handler
	apiRepo      user.ApiRepository
	apiClient    api.ApiClient
	dispatcher   webhook.Dispatcher
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

func (d *diContainer) Handler(ctx context.Context) server.Handler {
	if d.handler == nil {
		d.handler = server.NewHandler()
	}
	return d.handler
}

func (d *diContainer) StartHandler() bot.StartHandler {
	if d.startHandler == nil {
		d.startHandler = start.NewStartHandler(d.UserService())
	}
	return d.startHandler
}

func (d *diContainer) UserService() start.UserService {
	if d.userService == nil {
		d.userService = user.NewService(d.ApiRepository())
	}
	return d.userService
}

func (d *diContainer) ApiRepository() user.ApiRepository {
	if d.apiRepo == nil {
		d.apiRepo = api.NewRepository(d.ApiClient())
	}
	return d.apiRepo
}

func (d *diContainer) ApiClient() api.ApiClient {
	if d.apiClient == nil {
		d.apiClient = apiintegration.New(
			integrations.NewHTTPClient(d.cfg.Api.Timeout),
			d.cfg.Api.BaseURL,
			d.cfg.Security.MaxBotToken,
		)
	}
	return d.apiClient
}

func (d *diContainer) Dispatcher() webhook.Dispatcher {
	if d.dispatcher == nil {
		d.dispatcher = bot.NewDispatcher(d.StartHandler())
	}
	return d.dispatcher
}
