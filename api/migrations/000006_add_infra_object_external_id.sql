-- +goose Up
ALTER TABLE infra_objects
    ADD COLUMN external_id TEXT;

UPDATE infra_objects
SET external_id = id::text;

ALTER TABLE infra_objects
    ALTER COLUMN external_id SET NOT NULL,
    ADD CONSTRAINT infra_objects_external_id_not_blank CHECK (btrim(external_id) <> ''),
    ADD CONSTRAINT infra_objects_external_id_key UNIQUE (external_id),
    ALTER COLUMN id SET DEFAULT gen_random_uuid();

-- +goose Down
ALTER TABLE infra_objects
    ALTER COLUMN id DROP DEFAULT,
    DROP CONSTRAINT infra_objects_external_id_key,
    DROP CONSTRAINT infra_objects_external_id_not_blank,
    DROP COLUMN external_id;
