-- +goose Up
-- Depends on users and business_types created by earlier migrations.
CREATE TABLE tracked_locations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users (id),
    business_type_id UUID NOT NULL REFERENCES business_types (id),
    address TEXT NOT NULL CHECK (btrim(address) <> ''),
    location geometry(Point, 4326) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX tracked_locations_user_id_idx
    ON tracked_locations (user_id);

-- +goose StatementBegin
CREATE FUNCTION set_tracked_locations_updated_at()
RETURNS trigger
LANGUAGE plpgsql
AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$;
-- +goose StatementEnd

CREATE TRIGGER tracked_locations_set_updated_at
    BEFORE UPDATE ON tracked_locations
    FOR EACH ROW
    EXECUTE FUNCTION set_tracked_locations_updated_at();

-- +goose Down
DROP TABLE tracked_locations;
DROP FUNCTION set_tracked_locations_updated_at();
