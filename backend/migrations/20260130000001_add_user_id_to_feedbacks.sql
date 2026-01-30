-- +goose Up
-- +goose StatementBegin

-- Add user_id column to feedbacks table
ALTER TABLE feedback.feedbacks
    ADD COLUMN IF NOT EXISTS user_id UUID NOT NULL REFERENCES feedback.users (id) ON DELETE CASCADE DEFAULT uuid_generate_v4();

-- Add index for user_id lookups
CREATE INDEX IF NOT EXISTS feedback_feedbacks_user_id_idx ON feedback.feedbacks (user_id);
COMMENT ON INDEX feedback.feedback_feedbacks_user_id_idx IS 'Index for querying feedbacks by user ID';

-- Add composite index for user_id and created_at (common query pattern: user's recent feedbacks)
CREATE INDEX IF NOT EXISTS feedback_feedbacks_user_id_created_at_idx ON feedback.feedbacks (user_id, created_at DESC);
COMMENT ON INDEX feedback.feedback_feedbacks_user_id_created_at_idx IS 'Composite index for querying user feedbacks sorted by creation date';

-- Update comment on table to reflect the relationship
COMMENT ON COLUMN feedback.feedbacks.user_id IS 'Reference to the user who submitted the feedback';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS feedback_feedbacks_user_id_created_at_idx;
DROP INDEX IF EXISTS feedback_feedbacks_user_id_idx;
ALTER TABLE feedback.feedbacks
    DROP COLUMN IF EXISTS user_id;

-- +goose StatementEnd
