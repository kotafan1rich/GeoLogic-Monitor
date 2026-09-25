package domain

import "github.com/kotafan1rich/GeoLogic-Monitor/api/internal/errs"

const (
	maxLatitude  = 90.0
	minLatitude  = -90.0
	maxLongitude = 180.0
	minLongitude = -180.0
)

func validateLat(lat float64) bool {
	if lat < minLatitude || lat > maxLatitude {
		return false
	}
	return true
}

func validateLon(lon float64) bool {
	if lon < minLongitude || lon > maxLongitude {
		return false
	}
	return true
}

type GeoPoint struct {
	Lat float64
	Lon float64
}

func NewGeoPoint(lat, lon float64) (*GeoPoint, error) {
	if !validateLat(lat) {
		return nil, errs.ErrInvalidLat
	}
	if !validateLon(lon) {
		return nil, errs.ErrInvalidLon
	}
	return &GeoPoint{Lat: lat, Lon: lon}, nil
}

func (g *GeoPoint) UpdateLat(lat float64) error {
	if !validateLat(lat) {
		return errs.ErrInvalidLat
	}
	g.Lat = lat
	return nil
}

func (g *GeoPoint) UpdateLon(lon float64) error {
	if !validateLon(lon) {
		return errs.ErrInvalidLon
	}
	g.Lon = lon
	return nil
}
