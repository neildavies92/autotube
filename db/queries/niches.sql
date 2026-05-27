-- name: CreateNiche :one
INSERT INTO niches (
    name,
    summary,
    audience_regions,
    seed_keywords,
    status
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING id, name, summary, audience_regions, seed_keywords, status, created_at, updated_at;

-- name: GetNiche :one
SELECT id, name, summary, audience_regions, seed_keywords, status, created_at, updated_at
FROM niches
WHERE id = $1;

-- name: ListNiches :many
SELECT id, name, summary, audience_regions, seed_keywords, status, created_at, updated_at
FROM niches
ORDER BY name;

-- name: UpdateNicheStatus :one
UPDATE niches
SET status = $2
WHERE id = $1
RETURNING id, name, summary, audience_regions, seed_keywords, status, created_at, updated_at;
