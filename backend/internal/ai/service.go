package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/open-same/backend/internal/config"
)

// AIService provides AI-powered content generation and assistance
type AIService struct {
	config config.AIConfig
	client *http.Client
}

// NewAIService creates a new AI service instance
func NewAIService(cfg config.AIConfig) *AIService {
	return &AIService{
		config: cfg,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GenerateContentRequest represents a content generation request
type GenerateContentRequest struct {
	Prompt     string                 `json:"prompt"`
	Type       string                 `json:"type"`
	Length     string                 `json:"length,omitempty"`
	Style      string                 `json:"style,omitempty"`
	Tone       string                 `json:"tone,omitempty"`
	Language   string                 `json:"language,omitempty"`
	Context    string                 `json:"context,omitempty"`
	Parameters map[string]interface{} `json:"parameters,omitempty"`
}

// GenerateContentResponse represents the AI-generated content response
type GenerateContentResponse struct {
	Content     string                 `json:"content"`
	Title       string                 `json:"title,omitempty"`
	Description string                 `json:"description,omitempty"`
	Tags        []string               `json:"tags,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	Model       string                 `json:"model"`
	Usage       *Usage                 `json:"usage,omitempty"`
	Error       string                 `json:"error,omitempty"`
}

// Usage represents token usage information
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// OpenAIRequest represents OpenAI API request
type OpenAIRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
	Temperature float64   `json:"temperature,omitempty"`
	TopP        float64   `json:"top_p,omitempty"`
	Stream      bool      `json:"stream,omitempty"`
}

// OpenAIResponse represents OpenAI API response
type OpenAIResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage Usage `json:"usage"`
}

// Message represents a chat message
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// AnthropicRequest represents Anthropic API request
type AnthropicRequest struct {
	Model       string    `json:"model"`
	MaxTokens   int       `json:"max_tokens"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature,omitempty"`
	TopP        float64   `json:"top_p,omitempty"`
	Stream      bool      `json:"stream,omitempty"`
}

// AnthropicResponse represents Anthropic API response
type AnthropicResponse struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Role    string `json:"role"`
	Content []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"content"`
	Model   string `json:"model"`
	Usage   Usage  `json:"usage"`
	StopSeq string `json:"stop_seq,omitempty"`
}

// GenerateContent generates content using AI
func (s *AIService) GenerateContent(ctx context.Context, req GenerateContentRequest) (*GenerateContentResponse, error) {
	// Try OpenAI first if configured
	if s.config.OpenAIKey != "" {
		response, err := s.generateWithOpenAI(ctx, req)
		if err == nil {
			return response, nil
		}
		// Log error but continue to try other providers
		fmt.Printf("OpenAI generation failed: %v\n", err)
	}

	// Try Anthropic if configured
	if s.config.AnthropicKey != "" {
		response, err := s.generateWithAnthropic(ctx, req)
		if err == nil {
			return response, nil
		}
		// Log error but continue
		fmt.Printf("Anthropic generation failed: %v\n", err)
	}

	// Return error if no providers available
	return nil, fmt.Errorf("no AI providers configured or available")
}

// generateWithOpenAI generates content using OpenAI API
func (s *AIService) generateWithOpenAI(ctx context.Context, req GenerateContentRequest) (*GenerateContentResponse, error) {
	// Build system prompt based on content type
	systemPrompt := s.buildSystemPrompt(req)

	// Build user prompt
	userPrompt := s.buildUserPrompt(req)

	// Create OpenAI request
	openAIReq := OpenAIRequest{
		Model:       s.config.OpenAIModel,
		MaxTokens:   s.config.MaxTokens,
		Temperature: s.config.Temperature,
		Messages: []Message{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userPrompt},
		},
	}

	// Marshal request
	reqBody, err := json.Marshal(openAIReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal OpenAI request: %v", err)
	}

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, "POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %v", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+s.config.OpenAIKey)

	// Make request
	resp, err := s.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to make OpenAI request: %v", err)
	}
	defer resp.Body.Close()

	// Read response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read OpenAI response: %v", err)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("OpenAI API error: %s - %s", resp.Status, string(respBody))
	}

	// Parse response
	var openAIResp OpenAIResponse
	if err := json.Unmarshal(respBody, &openAIResp); err != nil {
		return nil, fmt.Errorf("failed to parse OpenAI response: %v", err)
	}

	// Extract content
	if len(openAIResp.Choices) == 0 {
		return nil, fmt.Errorf("no content generated by OpenAI")
	}

	content := openAIResp.Choices[0].Message.Content

	// Build response
	response := &GenerateContentResponse{
		Content: content,
		Model:   s.config.OpenAIModel,
		Usage:   &openAIResp.Usage,
	}

	// Extract title and description if possible
	response.Title = s.extractTitle(content, req.Type)
	response.Description = s.extractDescription(content, req.Type)
	response.Tags = s.extractTags(content, req.Type)

	return response, nil
}

// generateWithAnthropic generates content using Anthropic API
func (s *AIService) generateWithAnthropic(ctx context.Context, req GenerateContentRequest) (*GenerateContentResponse, error) {
	// Build system prompt
	systemPrompt := s.buildSystemPrompt(req)

	// Build user prompt
	userPrompt := s.buildUserPrompt(req)

	// Create Anthropic request
	anthropicReq := AnthropicRequest{
		Model:       s.config.AnthropicModel,
		MaxTokens:   s.config.MaxTokens,
		Temperature: s.config.Temperature,
		Messages: []Message{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userPrompt},
		},
	}

	// Marshal request
	reqBody, err := json.Marshal(anthropicReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Anthropic request: %v", err)
	}

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, "POST", "https://api.anthropic.com/v1/messages", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %v", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-api-key", s.config.AnthropicKey)
	httpReq.Header.Set("anthropic-version", "2023-06-01")

	// Make request
	resp, err := s.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to make Anthropic request: %v", err)
	}
	defer resp.Body.Close()

	// Read response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read Anthropic response: %v", err)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Anthropic API error: %s - %s", resp.Status, string(respBody))
	}

	// Parse response
	var anthropicResp AnthropicResponse
	if err := json.Unmarshal(respBody, &anthropicResp); err != nil {
		return nil, fmt.Errorf("failed to parse Anthropic response: %v", err)
	}

	// Extract content
	if len(anthropicResp.Content) == 0 {
		return nil, fmt.Errorf("no content generated by Anthropic")
	}

	content := anthropicResp.Content[0].Text

	// Build response
	response := &GenerateContentResponse{
		Content: content,
		Model:   s.config.AnthropicModel,
		Usage:   &anthropicResp.Usage,
	}

	// Extract title and description if possible
	response.Title = s.extractTitle(content, req.Type)
	response.Description = s.extractDescription(content, req.Type)
	response.Tags = s.extractTags(content, req.Type)

	return response, nil
}

// buildSystemPrompt builds a system prompt based on content type and parameters
func (s *AIService) buildSystemPrompt(req GenerateContentRequest) string {
	basePrompt := "You are an expert content creator. Generate high-quality, engaging content based on the user's request."

	switch req.Type {
	case "text":
		basePrompt += " Focus on creating well-structured, informative text."
	case "code":
		basePrompt += " Generate clean, well-commented, and efficient code."
	case "diagram":
		basePrompt += " Provide detailed descriptions for creating diagrams or visual content."
	case "document":
		basePrompt += " Create professional, well-formatted documents."
	case "template":
		basePrompt += " Generate reusable templates that can be easily customized."
	}

	if req.Style != "" {
		basePrompt += fmt.Sprintf(" Use a %s style.", req.Style)
	}

	if req.Tone != "" {
		basePrompt += fmt.Sprintf(" Maintain a %s tone.", req.Tone)
	}

	if req.Language != "" && req.Language != "en" {
		basePrompt += fmt.Sprintf(" Write in %s.", req.Language)
	}

	return basePrompt
}

// buildUserPrompt builds a user prompt from the request
func (s *AIService) buildUserPrompt(req GenerateContentRequest) string {
	prompt := req.Prompt

	if req.Length != "" {
		prompt += fmt.Sprintf("\n\nLength: %s", req.Length)
	}

	if req.Context != "" {
		prompt += fmt.Sprintf("\n\nContext: %s", req.Context)
	}

	return prompt
}

// extractTitle extracts a title from generated content
func (s *AIService) extractTitle(content, contentType string) string {
	// Simple title extraction - in production, you might want more sophisticated logic
	if len(content) > 100 {
		// Take first line or first 50 characters as title
		for i, char := range content {
			if char == '\n' && i > 10 {
				return content[:i]
			}
		}
		return content[:50] + "..."
	}
	return content
}

// extractDescription extracts a description from generated content
func (s *AIService) extractDescription(content, contentType string) string {
	// Simple description extraction
	if len(content) > 200 {
		return content[:200] + "..."
	}
	return content
}

// extractTags extracts tags from generated content
func (s *AIService) extractTags(content, contentType string) []string {
	// Simple tag extraction - in production, you might want AI-powered tag generation
	tags := []string{contentType}
	
	// Add some basic tags based on content type
	switch contentType {
	case "code":
		tags = append(tags, "programming", "development")
	case "text":
		tags = append(tags, "writing", "content")
	case "diagram":
		tags = append(tags, "visual", "design")
	case "document":
		tags = append(tags, "professional", "business")
	}
	
	return tags
}

// ImproveContent improves existing content using AI
func (s *AIService) ImproveContent(ctx context.Context, content, improvementType string) (*GenerateContentResponse, error) {
	prompt := fmt.Sprintf("Improve the following content for %s:\n\n%s", improvementType, content)
	
	req := GenerateContentRequest{
		Prompt: prompt,
		Type:   "text",
		Context: "Content improvement",
	}
	
	return s.GenerateContent(ctx, req)
}

// SummarizeContent summarizes content using AI
func (s *AIService) SummarizeContent(ctx context.Context, content, summaryType string) (*GenerateContentResponse, error) {
	prompt := fmt.Sprintf("Create a %s summary of the following content:\n\n%s", summaryType, content)
	
	req := GenerateContentRequest{
		Prompt: prompt,
		Type:   "text",
		Context: "Content summarization",
	}
	
	return s.GenerateContent(ctx, req)
}