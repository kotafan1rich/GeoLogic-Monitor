package geocoder

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type geocoderClient struct {
	httpClient *http.Client
	baseURL    string
}

var _ GeocoderClient = (*geocoderClient)(nil)

func New(client *http.Client, baseURL string) *geocoderClient {
	return &geocoderClient{
		httpClient: client,
		baseURL:    baseURL,
	}
}

func (c *geocoderClient) ParseFree(ctx context.Context, street string) (*Address, error) {
	query := url.Values{}
	query.Set("street", street)

	return get[Address](ctx, c, "/parse/free", query)
}

func (c *geocoderClient) Reverse(ctx context.Context, longitude, latitude float64) (*Geocode, error) {
	query := url.Values{}
	query.Set("x", fmt.Sprint(longitude))
	query.Set("y", fmt.Sprint(latitude))

	return get[Geocode](ctx, c, "/geocode/reverse", query)
}

func (c *geocoderClient) ParseEAS(ctx context.Context, street string) (*EASAddress, error) {
	query := url.Values{}
	query.Set("street", street)

	return get[EASAddress](ctx, c, "/parse/eas", query)
}

func (c *geocoderClient) Autocomplete(ctx context.Context, search string) (*Autocomplete, error) {
	query := url.Values{}
	query.Set("s", search)

	return get[Autocomplete](ctx, c, "/autocomplete/universal", query)
}

func get[T any](ctx context.Context, client *geocoderClient, path string, query url.Values) (*T, error) {
	requestURL := strings.TrimRight(client.baseURL, "/") + path + "?" + query.Encode()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create geocoder request: %w", err)
	}

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("execute geocoder request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("geocoder returned HTTP status %d", resp.StatusCode)
	}

	var result T
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode geocoder response: %w", err)
	}

	return &result, nil
}
