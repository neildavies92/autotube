DROP TRIGGER IF EXISTS workflow_runs_set_updated_at ON workflow_runs;
DROP TRIGGER IF EXISTS niches_set_updated_at ON niches;
DROP FUNCTION IF EXISTS set_updated_at();
DROP TABLE IF EXISTS workflow_events;
DROP TABLE IF EXISTS workflow_runs;
DROP TABLE IF EXISTS niches;
