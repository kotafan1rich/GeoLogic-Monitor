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

func (r *repository) FilterWalkingDistance(
	ctx context.Context,
	src *domain.GeoPoint,
	dst []*domain.InfraObject,
) ([]*domain.InfraObjectDistance, error) {
	if len(dst) == 0 {
		return []*domain.InfraObjectDistance{}, nil
	}

	distances, err := r.osrmClient.GetWalkingDistances(
		ctx,
		dto.ToCoordinate(src),
		dto.ToCoordinateSlice(dst),
	)
	if err != nil {
		return nil, err
	}
	result := make([]*domain.InfraObjectDistance, 0, len(distances))
	for indx, dist := range distances {
		if dist == nil {
			continue
		}
		if *dist <= float64(dst[indx].Type.MaxRadius) {
			result = append(result, &domain.InfraObjectDistance{
				Object:         dst[indx],
				DistanceMeters: *dist,
			})
		}
	}
	return result, nil
}
