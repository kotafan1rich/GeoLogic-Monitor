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
)

type diContainer struct {
	log   *slog.Logger
	dsc   *digitalspb.Client
	dsJob *job.DigitalSpb
	s     *scheduler.Scheduler
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

		log.Info("digital-spb http client initialized",
			slog.Int("max_idle_conns", cfg.Aggregator.DigitalSpb.MaxIdleConns),
			slog.Int("max_idle_conns_per_host", cfg.Aggregator.DigitalSpb.MaxIdleConnsPerHost),
			slog.Int("max_conns_per_host", cfg.Aggregator.DigitalSpb.MaxConnsPerHost),
			slog.Duration("request_timeout", cfg.Aggregator.DigitalSpb.RequestTimeout),
		)

		d.dsc = digitalspb.MustNew(
			httpClient,
			cfg.Aggregator.DigitalSpb.BaseURLMap,
			cfg.Aggregator.DigitalSpb.AttemptTimeout,
			cfg.Aggregator.DigitalSpb.MaxRetries,
		)

		log.Info("digital-spb client initialized",
			slog.Int("base_urls", len(cfg.Aggregator.DigitalSpb.BaseURLMap)),
			slog.Duration("attempt_timeout", cfg.Aggregator.DigitalSpb.AttemptTimeout),
			slog.Uint64("max_retries", uint64(cfg.Aggregator.DigitalSpb.MaxRetries)),
		)
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
		)

		d.Logger().Info(
			"job initialized",
			slog.String("job", d.dsJob.Name()),
			slog.Duration("interval", d.dsJob.Interval()),
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
