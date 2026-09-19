package query

const (
	Upsert = `
		INSERT INTO events (provider, external_id, location, date, info)
		VALUES ($1, $2, ST_GeomFromEWKT($3), $4, $5)
		ON CONFLICT (provider, external_id) DO UPDATE
		SET location = EXCLUDED.location,
			date = EXCLUDED.date,
			info = EXCLUDED.info
		RETURNING id, provider, external_id, ST_AsEWKB(location), date, info,
			notified_at, created_at, updated_at
	`

	GetByID = `
		SELECT id, provider, external_id, ST_AsEWKB(location), date, info,
			notified_at, created_at, updated_at
		FROM events
		WHERE id = $1
	`

	GetUnnotifiedByPeriod = `
		SELECT id, provider, external_id, ST_AsEWKB(location), date, info,
			notified_at, created_at, updated_at
		FROM events
		WHERE notified_at IS NULL
			AND date >= $1
			AND date < $2
		ORDER BY date ASC, id ASC
	`

	GetUnnotifiedNear = `
		SELECT id, provider, external_id, ST_AsEWKB(location), date, info,
			notified_at, created_at, updated_at
		FROM events
		WHERE notified_at IS NULL
			AND ST_DWithin(
				location::geography,
				ST_GeomFromEWKT($1)::geography,
				$2
			)
			AND ($3::timestamptz IS NULL OR date >= $3)
			AND ($4::timestamptz IS NULL OR date < $4)
		ORDER BY date ASC, id ASC
	`

	MarkNotified = `
		UPDATE events
		SET notified_at = COALESCE(notified_at, now())
		WHERE id = $1
		RETURNING notified_at
	`
)
