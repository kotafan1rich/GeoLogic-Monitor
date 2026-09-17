package query

const (
	Create = `
		INSERT INTO tracked_locations (user_id, location)
		VALUES ($1, ST_GeomFromEWKT($2))
		RETURNING id, user_id, ST_AsEWKB(location), created_at, updated_at
	`

	GetByID = `
		SELECT id, user_id, ST_AsEWKB(location), created_at, updated_at
		FROM tracked_locations
		WHERE id = $1
	`

	GetByUserID = `
		SELECT id, user_id, ST_AsEWKB(location), created_at, updated_at
		FROM tracked_locations
		WHERE user_id = $1
		ORDER BY created_at ASC
	`

	Delete = `
		DELETE FROM tracked_locations
		WHERE id = $1
		RETURNING id
	`
)
