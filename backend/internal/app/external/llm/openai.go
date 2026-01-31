package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/ktruedat/llm-feedback-analysis/internal/app/external"
	"github.com/ktruedat/llm-feedback-analysis/internal/domain/analysis"
	"github.com/ktruedat/llm-feedback-analysis/internal/domain/feedback"
	"github.com/ktruedat/llm-feedback-analysis/pkg/tracelog"
)

// OpenAIClient implements the external.LLMClient interface using OpenAI's Responses API.
type OpenAIClient struct {
	apiKey string
	model  string
	logger tracelog.TraceLogger
}

// NewOpenAIClient creates a new OpenAI client.
func NewOpenAIClient(apiKey string, model string, logger tracelog.TraceLogger) *OpenAIClient {
	return &OpenAIClient{
		apiKey: apiKey,
		model:  model,
		logger: logger,
	}
}

// APIResponse represents the response structure from OpenAI Responses API.
type APIResponse struct {
	ID     string       `json:"id"`
	Status string       `json:"status"`
	Output []OutputItem `json:"output"`
	Usage  Usage        `json:"usage"`
	Error  *APIError    `json:"error,omitempty"`
}

// OutputItem represents an output item in the API response.
type OutputItem struct {
	Type    string        `json:"type"`
	ID      string        `json:"id"`
	Content []ContentItem `json:"content"`
}

// ContentItem represents content within an output item.
type ContentItem struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// Usage represents token usage information.
type Usage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
	TotalTokens  int `json:"total_tokens"`
}

// APIError represents an error from the API.
type APIError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
}

// AnalysisResponse represents the structured JSON response from the LLM.
type AnalysisResponse struct {
	OverallSummary string          `json:"overall_summary"`
	Sentiment      string          `json:"sentiment"`
	KeyInsights    []string        `json:"key_insights"`
	Topics         []TopicResponse `json:"topics"`
}

// TopicResponse represents a topic in the LLM response.
type TopicResponse struct {
	TopicEnum   string   `json:"topic_enum"`
	Summary     string   `json:"summary"`
	FeedbackIDs []string `json:"feedback_ids"`
	Sentiment   string   `json:"sentiment"`
}

// AnalyzeFeedbacks performs LLM analysis on the given feedbacks.
func (c *OpenAIClient) AnalyzeFeedbacks(
	ctx context.Context,
	feedbacks []*feedback.Feedback,
	previousAnalysis *analysis.Analysis,
) (*external.AnalysisResult, error) {
	// Build the user payload with feedback data
	userPayload := c.buildUserPayload(feedbacks, previousAnalysis)

	// Build the request body
	requestBody, err := c.buildRequestBody(userPayload)
	if err != nil {
		return nil, fmt.Errorf("failed to build request body: %w", err)
	}

	// Make the HTTP request
	httpReq, err := http.NewRequestWithContext(
		ctx,
		"POST",
		"https://api.openai.com/v1/responses",
		bytes.NewReader(requestBody),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			c.logger.RecordSpanError(ctx, fmt.Errorf("failed to close response body: %w", err))
		}
	}(resp.Body)

	rawBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("OpenAI API error (HTTP %d): %s", resp.StatusCode, string(rawBody))
	}

	// Parse the API response
	var apiResp APIResponse
	if err := json.Unmarshal(rawBody, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse API response: %w", err)
	}

	// Check for API-level errors
	if apiResp.Error != nil {
		return nil, fmt.Errorf("OpenAI API error: %s (type: %s)", apiResp.Error.Message, apiResp.Error.Type)
	}

	// Extract the output text from the response
	outputText, err := c.extractOutputText(apiResp)
	if err != nil {
		return nil, fmt.Errorf("failed to extract output text: %w", err)
	}

	// Parse the structured JSON response
	var analysisResp AnalysisResponse
	if err := json.Unmarshal([]byte(outputText), &analysisResp); err != nil {
		return nil, fmt.Errorf("model returned invalid JSON or schema mismatch: %w (raw: %s)", err, outputText)
	}

	c.logger.Debug("parsed analysis response", "topics_count", len(analysisResp.Topics))

	// Convert to external.AnalysisResult
	convertedTopics := c.convertTopics(ctx, analysisResp.Topics)
	c.logger.Debug("converted topics", "topics_count", len(convertedTopics))

	result := &external.AnalysisResult{
		OverallSummary: analysisResp.OverallSummary,
		Sentiment:      analysis.Sentiment(analysisResp.Sentiment),
		KeyInsights:    analysisResp.KeyInsights,
		TokensUsed:     apiResp.Usage.TotalTokens,
		Topics:         convertedTopics,
	}

	return result, nil
}

// buildUserPayload creates the user payload with feedback data.
func (c *OpenAIClient) buildUserPayload(feedbacks []*feedback.Feedback, previousAnalysis *analysis.Analysis) Map {
	feedbackItems := make([]Map, 0, len(feedbacks))
	for _, fb := range feedbacks {
		feedbackItems = append(
			feedbackItems, Map{
				"id":      fb.ID().String(),
				"rating":  fb.Rating().Value(),
				"comment": fb.Comment().Value(),
			},
		)
	}

	payload := Map{
		"feedbacks": feedbackItems,
	}

	// Include previous analysis summary if available
	if previousAnalysis != nil {
		payload["previous_analysis"] = Map{
			"overall_summary": previousAnalysis.OverallSummary(),
			"sentiment":       string(previousAnalysis.Sentiment()),
			"key_insights":    previousAnalysis.KeyInsights(),
		}
	}

	return payload
}

// buildRequestBody builds the request body for the OpenAI API.
func (c *OpenAIClient) buildRequestBody(userPayload Map) ([]byte, error) {
	userJSON, err := json.Marshal(userPayload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal user payload: %w", err)
	}

	systemPrompt := c.buildSystemPrompt()

	requestBody := Map{
		"model": c.model,
		"input": []Map{
			{
				"role":    "system",
				"content": systemPrompt,
			},
			{
				"role":    "user",
				"content": string(userJSON),
			},
		},
		"text": Map{
			"format": Map{
				"type":   "json_schema",
				"name":   "feedback_analysis",
				"strict": true,
				"schema": AnalysisSchema(),
			},
		},
	}

	return json.Marshal(requestBody)
}

// buildSystemPrompt creates the system prompt for the LLM.
func (c *OpenAIClient) buildSystemPrompt() string {
	// Build the topics list with descriptions for the prompt
	topicsList := ""
	for i, topic := range analysis.AllTopics() {
		if i > 0 {
			topicsList += "\n\n"
		}
		topicsList += fmt.Sprintf("%d. %s (%s)\n%s", i+1, topic.DisplayName(), string(topic), topic.Description())
	}

	return fmt.Sprintf(
		`Your task is to analyze customer feedback and categorize it into predefined business topics.

AVAILABLE TOPICS:
%s

INSTRUCTIONS:
1. Analyze all feedback and provide:
   - An overall summary of all feedback
   - The overall sentiment (positive, mixed, or negative)
   - Key insights as bullet points

2. Categorize feedbacks into topics:
   - You MUST use one of the predefined topic enum values listed above
   - A single feedback can belong to multiple topics if it addresses multiple themes
   - For each topic, provide:
     * The topic_enum value (one of the predefined values)
     * A summary explaining why this feedback belongs to this topic and what specific aspects it addresses
     * The feedback IDs that belong to this topic
     * The sentiment for this specific topic

3. Important rules:
   - DO NOT create new topic names - only use the predefined topic enum values
   - Group similar feedback together under the most appropriate topic(s)
   - Be specific about which feedback IDs map to which topics
   - Provide clear, actionable insights`, topicsList,
	)
}

// extractOutputText extracts the output text from the API response.
func (c *OpenAIClient) extractOutputText(apiResp APIResponse) (string, error) {
	for _, item := range apiResp.Output {
		if item.Type == "message" {
			for _, content := range item.Content {
				if content.Type == "output_text" {
					return content.Text, nil
				}
			}
		}
	}
	return "", errors.New("no output_text found in API response")
}

// convertTopics converts TopicResponse to external.Topic.
func (c *OpenAIClient) convertTopics(ctx context.Context, topics []TopicResponse) []external.Topic {
	if len(topics) == 0 {
		return nil
	}

	result := make([]external.Topic, 0, len(topics))
	for i, topic := range topics {
		feedbackIDs := make([]uuid.UUID, 0, len(topic.FeedbackIDs))
		for _, idStr := range topic.FeedbackIDs {
			id, err := uuid.Parse(idStr)
			if err != nil {
				// Log invalid UUID but continue - this shouldn't happen with proper schema validation
				c.logger.Warning(
					"failed to parse feedback ID for topic",
					"feedback_id",
					idStr,
					"topic",
					topic.TopicEnum,
					"error",
					err.Error(),
				)
				c.logger.RecordSpanError(
					ctx,
					fmt.Errorf("failed to parse feedback ID '%s' for topic '%s': %w", idStr, topic.TopicEnum, err),
				)
				continue
			}
			feedbackIDs = append(feedbackIDs, id)
		}

		// Only add topic if it has at least some valid feedback IDs or if it's a valid topic
		// (topics without feedback IDs might still be valid if the LLM didn't assign any)
		// Parse topic enum
		topicValue := analysis.Topic(topic.TopicEnum)
		if !topicValue.IsValid() {
			c.logger.Warning("invalid topic enum from LLM", "topic_enum", topic.TopicEnum, "index", i)
			c.logger.RecordSpanError(ctx, fmt.Errorf("invalid topic enum '%s' from LLM response", topic.TopicEnum))
			continue
		}

		result = append(
			result, external.Topic{
				Topic:       topicValue,
				Summary:     topic.Summary,
				FeedbackIDs: feedbackIDs,
				Sentiment:   analysis.Sentiment(topic.Sentiment),
			},
		)
		c.logger.Debug(
			"converted topic",
			"index",
			i,
			"topic_enum",
			topic.TopicEnum,
			"feedback_ids_count",
			len(feedbackIDs),
		)
	}
	return result
}
