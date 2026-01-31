-- +goose Up
-- +goose StatementBegin

-- Create topic enum type
CREATE TYPE feedback.topic_enum AS ENUM (
    'product_functionality_features',
    'ui_ux',
    'performance_reliability',
    'usability_productivity',
    'security_privacy',
    'compatibility_integration',
    'developer_experience',
    'pricing_licensing',
    'customer_support_community',
    'installation_setup_deployment',
    'data_analytics_reporting',
    'localization_internationalization',
    'product_strategy_roadmap'
    );

COMMENT ON TYPE feedback.topic_enum IS 'Predefined business topics for categorizing feedback';

ALTER TABLE feedback.analysis_topics
    ADD COLUMN topic_enum feedback.topic_enum NOT NULL DEFAULT 'product_functionality_features';

ALTER TABLE feedback.analysis_topics
    DROP COLUMN topic_name;

-- Drop the old description column and add summary column
ALTER TABLE feedback.analysis_topics
    DROP COLUMN description,
    ADD COLUMN summary TEXT NOT NULL DEFAULT '';

-- Remove default from summary after adding it
ALTER TABLE feedback.analysis_topics
    ALTER COLUMN summary DROP DEFAULT;

-- Update comments
COMMENT ON COLUMN feedback.analysis_topics.topic_enum IS 'Predefined topic enum value';
COMMENT ON COLUMN feedback.analysis_topics.summary IS 'Summary of the analysis for this topic';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Add back topic_name and description columns
ALTER TABLE feedback.analysis_topics
    ADD COLUMN topic_name  VARCHAR(100),
    ADD COLUMN description TEXT;

-- Migrate enum values back to topic_name (using display names)
UPDATE feedback.analysis_topics
SET topic_name = CASE topic_enum::text
                     WHEN 'product_functionality_features' THEN 'Product Functionality & Features'
                     WHEN 'ui_ux' THEN 'UI / UX'
                     WHEN 'performance_reliability' THEN 'Performance & Reliability'
                     WHEN 'usability_productivity' THEN 'Usability & Productivity'
                     WHEN 'security_privacy' THEN 'Security & Privacy'
                     WHEN 'compatibility_integration' THEN 'Compatibility & Integration'
                     WHEN 'developer_experience' THEN 'Developer Experience'
                     WHEN 'pricing_licensing' THEN 'Pricing & Licensing'
                     WHEN 'customer_support_community' THEN 'Customer Support & Community'
                     WHEN 'installation_setup_deployment' THEN 'Installation, Setup & Deployment'
                     WHEN 'data_analytics_reporting' THEN 'Data & Analytics / Reporting'
                     WHEN 'localization_internationalization' THEN 'Localization & Internationalization'
                     WHEN 'product_strategy_roadmap' THEN 'Product Strategy & Roadmap'
                     ELSE 'Unknown'
    END;

-- Migrate summary back to description
UPDATE feedback.analysis_topics
SET description = summary;

-- Make topic_name and description NOT NULL
ALTER TABLE feedback.analysis_topics
    ALTER COLUMN topic_name SET NOT NULL,
    ALTER COLUMN description SET NOT NULL;

-- Drop topic_enum and summary columns
ALTER TABLE feedback.analysis_topics
    DROP COLUMN topic_enum,
    DROP COLUMN summary;

-- Drop the enum type
DROP TYPE IF EXISTS feedback.topic_enum;

-- +goose StatementEnd
