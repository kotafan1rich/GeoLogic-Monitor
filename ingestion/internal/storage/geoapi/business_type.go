package geoapi

import (
	"context"
	"fmt"
)

const businessTypesEndpoint = "internal/v1/business-types"

func (c *Client) ConnectWithInfra(ctx context.Context, ID string) (BusinessType, error) {
	const op = "geoapi.Client.ConnectWithInfra"

	var raw BusinessType

	if err := c.put(
		ctx, businessTypesEndpoint, InfraTypeID{InfraTypeID: ID}, &raw,
	); err != nil {
		return BusinessType{}, fmt.Errorf("%s: %w", op, err)
	}

	return raw, nil
}
