package query

const (
	Upsert = `
		INSERT INTO infra_types (slug, name, weight, max_radius)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (slug) DO UPDATE
		SET name = EXCLUDED.name,
			weight = EXCLUDED.weight,
			max_radius = EXCLUDED.max_radius
		RETURNING id, slug, name, weight, max_radius, created_at, updated_at
	`

	GetByID = `
		SELECT id, slug, name, weight, max_radius, created_at, updated_at
		FROM infra_types
		WHERE id = $1
	`
)
