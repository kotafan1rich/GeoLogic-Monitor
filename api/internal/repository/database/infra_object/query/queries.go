package query

const (
	Upsert = `
		INSERT INTO infra_objects (id, type_id, location, address, name)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (id) DO UPDATE
		SET type_id = EXCLUDED.type_id,
			location = EXCLUDED.location,
			address = EXCLUDED.address,
			name = EXCLUDED.name
		RETURNING id
	`

	GetByID = `
		SELECT infra_objects.id, infra_objects.type_id, infra_objects.location,
			infra_objects.address, infra_objects.name,
			infra_types.id, infra_types.slug, infra_types.name,
			infra_types.weight, infra_types.max_radius,
			infra_objects.created_at, infra_objects.updated_at
		FROM infra_objects
		JOIN infra_types ON infra_types.id = infra_objects.type_id
		WHERE infra_objects.id = $1
	`

	Near = `
		SELECT infra_objects.id, infra_objects.type_id, infra_objects.location,
			infra_objects.address, infra_objects.name,
			infra_types.id, infra_types.slug, infra_types.name,
			infra_types.weight, infra_types.max_radius,
			infra_objects.created_at, infra_objects.updated_at
		FROM infra_objects
		JOIN infra_types ON infra_types.id = infra_objects.type_id
		WHERE ST_DWithin(
			infra_objects.location::geography,
			$1::geometry::geography,
			infra_types.max_radius
		)
		ORDER BY infra_objects.id
	`
)
