CREATE TABLE IF NOT EXISTS metrics (
    id          BIGSERIAL PRIMARY KEY,
    name        TEXT NOT NULL,
    value       DOUBLE PRECISION NOT NULL,
    type        TEXT NOT NULL,
    labels      JSONB,
    timestamp   BIGINT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_metrics_name_timestamp
    ON metrics (name, timestamp DESC);