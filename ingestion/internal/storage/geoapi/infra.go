package geoapi

import (
	"context"
	"fmt"
)

const infraObjectEndpoint = "internal/v1/infra"

func (c *Client) PutInfraObject(ctx context.Context, obj InfraObjectInput) (InfraObject, error) {
	const op = "geoapi.Client.PutInfraObject"

	var raw InfraObject

	err := c.do(ctx, c.baseURL, infraObjectEndpoint, obj, &raw)
	if err != nil {
		return InfraObject{}, fmt.Errorf("%s: %w", op, err)
	}

	return raw, nil
}
