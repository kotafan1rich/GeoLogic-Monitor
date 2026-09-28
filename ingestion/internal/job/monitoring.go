package job

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"math"
	"slices"
	"strings"
	"sync/atomic"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/kotafan1rich/GeoLogic-Monitor/ingestion/internal/aggregator/twogis"
	"github.com/kotafan1rich/GeoLogic-Monitor/ingestion/internal/errs"
	"github.com/kotafan1rich/GeoLogic-Monitor/ingestion/internal/producer"
	"github.com/kotafan1rich/GeoLogic-Monitor/ingestion/internal/repository/postgresql/competitor"
	"github.com/kotafan1rich/GeoLogic-Monitor/ingestion/internal/repository/postgresql/geo"
	"github.com/kotafan1rich/GeoLogic-Monitor/ingestion/internal/storage/geoapi"
)

const (
	MonitoringName  = "monitoring"
	CompetitorsName = "competitors"
)

const (
	defaultFetchConcurrency = 2
	defaultRouteConcurrency = 4

	maxDestinations        = 99
	walkingSpeedMPS        = 5.0 / 3.6
	defaultCheckpointStart = 24 * time.Hour

	reasonSameCategory = "Та же категория бизнеса"
	reasonWalkable     = "Пешая доступность"
)

type Monitoring struct {
	log         *slog.Logger
	geo         *geoapi.Client
	competitors *twogis.Client
	types       *geoapi.TypeRegistry
	repo        *competitor.CompetitorRepo
	producer    *producer.Producer
	checkpoints *twogis.CheckpointsCache
	topic       string
	window      time.Duration
	bootstrap   time.Duration
	ttl         time.Duration
	fetchLimit  int
	routeLimit  int
}

func NewMonitoring(
	log *slog.Logger,
	geoClient *geoapi.Client,
	competitors *twogis.Client,
	types *geoapi.TypeRegistry,
	repo *competitor.CompetitorRepo,
	p *producer.Producer,
	checkpoints *twogis.CheckpointsCache,
	topic string,
	window time.Duration,
	bootstrap time.Duration,
	ttl time.Duration,
	fetchLimit int,
	routeLimit int,
) *Monitoring {
	if bootstrap <= 0 {
		bootstrap = defaultCheckpointStart
	}

	return &Monitoring{
		log:         log.With(slog.String("source", MonitoringName)),
		geo:         geoClient,
		competitors: competitors,
		types:       types,
		repo:        repo,
		producer:    p,
		checkpoints: checkpoints,
		topic:       topic,
		window:      window,
		bootstrap:   bootstrap,
		ttl:         ttl,
		fetchLimit:  fetchLimit,
		routeLimit:  routeLimit,
	}
}

type businessGroup struct {
	id         string
	def        geoapi.InfraTypeInput
	locations  []geoapi.TrackedLocation
	fetched    int
	mismatched int
	fresh      []twogis.Item
	checkpoint time.Time
	failed     atomic.Bool
}

func (m *Monitoring) Competitors(ctx context.Context) error {
	start := time.Now()

	if !m.types.HasBusinessTypes() {
		m.log.InfoContext(ctx, "competitors monitoring skipped: business types are not connected yet")

		return nil
	}

	locations, err := m.geo.TrackedLocations(ctx)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrLoadLocations, err)
	}

	// Мок-данные
	// locations, err := m.testTrackedLocations(ctx)
	// if err != nil {
	// 	return fmt.Errorf("%w: %v", ErrLoadLocations, err)
	// }

	groups := m.groupByBusinessType(ctx, locations)
	if len(groups) == 0 {
		m.log.InfoContext(ctx, "competitors monitoring skipped: no groups",
			slog.Int("locations", len(locations)),
		)

		return nil
	}

	collector := errs.NewCollector(errs.DefaultSampleSize)

	m.fetchCompetitors(ctx, groups, collector)
	notified := m.matchCompetitors(ctx, groups, collector)

	var fetched, mismatched, fresh, advanced int

	for _, g := range groups {
		fetched += g.fetched
		mismatched += g.mismatched
		fresh += len(g.fresh)

		if g.failed.Load() {
			continue
		}

		m.checkpoints.Set(g.def.Slug, g.checkpoint)
		advanced++
	}

	attrs := []any{
		slog.Int("locations", len(locations)),
		slog.Int("groups", len(groups)),
		slog.Int("fetched", fetched),
		slog.Int("type_mismatched", mismatched),
		slog.Int("fresh", fresh),
		slog.Int("notified", notified),
		slog.Int("checkpoints_advanced", advanced),
		slog.Int("failed", collector.Total()),
		slog.Duration("duration", time.Since(start)),
	}

	if n := collector.Total(); n > 0 {
		m.log.ErrorContext(ctx, "competitors monitoring finished with errors",
			append(attrs, collector.Attr())...)

		return fmt.Errorf("%w: %d error(s)", ErrMonitoringFailed, n)
	}

	m.log.InfoContext(ctx, "competitors monitoring finished", attrs...)

	return nil
}

func (m *Monitoring) Events(ctx context.Context) error {
	// TODO: События для оповещения
	return nil
}

func (m *Monitoring) groupByBusinessType(
	ctx context.Context, locations []geoapi.TrackedLocation,
) []*businessGroup {
	var (
		byID    = make(map[string]*businessGroup)
		groups  []*businessGroup
		unknown int
	)

	for _, loc := range locations {
		id := loc.BusinessTypeID.String()

		g, ok := byID[id]
		if !ok {
			def, found := m.types.BusinessType(id)
			if !found {
				unknown++
				continue
			}

			g = &businessGroup{id: id, def: def}
			byID[id] = g
			groups = append(groups, g)
		}

		g.locations = append(g.locations, loc)
	}

	if unknown > 0 {
		m.log.WarnContext(ctx, "tracked locations with unknown business type skipped",
			slog.Int("skipped", unknown),
		)
	}

	return groups
}

func (m *Monitoring) fetchCompetitors(
	ctx context.Context, groups []*businessGroup, collector *errs.Collector,
) {
	var g errgroup.Group

	g.SetLimit(m.fetchLimit)

	since := time.Now().Add(-m.window)

	for _, group := range groups {
		g.Go(func() error {
			items, err := m.competitors.ParseCompetitors(ctx, searchQuery(group.def), since)
			if err != nil {
				group.failed.Store(true)
				collector.Add(fmt.Errorf("%w [%s]: %v", ErrFetchCompetitors, group.def.Slug, err))

				return nil
			}

			group.fetched = len(items)
			m.freshCompetitors(ctx, group, items)

			return nil
		})
	}

	_ = g.Wait()
}

func (m *Monitoring) freshCompetitors(ctx context.Context, group *businessGroup, items []twogis.Item) {
	last, ok := m.checkpoints.Get(group.def.Slug)
	if !ok {
		last = time.Now().Add(-m.bootstrap)
	}

	now := time.Now()
	next := last
	fresh := make([]twogis.Item, 0)

	for _, item := range items {
		if item.Dates == nil || !item.Dates.CreatedAt.After(last) {
			continue
		}

		if item.Dates.CreatedAt.After(next) && !item.Dates.CreatedAt.After(now) {
			next = item.Dates.CreatedAt
		}

		if item.Point == nil {
			continue
		}

		if !sameBusinessType(group.def, item) {
			group.mismatched++

			m.log.DebugContext(ctx, "competitor skipped, business type mismatch",
				slog.String("business_type", group.def.Slug),
				slog.String("external_id", item.ID),
				slog.String("name", item.Name),
				slog.Any("rubrics", rubricNames(item)),
			)

			continue
		}

		fresh = append(fresh, item)
	}

	group.fresh, group.checkpoint = fresh, next

	m.log.DebugContext(ctx, "competitors fetched",
		slog.String("business_type", group.def.Slug),
		slog.Int("fetched", len(items)),
		slog.Int("type_mismatched", group.mismatched),
		slog.Int("fresh", len(fresh)),
		slog.Time("checkpoint_from", last),
		slog.Time("checkpoint_to", next),
	)
}

func sameBusinessType(def geoapi.InfraTypeInput, item twogis.Item) bool {
	if len(def.Rubrics) == 0 {
		return true
	}

	for _, r := range item.Rubrics {
		name := strings.ToLower(r.Name)

		for _, want := range def.Rubrics {
			if strings.Contains(name, strings.ToLower(want)) {
				return true
			}
		}
	}

	return false
}

func rubricNames(item twogis.Item) []string {
	names := make([]string, 0, len(item.Rubrics))
	for _, r := range item.Rubrics {
		names = append(names, r.Name)
	}

	return names
}

func (m *Monitoring) matchCompetitors(
	ctx context.Context, groups []*businessGroup, collector *errs.Collector,
) int {
	var (
		g        errgroup.Group
		notified atomic.Int64
	)

	g.SetLimit(m.routeLimit)

	for _, group := range groups {
		if group.failed.Load() || len(group.fresh) == 0 {
			continue
		}

		dests := make([]geoapi.GeoPoint, 0, len(group.locations))
		for _, loc := range group.locations {
			dests = append(dests, geoapi.GeoPoint{Lat: loc.Lat, Lon: loc.Lon})
		}

		for _, item := range group.fresh {
			g.Go(func() error {
				notified.Add(int64(m.matchCompetitor(ctx, group, item, dests, collector)))
				return nil
			})
		}
	}

	_ = g.Wait()

	return int(notified.Load())
}

func (m *Monitoring) matchCompetitor(
	ctx context.Context,
	group *businessGroup,
	item twogis.Item,
	dests []geoapi.GeoPoint,
	collector *errs.Collector,
) int {
	src := geoapi.GeoPoint{Lat: item.Point.Lat, Lon: item.Point.Lon}

	distances, err := m.walkingDistances(ctx, src, dests)
	if err != nil {
		group.failed.Store(true)
		collector.Add(fmt.Errorf("%w [%s:%s]: %v", ErrRouteCompetitor, group.def.Slug, item.ID, err))

		return 0
	}

	var notified int

	for i, d := range distances {
		if d == nil || *d > float64(group.def.MaxRadius) {
			continue
		}

		sent, err := m.notifyCompetitor(ctx, group, item, group.locations[i], *d)
		if err != nil {
			group.failed.Store(true)
			collector.Add(err)
			continue
		}

		if sent {
			notified++
		}
	}

	return notified
}

func (m *Monitoring) notifyCompetitor(
	ctx context.Context,
	group *businessGroup,
	item twogis.Item,
	loc geoapi.TrackedLocation,
	distance float64,
) (bool, error) {
	locationID := loc.ID.String()
	openedAt := item.Dates.CreatedAt

	alreadyNotified, err := m.repo.Add(
		ctx,
		locationID,
		item.ID,
		item.Name,
		group.id,
		item.AddressName,
		geo.GeoPoint{Lat: item.Point.Lat, Lon: item.Point.Lon},
		openedAt,
		m.ttl,
	)
	if err != nil {
		return false, fmt.Errorf("%w [%s:%s]: %v", ErrStoreCompetitor, locationID, item.ID, err)
	}

	if alreadyNotified {
		m.log.DebugContext(ctx, "competitor already notified, skipped",
			slog.String("tracked_location_id", locationID),
			slog.String("external_id", item.ID),
		)

		return false, nil
	}

	n := newNotification(producer.TypeCompetitorOpened, loc, producer.Subject{
		ExternalID: item.ID,
		Title:      item.Name,
		Category:   group.def.Name,
		Address:    item.AddressName,
		OpenedAt:   openedAt.Format(time.DateOnly),
	}, distance, []string{reasonSameCategory, reasonWalkable})

	if err := m.publish(ctx, n); err != nil {
		return false, fmt.Errorf("[%s:%s]: %w", locationID, item.ID, err)
	}

	if err := m.repo.UpdateStatus(ctx, locationID, item.ID, time.Now()); err != nil {
		return false, fmt.Errorf("%w [%s:%s]: %v", ErrMarkNotified, locationID, item.ID, err)
	}

	m.log.InfoContext(ctx, "competitor notification published",
		slog.String("tracked_location_id", locationID),
		slog.String("business_type", group.def.Slug),
		slog.String("external_id", item.ID),
		slog.String("name", item.Name),
		slog.Float64("distance_m", distance),
		slog.String("topic", m.topic),
	)

	return true, nil
}

func (m *Monitoring) walkingDistances(
	ctx context.Context, src geoapi.GeoPoint, dests []geoapi.GeoPoint,
) ([]*float64, error) {
	distances := make([]*float64, 0, len(dests))

	for chunk := range slices.Chunk(dests, maxDestinations) {
		part, err := m.geo.CalculateWalkingDistance(ctx, src, chunk)
		if err != nil {
			return nil, err
		}

		if len(part) != len(chunk) {
			return nil, fmt.Errorf("%w: got %d, want %d", ErrDistancesMismatch, len(part), len(chunk))
		}

		distances = append(distances, part...)
	}

	return distances, nil
}

func (m *Monitoring) publish(ctx context.Context, n producer.Notification) error {
	payload, err := json.Marshal(n)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrPublishNotification, err)
	}

	if err := m.producer.Produce(ctx, string(payload), m.topic); err != nil {
		return fmt.Errorf("%w: %v", ErrPublishNotification, err)
	}

	return nil
}

func newNotification(
	kind string,
	loc geoapi.TrackedLocation,
	subject producer.Subject,
	distance float64,
	reasons []string,
) producer.Notification {
	return producer.Notification{
		Type:      kind,
		MaxChatID: loc.MaxChatID,
		TrackedLocation: producer.TrackedLocation{
			ID:             loc.ID.String(),
			Name:           loc.Name,
			Address:        loc.Address,
			BusinessTypeID: loc.BusinessTypeID.String(),
		},
		Subject: subject,
		Route: producer.Route{
			DistanceMeters:  distance,
			DurationSeconds: math.Round(distance / walkingSpeedMPS),
		},
		Assessment: producer.Assessment{
			Reasons:         reasons,
			Recommendations: []string{},
		},
	}
}

func searchQuery(def geoapi.InfraTypeInput) string {
	if def.Query != "" {
		return def.Query
	}

	return def.Name
}
