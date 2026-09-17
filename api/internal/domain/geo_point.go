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

func validateLong(long float64) bool {
	if long < minLongitude || long > maxLongitude {
		return false
	}
	return true
}

type GeoPoint struct {
	Lat float64
	Lng float64
}

func NewGeoPoint(lat, lng float64) (*GeoPoint, error) {
	if !validateLat(lat) {
		return nil, errs.ErrInvalidLat
	}
	if !validateLong(lng) {
		return nil, errs.ErrInvalidLong
	}
	return &GeoPoint{Lat: lat, Lng: lng}, nil
}

func (g *GeoPoint) UpdateLat(lat float64) error {
	if !validateLat(lat) {
		return errs.ErrInvalidLat
	}
	g.Lat = lat
	return nil
}

func (g *GeoPoint) UpdateLong(lng float64) error {
	if !validateLong(lng) {
		return errs.ErrInvalidLong
	}
	g.Lng = lng
	return nil
}
