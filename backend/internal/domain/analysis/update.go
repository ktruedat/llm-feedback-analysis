package analysis

import (
	"time"

	"github.com/ktruedat/llm-feedback-analysis/pkg/optional"
)

// UpdatableFields represents fields that can be updated in an analysis.
type UpdatableFields struct {
	Results       optional.Optional[*UpdatedResults]
	FailureReason optional.Optional[string]
	Status        Status
	CompletedAt   time.Time
}

type UpdatedResults struct {
	OverallSummary string
	Sentiment      Sentiment
	KeyInsights    []string
	Tokens         int
}
