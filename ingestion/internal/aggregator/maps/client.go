package maps

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strings"

	a "github.com/kotafan1rich/GeoLogic-Monitor/ingestion/internal/aggregator"
	"github.com/kotafan1rich/GeoLogic-Monitor/ingestion/internal/infra"
)

const (
	isoRegion = "RU-SPE"

	timeoutSec = 900
)

type Client struct {
	BaseURL a.URL
	log     *slog.Logger
	req     *infra.Requester
}

func New(log *slog.Logger, req *infra.Requester, baseURL string) (*Client, error) {
	if log == nil {
		return nil, ErrInvalidLogger
	}

	if req == nil {
		return nil, ErrInvalidRequester
	}

	if baseURL == "" {
		return nil, ErrInvalidURL
	}

	parsedURL, err := a.NewURL(baseURL)
	if err != nil {
		return nil, ErrParseURL
	}

	return &Client{
		BaseURL: parsedURL,
		log:     log,
		req:     req,
	}, nil
}

func MustNew(log *slog.Logger, req *infra.Requester, baseURL string) *Client {
	const op = "maps.MustNew"

	c, err := New(log, req, baseURL)
	if err != nil {
		panic(fmt.Sprintf("%s: failed to init maps-mail-ru client: %v", op, err))
	}

	return c
}

func (c *Client) do(
	ctx context.Context, baseURL *url.URL, endpoint string, query url.Values, out any,
) error {
	return c.req.Form(ctx, http.MethodPost, infra.Target(baseURL, endpoint, query), query, out)
}

func (c *Client) buildQuery() url.Values {
	amenity := []string{
		"cafe",
		"restaurant",
		"fast_food",
		"food_court",
		"pharmacy",
		"dentist",
		"veterinary",
		"car_wash",
		"parcel_locker",
	}
	healthcare := []string{"pharmacy", "dentist", "veterinary"}

	var b strings.Builder
	fmt.Fprintf(&b, "[out:json][timeout:%d];\n", timeoutSec)
	fmt.Fprintf(&b, "area[\"ISO3166-2\"=%q]->.a;\n(\n", isoRegion)
	b.WriteString("  nwr(area.a)[\"shop\"];\n")
	b.WriteString("  nwr(area.a)[\"craft\"];\n")
	fmt.Fprintf(&b, "  nwr(area.a)[\"amenity\"~\"^(%s)$\"];\n", strings.Join(amenity, "|"))
	fmt.Fprintf(&b, "  nwr(area.a)[\"healthcare\"~\"^(%s)$\"];\n", strings.Join(healthcare, "|"))
	b.WriteString(");\n")
	b.WriteString("out center meta;\n")

	return url.Values{"data": {b.String()}}
}
