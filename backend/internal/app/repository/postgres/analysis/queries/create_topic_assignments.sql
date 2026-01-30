-- name: CreateTopicAssignment :exec
INSERT INTO feedback.feedback_topic_assignments (
    id,
    analysis_id,
    feedback_id,
    topic_id,
    created_at
) VALUES (
    $1,  -- id
    $2,  -- analysis_id
    $3,  -- feedback_id
    $4,  -- topic_id
    $5   -- created_at
)
ON CONFLICT (analysis_id, feedback_id, topic_id) DO NOTHING;
