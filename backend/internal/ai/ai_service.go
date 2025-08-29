package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/open-same/backend/internal/config"
	"github.com/open-same/backend/internal/models"
)

// AIService provides AI-powered content generation and assistance
type AIService struct {
	config     *config.Config
	openAI     *OpenAIClient
	anthropic  *AnthropicClient
	localLLM   *LocalLLMClient
	rateLimiter *RateLimiter
}

// ContentGenerationRequest represents a request for AI content generation
type ContentGenerationRequest struct {
	Type        string                 `json:"type"`        // document, code, diagram, etc.
	Prompt      string                 `json:"prompt"`
	Context     string                 `json:"context"`
	Style       string                 `json:"style"`
	Length      int                    `json:"length"`
	Language    string                 `json:"language"`
	Metadata    map[string]interface{} `json:"metadata"`
	UserID      string                 `json:"user_id"`
	CollaborationID string             `json:"collaboration_id"`
}

// ContentGenerationResponse represents AI-generated content
type ContentGenerationResponse struct {
	Content     string                 `json:"content"`
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	Metadata    map[string]interface{} `json:"metadata"`
	Model       string                 `json:"model"`
	Tokens      int                    `json:"tokens"`
	Cost        float64                `json:"cost"`
	Latency     time.Duration          `json:"latency"`
}

// AISuggestion represents AI-powered suggestions
type AISuggestion struct {
	Type        string `json:"type"`        // completion, improvement, correction
	Content     string `json:"content"`
	Confidence  float64 `json:"confidence"`
	Explanation string `json:"explanation"`
}

// AITemplate represents AI-generated templates
type AITemplate struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Content     string                 `json:"content"`
	Type        string                 `json:"type"`
	Category    string                 `json:"category"`
	Tags        []string               `json:"tags"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// NewAIService creates a new AI service instance
func NewAIService(cfg *config.Config) *AIService {
	service := &AIService{
		config: cfg,
		rateLimiter: NewRateLimiter(cfg.AI.RateLimit, cfg.AI.MaxConcurrentRequests),
	}

	// Initialize OpenAI client if configured
	if cfg.AI.OpenAI.APIKey != "" {
		service.openAI = NewOpenAIClient(cfg.AI.OpenAI)
	}

	// Initialize Anthropic client if configured
	if cfg.AI.Anthropic.APIKey != "" {
		service.anthropic = NewAnthropicClient(cfg.AI.Anthropic)
	}

	// Initialize local LLM client if enabled
	if cfg.AI.LocalLLM.Enabled {
		service.localLLM = NewLocalLLMClient(cfg.AI.LocalLLM)
	}

	return service
}

// GenerateContent generates AI-powered content based on the request
func (s *AIService) GenerateContent(ctx context.Context, req *ContentGenerationRequest) (*ContentGenerationResponse, error) {
	// Check rate limits
	if !s.rateLimiter.Allow() {
		return nil, fmt.Errorf("rate limit exceeded")
	}

	start := time.Now()
	
	// Select the best AI model based on request type and availability
	model, err := s.selectBestModel(req)
	if err != nil {
		return nil, fmt.Errorf("failed to select AI model: %w", err)
	}

	var response *ContentGenerationResponse

	// Generate content using the selected model
	switch model {
	case "openai":
		response, err = s.generateWithOpenAI(ctx, req)
	case "anthropic":
		response, err = s.generateWithAnthropic(ctx, req)
	case "local":
		response, err = s.generateWithLocalLLM(ctx, req)
	default:
		return nil, fmt.Errorf("unsupported AI model: %s", model)
	}

	if err != nil {
		// Try fallback model if enabled
		if s.config.AI.Fallback.Enabled {
			log.Printf("Primary AI model failed, trying fallback: %v", err)
			response, err = s.generateWithFallback(ctx, req)
			if err != nil {
				return nil, fmt.Errorf("both primary and fallback AI models failed: %w", err)
			}
		} else {
			return nil, fmt.Errorf("AI content generation failed: %w", err)
		}
	}

	// Calculate latency
	response.Latency = time.Since(start)

	// Log the generation for analytics
	s.logGeneration(req, response)

	return response, nil
}

// GenerateSuggestions generates AI-powered suggestions for existing content
func (s *AIService) GenerateSuggestions(ctx context.Context, content *models.Content, userID string) ([]*AISuggestion, error) {
	if !s.rateLimiter.Allow() {
		return nil, fmt.Errorf("rate limit exceeded")
	}

	suggestions := []*AISuggestion{}

	// Generate completion suggestions
	if completion, err := s.generateCompletionSuggestion(ctx, content, userID); err == nil {
		suggestions = append(suggestions, completion)
	}

	// Generate improvement suggestions
	if improvement, err := s.generateImprovementSuggestion(ctx, content, userID); err == nil {
		suggestions = append(suggestions, improvement)
	}

	// Generate correction suggestions
	if correction, err := s.generateCorrectionSuggestion(ctx, content, userID); err == nil {
		suggestions = append(suggestions, correction)
	}

	return suggestions, nil
}

// GenerateTemplate generates AI-powered templates
func (s *AIService) GenerateTemplate(ctx context.Context, templateType, category string, userID string) (*AITemplate, error) {
	if !s.rateLimiter.Allow() {
		return nil, fmt.Errorf("rate limit exceeded")
	}

	prompt := fmt.Sprintf("Generate a %s template for %s category. Include placeholders and examples.", templateType, category)

	req := &ContentGenerationRequest{
		Type:     templateType,
		Prompt:   prompt,
		Style:    "template",
		Length:   500,
		UserID:   userID,
		Metadata: map[string]interface{}{"category": category},
	}

	response, err := s.GenerateContent(ctx, req)
	if err != nil {
		return nil, err
	}

	template := &AITemplate{
		Name:        fmt.Sprintf("%s Template", strings.Title(templateType)),
		Description: fmt.Sprintf("AI-generated %s template for %s", templateType, category),
		Content:     response.Content,
		Type:        templateType,
		Category:    category,
		Tags:        []string{templateType, category, "ai-generated"},
		Metadata:    response.Metadata,
	}

	return template, nil
}

// selectBestModel selects the best AI model based on request type and availability
func (s *AIService) selectBestModel(req *ContentGenerationRequest) (string, error) {
	// Priority order: OpenAI > Anthropic > Local LLM
	
	if s.openAI != nil && s.openAI.IsAvailable() {
		return "openai", nil
	}
	
	if s.anthropic != nil && s.anthropic.IsAvailable() {
		return "anthropic", nil
	}
	
	if s.localLLM != nil && s.localLLM.IsAvailable() {
		return "local", nil
	}
	
	return "", fmt.Errorf("no AI models available")
}

// generateWithOpenAI generates content using OpenAI
func (s *AIService) generateWithOpenAI(ctx context.Context, req *ContentGenerationRequest) (*ContentGenerationResponse, error) {
	if s.openAI == nil {
		return nil, fmt.Errorf("OpenAI client not configured")
	}

	response, err := s.openAI.GenerateContent(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("OpenAI generation failed: %w", err)
	}

	return response, nil
}

// generateWithAnthropic generates content using Anthropic Claude
func (s *AIService) generateWithAnthropic(ctx context.Context, req *ContentGenerationRequest) (*ContentGenerationResponse, error) {
	if s.anthropic == nil {
		return nil, fmt.Errorf("Anthropic client not configured")
	}

	response, err := s.anthropic.GenerateContent(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("Anthropic generation failed: %w", err)
	}

	return response, nil
}

// generateWithLocalLLM generates content using local LLM
func (s *AIService) generateWithLocalLLM(ctx context.Context, req *ContentGenerationRequest) (*ContentGenerationResponse, error) {
	if s.localLLM == nil {
		return nil, fmt.Errorf("Local LLM client not configured")
	}

	response, err := s.localLLM.GenerateContent(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("Local LLM generation failed: %w", err)
	}

	return response, nil
}

// generateWithFallback generates content using the fallback model
func (s *AIService) generateWithFallback(ctx context.Context, req *ContentGenerationRequest) (*ContentGenerationResponse, error) {
	// Use OpenAI with fallback model if available
	if s.openAI != nil {
		req.Metadata["fallback_model"] = s.config.AI.Fallback.Model
		return s.generateWithOpenAI(ctx, req)
	}

	return nil, fmt.Errorf("no fallback model available")
}

// generateCompletionSuggestion generates completion suggestions
func (s *AIService) generateCompletionSuggestion(ctx context.Context, content *models.Content, userID string) (*AISuggestion, error) {
	prompt := fmt.Sprintf("Complete the following %s content naturally:\n\n%s", content.Type, content.Content)
	
	req := &ContentGenerationRequest{
		Type:     content.Type,
		Prompt:   prompt,
		Context:  content.Content,
		Length:   200,
		UserID:   userID,
		Metadata: map[string]interface{}{"suggestion_type": "completion"},
	}

	response, err := s.GenerateContent(ctx, req)
	if err != nil {
		return nil, err
	}

	return &AISuggestion{
		Type:        "completion",
		Content:     response.Content,
		Confidence:  0.85,
		Explanation: "AI-generated completion suggestion based on existing content",
	}, nil
}

// generateImprovementSuggestion generates improvement suggestions
func (s *AIService) generateImprovementSuggestion(ctx context.Context, content *models.Content, userID string) (*AISuggestion, error) {
	prompt := fmt.Sprintf("Suggest improvements for the following %s content:\n\n%s", content.Type, content.Content)
	
	req := &ContentGenerationRequest{
		Type:     content.Type,
		Prompt:   prompt,
		Context:  content.Content,
		Length:   300,
		UserID:   userID,
		Metadata: map[string]interface{}{"suggestion_type": "improvement"},
	}

	response, err := s.GenerateContent(ctx, req)
	if err != nil {
		return nil, err
	}

	return &AISuggestion{
		Type:        "improvement",
		Content:     response.Content,
		Confidence:  0.80,
		Explanation: "AI-generated improvement suggestions for better content quality",
	}, nil
}

// generateCorrectionSuggestion generates correction suggestions
func (s *AIService) generateCorrectionSuggestion(ctx context.Context, content *models.Content, userID string) (*AISuggestion, error) {
	prompt := fmt.Sprintf("Identify and correct any errors in the following %s content:\n\n%s", content.Type, content.Content)
	
	req := &ContentGenerationRequest{
		Type:     content.Type,
		Prompt:   prompt,
		Context:  content.Content,
		Length:   250,
		UserID:   userID,
		Metadata: map[string]interface{}{"suggestion_type": "correction"},
	}

	response, err := s.GenerateContent(ctx, req)
	if err != nil {
		return nil, err
	}

	return &AISuggestion{
		Type:        "correction",
		Content:     response.Content,
		Confidence:  0.90,
		Explanation: "AI-generated corrections for grammar, spelling, and factual errors",
	}, nil
}

// logGeneration logs the AI generation for analytics and monitoring
func (s *AIService) logGeneration(req *ContentGenerationRequest, response *ContentGenerationResponse) {
	logData := map[string]interface{}{
		"user_id":           req.UserID,
		"collaboration_id":  req.CollaborationID,
		"content_type":      req.Type,
		"model":             response.Model,
		"tokens":            response.Tokens,
		"cost":              response.Cost,
		"latency":           response.Latency.String(),
		"timestamp":         time.Now().UTC(),
	}

	logBytes, _ := json.Marshal(logData)
	log.Printf("AI Generation: %s", string(logBytes))
}

// GetAvailableModels returns available AI models
func (s *AIService) GetAvailableModels() []string {
	models := []string{}
	
	if s.openAI != nil && s.openAI.IsAvailable() {
		models = append(models, "openai")
	}
	
	if s.anthropic != nil && s.anthropic.IsAvailable() {
		models = append(models, "anthropic")
	}
	
	if s.localLLM != nil && s.localLLM.IsAvailable() {
		models = append(models, "local")
	}
	
	return models
}

// GetModelStatus returns the status of all AI models
func (s *AIService) GetModelStatus() map[string]interface{} {
	status := make(map[string]interface{})
	
	if s.openAI != nil {
		status["openai"] = map[string]interface{}{
			"available": s.openAI.IsAvailable(),
			"model":     s.config.AI.OpenAI.Model,
		}
	}
	
	if s.anthropic != nil {
		status["anthropic"] = map[string]interface{}{
			"available": s.anthropic.IsAvailable(),
			"model":     s.config.AI.Anthropic.Model,
		}
	}
	
	if s.localLLM != nil {
		status["local"] = map[string]interface{}{
			"available": s.localLLM.IsAvailable(),
			"model":     s.config.AI.LocalLLM.Model,
		}
	}
	
	return status
}