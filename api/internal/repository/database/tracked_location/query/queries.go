package query

const (
	Create = `
		INSERT INTO tracked_locations (user_id, business_type_id, address, location)
		VALUES ($1, $2, $3, ST_GeomFromEWKT($4))
		RETURNING id, user_id, business_type_id, address,
			ST_AsEWKB(location), created_at, updated_at
	`

	GetByID = `
		SELECT id, user_id, business_type_id, address,
			ST_AsEWKB(location), created_at, updated_at
		FROM tracked_locations
		WHERE id = $1
	`

	GetByUserID = `
		SELECT id, user_id, business_type_id, address,
			ST_AsEWKB(location), created_at, updated_at
		FROM tracked_locations
		WHERE user_id = $1
		ORDER BY created_at ASC
	`

	GetAllForMonitoring = `
		SELECT tracked_locations.id, tracked_locations.user_id,
			tracked_locations.business_type_id, tracked_locations.address,
			ST_AsEWKB(tracked_locations.location), users.max_chat_id,
			tracked_locations.created_at, tracked_locations.updated_at
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
