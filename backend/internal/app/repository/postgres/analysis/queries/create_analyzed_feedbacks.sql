-- name: CreateAnalyzedFeedback :exec
INSERT INTO feedback.analyzed_feedbacks (
    analysis_id,
    feedback_id,
    created_at
) VALUES (
    $1,  -- analysis_id
    $2,  -- feedback_id
    $3   -- created_at
)
ON CONFLICT (analysis_id, feedback_id) DO NOTHING;

-- name: GetFeedbackIDsByAnalysisID :many
SELECT feedback_id FROM feedback.analyzed_feedbacks
WHERE analysis_id = $1;
