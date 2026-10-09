CREATE TABLE outbox (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    topic       TEXT NOT NULL,
    payload     JSONB NOT NULL,
    processed   BOOL NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_outbox_unprocessed ON outbox (processed, created_at) WHERE processed = FALSE;
