-- name: CreateChirp :one
INSERT INTO chirps (body, created_at, updated_at, user_id)
VALUES (
    $1,
    NOW(),
    NOW(),
    $2
)
RETURNING *;

-- name: GetChirps :many
SELECT * FROM chirps
ORDER BY created_at ASC;

-- name: GetChirpById :one
SELECT * FROM chirps
WHERE id = $1;

-- name: DeleteChirpById :exec
DELETE FROM chirps
WHERE id = $1;