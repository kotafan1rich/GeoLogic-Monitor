package app

import (
	"context"
	"fmt"
	"os"

	maxbot "github.com/max-messenger/max-bot-api-client-go/v2"

	"github.com/kotafan1rich/GeoLogic-Monitor/bot/internal/config"
	"github.com/kotafan1rich/GeoLogic-Monitor/bot/internal/handler/bot"
	"github.com/kotafan1rich/GeoLogic-Monitor/bot/internal/handler/bot/start"
	"github.com/kotafan1rich/GeoLogic-Monitor/bot/internal/handler/http"
	"github.com/kotafan1rich/GeoLogic-Monitor/bot/internal/handler/http/webhook"
	"github.com/kotafan1rich/GeoLogic-Monitor/bot/internal/integrations"
	apiintegration "github.com/kotafan1rich/GeoLogic-Monitor/bot/internal/integrations/api"
	"github.com/kotafan1rich/GeoLogic-Monitor/bot/internal/logger"
	"github.com/kotafan1rich/GeoLogic-Monitor/bot/internal/repository"
	"github.com/kotafan1rich/GeoLogic-Monitor/bot/internal/repository/api"
	"github.com/kotafan1rich/GeoLogic-Monitor/bot/internal/repository/max"
	"github.com/kotafan1rich/GeoLogic-Monitor/bot/internal/server"
	"github.com/kotafan1rich/GeoLogic-Monitor/bot/internal/service"
	"github.com/kotafan1rich/GeoLogic-Monitor/bot/internal/service/onboarding"
)

type diContainer struct {
	cfg            *config.Config
	log            *logger.Logger
	startHandler   bot.StartHandler
	userService    service.OnboardingService
	webhookHandler http.WebhookHandler
	handler        server.Handler
	apiRepo        repository.ApiRepository
	apiClient      api.ApiClient
	messenger      repository.Messenger
	maxClient      *maxbot.Api
	dispatcher     webhook.Dispatcher
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
		d.handler = server.NewHandler(
			d.MAXClient(),
			d.WebhookHandler(),
			d.cfg.Security.WebhookSecret,
		)
	}
	return d.handler
}

func (d *diContainer) StartHandler() bot.StartHandler {
	if d.startHandler == nil {
		d.startHandler = start.NewStartHandler(
			d.Log(),
			d.UserService(),
		)
	}
	return d.startHandler
}

func (d *diContainer) UserService() service.OnboardingService {
	if d.userService == nil {
		d.userService = onboarding.NewService(
			d.ApiRepository(),
			d.Messenger(),
		)
	}
	return d.userService
}

func (d *diContainer) WebhookHandler() http.WebhookHandler {
	if d.webhookHandler == nil {
		d.webhookHandler = webhook.NewHandler(
			d.Dispatcher(),
		)
	}
	return d.webhookHandler
}

func (d *diContainer) ApiRepository() repository.ApiRepository {
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
			d.cfg.Security.BotServiceToken,
		)
	}
	return d.apiClient
}

func (d *diContainer) Messenger() repository.Messenger {
	if d.messenger == nil {
		d.messenger = max.NewRepository(d.MAXClient())
	}
	return d.messenger
}

func (d *diContainer) MAXClient() *maxbot.Api {
	if d.maxClient == nil {
		client, err := maxbot.NewApi(d.cfg.Security.MaxBotToken)
		if err != nil {
			panic(fmt.Errorf("create MAX client: %w", err))
		}

		d.maxClient = client
	}

	return d.maxClient
}

func (d *diContainer) Dispatcher() webhook.Dispatcher {
	if d.dispatcher == nil {
		d.dispatcher = bot.NewDispatcher(d.StartHandler())
	}
	return d.dispatcher
}
