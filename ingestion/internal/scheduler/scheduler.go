package scheduler

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/go-co-op/gocron/v2"
)

type Job interface {
	Name() string
	Schedule() string
	Run(ctx context.Context)
}

type Scheduler struct {
	name  string
	log   *slog.Logger
	inner gocron.Scheduler
}

func New(name string, log *slog.Logger) (*Scheduler, error) {
	if log == nil {
		return nil, ErrInvalidLogger
	}

	s, err := gocron.NewScheduler()
	if err != nil {
		return nil, ErrInitCronScheduler
	}

	return &Scheduler{name: name, log: log, inner: s}, nil
}

func MustNew(name string, log *slog.Logger) *Scheduler {
	const op = "scheduler.MustNew"

	s, err := New(name, log)
	if err != nil {
		panic(fmt.Sprintf("%s: failed to init scheduler: %v", op, err))
	}

	return s
}

func (s *Scheduler) Register(ctx context.Context, job Job) error {
	const op = "scheduler.Scheduler.Register"

	j, err := s.inner.NewJob(
		gocron.CronJob(job.Schedule(), false),
		gocron.NewTask(job.Run),
		gocron.WithContext(ctx),
		gocron.WithName(job.Name()),
		gocron.WithSingletonMode(gocron.LimitModeReschedule),
		gocron.WithStartAt(gocron.WithStartImmediately()),
	)
	if err != nil {
		return fmt.Errorf("%s: %w: %v", op, ErrRegisterCronJob, err)
	}

	s.log.InfoContext(
		ctx,
		"job registered",
		slog.String("scheduler", s.name),
		slog.String("job", job.Name()),
		slog.String("job_id", j.ID().String()),
		slog.String("schedule", job.Schedule()),
	)

	return nil
}

func (s *Scheduler) Start() {
	s.log.Info("scheduler started", slog.String("scheduler", s.name))
	s.inner.Start()
}

func (s *Scheduler) Shutdown() error {
	const op = "scheduler.Scheduler.Shutdown"

	s.log.Info("scheduler is stopping", slog.String("scheduler", s.name))

	if err := s.inner.Shutdown(); err != nil {
		return fmt.Errorf("%s: %w: %v", op, ErrShutdownCron, err)
	}

	s.log.Info("scheduler stopped", slog.String("scheduler", s.name))

	return nil
}

func (s *Scheduler) Name() string {
	return s.name
}
