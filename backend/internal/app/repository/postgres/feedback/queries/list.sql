-- name: ListFeedbacks :many
SELECT * FROM feedback.feedbacks
WHERE deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;
