package query

const (
	Upsert = `
		INSERT INTO users (max_user_id, max_chat_id)
		VALUES ($1, $2)
		ON CONFLICT (max_user_id) DO UPDATE
		SET max_chat_id = EXCLUDED.max_chat_id
		RETURNING id, max_user_id, max_chat_id, created_at, updated_at
	`

	GetByID = `
		SELECT id, max_user_id, max_chat_id, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	GetByMaxUserID = `
		SELECT id, max_user_id, max_chat_id, created_at, updated_at
		FROM users
		WHERE max_user_id = $1
	`
)
