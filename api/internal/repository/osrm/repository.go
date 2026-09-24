package osrm

import (
	"context"

	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/domain"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/integrations/osrm"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/repository/osrm/dto"
)

type OSRMClient interface {
	GetWalkingDistances(ctx context.Context, src *osrm.Coordinate, dst []*osrm.Coordinate) ([]*float64, error)
}

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

	distances, err := r.osrmClient.GetWalkingDistances(
		ctx,
		dto.ToCoordinate(src),
		dto.ToCoordinateSlice(dst),
	)
	if err != nil {
		return nil, err
	}
	return distances, nil
}
