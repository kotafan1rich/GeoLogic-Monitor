package maps

import (
	"context"
	"fmt"

	"github.com/kotafan1rich/GeoLogic-Monitor/ingestion/internal/storage/geoapi"
)

const interpreterEndpoint = "interpreter"

func (c *Client) GetBusinessInfra(ctx context.Context) ([]geoapi.InfraObjectInput, error) {
	const op = "maps.Client.GetBusinessInfra"

	var raw Response

	query := c.buildQuery()
	err := c.do(ctx, c.BaseURL, interpreterEndpoint, query, &raw)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return MapToInfraObjects(raw.Elements), nil
}
