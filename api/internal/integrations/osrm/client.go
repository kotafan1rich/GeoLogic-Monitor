package osrm

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

)

type osrmClient struct {
	httpClient *http.Client
	baseUrl    string
}

func New(client *http.Client, baseUrl string) *osrmClient {
	return &osrmClient{
		httpClient: client,
		baseUrl:    baseUrl,
	}
}

func (o *osrmClient) GetWalkingDistances(ctx context.Context, src *Coordinate, dst []*Coordinate) ([]*float64, error) {
	path := "/table/v1/foot/" + joinCoordinates(src, dst)

	query := url.Values{}
	query.Set("sources", "0")
	query.Set("destinations", joinIndices(len(dst)))
	query.Set("annotations", "distance")
	query.Set("skip_waypoints", "true")

	requestURL := strings.TrimRight(o.baseUrl, "/") +
		path +
		"?" +
		query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create OSRM request: %w", err)
	}

	resp, err := o.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("execute OSRM request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("OSRM returned HTTP status %d", resp.StatusCode)
	}

	var response tableResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("decode OSRM response: %w", err)
	}

	if response.Code != "Ok" {
		return nil, fmt.Errorf("OSRM returned code %q", response.Code)
	}

	if len(response.Distances) != 1 {
		return nil, fmt.Errorf(
			"unexpected OSRM distance matrix rows: %d",
			len(response.Distances),
		)
	}

	if len(response.Distances[0]) != len(dst) {
		return nil, fmt.Errorf(
			"unexpected OSRM distance count: got %d, want %d",
			len(response.Distances[0]),
			len(dst),
		)
	}

	return response.Distances[0], nil
}

func formatCoordinate(point *Coordinate) string {
	return strconv.FormatFloat(point.Lon, 'f', -1, 64) +
		"," +
		strconv.FormatFloat(point.Lat, 'f', -1, 64)
}

func joinCoordinates(src *Coordinate, dst []*Coordinate) string {
	coordinates := make([]string, 0, len(dst)+1)
	coordinates = append(coordinates, formatCoordinate(src))

	for _, destination := range dst {
		coordinates = append(coordinates, formatCoordinate(destination))
	}

	return strings.Join(coordinates, ";")
}

func joinIndices(count int) string {
	indices := make([]string, count)
	for i := range count {
		indices[i] = strconv.Itoa(i + 1)
	}
	return strings.Join(indices, ";")
}
