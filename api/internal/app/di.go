package app

import (
	"context"
	"net/http"
	"os"

	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/api"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/config"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/database"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/database/postgres"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/handler"
	businesstypehandler "github.com/kotafan1rich/GeoLogic-Monitor/api/internal/handler/business_type"
	eventhandler "github.com/kotafan1rich/GeoLogic-Monitor/api/internal/handler/event"
	geocodinghandler "github.com/kotafan1rich/GeoLogic-Monitor/api/internal/handler/geocoding"
	healthhandler "github.com/kotafan1rich/GeoLogic-Monitor/api/internal/handler/health"
	infrahandler "github.com/kotafan1rich/GeoLogic-Monitor/api/internal/handler/infra"
	trackedlocationhandler "github.com/kotafan1rich/GeoLogic-Monitor/api/internal/handler/tracked_location"
	userhandler "github.com/kotafan1rich/GeoLogic-Monitor/api/internal/handler/user"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/integrations/geocoder"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/integrations/osrm"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/logger"
	businesstyperepository "github.com/kotafan1rich/GeoLogic-Monitor/api/internal/repository/database/business_type"
	eventrepository "github.com/kotafan1rich/GeoLogic-Monitor/api/internal/repository/database/event"
	infraobjectrepository "github.com/kotafan1rich/GeoLogic-Monitor/api/internal/repository/database/infra_object"
	infratyperepository "github.com/kotafan1rich/GeoLogic-Monitor/api/internal/repository/database/infra_type"
	ratingrepository "github.com/kotafan1rich/GeoLogic-Monitor/api/internal/repository/database/rating"
	trackedlocationrepository "github.com/kotafan1rich/GeoLogic-Monitor/api/internal/repository/database/tracked_location"
	userrepository "github.com/kotafan1rich/GeoLogic-Monitor/api/internal/repository/database/user"
	geocoderrepository "github.com/kotafan1rich/GeoLogic-Monitor/api/internal/repository/geocoder"
	osrmrepository "github.com/kotafan1rich/GeoLogic-Monitor/api/internal/repository/osrm"
	businesstypeservice "github.com/kotafan1rich/GeoLogic-Monitor/api/internal/service/business_type"
	calculateservice "github.com/kotafan1rich/GeoLogic-Monitor/api/internal/service/calculate"
	eventservice "github.com/kotafan1rich/GeoLogic-Monitor/api/internal/service/event"
	geocodingservice "github.com/kotafan1rich/GeoLogic-Monitor/api/internal/service/geocoding"
	infraservice "github.com/kotafan1rich/GeoLogic-Monitor/api/internal/service/infra"
	ratingservice "github.com/kotafan1rich/GeoLogic-Monitor/api/internal/service/rating"
	ratinghistoryservice "github.com/kotafan1rich/GeoLogic-Monitor/api/internal/service/rating_history"
	trackedlocationservice "github.com/kotafan1rich/GeoLogic-Monitor/api/internal/service/tracked_location"
	userservice "github.com/kotafan1rich/GeoLogic-Monitor/api/internal/service/user"
)

type userService interface {
	userhandler.UserService
	trackedlocationhandler.UserService
}

type businessTypeService interface {
	businesstypehandler.BusinessTypeService
	trackedlocationservice.BusinessTypeService
}

type infraService interface {
	infrahandler.InfraService
	trackedlocationservice.InfraService
}

type trackedLocationService interface {
	ratinghistoryservice.TrackedLocationService
	trackedlocationhandler.TrackedLocationService
}

type ratingRepository interface {
	ratingservice.RatingRepository
	ratinghistoryservice.RatingRepository
}

type diContainer struct {
	cfg                       *config.Config
	db                        database.DBTX
	txManager                 database.TxManager
	log                       *logger.Logger
	businessTypeRepository    businesstypeservice.Repository
	eventRepository           eventservice.EventRepository
	infraObjectRepository     infraservice.InfraRepository
	infraTypeRepository       infraservice.InfraTypeRepository
	ratingRepository          ratingRepository
	trackedLocationRepository trackedlocationservice.TrackedLocationRepository
	userRepository            userservice.UserRepository
	osrmRepository            trackedlocationservice.OSRMRepository
	geocoderRepository        geocodingservice.Repository
	geocodingService          geocodinghandler.GeocodingService
	ratingService             ratingservice.Service
	ratingHistoryService      ratinghistoryservice.Service
	trackedLocationService    trackedLocationService
	userService               userService
	businessTypeService       businessTypeService
	infraTypeService          infrahandler.InfraTypeService
	infraService              infraService
	healthHandler             handler.HealthHandler
	userHandler               handler.UserHandler
	geocodingHandler          handler.GeocodingHandler
	trackedLocationHandler    handler.TrackedLocationHandler
	businessTypeHandler       handler.BusinessTypeHandler
	infraHandler              handler.InfraHandler
	eventService              eventhandler.EventService
	eventHandler              handler.EventHandler
	handler                   api.Handler
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
			cfg.MinIdleConns,
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

func (d *diContainer) BusinessTypeRepository(ctx context.Context) businesstypeservice.Repository {
	if d.businessTypeRepository == nil {
		d.businessTypeRepository = businesstyperepository.NewRepository(d.DB(ctx))
	}
	return d.businessTypeRepository
}

func (d *diContainer) EventRepository(ctx context.Context) eventservice.EventRepository {
	if d.eventRepository == nil {
		d.eventRepository = eventrepository.NewRepository(d.DB(ctx))
	}
	return d.eventRepository
}

func (d *diContainer) InfraObjectRepository(ctx context.Context) infraservice.InfraRepository {
	if d.infraObjectRepository == nil {
		d.infraObjectRepository = infraobjectrepository.NewRepository(d.DB(ctx))
	}
	return d.infraObjectRepository
}

func (d *diContainer) InfraTypeRepository(ctx context.Context) infraservice.InfraTypeRepository {
	if d.infraTypeRepository == nil {
		d.infraTypeRepository = infratyperepository.NewRepository(d.DB(ctx))
	}
	return d.infraTypeRepository
}

func (d *diContainer) RatingRepository(ctx context.Context) ratingRepository {
	if d.ratingRepository == nil {
		d.ratingRepository = ratingrepository.NewRepository(d.DB(ctx))
	}
	return d.ratingRepository
}

func (d *diContainer) TrackedLocationRepository(ctx context.Context) trackedlocationservice.TrackedLocationRepository {
	if d.trackedLocationRepository == nil {
		d.trackedLocationRepository = trackedlocationrepository.NewRepository(d.DB(ctx))
	}
	return d.trackedLocationRepository
}

func (d *diContainer) UserRepository(ctx context.Context) userservice.UserRepository {
	if d.userRepository == nil {
		d.userRepository = userrepository.NewRepository(d.DB(ctx))
	}
	return d.userRepository
}

func (d *diContainer) OSRMRepository() trackedlocationservice.OSRMRepository {
	if d.osrmRepository == nil {
		client := &http.Client{Timeout: d.cfg.OSRM.Timeout}
		osrmClient := osrm.New(client, d.cfg.OSRM.BaseURL)
		d.osrmRepository = osrmrepository.New(osrmClient)
	}
	return d.osrmRepository
}

func (d *diContainer) GeocoderRepository() geocodingservice.Repository {
	if d.geocoderRepository == nil {
		client := &http.Client{Timeout: d.cfg.Geocoder.Timeout}
		geocoderClient := geocoder.New(client, d.cfg.Geocoder.BaseURL)
		d.geocoderRepository = geocoderrepository.New(geocoderClient)
	}
	return d.geocoderRepository
}

func (d *diContainer) GeocodingService() geocodinghandler.GeocodingService {
	if d.geocodingService == nil {
		d.geocodingService = geocodingservice.NewService(d.GeocoderRepository(), d.Log())
	}
	return d.geocodingService
}

func (d *diContainer) RatingService(ctx context.Context) ratingservice.Service {
	if d.ratingService == nil {
		d.ratingService = ratingservice.NewService(
			d.Log(),
			calculateservice.NewFormulaCalculator(),
			d.RatingRepository(ctx),
		)
	}
	return d.ratingService
}

func (d *diContainer) RatingHistoryService(ctx context.Context) ratinghistoryservice.Service {
	if d.ratingHistoryService == nil {
		d.ratingHistoryService = ratinghistoryservice.NewService(
			d.Log(),
			d.RatingRepository(ctx),
			d.TrackedLocationService(ctx),
		)
	}
	return d.ratingHistoryService
}

func (d *diContainer) UserService(ctx context.Context) userService {
	if d.userService == nil {
		d.userService = userservice.NewUserService(d.Log(), d.UserRepository(ctx))
	}
	return d.userService
}

func (d *diContainer) HealthHandler(_ context.Context) handler.HealthHandler {
	if d.healthHandler == nil {
		d.healthHandler = healthhandler.New()
	}
	return d.healthHandler
}

func (d *diContainer) UserHandler(ctx context.Context) handler.UserHandler {
	if d.userHandler == nil {
		d.userHandler = userhandler.New(d.UserService(ctx))
	}
	return d.userHandler
}

func (d *diContainer) GeocodingHandler() handler.GeocodingHandler {
	if d.geocodingHandler == nil {
		d.geocodingHandler = geocodinghandler.New(d.GeocodingService())
	}
	return d.geocodingHandler
}

func (d *diContainer) TrackedLocationService(ctx context.Context) trackedLocationService {
	if d.trackedLocationService == nil {
		d.trackedLocationService = trackedlocationservice.NewTrackedLocationService(
			d.TrackedLocationRepository(ctx),
			d.OSRMRepository(),
			d.InfraService(ctx),
			d.BusinessTypeService(ctx),
			d.RatingService(ctx),
			d.TxManager(ctx),
			d.Log(),
		)
	}
	return d.trackedLocationService
}

func (d *diContainer) TrackedLocationHandler(ctx context.Context) handler.TrackedLocationHandler {
	if d.trackedLocationHandler == nil {
		d.trackedLocationHandler = trackedlocationhandler.New(
			d.UserService(ctx),
			d.TrackedLocationService(ctx),
			d.RatingHistoryService(ctx),
		)
	}
	return d.trackedLocationHandler
}

func (d *diContainer) BusinessTypeService(ctx context.Context) businessTypeService {
	if d.businessTypeService == nil {
		d.businessTypeService = businesstypeservice.NewService(d.Log(), d.BusinessTypeRepository(ctx))
	}
	return d.businessTypeService
}

func (d *diContainer) BusinessTypeHandler(ctx context.Context) handler.BusinessTypeHandler {
	if d.businessTypeHandler == nil {
		d.businessTypeHandler = businesstypehandler.New(d.BusinessTypeService(ctx))
	}
	return d.businessTypeHandler
}

func (d *diContainer) InfraTypeService(ctx context.Context) infrahandler.InfraTypeService {
	if d.infraTypeService == nil {
		d.infraTypeService = infraservice.NewTypeService(d.Log(), d.InfraTypeRepository(ctx))
	}
	return d.infraTypeService
}

func (d *diContainer) InfraService(ctx context.Context) infraService {
	if d.infraService == nil {
		d.infraService = infraservice.NewInfraService(d.Log(), d.InfraObjectRepository(ctx))
	}
	return d.infraService
}

func (d *diContainer) InfraHandler(ctx context.Context) handler.InfraHandler {
	if d.infraHandler == nil {
		d.infraHandler = infrahandler.New(d.InfraTypeService(ctx), d.InfraService(ctx))
	}
	return d.infraHandler
}

func (d *diContainer) EventService(ctx context.Context) eventhandler.EventService {
	if d.eventService == nil {
		d.eventService = eventservice.NewEventService(d.Log(), d.EventRepository(ctx))
	}
	return d.eventService
}

func (d *diContainer) EventHandler(ctx context.Context) handler.EventHandler {
	if d.eventHandler == nil {
		d.eventHandler = eventhandler.New(d.EventService(ctx))
	}
	return d.eventHandler
}

func (d *diContainer) Handler(ctx context.Context) api.Handler {
	if d.handler == nil {
		d.handler = api.NewHandler(
			d.HealthHandler(ctx),
			d.UserHandler(ctx),
			d.GeocodingHandler(),
			d.TrackedLocationHandler(ctx),
			d.BusinessTypeHandler(ctx),
			d.InfraHandler(ctx),
			d.EventHandler(ctx),
			d.cfg.Security.MaxBotToken,
			d.cfg.Security.MiniAppInitDataMaxAge,
			d.cfg.Security.BotServiceToken,
			d.cfg.Security.IngestionServiceToken,
			d.cfg.Security.CORSAllowedOrigin,
		)
	}
	return d.handler
}
