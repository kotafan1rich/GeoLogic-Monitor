package job

import "errors"

var (
	ErrResolveType         = errors.New("failed to resolve infra type id")
	ErrStoreFailed         = errors.New("failed to store data")
	ErrConnectBusiness     = errors.New("failed to connect business type with infra")
	ErrLoadLocations       = errors.New("failed to load tracked locations")
	ErrLoadBusinessTypes   = errors.New("failed to load business types")
	ErrFetchCompetitors    = errors.New("failed to fetch competitors")
	ErrRouteCompetitor     = errors.New("failed to calculate walking distances")
	ErrDistancesMismatch   = errors.New("distances count mismatch")
	ErrStoreCompetitor     = errors.New("failed to store competitor")
	ErrPublishNotification = errors.New("failed to publish notification")
	ErrMarkNotified        = errors.New("failed to mark competitor notified")
	ErrMonitoringFailed    = errors.New("monitoring finished with errors")
)
