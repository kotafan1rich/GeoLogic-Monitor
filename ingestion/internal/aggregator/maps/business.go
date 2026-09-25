package maps

import (
	"context"
)

const interpreterEndpoint = "interpreter"

func (c *Client) GetBusinessInfra(ctx context.Context) ([]Element, error) {
	var raw Response

	query := c.buildQuery()
	err := c.do(ctx, c.BaseURL, interpreterEndpoint, query, &raw)
	if err != nil {
		return nil, err
	}

	return raw.Elements, nil
}
