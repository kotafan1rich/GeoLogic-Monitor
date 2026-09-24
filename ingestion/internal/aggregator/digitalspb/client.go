package digitalspb

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"

	a "github.com/kotafan1rich/GeoLogic-Monitor/ingestion/internal/aggregator"
	"github.com/kotafan1rich/GeoLogic-Monitor/ingestion/internal/infra"
)

const maxBodySize = infra.MaxBodySize

type Client struct {
	BaseURLs map[string]a.URL
	log      *slog.Logger
	req      *infra.Requester
}

func New(log *slog.Logger, req *infra.Requester, baseURLs map[string]string) (*Client, error) {
	if log == nil {
		return nil, ErrInvalidLogger
	}

	if req == nil {
		return nil, ErrInvalidRequester
	}

	if len(baseURLs) == 0 {
		return nil, ErrInvalidURLMap
	}

	parsedBaseURLs := make(map[string]a.URL, len(baseURLs))

	for source, raw := range baseURLs {
		src, err := a.NewURL(raw)
		if err != nil {
			return nil, fmt.Errorf("%w [%s]: %v", ErrInvalidURL, source, err)
		}

		parsedBaseURLs[source] = src
	}

	return &Client{
		BaseURLs: parsedBaseURLs,
		log:      log,
		req:      req,
	}, nil
}

func MustNew(log *slog.Logger, req *infra.Requester, baseURLs map[string]string) *Client {
	const op = "digitalspb.MustNew"

	c, err := New(log, req, baseURLs)
	if err != nil {
		panic(fmt.Sprintf("%s: failed to init digital-spb client: %v", op, err))
	}

	return c
}

func (c *Client) do(
	ctx context.Context, baseURL *url.URL, endpoint string, query url.Values, out any,
) error {
	return c.req.JSON(ctx, http.MethodGet, infra.Target(baseURL, endpoint, query), nil, out)
}
