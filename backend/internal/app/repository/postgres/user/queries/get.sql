-- name: GetUserByID :one
SELECT * FROM feedback.users
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM feedback.users
WHERE email = $1;
