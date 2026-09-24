package scheduler

import (
	"context"
	"fmt"
	"sync"
	"time"
	_ "time/tzdata"

	"github.com/go-co-op/gocron/v2"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/logger"
)

type RatingService interface {
	RecalculateAll(context.Context) error
}

type Rating struct {
	scheduler gocron.Scheduler
	service   RatingService
	log       *logger.Logger
	ctx       context.Context
	cancel    context.CancelFunc
	mu        sync.Mutex
	stopping  bool
	active    sync.WaitGroup
}

func NewRating(ctx context.Context, schedule string, service RatingService, log *logger.Logger) (*Rating, error) {
	location, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		return nil, err
	}
	s, err := gocron.NewScheduler(gocron.WithLocation(location))
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithCancel(ctx)
	r := &Rating{scheduler: s, service: service, log: log, ctx: ctx, cancel: cancel}
	_, err = s.NewJob(
		gocron.CronJob(schedule, false),
		gocron.NewTask(r.run),
		gocron.WithName("rating-recalculation"),
		gocron.WithSingletonMode(gocron.LimitModeReschedule),
	)
	if err != nil {
		cancel()
		_ = s.Shutdown()
		return nil, fmt.Errorf("configure rating recalculation: %w", err)
	}
	return r, nil
}

func (r *Rating) Start() {
	r.scheduler.Start()
}

func (r *Rating) run() {
	r.mu.Lock()
	if r.stopping || r.ctx.Err() != nil {
		r.mu.Unlock()
		return
	}
	r.active.Add(1)
	r.mu.Unlock()
	defer r.active.Done()

	if err := r.service.RecalculateAll(r.ctx); err != nil && r.ctx.Err() == nil {
		r.log.ErrorContext(r.ctx, "rating recalculation failed", "err", err)
	}
}

func (r *Rating) Shutdown() error {
	r.mu.Lock()
	r.stopping = true
	r.cancel()
	r.mu.Unlock()
	err := r.scheduler.Shutdown()
	// Even if gocron's stop timeout expires, dependencies must remain open
	// until the running service call has finished using them.
	r.active.Wait()
	return err
}
