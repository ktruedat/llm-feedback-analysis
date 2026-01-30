package analysis

// Sentiment represents the sentiment of an analysis.
type Sentiment string

const (
	SentimentPositive Sentiment = "positive"
	SentimentMixed    Sentiment = "mixed"
	SentimentNegative Sentiment = "negative"
)

// String returns the string representation of the sentiment.
func (s Sentiment) String() string {
	return string(s)
}

// IsValid checks if the sentiment is valid.
func (s Sentiment) IsValid() bool {
	switch s {
	case SentimentPositive, SentimentMixed, SentimentNegative:
		return true
	default:
		return false
	}
}
