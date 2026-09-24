package geocoder

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const suggestionsLimit = 5

type geocoderClient struct {
	httpClient *http.Client
	baseURL    string
	apiKey     string
}

var _ GeocoderClient = (*geocoderClient)(nil)

func New(httpClient *http.Client, baseURL, apiKey string) *geocoderClient {
	return &geocoderClient{
		httpClient: httpClient,
		baseURL:    strings.TrimRight(baseURL, "/"),
		apiKey:     apiKey,
	}
}

func (c *geocoderClient) Autocomplete(ctx context.Context, search string) ([]Autocomplete, error) {
	request := suggestRequest{
		Query: search,
		Count: suggestionsLimit,
		Locations: []suggestLocation{{
			City: "Санкт-Петербург",
		}},
		RestrictValue: true,
	}
	var response suggestionsResponse
	if err := c.post(ctx, "/suggest/address", request, &response); err != nil {
		return nil, err
	}

	result := make([]Autocomplete, 0, len(response.Suggestions))
	for _, suggestion := range response.Suggestions {
		result = append(result, Autocomplete{
			Name:      suggestion.Value,
			Latitude:  suggestion.Data.GeoLat,
			Longitude: suggestion.Data.GeoLon,
		})
	}
	return result, nil
}

func (c *geocoderClient) Reverse(ctx context.Context, longitude, latitude float64) (*Geocode, error) {
	request := geolocateRequest{Lat: latitude, Lon: longitude, Count: 1}
	var response suggestionsResponse
	if err := c.post(ctx, "/geolocate/address", request, &response); err != nil {
		return nil, err
	}
	if len(response.Suggestions) == 0 {
		return nil, fmt.Errorf("DaData address not found")
	}

	suggestion := response.Suggestions[0]
	return &Geocode{
		Address: suggestion.Value,
		Center:  GeocodeCenter{X: suggestion.Data.GeoLon, Y: suggestion.Data.GeoLat},
	}, nil
}

func (c *geocoderClient) post(ctx context.Context, path string, body, result any) error {
	payload, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("encode DaData request: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+path, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("create DaData request: %w", err)
	}
	req.Header.Set("Authorization", "Token "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	response, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("execute DaData request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("DaData returned HTTP status %d", response.StatusCode)
	}
	if err := json.NewDecoder(response.Body).Decode(result); err != nil {
		return fmt.Errorf("decode DaData response: %w", err)
	}
	return nil
}
