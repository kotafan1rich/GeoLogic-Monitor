package domain

import (
	"time"

	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/errs"
)

type CalculatedRating struct {
	Value        float64
	CalculatedAt time.Time
}

func NewCalculatedRating(value float64, calculatedAt time.Time) (*CalculatedRating, error) {
	if value < 0 || value > 10 {
		return nil, errs.ErrInvalidRatingValue
	}
	return &CalculatedRating{
		Value:        value,
		CalculatedAt: calculatedAt,
	}, nil
}
