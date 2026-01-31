-- +goose Up
-- +goose StatementBegin

-- Enable UUID extension for generating UUIDs
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE SCHEMA IF NOT EXISTS feedback;
COMMENT ON SCHEMA feedback IS 'Schema for storing user feedback submissions';

-- Create feedback table
CREATE TABLE IF NOT EXISTS feedback.feedbacks
(
    id         UUID PRIMARY KEY                  DEFAULT uuid_generate_v4() NOT NULL,
    rating     INTEGER                          NOT NULL CHECK (rating >= 1 AND rating <= 5),
    comment    TEXT                             NOT NULL CHECK (LENGTH(comment) >= 1 AND LENGTH(comment) <= 1000),
    created_at TIMESTAMP                        NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP                        NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP                        NULL
);
COMMENT ON TABLE feedback.feedbacks IS 'Stores user feedback submissions with ratings and comments';
COMMENT ON COLUMN feedback.feedbacks.id IS 'Unique identifier for the feedback submission';
COMMENT ON COLUMN feedback.feedbacks.rating IS 'Rating value from 1 to 5 stars';
COMMENT ON COLUMN feedback.feedbacks.comment IS 'Free-text feedback comment (1-1000 characters)';
COMMENT ON COLUMN feedback.feedbacks.created_at IS 'Timestamp when the feedback was submitted';
COMMENT ON COLUMN feedback.feedbacks.updated_at IS 'Timestamp when the feedback was last updated';
COMMENT ON COLUMN feedback.feedbacks.deleted_at IS 'Timestamp when the feedback was soft-deleted (NULL if not deleted)';

-- Indexes for common query patterns
CREATE INDEX IF NOT EXISTS feedback_feedbacks_created_at_idx ON feedback.feedbacks (created_at DESC);
COMMENT ON INDEX feedback.feedback_feedbacks_created_at_idx IS 'Index for querying recent feedback submissions';

CREATE INDEX IF NOT EXISTS feedback_feedbacks_rating_idx ON feedback.feedbacks (rating);
COMMENT ON INDEX feedback.feedback_feedbacks_rating_idx IS 'Index for rating-based queries and statistics';

CREATE INDEX IF NOT EXISTS feedback_feedbacks_deleted_at_idx ON feedback.feedbacks (deleted_at) WHERE deleted_at IS NULL;
COMMENT ON INDEX feedback.feedback_feedbacks_deleted_at_idx IS 'Partial index for filtering active (non-deleted) feedback';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS feedback.feedbacks;
DROP SCHEMA IF EXISTS feedback;

-- +goose StatementEnd
