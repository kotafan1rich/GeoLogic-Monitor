CREATE EXTENSION IF NOT EXISTS postgis;

CREATE TABLE IF NOT EXISTS competitors (
    tracked_location_id UUID NOT NULL,
    external_id VARCHAR(255) NOT NULL,
    competitor_name VARCHAR(255) NOT NULL,
    type_id VARCHAR(255) NOT NULL,
    competitor_address VARCHAR(255) NOT NULL,
    competitor_location GEOMETRY(Point, 4326) NOT NULL,
    opened_at TIMESTAMP NOT NULL,
    notified_at TIMESTAMP,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    PRIMARY KEY (tracked_location_id, external_id)
);

---- create above / drop below ----

DROP TABLE IF EXISTS competitors;