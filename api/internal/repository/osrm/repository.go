package osrm

import (
	"context"
	"fmt"

	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/domain"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/integrations/osrm"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/repository/osrm/dto"
)

type OSRMClient interface {
	GetWalkingDistances(ctx context.Context, src *osrm.Coordinate, dst []*osrm.Coordinate) ([]*float64, error)
}

// One source plus at most 99 targets fits OSRM's 100-coordinate table limit.
const maxWalkingDestinations = 99

type repository struct {
	osrmClient OSRMClient
}

func New(client OSRMClient) *repository {
	return &repository{
		osrmClient: client,
	}
}

func (r *repository) WalkingDistances(
	ctx context.Context,
	src *domain.GeoPoint,
	dst []*domain.GeoPoint,
) ([]*float64, error) {
	if len(dst) == 0 {
		return []*float64{}, nil
	}

	distances := make([]*float64, 0, len(dst))
	for start := 0; start < len(dst); start += maxWalkingDestinations {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		end := min(start+maxWalkingDestinations, len(dst))
		batch, err := r.osrmClient.GetWalkingDistances(ctx, dto.ToCoordinate(src), dto.ToCoordinateSlice(dst[start:end]))
		if err != nil {
			return nil, err
		}
		if len(batch) != end-start {
			return nil, fmt.Errorf("unexpected OSRM batch distance count")
		}
		distances = append(distances, batch...)
	}
	return distances, nil
}
