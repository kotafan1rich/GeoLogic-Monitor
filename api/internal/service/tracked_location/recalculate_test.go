package trackedlocation

import (
	"context"
	"errors"
	"testing"
	"uuid"

	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/domain"
	domainerrs "github.com/kotafan1rich/GeoLogic-Monitor/api/internal/errs"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/service/calculate"
	ratingservice "github.com/kotafan1rich/GeoLogic-Monitor/api/internal/service/rating"
)

func TestCreateAndRecalculateAppendRatings(t *testing.T) {
	repo := &monitoringRepository{}
	s := newTestService(repo)
	history := &ratingHistoryRepository{}
	s.ratingService = ratingservice.NewService(testLogger(), calculate.NewFormulaCalculator(), history)

	point, err := s.Create(context.Background(), uuid.New(), uuid.New(), "name", "address", 59.93, 30.32)
	if err != nil {
		t.Fatal(err)
	}
	if len(history.rows) != 1 || history.rows[0].Value != point.Value || point.Value != 5 {
		t.Fatalf("first rating was not saved: %+v", history.rows)
	}
	repo.locations = []domain.MonitoringLocation{{TrackedLocation: domain.TrackedLocation{
		ID: point.ID, BusinessTypeID: point.BusinessTypeID, GeoPoint: point.GeoPoint,
	}}}
	for range 2 {
		if err := s.Recalculate(context.Background(), point.ID); err != nil {
			t.Fatal(err)
		}
	}
	if len(history.rows) != 3 {
		t.Fatalf("history length = %d, want 3 even for unchanged ratings", len(history.rows))
	}
	for i, row := range history.rows {
		if row.TrackedLocationID != point.ID || row.Value != point.Value || row.CalculatedAt.IsZero() {
			t.Fatalf("unexpected history row %d: %+v", i, row)
		}
	}
	if repo.createCalls != 1 {
		t.Fatal("recalculation must not create a tracked location")
	}
}

func TestRecalculateAllContinuesAfterFailure(t *testing.T) {
	for _, failure := range []error{errors.New("calculation unavailable"), domainerrs.ErrTrackedLocationNotFound} {
		t.Run(failure.Error(), func(t *testing.T) {
			repo := monitoringPoints()
			s := newTestService(repo)
			history := &ratingHistoryRepository{failID: repo.locations[0].ID, failure: failure}
			s.ratingService = ratingservice.NewService(testLogger(), calculate.NewFormulaCalculator(), history)
			if !errors.Is(failure, domainerrs.ErrTrackedLocationNotFound) {
				s.infraService = &failingOnceInfra{failure: failure}
			}
			if err := s.RecalculateAll(context.Background()); err != nil {
				t.Fatal(err)
			}
			if repo.reads != 2 || len(history.rows) != 1 || history.rows[0].TrackedLocationID != repo.locations[1].ID {
				t.Fatalf("failed/deleted point prevented later save: reads=%d, rows=%+v", repo.reads, history.rows)
			}
		})
	}
}

func TestRecalculateAllStopsOnCancellation(t *testing.T) {
	repo := monitoringPoints()
	s := newTestService(repo)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	history := &ratingHistoryRepository{}
	s.ratingService = ratingservice.NewService(testLogger(), calculate.NewFormulaCalculator(), history)
	s.infraService = cancellingInfra{cancel: cancel}
	if err := s.RecalculateAll(ctx); !errors.Is(err, context.Canceled) {
		t.Fatalf("got %v, want context.Canceled", err)
	}
	if repo.reads != 1 || len(history.rows) != 0 {
		t.Fatalf("work continued after cancellation: reads=%d, rows=%d", repo.reads, len(history.rows))
	}
}

func TestCreateKeepsFirstRatingInTransaction(t *testing.T) {
	for _, failure := range []error{nil, errors.New("rating insert failed")} {
		t.Run("save", func(t *testing.T) {
			tx := &checkedTransaction{}
			repo := &transactionLocationRepository{tx: tx, t: t}
			s := newTestService(repo)
			s.txManager = tx
			history := &ratingHistoryRepository{failure: failure, beforeSave: func(ctx context.Context) {
				if !tx.active || ctx != tx.ctx {
					t.Fatal("rating saved outside location transaction")
				}
			}}
			s.ratingService = ratingservice.NewService(testLogger(), calculate.NewFormulaCalculator(), history)
			_, err := s.Create(context.Background(), uuid.New(), uuid.New(), "name", "address", 59.93, 30.32)
			if !errors.Is(err, failure) || tx.committed != (failure == nil) {
				t.Fatalf("transaction result: err=%v, committed=%v", err, tx.committed)
			}
		})
	}
}

type monitoringRepository struct {
	fakeTrackedLocationRepository
	locations []domain.MonitoringLocation
	reads     int
}

func monitoringPoints() *monitoringRepository {
	return &monitoringRepository{locations: []domain.MonitoringLocation{
		{TrackedLocation: domain.TrackedLocation{ID: uuid.New(), BusinessTypeID: uuid.New()}},
		{TrackedLocation: domain.TrackedLocation{ID: uuid.New(), BusinessTypeID: uuid.New()}},
	}}
}

func (r *monitoringRepository) GetAllForMonitoring(context.Context) ([]domain.MonitoringLocation, error) {
	return r.locations, nil
}

func (r *monitoringRepository) GetByID(_ context.Context, id uuid.UUID) (*domain.TrackedLocation, error) {
	r.reads++
	for _, location := range r.locations {
		if location.ID == id {
			return &location.TrackedLocation, nil
		}
	}
	return nil, domainerrs.ErrTrackedLocationNotFound
}

type ratingHistoryRepository struct {
	rows       []domain.LocationRating
	failID     uuid.UUID
	failure    error
	beforeSave func(context.Context)
}

func (r *ratingHistoryRepository) Create(ctx context.Context, row *domain.LocationRating) (*domain.LocationRating, error) {
	if r.beforeSave != nil {
		r.beforeSave(ctx)
	}
	if row.TrackedLocationID == r.failID && r.failure != nil {
		return nil, r.failure
	}
	r.rows = append(r.rows, *row)
	return row, nil
}

type cancellingInfra struct{ cancel context.CancelFunc }

type failingOnceInfra struct{ failure error }

func (f *failingOnceInfra) Near(context.Context, *domain.GeoPoint) ([]*domain.InfraObject, error) {
	err := f.failure
	f.failure = nil
	return nil, err
}

func (c cancellingInfra) Near(context.Context, *domain.GeoPoint) ([]*domain.InfraObject, error) {
	c.cancel()
	return nil, context.Canceled
}

type checkedTransaction struct {
	active, committed bool
	ctx               context.Context
}

func (tx *checkedTransaction) WithTx(ctx context.Context, fn func(context.Context) error) error {
	tx.ctx = context.WithValue(ctx, tx, true)
	tx.active = true
	err := fn(tx.ctx)
	tx.active = false
	tx.committed = err == nil
	return err
}

type transactionLocationRepository struct {
	fakeTrackedLocationRepository
	tx *checkedTransaction
	t  *testing.T
}

func (r *transactionLocationRepository) Create(ctx context.Context, location *domain.TrackedLocation) (*domain.TrackedLocation, error) {
	if !r.tx.active || ctx != r.tx.ctx {
		r.t.Fatal("location saved outside transaction")
	}
	return r.fakeTrackedLocationRepository.Create(ctx, location)
}
