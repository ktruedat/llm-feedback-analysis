-- name: CreateFeedback :one
INSERT INTO feedback.feedbacks (
    id,
    rating,
    comment,
    created_at,
    updated_at,
    deleted_at
) VALUES (
    $1, -- id
    $2, -- rating
    $3, -- comment
    $4, -- created_at
    $5, -- updated_at
    $6  -- deleted_at
)
RETURNING *;
