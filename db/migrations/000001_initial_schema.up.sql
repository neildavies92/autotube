CREATE TABLE niches (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    summary TEXT NOT NULL DEFAULT '',
    audience_regions TEXT[] NOT NULL DEFAULT ARRAY['US', 'GB']::TEXT[],
    seed_keywords TEXT[] NOT NULL DEFAULT ARRAY[]::TEXT[],
    status TEXT NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT niches_status_check CHECK (status IN ('active', 'paused', 'archived'))
);

CREATE TABLE workflow_runs (
    id BIGSERIAL PRIMARY KEY,
    niche_id BIGINT REFERENCES niches(id) ON DELETE SET NULL,
    workflow_type TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'queued',
    input JSONB NOT NULL DEFAULT '{}'::JSONB,
    output JSONB,
    error_message TEXT,
    started_at TIMESTAMPTZ,
    finished_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT workflow_runs_status_check CHECK (status IN ('queued', 'running', 'succeeded', 'failed', 'canceled')),
    CONSTRAINT workflow_runs_type_check CHECK (workflow_type <> '')
);

CREATE TABLE workflow_events (
    id BIGSERIAL PRIMARY KEY,
    workflow_run_id BIGINT NOT NULL REFERENCES workflow_runs(id) ON DELETE CASCADE,
    event_type TEXT NOT NULL,
    message TEXT NOT NULL DEFAULT '',
    payload JSONB NOT NULL DEFAULT '{}'::JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT workflow_events_type_check CHECK (event_type <> '')
);

CREATE INDEX niches_status_idx ON niches(status);
CREATE INDEX workflow_runs_niche_id_idx ON workflow_runs(niche_id);
CREATE INDEX workflow_runs_status_idx ON workflow_runs(status);
CREATE INDEX workflow_runs_type_created_at_idx ON workflow_runs(workflow_type, created_at DESC);
CREATE INDEX workflow_events_run_created_at_idx ON workflow_events(workflow_run_id, created_at);

CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER niches_set_updated_at
BEFORE UPDATE ON niches
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER workflow_runs_set_updated_at
BEFORE UPDATE ON workflow_runs
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();
