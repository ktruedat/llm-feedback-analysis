-- +goose Up
-- +goose StatementBegin

-- Create sentiment enum type
CREATE TYPE feedback.sentiment AS ENUM ('positive', 'mixed', 'negative');

-- Create analysis status enum type
CREATE TYPE feedback.analysis_status AS ENUM ('processing', 'success', 'failed');

-- ============================================================================
-- TABLE 1: analyses
-- PURPOSE: Store snapshots of AI analysis at different points in time
-- WHY: We don't reanalyze everything from scratch, we build on previous analysis
-- ============================================================================
CREATE TABLE IF NOT EXISTS feedback.analyses
(
    id                    UUID PRIMARY KEY                  DEFAULT uuid_generate_v4() NOT NULL,
    previous_analysis_id  UUID REFERENCES feedback.analyses (id),

    -- Period covered in the analysis
    period_start          TIMESTAMP                NOT NULL,
    period_end            TIMESTAMP                NOT NULL,

    feedback_count        INTEGER                  NOT NULL,              -- Total feedbacks in this analysis
    new_feedback_count    INTEGER,                                        -- NEW feedbacks since last analysis

    -- Actual LLM generated content
    overall_summary       TEXT                     NOT NULL,              -- Human-readable summary of all feedback
    sentiment             feedback.sentiment       NOT NULL,              -- Overall sentiment (positive/mixed/negative)
    key_insights          TEXT[]                   NOT NULL DEFAULT '{}', -- Array of bullet points (key takeaways)

    -- Metadata
    model                 VARCHAR(50)              NOT NULL,              -- e.g., "gpt-5-mini"
    tokens                INTEGER                  NOT NULL,              -- Total tokens consumed
    analysis_duration_ms  INTEGER                  NOT NULL,              -- How long analysis took

    status                feedback.analysis_status NOT NULL DEFAULT 'processing',
    failure_reason        TEXT                     NULL,                  -- If failed, store failure reason

    created_at            TIMESTAMP                NOT NULL DEFAULT NOW(),
    completed_at          TIMESTAMP                NULL
);

COMMENT ON TABLE feedback.analyses IS 'Stores snapshots of AI analysis at different points in time';
COMMENT ON COLUMN feedback.analyses.id IS 'Unique identifier for the analysis';
COMMENT ON COLUMN feedback.analyses.period_start IS 'Start timestamp of the period covered by this analysis';
COMMENT ON COLUMN feedback.analyses.period_end IS 'End timestamp of the period covered by this analysis';
COMMENT ON COLUMN feedback.analyses.feedback_count IS 'Total number of feedbacks included in this analysis';
COMMENT ON COLUMN feedback.analyses.new_feedback_count IS 'Number of new feedbacks since the previous analysis';
COMMENT ON COLUMN feedback.analyses.previous_analysis_id IS 'Reference to the previous analysis (for incremental updates)';
COMMENT ON COLUMN feedback.analyses.overall_summary IS 'Human-readable summary of all feedback in this analysis';
COMMENT ON COLUMN feedback.analyses.sentiment IS 'Overall sentiment analysis (positive/mixed/negative)';
COMMENT ON COLUMN feedback.analyses.key_insights IS 'Array of key insights/takeaways from the analysis';
COMMENT ON COLUMN feedback.analyses.model IS 'LLM model used for this analysis (e.g., gpt-5-mini)';
COMMENT ON COLUMN feedback.analyses.tokens IS 'Total tokens consumed during analysis';
COMMENT ON COLUMN feedback.analyses.analysis_duration_ms IS 'Analysis duration in milliseconds';
COMMENT ON COLUMN feedback.analyses.status IS 'Analysis status (processing/success/failed)';
COMMENT ON COLUMN feedback.analyses.failure_reason IS 'Failure reason if analysis failed';
COMMENT ON COLUMN feedback.analyses.created_at IS 'Timestamp when the analysis was created';
COMMENT ON COLUMN feedback.analyses.completed_at IS 'Timestamp when the analysis was completed (NULL if not completed)';

CREATE INDEX IF NOT EXISTS feedback_analyses_previous_analysis_id_idx ON feedback.analyses (previous_analysis_id);
COMMENT ON INDEX feedback.feedback_analyses_previous_analysis_id_idx IS 'Index for traversing analysis chain';


-- ============================================================================
-- TABLE 2: analyzed_feedbacks
-- PURPOSE: Map which feedbacks were analyzed in which analysis (many-to-many)
-- WHY: Normalized relationship table for efficient querying of analyzed/unanalyzed feedbacks
-- ============================================================================
CREATE TABLE IF NOT EXISTS feedback.analyzed_feedbacks
(
    analysis_id UUID      NOT NULL REFERENCES feedback.analyses (id) ON DELETE CASCADE,
    feedback_id UUID      NOT NULL REFERENCES feedback.feedbacks (id) ON DELETE CASCADE,

    created_at  TIMESTAMP NOT NULL DEFAULT NOW(),

    PRIMARY KEY (analysis_id, feedback_id)
);

COMMENT ON TABLE feedback.analyzed_feedbacks IS 'Maps feedbacks to analyses (many-to-many relationship)';
COMMENT ON COLUMN feedback.analyzed_feedbacks.analysis_id IS 'Reference to the analysis';
COMMENT ON COLUMN feedback.analyzed_feedbacks.feedback_id IS 'Reference to the feedback that was analyzed';
COMMENT ON COLUMN feedback.analyzed_feedbacks.created_at IS 'Timestamp when the feedback was analyzed';

CREATE INDEX IF NOT EXISTS feedback_analyzed_feedbacks_analysis_id_idx ON feedback.analyzed_feedbacks (analysis_id);
COMMENT ON INDEX feedback.feedback_analyzed_feedbacks_analysis_id_idx IS 'Index for querying feedbacks by analysis';
CREATE INDEX IF NOT EXISTS feedback_analyzed_feedbacks_feedback_id_idx ON feedback.analyzed_feedbacks (feedback_id);
COMMENT ON INDEX feedback.feedback_analyzed_feedbacks_feedback_id_idx IS 'Index for querying analyses by feedback';


-- ============================================================================
-- TABLE 3: analysis_topics
-- PURPOSE: Store topics/themes identified by AI (e.g., "UI Issues", "Performance")
-- WHY: Topics are the main clustering mechanism - group similar feedback together
-- ============================================================================
CREATE TABLE IF NOT EXISTS feedback.analysis_topics
(
    id               UUID PRIMARY KEY      DEFAULT uuid_generate_v4() NOT NULL,
    analysis_id      UUID         NOT NULL REFERENCES feedback.analyses (id) ON DELETE CASCADE,

    topic_name       VARCHAR(100) NOT NULL, -- e.g., "App Crashes"
    description      TEXT          NOT NULL, -- AI explanation of this topic

    -- Metrics
    feedback_count   INTEGER      NOT NULL, -- How many feedbacks belong to this topic
    sentiment        feedback.sentiment NOT NULL, -- Sentiment for this topic (positive/mixed/negative)

    created_at       TIMESTAMP    NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMP    NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE feedback.analysis_topics IS 'Stores topics/themes identified by AI analysis';
COMMENT ON COLUMN feedback.analysis_topics.id IS 'Unique identifier for the topic';
COMMENT ON COLUMN feedback.analysis_topics.analysis_id IS 'Reference to the analysis this topic belongs to';
COMMENT ON COLUMN feedback.analysis_topics.topic_name IS 'Name of the topic (e.g., Mobile App Crashes)';
COMMENT ON COLUMN feedback.analysis_topics.description IS 'AI-generated explanation of this topic';
COMMENT ON COLUMN feedback.analysis_topics.feedback_count IS 'Number of feedbacks belonging to this topic';
COMMENT ON COLUMN feedback.analysis_topics.sentiment IS 'Sentiment for this topic (positive/mixed/negative)';

CREATE INDEX IF NOT EXISTS feedback_analysis_topics_analysis_id_idx ON feedback.analysis_topics (analysis_id);
COMMENT ON INDEX feedback.feedback_analysis_topics_analysis_id_idx IS 'Index for querying topics by analysis';


-- ============================================================================
-- TABLE 4: feedback_topic_assignments
-- PURPOSE: Map which feedbacks belong to which topics (many-to-many)
-- WHY: A single feedback might relate to multiple topics
-- ============================================================================
CREATE TABLE IF NOT EXISTS feedback.feedback_topic_assignments
(
    id          UUID PRIMARY KEY   DEFAULT uuid_generate_v4() NOT NULL,
    analysis_id UUID      NOT NULL REFERENCES feedback.analyses (id) ON DELETE CASCADE,
    feedback_id UUID      NOT NULL REFERENCES feedback.feedbacks (id) ON DELETE CASCADE,
    topic_id    UUID      NOT NULL REFERENCES feedback.analysis_topics (id) ON DELETE CASCADE,

    created_at  TIMESTAMP NOT NULL DEFAULT NOW(),

    UNIQUE (analysis_id, feedback_id, topic_id)
);

COMMENT ON TABLE feedback.feedback_topic_assignments IS 'Maps feedbacks to topics (many-to-many relationship)';
COMMENT ON COLUMN feedback.feedback_topic_assignments.id IS 'Unique identifier for the assignment';
COMMENT ON COLUMN feedback.feedback_topic_assignments.analysis_id IS 'Reference to the analysis this assignment belongs to';
COMMENT ON COLUMN feedback.feedback_topic_assignments.feedback_id IS 'Reference to the feedback being assigned';
COMMENT ON COLUMN feedback.feedback_topic_assignments.topic_id IS 'Reference to the topic being assigned to';

CREATE INDEX IF NOT EXISTS feedback_feedback_topic_assignments_analysis_id_idx ON feedback.feedback_topic_assignments (analysis_id);
COMMENT ON INDEX feedback.feedback_feedback_topic_assignments_analysis_id_idx IS 'Index for querying assignments by analysis';
CREATE INDEX IF NOT EXISTS feedback_feedback_topic_assignments_feedback_id_idx ON feedback.feedback_topic_assignments (feedback_id);
COMMENT ON INDEX feedback.feedback_feedback_topic_assignments_feedback_id_idx IS 'Index for querying assignments by feedback';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS feedback.feedback_topic_assignments CASCADE;
DROP TABLE IF EXISTS feedback.analysis_topics CASCADE;
DROP TABLE IF EXISTS feedback.analyzed_feedbacks CASCADE;
DROP TABLE IF EXISTS feedback.analyses CASCADE;

DROP TYPE IF EXISTS feedback.analysis_status;
DROP TYPE IF EXISTS feedback.sentiment;

-- +goose StatementEnd
