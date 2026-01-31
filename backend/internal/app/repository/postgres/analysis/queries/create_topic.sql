-- name: CreateTopicAnalysis :one
INSERT INTO feedback.analysis_topics (
    id,
    analysis_id,
    topic_enum,
    summary,
    feedback_count,
    sentiment,
    created_at,
    updated_at
) VALUES (
    $1,  -- id
    $2,  -- analysis_id
    $3,  -- topic_enum
    $4,  -- summary
    $5,  -- feedback_count
    $6,  -- sentiment
    $7,  -- created_at
    $8   -- updated_at
)
RETURNING *;
