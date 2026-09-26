package competitor

const (
	AddCompetitorQuery = `
		INSERT INTO competitors (
		tracked_location_id, external_id, competitor_name, type_id, 
		competitor_address, competitor_location, opened_at, notified_at,
		expires_at, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NULL, $8, NOW(), NOW())
		ON CONFLICT (tracked_location_id, external_id) DO NOTHING
	`

	ChangeStatusQuery = `
		UPDATE competitors
		SET
			notified_at = $3,
			updated_at = NOW()
		WHERE 
			tracked_location_id = $1
			AND external_id = $2
			AND notified_at IS NULL
	`
)
