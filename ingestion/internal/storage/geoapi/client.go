package geoapi

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/kotafan1rich/GeoLogic-Monitor/ingestion/internal/infra"
)

type Client struct {
	baseURL *url.URL
	write   *infra.Requester
	geocode *infra.Requester
}

func New(baseURL string, write, geocode *infra.Requester) (*Client, error) {
	if write == nil || geocode == nil {
		return nil, ErrInvalidRequester
	}

	parsed, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrParseURL, err)
	}

	if !validateURL(parsed) {
		return nil, fmt.Errorf("%w [%s]", ErrInvalidURL, baseURL)
	}

	return &Client{baseURL: parsed, write: write, geocode: geocode}, nil
}

func MustNew(baseURL string, write, geocode *infra.Requester) *Client {
	const op = "geoapi.MustNew"

	c, err := New(baseURL, write, geocode)
	if err != nil {
		panic(fmt.Sprintf("%s: failed to init geo-api client: %v", op, err))
	}

	return c
}

func (c *Client) put(ctx context.Context, endpoint string, body, out any) error {
	return c.write.JSON(ctx, http.MethodPut, infra.Target(c.baseURL, endpoint, nil), body, out)
}

func (c *Client) get(ctx context.Context, endpoint string, query url.Values, out any) error {
	return c.geocode.JSON(ctx, http.MethodGet, infra.Target(c.baseURL, endpoint, query), nil, out)
}

func validateURL(u *url.URL) bool {
	return u != nil && u.Host != "" && u.Scheme != ""
}
