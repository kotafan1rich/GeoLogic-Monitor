package query

const (
	Create = `
		INSERT INTO rating_history (id, tracked_location_id, value, calculated_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, tracked_location_id, value, calculated_at, created_at, updated_at
	`

	CreateMany = `
		INSERT INTO rating_history (id, tracked_location_id, value, calculated_at)
		SELECT id, tracked_location_id, value, calculated_at
		FROM unnest(
			$1::uuid[],
			$2::uuid[],
			$3::double precision[],
			$4::timestamptz[]
		) AS ratings(id, tracked_location_id, value, calculated_at)
		RETURNING id, tracked_location_id, value, calculated_at, created_at, updated_at
	`

	GetHistory = `
		SELECT id, tracked_location_id, value, calculated_at, created_at, updated_at
		FROM rating_history
		WHERE tracked_location_id = $1
			AND calculated_at >= (
				(now() AT TIME ZONE 'Europe/Moscow' - make_interval(months => $2))
				AT TIME ZONE 'Europe/Moscow'
			)
			AND calculated_at <= now()
		ORDER BY calculated_at ASC, id ASC
	`
)
