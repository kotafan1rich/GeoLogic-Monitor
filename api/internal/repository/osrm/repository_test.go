package osrm

import (
	"context"
	"errors"
	"testing"

	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/domain"
	integration "github.com/kotafan1rich/GeoLogic-Monitor/api/internal/integrations/osrm"
)

func TestWalkingDistancesBatchesAndPreservesOrder(t *testing.T) {
	client := &batchClient{}
	dst := make([]*domain.GeoPoint, 205)
	for i := range dst {
		dst[i] = &domain.GeoPoint{Lon: float64(i)}
	}
	distances, err := New(client).WalkingDistances(context.Background(), &domain.GeoPoint{}, dst)
	if err != nil {
		t.Fatal(err)
	}
	if client.calls != 3 || len(distances) != 205 {
		t.Fatal("wrong batch count")
	}
	for i, d := range distances {
		if i == 100 {
			if d != nil {
				t.Fatal("null route lost")
			}
			continue
		}
		if d == nil || *d != float64(i) {
			t.Fatalf("wrong distance at %d", i)
		}
	}
	client = &batchClient{failAt: 2}
	if result, err := New(client).WalkingDistances(context.Background(), &domain.GeoPoint{}, dst); err == nil || result != nil {
		t.Fatal("partial result leaked after failure")
	}
	ctx, cancel := context.WithCancel(context.Background())
	client = &batchClient{after: cancel}
	if _, err := New(client).WalkingDistances(ctx, &domain.GeoPoint{}, dst); !errors.Is(err, context.Canceled) || client.calls != 1 {
		t.Fatal("batching ignored cancellation")
	}
}

type batchClient struct {
	calls  int
	failAt int
	after  func()
}

func (c *batchClient) GetWalkingDistances(_ context.Context, _ *integration.Coordinate, dst []*integration.Coordinate) ([]*float64, error) {
	c.calls++
	if len(dst) > 99 {
		return nil, errors.New("too many coordinates")
	}
	if c.calls == c.failAt {
		return nil, errors.New("provider failed")
	}
	result := make([]*float64, len(dst))
	for i, p := range dst {
		if p.Lon == 100 {
			continue
		}
		d := p.Lon
		result[i] = &d
	}
	if c.after != nil {
		c.after()
	}
	return result, nil
}
