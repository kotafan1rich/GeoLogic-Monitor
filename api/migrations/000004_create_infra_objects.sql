-- +goose Up
CREATE TABLE infra_objects (
    id UUID PRIMARY KEY,
    type_id UUID NOT NULL REFERENCES infra_types (id),
    location geometry(Point, 4326) NOT NULL,
    address TEXT NOT NULL CHECK (btrim(address) <> ''),
    name TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX infra_objects_type_id_idx
    ON infra_objects (type_id);

CREATE INDEX infra_objects_location_geography_idx
    ON infra_objects USING gist ((location::geography));

-- +goose StatementBegin
CREATE FUNCTION set_infra_objects_updated_at()
RETURNS trigger
LANGUAGE plpgsql
AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$;
-- +goose StatementEnd

CREATE TRIGGER infra_objects_set_updated_at
    BEFORE UPDATE ON infra_objects
    FOR EACH ROW
    EXECUTE FUNCTION set_infra_objects_updated_at();

-- +goose Down
DROP TABLE infra_objects;
DROP FUNCTION set_infra_objects_updated_at();
