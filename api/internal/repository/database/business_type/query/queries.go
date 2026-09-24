package query

const (
	Upsert = `
		INSERT INTO business_types (infra_type_id)
		VALUES ($1)
		ON CONFLICT (infra_type_id) DO UPDATE
		SET infra_type_id = EXCLUDED.infra_type_id
		RETURNING id
	`

	GetByID = `
		SELECT business_types.id, business_types.infra_type_id,
			infra_types.id, infra_types.slug, infra_types.name,
			infra_types.weight, infra_types.max_radius,
			business_types.created_at, business_types.updated_at
		FROM business_types
		JOIN infra_types ON infra_types.id = business_types.infra_type_id
		WHERE business_types.id = $1
	`

	GetAll = `
		SELECT business_types.id, business_types.infra_type_id,
			infra_types.id, infra_types.slug, infra_types.name,
			infra_types.weight, infra_types.max_radius,
			business_types.created_at, business_types.updated_at
		FROM business_types
		JOIN infra_types ON infra_types.id = business_types.infra_type_id
		ORDER BY business_types.id
	`

	Delete = `
		DELETE FROM business_types
		WHERE id = $1
		RETURNING id
	`
)
