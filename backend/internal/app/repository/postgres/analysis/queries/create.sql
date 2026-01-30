-- name: CreateAnalysis :one
INSERT INTO feedback.analyses (
    id,
    previous_analysis_id,
    period_start,
    period_end,
    feedback_count,
    new_feedback_count,
    overall_summary,
    sentiment,
    key_insights,
    model,
    tokens,
    analysis_duration_ms,
    status,
    failure_reason,
    created_at,
    completed_at
) VALUES (
    $1,  -- id
    $2,  -- previous_analysis_id (nullable)
    $3,  -- period_start
    $4,  -- period_end
    $5,  -- feedback_count
    $6,  -- new_feedback_count (nullable)
    $7,  -- overall_summary
    $8,  -- sentiment
    $9,  -- key_insights
    $10, -- model
    $11, -- tokens
    $12, -- analysis_duration_ms
    $13, -- status
    $14, -- failure_reason (nullable)
    $15, -- created_at
    $16  -- completed_at (nullable)
)
RETURNING *;
