-- name: CreateFeedback :one
INSERT INTO feedback.feedbacks (
    id,
    user_id,
    rating,
    comment,
    created_at,
    updated_at,
    deleted_at
) VALUES (
    $1, -- id
    $2, -- user_id
    $3, -- rating
    $4, -- comment
    $5, -- created_at
    $6, -- updated_at
    $7  -- deleted_at
)
RETURNING *;
