-- +goose Up
ALTER TABLE tracked_locations
    ADD COLUMN name TEXT;

UPDATE tracked_locations
SET name = address;

ALTER TABLE tracked_locations
    ALTER COLUMN name SET NOT NULL,
    ADD CONSTRAINT tracked_locations_name_not_blank_check CHECK (btrim(name) <> '');

-- +goose Down
ALTER TABLE tracked_locations
    DROP COLUMN name;
