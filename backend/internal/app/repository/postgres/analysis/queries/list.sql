-- name: ListAnalyses :many
SELECT * FROM feedback.analyses
ORDER BY created_at DESC;
