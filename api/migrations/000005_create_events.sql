-- +goose Up
CREATE TABLE events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    provider TEXT NOT NULL CHECK (btrim(provider) <> ''),
    external_id TEXT NOT NULL CHECK (btrim(external_id) <> ''),
    location geometry(Point, 4326) NOT NULL,
    date TIMESTAMPTZ NOT NULL,
    info TEXT,
    notified_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (provider, external_id)
);

CREATE INDEX events_date_idx
    ON events (date);

CREATE INDEX events_location_geography_idx
    ON events USING gist ((location::geography));

-- +goose StatementBegin
CREATE FUNCTION set_events_updated_at()
RETURNS trigger
LANGUAGE plpgsql
AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$;
-- +goose StatementEnd

CREATE TRIGGER events_set_updated_at
    BEFORE UPDATE ON events
    FOR EACH ROW
    EXECUTE FUNCTION set_events_updated_at();

-- +goose Down
DROP TABLE events;
DROP FUNCTION set_events_updated_at();
