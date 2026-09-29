CREATE INDEX IF NOT EXISTS expires_at_idx ON competitors (expires_at);

---- create above / drop below ----

DROP INDEX IF EXISTS expires_at_idx;
