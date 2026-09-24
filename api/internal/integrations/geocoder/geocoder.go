package geocoder

import "context"

type GeocoderClient interface {
	Reverse(ctx context.Context, longitude, latitude float64) (*Geocode, error)
	Autocomplete(ctx context.Context, search string) ([]Autocomplete, error)
}
