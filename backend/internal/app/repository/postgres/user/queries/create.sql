-- name: CreateUser :one
INSERT INTO feedback.users (
    id,
    email,
    password_hash,
    roles,
    status,
    created_at,
    updated_at,
    deleted_at
) VALUES (
    $1, -- id
    $2, -- email
    $3, -- password_hash
    $4, -- roles
    $5, -- status
    $6, -- created_at
    $7, -- updated_at
    $8  -- deleted_at
)
RETURNING *;
