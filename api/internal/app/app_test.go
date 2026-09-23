package app

import (
	"context"
	"io"
	"net/http"
	"testing"
	"testing/synctest"
	"time"

	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/database"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/logger"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/scheduler"
)

func TestShutdownWaitsForRecalculationBeforeClosingDatabase(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		log := logger.New(logger.LevelError, logger.FormatText, false, io.Discard)
		db := &shutdownDatabase{}
		release := make(chan struct{})
		finished := false
		job := shutdownRatingFunc(func(ctx context.Context) error {
			<-ctx.Done()
			<-release
			if db.closed {
				t.Error("database closed while job was running")
			}
			finished = true
			return ctx.Err()
		})
		r, err := scheduler.NewRating(context.Background(), "* * * * *", job, log)
		if err != nil {
			t.Fatal(err)
		}
		a := &App{diContainer: &diContainer{db: db, log: log}, httpServer: &http.Server{}, ratingScheduler: r}
		r.Start()
		synctest.Wait()
		time.Sleep(time.Minute)
		synctest.Wait()
		stopped := make(chan error, 1)
		go func() { stopped <- a.gracefullShutdown() }()
		synctest.Wait()
		if db.closed {
			t.Fatal("database closed before job finished")
		}
		close(release)
		synctest.Wait()
		if err := <-stopped; err != nil {
			t.Fatal(err)
		}
		if !finished || !db.closed {
			t.Fatal("shutdown did not finish job and close database")
		}
	})
}

type shutdownDatabase struct {
	database.DBTX
	closed bool
}

func (d *shutdownDatabase) Close() { d.closed = true }

type shutdownRatingFunc func(context.Context) error

func (f shutdownRatingFunc) RecalculateAll(ctx context.Context) error { return f(ctx) }
