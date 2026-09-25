package twogis

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/kotafan1rich/GeoLogic-Monitor/ingestion/internal/infra"
)

const (
	maxBodySize = infra.MaxBodySize

	dateLayout = "2000-01-01"

	fieldsList = "items.point,items.dates,items.rubrics"
	sortOrder  = "creation_time"
)

type Client struct {
	baseURL *url.URL
	apiKey  string
	log     *slog.Logger
	req     *infra.Requester
}

func New(log *slog.Logger, req *infra.Requester, baseURL, apiKey string) (*Client, error) {
	if log == nil {
		return nil, ErrInvalidLogger
	}

	if req == nil {
		return nil, ErrInvalidRequester
	}

	if baseURL == "" {
		return nil, ErrInvalidURL
	}

	if apiKey == "" {
		return nil, ErrInvalidAPIKey
	}

	parsedBaseURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, ErrParseURL
	}

	return &Client{
		baseURL: parsedBaseURL,
		apiKey:  apiKey,
		log:     log,
		req:     req,
	}, nil
}

func MustNew(log *slog.Logger, req *infra.Requester, baseURL, apiKey string) *Client {
	const op = "digitalspb.MustNew"

	c, err := New(log, req, baseURL, apiKey)
	if err != nil {
		panic(fmt.Sprintf("%s: failed to init 2gis client: %v", op, err))
	}

	return c
}

func (c *Client) do(ctx context.Context, endpoint string, query url.Values, out any) error {
	return c.req.JSON(ctx, http.MethodGet, infra.Target(c.baseURL, endpoint, query), nil, out)
}

func (c *Client) buildQuery(category string, since time.Time) url.Values {
	return url.Values{
		"q":                 {category},
		"key":               {c.apiKey},
		"city_id":           {spbCityID},
		"opened_after_date": {since.Format(dateLayout)},
		"sort":              {sortOrder},
		"fields":            {fieldsList},
		"page_size":         {strconv.Itoa(recordsPerPage)},
	}
}
