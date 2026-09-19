package osrm

import "context"

type OSRMClient interface {
	GetWalkingDistances(ctx context.Context, src *Coordinate, dst []*Coordinate) ([]*float64, error)
}
