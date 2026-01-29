-- name: DeleteFeedback :execrows
UPDATE feedback.feedbacks
SET deleted_at = NOW(),
    updated_at = NOW()
WHERE id = $1
  AND deleted_at IS NULL;
