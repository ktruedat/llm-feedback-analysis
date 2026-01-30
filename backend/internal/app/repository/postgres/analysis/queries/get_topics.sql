-- name: GetTopicsByAnalysisID :many
SELECT * FROM feedback.analysis_topics
WHERE analysis_id = $1
ORDER BY created_at DESC;

-- name: GetFeedbackIDsByTopicID :many
SELECT feedback_id FROM feedback.feedback_topic_assignments
WHERE topic_id = $1;
