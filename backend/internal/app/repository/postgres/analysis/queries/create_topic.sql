-- name: CreateTopic :one
INSERT INTO feedback.analysis_topics (
    id,
    analysis_id,
    topic_name,
    description,
    feedback_count,
    sentiment,
    created_at,
    updated_at
) VALUES (
    $1,  -- id
    $2,  -- analysis_id
    $3,  -- topic_name
    $4,  -- description
    $5,  -- feedback_count
    $6,  -- sentiment
    $7,  -- created_at
    $8   -- updated_at
)
RETURNING *;
