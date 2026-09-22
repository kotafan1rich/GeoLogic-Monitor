package app

import (
	"context"
	"log/slog"
	"os"

	"github.com/kotafan1rich/GeoLogic-Monitor/ingestion/internal/aggregator/digitalspb"
	"github.com/kotafan1rich/GeoLogic-Monitor/ingestion/internal/closer"
	"github.com/kotafan1rich/GeoLogic-Monitor/ingestion/internal/config"
	"github.com/kotafan1rich/GeoLogic-Monitor/ingestion/internal/infra"
	"github.com/kotafan1rich/GeoLogic-Monitor/ingestion/internal/job"
	"github.com/kotafan1rich/GeoLogic-Monitor/ingestion/internal/logger"
	"github.com/kotafan1rich/GeoLogic-Monitor/ingestion/internal/scheduler"
	"github.com/kotafan1rich/GeoLogic-Monitor/ingestion/internal/storage/geoapi"
)

type diContainer struct {
	log   *slog.Logger
	dsc   *digitalspb.Client
	dsJob *job.DigitalSpb
	s     *scheduler.Scheduler
	gc    *geoapi.Client
	tr    *geoapi.TypeRegistry
}

func newDIContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) Logger() *slog.Logger {
	if d.log == nil {
		cfg := config.Get()
		d.log = logger.MustNew(cfg.Logger.Level, cfg.Logger.Format, os.Stdout)

		d.log.Info("logger initialized",
			slog.String("level", cfg.Logger.Level),
			slog.String("format", cfg.Logger.Format),
		)
	}
	return d.log
}

func (d *diContainer) DSClient() *digitalspb.Client {
	if d.dsc == nil {
		cfg := config.Get()
		log := d.Logger()

		httpClient := infra.NewDigitalSpbHTTPClient(
			cfg.Aggregator.DigitalSpb.MaxIdleConns,
			cfg.Aggregator.DigitalSpb.MaxIdleConnsPerHost,
			cfg.Aggregator.DigitalSpb.MaxConnsPerHost,
			cfg.Aggregator.DigitalSpb.RequestTimeout,
		)

		log.Info("digital-spb http client initialized")

		d.dsc = digitalspb.MustNew(
			httpClient,
			cfg.Aggregator.DigitalSpb.BaseURLMap,
			cfg.Aggregator.DigitalSpb.AttemptTimeout,
			cfg.Aggregator.DigitalSpb.MaxRetries,
		)

		log.Info("digital-spb client initialized")
	}
	return d.dsc
}

func (d *diContainer) DSJob() *job.DigitalSpb {
	if d.dsJob == nil {
		cfg := config.Get()

		d.dsJob = job.NewDigitalSpb(
			d.Logger(),
			d.DSClient(),
			cfg.Aggregator.DigitalSpb.JobInterval,
			cfg.Aggregator.DigitalSpb.StaticFiles,
			d.GeoApiClient(),
			d.TypeRegistry(),
		)

		d.Logger().Info(
			"job initialized",
			slog.String("job", d.dsJob.Name()),
			slog.Duration("interval", d.dsJob.Interval()),
			slog.Int("static_files", len(cfg.Aggregator.DigitalSpb.StaticFiles)),
		)
	}
	return d.dsJob
}

// TODO: Впоследствии дополнить 2ГИС (будет слайс / мапа с планировщиками под digitalspb и 2ГИС)
func (d *diContainer) Scheduler() *scheduler.Scheduler {
	if d.s == nil {
		s := scheduler.MustNew(d.DSJob().Name(), d.Logger())

		d.Logger().Info(
			"scheduler initialized",
			slog.String("scheduler", s.Name()),
		)

		closer.Add("scheduler:"+s.Name(), func(context.Context) error {
			return s.Shutdown()
		})

		d.s = s
	}
	return d.s
}

func (d *diContainer) GeoApiClient() *geoapi.Client {
	if d.gc == nil {
		cfg := config.Get()
		log := d.Logger()

		httpClient := infra.MustNewGeoAPIHTTPClient(
			cfg.GeoApi.AuthToken,
			cfg.GeoApi.MaxIdleConns,
			cfg.GeoApi.MaxIdleConnsPerHost,
			cfg.GeoApi.MaxConnsPerHost,
			cfg.GeoApi.RequestTimeout,
		)

		log.Info("geo-api http client initialized")

		d.gc = geoapi.MustNew(
			httpClient,
			cfg.GeoApi.URL,
			cfg.GeoApi.AttemptTimeout,
			cfg.GeoApi.MaxRetries,
		)

		log.Info("geo-api client initialized")
	}
	return d.gc
}

func (d *diContainer) TypeRegistry() *geoapi.TypeRegistry {
	if d.tr == nil {
		cfg := config.Get()

		defs := make([]geoapi.InfraTypeInput, 0, len(cfg.GeoApi.InfraTypes))
		for _, t := range cfg.GeoApi.InfraTypes {
			defs = append(defs, geoapi.InfraTypeInput{
				Slug:      t.Slug,
				Name:      t.Name,
				Weight:    t.Weight,
				MaxRadius: t.MaxRadius,
			})
		}

		d.tr = geoapi.MustNewTypeRegistry(d.GeoApiClient(), defs)

		d.Logger().Info("infra type registry initialized", slog.Int("types", len(defs)))
	}
	return d.tr
}
