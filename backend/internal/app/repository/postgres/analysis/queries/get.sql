-- name: GetAnalysisByID :one
SELECT * FROM feedback.analyses
WHERE id = $1;

-- name: GetLatestAnalysis :one
SELECT * FROM feedback.analyses
ORDER BY created_at DESC
LIMIT 1;
