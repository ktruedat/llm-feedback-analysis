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
				"description": "Array of topics/themes identified in the feedback. You MUST use one of the predefined topic enum values.",
				"items": Map{
					"type": "object",
					"properties": Map{
						"topic_enum": Map{
							"type": "string",
							"enum": []any{
								"product_functionality_features",
								"ui_ux",
								"performance_reliability",
								"usability_productivity",
								"security_privacy",
								"compatibility_integration",
								"developer_experience",
								"pricing_licensing",
								"customer_support_community",
								"installation_setup_deployment",
								"data_analytics_reporting",
								"localization_internationalization",
								"product_strategy_roadmap",
							},
							"description": "The predefined topic enum value that best categorizes this feedback",
						},
						"summary": Map{
							"type":        "string",
							"description": "Summary of the analysis for this topic - explaining why this feedback belongs to this topic and what specific aspects it addresses",
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
					"required":             []any{"topic_enum", "summary", "feedback_ids", "sentiment"},
					"additionalProperties": false,
				},
			},
		},
		"required":             []any{"overall_summary", "sentiment", "key_insights", "topics"},
		"additionalProperties": false,
	}
}
