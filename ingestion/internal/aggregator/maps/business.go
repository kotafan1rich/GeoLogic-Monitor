package maps

import (
	"context"
	"fmt"

	a "github.com/kotafan1rich/GeoLogic-Monitor/ingestion/internal/aggregator"
	"github.com/kotafan1rich/GeoLogic-Monitor/ingestion/internal/storage/geoapi"
)

const interpreterEndpoint = "interpreter"

func (c *Client) GetBusinessInfra(ctx context.Context, src a.URL) ([]geoapi.InfraObjectInput, error) {
	const op = "maps.Client.GetBusinessInfra"

	if err := src.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var raw Response

	query := c.buildQuery()

	err := c.do(ctx, src.URL(), interpreterEndpoint, query, &raw)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return MapToInfraObjects(raw.Elements), nil
}
