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
	userhandler "github.com/kotafan1rich/GeoLogic-Monitor/api/internal/handler/user"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/integrations/geocoder"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/integrations/osrm"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/logger"
	businesstyperepository "github.com/kotafan1rich/GeoLogic-Monitor/api/internal/repository/database/business_type"
	eventrepository "github.com/kotafan1rich/GeoLogic-Monitor/api/internal/repository/database/event"
	infraobjectrepository "github.com/kotafan1rich/GeoLogic-Monitor/api/internal/repository/database/infra_object"
	infratyperepository "github.com/kotafan1rich/GeoLogic-Monitor/api/internal/repository/database/infra_type"
	trackedlocationrepository "github.com/kotafan1rich/GeoLogic-Monitor/api/internal/repository/database/tracked_location"
	userrepository "github.com/kotafan1rich/GeoLogic-Monitor/api/internal/repository/database/user"
	geocoderrepository "github.com/kotafan1rich/GeoLogic-Monitor/api/internal/repository/geocoder"
	osrmrepository "github.com/kotafan1rich/GeoLogic-Monitor/api/internal/repository/osrm"
	businesstypeservice "github.com/kotafan1rich/GeoLogic-Monitor/api/internal/service/business_type"
	eventservice "github.com/kotafan1rich/GeoLogic-Monitor/api/internal/service/event"
	geocodingservice "github.com/kotafan1rich/GeoLogic-Monitor/api/internal/service/geocoding"
	infraservice "github.com/kotafan1rich/GeoLogic-Monitor/api/internal/service/infra"
	trackedlocationservice "github.com/kotafan1rich/GeoLogic-Monitor/api/internal/service/tracked_location"
	userservice "github.com/kotafan1rich/GeoLogic-Monitor/api/internal/service/user"
)

type diContainer struct {
	cfg                       *config.Config
	db                        database.DBTX
	txManager                 database.TxManager
	log                       *logger.Logger
	businessTypeRepository    businesstypeservice.Repository
	eventRepository           eventservice.EventRepository
	infraObjectRepository     infraservice.InfraRepository
	infraTypeRepository       infraservice.InfraTypeRepository
	trackedLocationRepository trackedlocationservice.TrackedLocationRepository
	userRepository            userservice.UserRepository
	osrmRepository            trackedlocationservice.OSRMRepository
	geocoderRepository        geocodingservice.Repository
	geocodingService          geocodingservice.Service
	userService               userhandler.UserService
	userHandler               handler.UserHandler
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

func (d *diContainer) GeocodingService() geocodingservice.Service {
	if d.geocodingService == nil {
		d.geocodingService = geocodingservice.NewService(d.GeocoderRepository(), d.Log())
	}
	return d.geocodingService
}

func (d *diContainer) UserService(ctx context.Context) userhandler.UserService {
	if d.userService == nil {
		d.userService = userservice.NewUserService(d.Log(), d.UserRepository(ctx))
	}
	return d.userService
}

func (d *diContainer) UserHandler(ctx context.Context) handler.UserHandler {
	if d.userHandler == nil {
		d.userHandler = userhandler.New(d.UserService(ctx))
	}
	return d.userHandler
}

func (d *diContainer) Handler(ctx context.Context) api.Handler {
	if d.handler == nil {
		d.handler = api.NewHandler(d.UserHandler(ctx), d.cfg.BotServiceToken)
	}
	return d.handler
}
