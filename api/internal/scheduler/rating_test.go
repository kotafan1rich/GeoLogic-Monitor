package scheduler

import (
	"context"
	"io"
	"testing"
	"testing/synctest"
	"time"

	"github.com/caarlos0/env/v9"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/config"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/logger"
)

func TestMonthlyScheduleDoesNotRunAtStartup(t *testing.T) {
	t.Setenv("RATING_RECALC_CRON", "")
	synctest.Test(t, func(t *testing.T) {
		var cfg config.Rating
		if err := env.Parse(&cfg); err != nil {
			t.Fatal(err)
		}
		if cfg.RecalcCron != "0 3 1 * *" {
			t.Fatalf("default schedule = %q", cfg.RecalcCron)
		}
		calls := 0
		r := newTestRating(t, cfg.RecalcCron, func(context.Context) error { calls++; return nil })
		defer r.Shutdown()
		r.Start()
		synctest.Wait()
		if calls != 0 {
			t.Fatal("recalculation ran at startup")
		}
		next, err := r.scheduler.Jobs()[0].NextRun()
		if err != nil {
			t.Fatal(err)
		}
		moscow, _ := time.LoadLocation("Europe/Moscow")
		local := next.In(moscow)
		if local.Day() != 1 || local.Hour() != 3 || local.Minute() != 0 || !next.After(time.Now()) {
			t.Fatalf("unexpected next monthly run: %v", local)
		}
		time.Sleep(time.Until(next) - time.Nanosecond)
		synctest.Wait()
		if calls != 0 {
			t.Fatal("recalculation ran before cron time")
		}
		time.Sleep(time.Nanosecond)
		synctest.Wait()
		if calls != 1 {
			t.Fatalf("scheduled runs = %d, want 1", calls)
		}
	})
}

func TestOverlappingRunsAreSkipped(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		calls := 0
		release := make(chan struct{})
		r := newTestRating(t, "* * * * *", func(ctx context.Context) error {
			calls++
			select {
			case <-release:
			case <-ctx.Done():
			}
			return nil
		})
		defer r.Shutdown()
		r.Start()
		synctest.Wait()
		time.Sleep(time.Minute)
		synctest.Wait()
		if calls != 1 {
			t.Fatalf("runs = %d, want 1", calls)
		}
		time.Sleep(2 * time.Minute)
		synctest.Wait()
		if calls != 1 {
			t.Fatalf("overlapping runs were not skipped: %d", calls)
		}
		close(release)
		synctest.Wait()
		if calls != 1 {
			t.Fatal("skipped runs were queued")
		}
		time.Sleep(time.Minute)
		synctest.Wait()
		if calls != 2 {
			t.Fatalf("next scheduled run missing: %d", calls)
		}
	})
}

func TestShutdownCancelsAndWaitsForJob(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		calls := 0
		cancelled := make(chan struct{})
		release := make(chan struct{})
		r := newTestRating(t, "* * * * *", func(ctx context.Context) error {
			calls++
			<-ctx.Done()
			close(cancelled)
			<-release
			return ctx.Err()
		})
		r.Start()
		synctest.Wait()
		time.Sleep(time.Minute)
		synctest.Wait()
		stopped := make(chan error, 1)
		go func() { stopped <- r.Shutdown() }()
		synctest.Wait()
		select {
		case <-cancelled:
		default:
			t.Fatal("job context was not cancelled")
		}
		select {
		case <-stopped:
			t.Fatal("shutdown returned while job still used dependencies")
		default:
		}
		close(release)
		synctest.Wait()
		if err := <-stopped; err != nil {
			t.Fatal(err)
		}
		time.Sleep(2 * time.Minute)
		synctest.Wait()
		if calls != 1 {
			t.Fatalf("job ran after shutdown: %d", calls)
		}
	})
}

func TestInvalidSchedule(t *testing.T) {
	if _, err := NewRating(context.Background(), "invalid", ratingFunc(func(context.Context) error { return nil }), schedulerLogger()); err == nil {
		t.Fatal("invalid cron accepted")
	}
}

type ratingFunc func(context.Context) error

func (f ratingFunc) RecalculateAll(ctx context.Context) error { return f(ctx) }

func newTestRating(t *testing.T, schedule string, fn ratingFunc) *Rating {
	t.Helper()
	r, err := NewRating(context.Background(), schedule, fn, schedulerLogger())
	if err != nil {
		t.Fatal(err)
	}
	return r
}

func schedulerLogger() *logger.Logger {
	return logger.New(logger.LevelError, logger.FormatText, false, io.Discard)
}
