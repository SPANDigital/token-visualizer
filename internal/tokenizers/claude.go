package tokenizers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/spandigital/token-visualizer/internal/cache"
)

// ClaudeTokenizer implements the Tokenizer interface using Anthropic's Token Counting API
type ClaudeTokenizer struct {
	model  string
	apiKey string
	client *http.Client
	cache  *cache.Cache
}

type claudeTokenCountRequest struct {
	Model    string    `json:"model"`
	Messages []message `json:"messages"`
}

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type claudeTokenCountResponse struct {
	InputTokens int `json:"input_tokens"`
}

// NewClaudeTokenizer creates a new Claude tokenizer
// model should be one of: "claude-3-5-sonnet-20241022", "claude-3-5-haiku-20241022", etc.
func NewClaudeTokenizer(model string, useCache bool) (*ClaudeTokenizer, error) {
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("ANTHROPIC_API_KEY environment variable not set")
	}

	var c *cache.Cache
	var err error
	if useCache {
		c, err = cache.NewCache("")
		if err != nil {
			return nil, fmt.Errorf("failed to initialize cache: %w", err)
		}
	}

	return &ClaudeTokenizer{
		model:  model,
		apiKey: apiKey,
		client: &http.Client{},
		cache:  c,
	}, nil
}

// Name returns the name of this tokenizer
func (c *ClaudeTokenizer) Name() string {
	return fmt.Sprintf("Claude (%s)", c.model)
}

// Encode is limited for Claude - we can only get token count via API, not individual tokens
func (c *ClaudeTokenizer) Encode(ctx context.Context, text string) (*TokenizationResult, error) {
	count, err := c.CountTokens(ctx, text)
	if err != nil {
		return nil, err
	}

	// Claude API doesn't provide individual tokens, so we return a single "token" representing the text
	return &TokenizationResult{
		Tokens: []Token{
			{
				Text:  text,
				ID:    -1, // No token ID available
				Start: 0,
				End:   len(text),
			},
		},
		TotalCount: count,
		Text:       text,
		Model:      c.model,
	}, nil
}

// CountTokens returns the token count using Anthropic's API
func (c *ClaudeTokenizer) CountTokens(ctx context.Context, text string) (int, error) {
	// Check cache first
	cacheKey := fmt.Sprintf("claude:%s:%s", c.model, text)
	if c.cache != nil {
		var count int
		if err := c.cache.Get(cacheKey, &count); err == nil {
			return count, nil
		}
	}

	// Prepare request
	reqBody := claudeTokenCountRequest{
		Model: c.model,
		Messages: []message{
			{
				Role:    "user",
				Content: text,
			},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.anthropic.com/v1/messages/count_tokens", bytes.NewBuffer(jsonData))
	if err != nil {
		return 0, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("x-api-key", c.apiKey)
	req.Header.Set("content-type", "application/json")
	req.Header.Set("anthropic-version", "2023-06-01")
	req.Header.Set("anthropic-beta", "token-counting-2024-11-01")

	// Send request
	resp, err := c.client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("failed to send request: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return 0, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var response claudeTokenCountResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return 0, fmt.Errorf("failed to decode response: %w", err)
	}

	// Cache the result
	if c.cache != nil {
		_ = c.cache.Set(cacheKey, response.InputTokens)
	}

	return response.InputTokens, nil
}

// SupportsTokenIDs returns false (Claude API doesn't provide token IDs)
func (c *ClaudeTokenizer) SupportsTokenIDs() bool {
	return false
}

// SupportsDecoding returns false (Claude API doesn't provide token decoding)
func (c *ClaudeTokenizer) SupportsDecoding() bool {
	return false
}
