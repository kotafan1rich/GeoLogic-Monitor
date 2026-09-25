package maps

import (
	"context"
	"net/url"
)

const interpreterEndpoint = "interpreter"

func (c *Client) GetBusinessInfra(ctx context.Context) ([]Element, error) {
	var raw Response

	q := c.buildQuery()
	err := c.do(ctx, c.BaseURL, interpreterEndpoint, url.Values{"data": {q}}, &raw)
	if err != nil {
		return nil, err
	}

	return raw.Elements, nil
}
