package llm

type Map = map[string]any

// AnalysisSchema creates a JSON schema for structured output from the LLM analysis.
// The schema defines the expected structure of the analysis response.
func AnalysisSchema() Map {
	return Map{
		"type": "object",
		"properties": Map{
			"overall_summary": Map{
				"type":        "string",
				"description": "Human-readable summary of all feedback in this analysis",
			},
			"sentiment": Map{
				"type":        "string",
				"enum":        []any{"positive", "mixed", "negative"},
				"description": "Overall sentiment analysis of all feedback",
			},
			"key_insights": Map{
				"type":        "array",
				"description": "Array of key insights/takeaways from the analysis",
				"items": Map{
					"type": "string",
				},
			},
			"topics": Map{
				"type":        "array",
				"description": "Array of topics/themes identified in the feedback",
				"items": Map{
					"type": "object",
					"properties": Map{
						"name": Map{
							"type":        "string",
							"description": "Name of the topic (e.g., 'App Crashes', 'UI Issues')",
						},
						"description": Map{
							"type":        "string",
							"description": "Detailed description of this topic",
						},
						"feedback_ids": Map{
							"type":        "array",
							"description": "Array of feedback IDs that belong to this topic",
							"items": Map{
								"type": "string",
							},
						},
						"sentiment": Map{
							"type":        "string",
							"enum":        []any{"positive", "mixed", "negative"},
							"description": "Sentiment for this specific topic",
						},
					},
					"required":             []any{"name", "description", "feedback_ids", "sentiment"},
					"additionalProperties": false,
				},
			},
		},
		"required":             []any{"overall_summary", "sentiment", "key_insights", "topics"},
		"additionalProperties": false,
	}
}
