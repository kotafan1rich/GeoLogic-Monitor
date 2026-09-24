package geoapi

import (
	"context"
	"fmt"
)

const infraObjectEndpoint = "internal/v1/infra"

func (c *Client) PutInfraObject(ctx context.Context, obj InfraObjectInput) (InfraObject, error) {
	const op = "geoapi.Client.PutInfraObject"

	var raw InfraObject

	if err := c.put(ctx, infraObjectEndpoint, obj, &raw); err != nil {
		return InfraObject{}, fmt.Errorf("%s: %w", op, err)
	}

	return raw, nil
}
