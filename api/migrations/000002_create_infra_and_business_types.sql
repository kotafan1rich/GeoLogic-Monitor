-- +goose Up
CREATE TABLE infra_types (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    slug TEXT NOT NULL UNIQUE CHECK (btrim(slug) <> ''),
    name TEXT NOT NULL CHECK (btrim(name) <> ''),
    weight DOUBLE PRECISION NOT NULL CHECK (weight > 0),
    max_radius INTEGER NOT NULL CHECK (max_radius BETWEEN 1 AND 65535),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- +goose StatementBegin
CREATE FUNCTION set_infra_types_updated_at()
RETURNS trigger
LANGUAGE plpgsql
AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$;
-- +goose StatementEnd

CREATE TRIGGER infra_types_set_updated_at
    BEFORE UPDATE ON infra_types
    FOR EACH ROW
    EXECUTE FUNCTION set_infra_types_updated_at();

CREATE TABLE business_types (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    infra_type_id UUID NOT NULL UNIQUE REFERENCES infra_types (id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- +goose StatementBegin
CREATE FUNCTION set_business_types_updated_at()
RETURNS trigger
LANGUAGE plpgsql
AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$;
-- +goose StatementEnd

CREATE TRIGGER business_types_set_updated_at
    BEFORE UPDATE ON business_types
    FOR EACH ROW
    EXECUTE FUNCTION set_business_types_updated_at();

-- +goose Down
DROP TABLE business_types;
DROP FUNCTION set_business_types_updated_at();
DROP TABLE infra_types;
DROP FUNCTION set_infra_types_updated_at();
