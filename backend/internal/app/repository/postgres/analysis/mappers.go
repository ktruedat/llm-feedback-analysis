package analysis

import (
	"github.com/ktruedat/llm-feedback-analysis/internal/app/repository/postgres/analysis/sqlc"
	"github.com/ktruedat/llm-feedback-analysis/internal/domain/analysis"
)

// mapSQLCAnalysisToDomain maps a SQLC analysis model to a domain analysis entity.
func mapSQLCAnalysisToDomain(sqlcAnalysis sqlc.Analysis) *analysis.Analysis {
	builder := analysis.NewBuilder().
		WithID(sqlcAnalysis.ID).
		WithPeriod(sqlcAnalysis.PeriodStart, sqlcAnalysis.PeriodEnd).
		WithFeedbackCount(int(sqlcAnalysis.FeedbackCount)).
		WithOverallSummary(sqlcAnalysis.OverallSummary).
		WithSentiment(analysis.Sentiment(sqlcAnalysis.Sentiment)).
		WithKeyInsights(sqlcAnalysis.KeyInsights).
		WithModel(sqlcAnalysis.Model).
		WithTokens(int(sqlcAnalysis.Tokens)).
		WithAnalysisDurationMs(int(sqlcAnalysis.AnalysisDurationMs)).
		WithStatus(analysis.Status(sqlcAnalysis.Status)).
		WithCreatedAt(sqlcAnalysis.CreatedAt)

	// Handle optional fields (nullable fields use pointers)
	if sqlcAnalysis.PreviousAnalysisID != nil {
		builder.WithPreviousAnalysisID(*sqlcAnalysis.PreviousAnalysisID)
	}

	if sqlcAnalysis.NewFeedbackCount != nil {
		builder.WithNewFeedbackCount(int(*sqlcAnalysis.NewFeedbackCount))
	}

	if sqlcAnalysis.FailureReason != nil {
		builder.WithFailureReason(*sqlcAnalysis.FailureReason)
	}

	if sqlcAnalysis.CompletedAt != nil {
		builder.WithCompletedAt(*sqlcAnalysis.CompletedAt)
	}

	return builder.BuildUnchecked()
}

// mapSQLCTopicToDomain maps a SQLC topic model to a domain topic analysis entity.
func mapSQLCTopicToDomain(sqlcTopic sqlc.Topic) *analysis.TopicAnalysis {
	builder := analysis.NewTopicAnalysisBuilder().
		WithID(sqlcTopic.ID).
		WithAnalysisID(sqlcTopic.AnalysisID).
		WithTopic(analysis.Topic(sqlcTopic.TopicEnum)).
		WithSummary(sqlcTopic.Summary).
		WithFeedbackCount(int(sqlcTopic.FeedbackCount)).
		WithSentiment(analysis.Sentiment(sqlcTopic.Sentiment)).
		WithCreatedAt(sqlcTopic.CreatedAt).
		WithUpdatedAt(sqlcTopic.UpdatedAt)

	return builder.BuildUnchecked()
}
