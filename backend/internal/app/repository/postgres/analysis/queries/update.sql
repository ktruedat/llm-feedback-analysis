-- name: UpdateAnalysis :exec
UPDATE feedback.analyses
SET
    overall_summary = $2,
    sentiment = $3,
    key_insights = $4,
    tokens = $5,
    status = $6,
    failure_reason = $7,
    completed_at = $8
WHERE id = $1;
