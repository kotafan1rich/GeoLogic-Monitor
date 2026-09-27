package infraobject

import (
	"context"
	"testing"

	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/database"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/domain"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/repository/database/infra_object/query"
	basemodel "github.com/kotafan1rich/GeoLogic-Monitor/api/internal/repository/database/model"
)

func TestNearUsesSingleQueryWithConfiguredRadius(t *testing.T) {
	db := &queryRecorder{}
	point := &domain.GeoPoint{Lat: 59.93, Lon: 30.32}
	result, err := NewRepository(db).Near(context.Background(), point, 800)
	if err != nil {
		t.Fatal(err)
	}
	if db.calls != 1 || db.sql != query.Near || len(db.args) != 2 || db.args[0] != basemodel.GeoPoint(*point) || db.args[1] != float64(800) || len(result) != 0 {
		t.Fatalf("unexpected spatial query: %+v", db)
	}
	if !db.rows.closed {
		t.Fatal("rows not closed")
	}
}

type queryRecorder struct {
	database.DBTX
	calls int
	sql   string
	args  []any
	rows  emptyRows
}

func (d *queryRecorder) Query(_ context.Context, sql string, args ...any) (database.Rows, error) {
	d.calls++
	d.sql = sql
	d.args = args
	return &d.rows, nil
}

type emptyRows struct{ closed bool }

func (*emptyRows) Next() bool        { return false }
func (*emptyRows) Scan(...any) error { return nil }
func (r *emptyRows) Close()          { r.closed = true }
func (*emptyRows) Err() error        { return nil }
