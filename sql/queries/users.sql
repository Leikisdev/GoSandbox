-- name: CreateUser :one
INSERT INTO users (created_at, updated_at, email, hashed_password)
VALUES (
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;

-- name: DeleteUsers :exec
DELETE FROM users;

-- name: LoginUser :one
SELECT * FROM users
WHERE email = $1;

-- name: UpdateUser :one
UPDATE users
SET email = $2, hashed_password = $3
WHERE id = $1
RETURNING *;