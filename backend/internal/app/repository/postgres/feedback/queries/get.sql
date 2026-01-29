-- name: GetFeedback :one
SELECT * FROM feedback.feedbacks
WHERE id = $1;
