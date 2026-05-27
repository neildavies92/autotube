-- name: CreateWorkflowRun :one
INSERT INTO workflow_runs (
    niche_id,
    workflow_type,
    status,
    input
) VALUES (
    $1,
    $2,
    $3,
    $4
)
RETURNING id, niche_id, workflow_type, status, input, output, error_message, started_at, finished_at, created_at, updated_at;

-- name: GetWorkflowRun :one
SELECT id, niche_id, workflow_type, status, input, output, error_message, started_at, finished_at, created_at, updated_at
FROM workflow_runs
WHERE id = $1;

-- name: ListWorkflowRunsForNiche :many
SELECT id, niche_id, workflow_type, status, input, output, error_message, started_at, finished_at, created_at, updated_at
FROM workflow_runs
WHERE niche_id = $1
ORDER BY created_at DESC;

-- name: UpdateWorkflowRunStatus :one
UPDATE workflow_runs
SET
    status = $2,
    output = $3,
    error_message = $4,
    started_at = COALESCE(started_at, $5),
    finished_at = $6
WHERE id = $1
RETURNING id, niche_id, workflow_type, status, input, output, error_message, started_at, finished_at, created_at, updated_at;

-- name: AddWorkflowEvent :one
INSERT INTO workflow_events (
    workflow_run_id,
    event_type,
    message,
    payload
) VALUES (
    $1,
    $2,
    $3,
    $4
)
RETURNING id, workflow_run_id, event_type, message, payload, created_at;

-- name: ListWorkflowEvents :many
SELECT id, workflow_run_id, event_type, message, payload, created_at
FROM workflow_events
WHERE workflow_run_id = $1
ORDER BY created_at, id;
