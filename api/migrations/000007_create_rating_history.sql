-- +goose Up
CREATE TABLE rating_history (
    id UUID PRIMARY KEY,
    tracked_location_id UUID NOT NULL REFERENCES tracked_locations (id) ON DELETE CASCADE,
    value NUMERIC(2, 1) NOT NULL CHECK (value BETWEEN 0.1 AND 9.9),
    calculated_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX rating_history_tracked_location_id_calculated_at_idx
    ON rating_history (tracked_location_id, calculated_at);

-- +goose StatementBegin
CREATE FUNCTION set_rating_history_updated_at()
RETURNS trigger
LANGUAGE plpgsql
AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$;
-- +goose StatementEnd

CREATE TRIGGER rating_history_set_updated_at
    BEFORE UPDATE ON rating_history
    FOR EACH ROW
    EXECUTE FUNCTION set_rating_history_updated_at();

-- +goose Down
DROP TABLE rating_history;
DROP FUNCTION set_rating_history_updated_at();
