package geoapi

import (
	"context"
	"fmt"
	"net/url"
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

func (c *Client) BusinessTypes(ctx context.Context) ([]BusinessType, error) {
	const op = "geoapi.Client.BusinessTypes"

	var raw []BusinessType

	if err := c.get(
		ctx, businessTypesEndpoint, url.Values{}, &raw,
	); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return raw, nil
}
