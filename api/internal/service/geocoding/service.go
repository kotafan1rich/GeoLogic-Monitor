package geocoding

import (
	"context"
	"errors"
	"log/slog"
	"strings"
	"unicode/utf8"

	apperrs "github.com/kotafan1rich/GeoLogic-Monitor/api/internal/errs/app"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/integrations/geocoder"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/logger"
)

const maxQueryLength = 500

var (
	errEmptyQuery   = errors.New("query is required")
	errQueryTooLong = errors.New("query is too long")
)

type Client interface {
	Autocomplete(ctx context.Context, search string) (*geocoder.Autocomplete, error)
}

type Suggestion struct {
	Address string
	Lat     float64
	Lon     float64
}

type Service interface {
	Suggest(ctx context.Context, query string) ([]Suggestion, error)
}

type service struct {
	client Client
	log    *logger.Logger
}

func NewService(client Client, log *logger.Logger) Service {
	return &service{
		client: client,
		log:    log,
	}
}

func (s *service) Suggest(ctx context.Context, query string) ([]Suggestion, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return nil, apperrs.ValidationError(errEmptyQuery)
	}
	if utf8.RuneCountInString(query) > maxQueryLength {
		return nil, apperrs.ValidationError(errQueryTooLong)
	}

	result, err := s.client.Autocomplete(ctx, query)
	if err != nil {
		s.log.ErrorContext(ctx,
			"failed to get geocoding suggestions",
			slog.String("error", err.Error()),
		)
		return nil, apperrs.Wrap(err, apperrs.ErrProviderUnavailable)
	}
	if result == nil {
		return []Suggestion{}, nil
	}

	return []Suggestion{{
		Address: formatAddress(result.Name, result.BuildingName),
		Lat:     result.Latitude,
		Lon:     result.Longitude,
	}}, nil
}

func formatAddress(name, buildingName string) string {
	name = strings.TrimSpace(name)
	buildingName = strings.TrimSpace(buildingName)
	if name == "" {
		return buildingName
	}
	if buildingName == "" || strings.Contains(name, buildingName) {
		return name
	}
	return name + ", " + buildingName
}
