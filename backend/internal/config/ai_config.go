package config

import (
	"time"
)

// AIConfig represents the AI services configuration
type AIConfig struct {
	OpenAI    OpenAIConfig    `json:"openai"`
	Anthropic AnthropicConfig `json:"anthropic"`
	LocalLLM  LocalLLMConfig  `json:"local_llm"`
	Fallback  FallbackConfig  `json:"fallback"`
	RateLimit float64         `json:"rate_limit"`
	MaxConcurrentRequests int `json:"max_concurrent_requests"`
}

// OpenAIConfig represents OpenAI API configuration
type OpenAIConfig struct {
	APIKey       string        `json:"api_key"`
	Model        string        `json:"model"`
	MaxTokens    int           `json:"max_tokens"`
	Temperature  float64       `json:"temperature"`
	Timeout      time.Duration `json:"timeout"`
	Organization string        `json:"organization"`
}

// AnthropicConfig represents Anthropic Claude API configuration
type AnthropicConfig struct {
	APIKey      string        `json:"api_key"`
	Model       string        `json:"model"`
	MaxTokens   int           `json:"max_tokens"`
	Timeout     time.Duration `json:"timeout"`
	Version     string        `json:"version"`
}

// LocalLLMConfig represents local LLM configuration
type LocalLLMConfig struct {
	Enabled bool   `json:"enabled"`
	URL     string `json:"url"`
	Model   string `json:"model"`
	Timeout time.Duration `json:"timeout"`
}

// FallbackConfig represents fallback AI model configuration
type FallbackConfig struct {
	Enabled bool   `json:"enabled"`
	Model   string `json:"model"`
}

// LoadAIConfig loads AI configuration from environment variables
func LoadAIConfig() *AIConfig {
	return &AIConfig{
		OpenAI: OpenAIConfig{
			APIKey:       getEnv("OPENAI_API_KEY", ""),
			Model:        getEnv("OPENAI_MODEL", "gpt-4"),
			MaxTokens:    getEnvAsInt("OPENAI_MAX_TOKENS", 4000),
			Temperature:  getEnvAsFloat("OPENAI_TEMPERATURE", 0.7),
			Timeout:      getEnvAsDuration("OPENAI_TIMEOUT", 30*time.Second),
			Organization: getEnv("OPENAI_ORGANIZATION", ""),
		},
		Anthropic: AnthropicConfig{
			APIKey:    getEnv("ANTHROPIC_API_KEY", ""),
			Model:     getEnv("ANTHROPIC_MODEL", "claude-3-sonnet-20240229"),
			MaxTokens: getEnvAsInt("ANTHROPIC_MAX_TOKENS", 4000),
			Timeout:   getEnvAsDuration("ANTHROPIC_TIMEOUT", 30*time.Second),
			Version:   getEnv("ANTHROPIC_VERSION", "2023-06-01"),
		},
		LocalLLM: LocalLLMConfig{
			Enabled: getEnv("LOCAL_LLM_ENABLED", "false") == "true",
			URL:     getEnv("LOCAL_LLM_URL", "http://localhost:11434"),
			Model:   getEnv("LOCAL_LLM_MODEL", "llama2:13b"),
			Timeout: getEnvAsDuration("LOCAL_LLM_TIMEOUT", 60*time.Second),
		},
		Fallback: FallbackConfig{
			Enabled: getEnv("AI_FALLBACK_ENABLED", "true") == "true",
			Model:   getEnv("AI_FALLBACK_MODEL", "gpt-3.5-turbo"),
		},
		RateLimit:             getEnvAsFloat("AI_RATE_LIMIT", 50.0),
		MaxConcurrentRequests: getEnvAsInt("AI_MAX_CONCURRENT_REQUESTS", 10),
	}
}