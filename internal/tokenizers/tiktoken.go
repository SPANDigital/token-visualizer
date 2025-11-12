package tokenizers

import (
	"context"
	"fmt"

	"github.com/pkoukk/tiktoken-go"
)

// TikTokenizer implements the Tokenizer interface using tiktoken-go for OpenAI models
type TikTokenizer struct {
	encoding string
	encoder  *tiktoken.Tiktoken
}

// NewTikTokenizer creates a new tiktoken-based tokenizer
// encoding should be one of: "cl100k_base" (GPT-4, GPT-3.5), "o200k_base" (GPT-4o), "p50k_base" (Codex), "r50k_base" (GPT-3)
func NewTikTokenizer(encoding string) (*TikTokenizer, error) {
	enc, err := tiktoken.GetEncoding(encoding)
	if err != nil {
		return nil, fmt.Errorf("failed to get tiktoken encoding %s: %w", encoding, err)
	}

	return &TikTokenizer{
		encoding: encoding,
		encoder:  enc,
	}, nil
}

// Name returns the name of this tokenizer
func (t *TikTokenizer) Name() string {
	return fmt.Sprintf("OpenAI (%s)", t.encoding)
}

// Encode converts text into tokens
func (t *TikTokenizer) Encode(ctx context.Context, text string) (*TokenizationResult, error) {
	// Encode the text to get token IDs
	tokenIDs := t.encoder.Encode(text, nil, nil)

	// Build token list with details
	tokens := make([]Token, 0, len(tokenIDs))
	currentPos := 0

	for _, id := range tokenIDs {
		// Decode this single token to get its text
		tokenText := t.encoder.Decode([]int{id})

		token := Token{
			Text:  tokenText,
			ID:    id,
			Start: currentPos,
			End:   currentPos + len(tokenText),
		}
		tokens = append(tokens, token)
		currentPos += len(tokenText)
	}

	return &TokenizationResult{
		Tokens:     tokens,
		TotalCount: len(tokenIDs),
		Text:       text,
		Model:      t.encoding,
	}, nil
}

// CountTokens returns just the token count (optimization)
func (t *TikTokenizer) CountTokens(ctx context.Context, text string) (int, error) {
	tokenIDs := t.encoder.Encode(text, nil, nil)
	return len(tokenIDs), nil
}

// SupportsTokenIDs returns true
func (t *TikTokenizer) SupportsTokenIDs() bool {
	return true
}

// SupportsDecoding returns true
func (t *TikTokenizer) SupportsDecoding() bool {
	return true
}
