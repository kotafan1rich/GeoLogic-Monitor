package geocoder

import "context"

type GeocoderClient interface {
	ParseFree(ctx context.Context, street string) (*Address, error)
	Reverse(ctx context.Context, longitude, latitude float64) (*Geocode, error)
	ParseEAS(ctx context.Context, street string) (*EASAddress, error)
	Autocomplete(ctx context.Context, search string) ([]Autocomplete, error)
}
