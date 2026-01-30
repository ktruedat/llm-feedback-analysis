-- +goose Up
-- +goose StatementBegin

-- Create users table in feedback schema
CREATE TABLE IF NOT EXISTS feedback.users
(
    id         UUID PRIMARY KEY                  DEFAULT uuid_generate_v4() NOT NULL,
    email      VARCHAR(254)                     NOT NULL UNIQUE,
    password_hash TEXT                         NOT NULL,
    roles      TEXT[]                          NOT NULL DEFAULT ARRAY['user']::TEXT[],
    status     VARCHAR(20)                      NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'inactive', 'suspended')),
    created_at TIMESTAMP                        NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP                        NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP                        NULL
);
COMMENT ON TABLE feedback.users IS 'Stores user accounts for authentication and authorization';
COMMENT ON COLUMN feedback.users.id IS 'Unique identifier for the user';
COMMENT ON COLUMN feedback.users.email IS 'User email address (unique, normalized to lowercase)';
COMMENT ON COLUMN feedback.users.password_hash IS 'Hashed password (never store plain text)';
COMMENT ON COLUMN feedback.users.roles IS 'Array of user roles (e.g., ["user", "admin"])';
COMMENT ON COLUMN feedback.users.status IS 'User account status: active, inactive, or suspended';
COMMENT ON COLUMN feedback.users.created_at IS 'Timestamp when the user account was created';
COMMENT ON COLUMN feedback.users.updated_at IS 'Timestamp when the user account was last updated';
COMMENT ON COLUMN feedback.users.deleted_at IS 'Timestamp when the user account was soft-deleted (NULL if not deleted)';

-- Indexes for common query patterns
CREATE INDEX IF NOT EXISTS feedback_users_email_idx ON feedback.users (email);
COMMENT ON INDEX feedback.feedback_users_email_idx IS 'Index for email-based lookups (authentication)';

CREATE INDEX IF NOT EXISTS feedback_users_status_idx ON feedback.users (status) WHERE deleted_at IS NULL;
COMMENT ON INDEX feedback.feedback_users_status_idx IS 'Partial index for querying active users by status';

CREATE INDEX IF NOT EXISTS feedback_users_deleted_at_idx ON feedback.users (deleted_at) WHERE deleted_at IS NULL;
COMMENT ON INDEX feedback.feedback_users_deleted_at_idx IS 'Partial index for filtering active (non-deleted) users';

CREATE INDEX IF NOT EXISTS feedback_users_created_at_idx ON feedback.users (created_at DESC);
COMMENT ON INDEX feedback.feedback_users_created_at_idx IS 'Index for querying users by creation date';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS feedback.users;

-- +goose StatementEnd
