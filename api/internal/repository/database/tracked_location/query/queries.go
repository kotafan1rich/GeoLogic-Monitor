package query

const (
	Create = `
		WITH created_location AS (
			INSERT INTO tracked_locations (user_id, business_type_id, address, location)
			VALUES ($1, $2, $3, ST_GeomFromEWKT($4))
			RETURNING id, user_id, business_type_id, address, location,
				created_at, updated_at
		)
		SELECT created_location.id, created_location.user_id,
			created_location.business_type_id, created_location.address,
			ST_AsEWKB(created_location.location), created_location.created_at,
			created_location.updated_at, users.id, users.max_user_id,
			users.max_chat_id, users.created_at, users.updated_at
		FROM created_location
		JOIN users ON users.id = created_location.user_id
	`

	GetByID = `
		SELECT tracked_locations.id, tracked_locations.user_id,
			tracked_locations.business_type_id, tracked_locations.address,
			ST_AsEWKB(tracked_locations.location), tracked_locations.created_at,
			tracked_locations.updated_at, users.id, users.max_user_id,
			users.max_chat_id, users.created_at, users.updated_at
		FROM tracked_locations
		JOIN users ON users.id = tracked_locations.user_id
		WHERE tracked_locations.id = $1
	`

	GetByUserID = `
		SELECT tracked_locations.id, tracked_locations.user_id,
			tracked_locations.business_type_id, tracked_locations.address,
			ST_AsEWKB(tracked_locations.location), tracked_locations.created_at,
			tracked_locations.updated_at, users.id, users.max_user_id,
			users.max_chat_id, users.created_at, users.updated_at,
			latest_rating.value, latest_rating.calculated_at
		FROM tracked_locations
		JOIN users ON users.id = tracked_locations.user_id
		LEFT JOIN LATERAL (
			SELECT value, calculated_at
			FROM rating_history
			WHERE tracked_location_id = tracked_locations.id
			ORDER BY calculated_at DESC, id DESC
			LIMIT 1
		) AS latest_rating ON TRUE
		WHERE tracked_locations.user_id = $1
		ORDER BY tracked_locations.created_at ASC
	`

	GetAllForMonitoring = `
		SELECT tracked_locations.id, tracked_locations.user_id,
			tracked_locations.business_type_id, tracked_locations.address,
			ST_AsEWKB(tracked_locations.location), tracked_locations.created_at,
			tracked_locations.updated_at, users.id, users.max_user_id,
			users.max_chat_id, users.created_at, users.updated_at
		FROM tracked_locations
		JOIN users ON users.id = tracked_locations.user_id
		ORDER BY tracked_locations.created_at ASC
	`

	Delete = `
		DELETE FROM tracked_locations
		WHERE id = $1
		RETURNING id
	`
)
