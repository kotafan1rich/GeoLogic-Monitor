package geocoding

import (
	"context"
	"errors"
	"log/slog"
	"strings"
	"unicode/utf8"

	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/domain"
	apperrs "github.com/kotafan1rich/GeoLogic-Monitor/api/internal/errs/app"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/logger"
)

const maxQueryLength = 500

var (
	errEmptyQuery   = errors.New("query is required")
	errQueryTooLong = errors.New("query is too long")
)

type Repository interface {
	Suggestions(ctx context.Context, query string) ([]domain.Address, error)
	Reverse(ctx context.Context, lat, lon float64) (*domain.Address, error)
}

type service struct {
	repo Repository
	log  *logger.Logger
}

func NewService(repo Repository, log *logger.Logger) *service {
	return &service{
		repo: repo,
		log:  log,
	}
}

func (s *service) Suggest(ctx context.Context, query string) ([]domain.Address, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return nil, apperrs.ValidationError(errEmptyQuery)
	}
	if utf8.RuneCountInString(query) > maxQueryLength {
		return nil, apperrs.ValidationError(errQueryTooLong)
	}

	result, err := s.repo.Suggestions(ctx, query)
	if err != nil {
		s.log.ErrorContext(ctx,
			"failed to get geocoding suggestions",
			slog.String("error", err.Error()),
		)
		return nil, apperrs.Wrap(err, apperrs.ErrProviderUnavailable)
	}
	if result == nil {
		return []domain.Address{}, nil
	}

	return result, nil
}

func (s *service) Reverse(ctx context.Context, geoPoint *domain.GeoPoint) (*domain.Address, error) {
	address, err := s.repo.Reverse(ctx, geoPoint.Lat, geoPoint.Lon)
	if err != nil {
		s.log.ErrorContext(ctx,
			"failed to get address from lat and lon",
			slog.String("error", err.Error()),
		)
		return nil, apperrs.Wrap(err, apperrs.ErrProviderUnavailable)
	}
	return address, nil
}
