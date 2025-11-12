package tokenizers

import (
	"context"
	"fmt"
	"os"

	"github.com/lwch/sentencepiece"
)

// LLaMATokenizer implements the Tokenizer interface using sentencepiece for LLaMA models
type LLaMATokenizer struct {
	model     *sentencepiece.Model
	modelPath string
}

// NewLLaMATokenizer creates a new LLaMA tokenizer
// modelPath should point to the tokenizer.model file
func NewLLaMATokenizer(modelPath string) (*LLaMATokenizer, error) {
	if modelPath == "" {
		return nil, fmt.Errorf("tokenizer model path is required")
	}

	if _, err := os.Stat(modelPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("tokenizer model file not found: %s", modelPath)
	}

	model, err := sentencepiece.Load(modelPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load sentencepiece model: %w", err)
	}

	return &LLaMATokenizer{
		model:     model,
		modelPath: modelPath,
	}, nil
}

// Name returns the name of this tokenizer
func (l *LLaMATokenizer) Name() string {
	return "LLaMA (sentencepiece)"
}

// Encode converts text into tokens
func (l *LLaMATokenizer) Encode(ctx context.Context, text string) (*TokenizationResult, error) {
	// Encode text to token IDs
	tokenIDs := l.model.Encode(text, true, true)

	// Build token list with details
	tokens := make([]Token, 0, len(tokenIDs))
	currentPos := 0

	for _, id := range tokenIDs {
		// Decode this single token to get its text
		tokenText := l.model.Decode([]uint64{id})

		token := Token{
			Text:  tokenText,
			ID:    int(id),
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
		Model:      "LLaMA",
	}, nil
}

// CountTokens returns just the token count
func (l *LLaMATokenizer) CountTokens(ctx context.Context, text string) (int, error) {
	tokenIDs := l.model.Encode(text, true, true)
	return len(tokenIDs), nil
}

// SupportsTokenIDs returns true
func (l *LLaMATokenizer) SupportsTokenIDs() bool {
	return true
}

// SupportsDecoding returns true
func (l *LLaMATokenizer) SupportsDecoding() bool {
	return true
}
