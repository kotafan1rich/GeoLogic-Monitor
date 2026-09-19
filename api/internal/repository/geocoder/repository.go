package geocoder

import (
	"context"
	"strings"

	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/domain"
	geocoderintegration "github.com/kotafan1rich/GeoLogic-Monitor/api/internal/integrations/geocoder"
)

type Client interface {
	Autocomplete(ctx context.Context, search string) (*geocoderintegration.Autocomplete, error)
}

type repository struct {
	client Client
}

func New(client Client) *repository {
	return &repository{client: client}
}

func (r *repository) Suggestions(ctx context.Context, query string) ([]domain.Address, error) {
	result, err := r.client.Autocomplete(ctx, query)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return []domain.Address{}, nil
	}

	return []domain.Address{{
		Address: formatAddress(result.Name, result.BuildingName),
		Lat:     result.Latitude,
		Lon:     result.Longitude,
	}}, nil
}

func formatAddress(name, buildingName string) string {
	name = strings.TrimSpace(name)
	buildingName = strings.TrimSpace(buildingName)
	if name == "" {
		return buildingName
	}
	if buildingName == "" || strings.Contains(name, buildingName) {
		return name
	}
	return name + ", " + buildingName
}
