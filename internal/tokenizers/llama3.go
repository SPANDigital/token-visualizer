package tokenizers

import (
	"context"
	"fmt"
	"os"

	"github.com/sugarme/tokenizer"
)

// LLaMA3Tokenizer implements the Tokenizer interface using HuggingFace tokenizers
// for LLaMA 3.0, 3.1, 3.2, and 3.3 models.
type LLaMA3Tokenizer struct {
	tokenizer *tokenizer.Tokenizer
	modelName string
}

// NewLLaMA3Tokenizer creates a new LLaMA 3+ tokenizer from a tokenizer.json file.
// The tokenizerPath should point to a tokenizer.json file downloaded from HuggingFace.
func NewLLaMA3Tokenizer(tokenizerPath string) (*LLaMA3Tokenizer, error) {
	// Verify the file exists
	if _, err := os.Stat(tokenizerPath); err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("tokenizer file not found: %s", tokenizerPath)
		}
		return nil, fmt.Errorf("failed to access tokenizer file: %w", err)
	}

	// Load the tokenizer from the JSON file
	tk := tokenizer.NewTokenizerFromFile(tokenizerPath)
	if tk == nil {
		return nil, fmt.Errorf("failed to load LLaMA 3 tokenizer from %s", tokenizerPath)
	}

	return &LLaMA3Tokenizer{
		tokenizer: tk,
		modelName: "llama3",
	}, nil
}

// Name returns the human-readable name of this tokenizer.
func (t *LLaMA3Tokenizer) Name() string {
	return t.modelName
}

// Encode tokenizes the input text and returns a TokenizationResult.
func (t *LLaMA3Tokenizer) Encode(ctx context.Context, text string) (*TokenizationResult, error) {
	// Create input sequence
	input := tokenizer.NewInputSequence(text)
	encodeInput := tokenizer.NewSingleEncodeInput(input)

	// Encode the text
	encoding, err := t.tokenizer.Encode(encodeInput, false)
	if err != nil {
		return nil, fmt.Errorf("failed to encode text: %w", err)
	}

	// Get token IDs and strings
	tokenIDs := encoding.GetIds()
	tokenStrings := encoding.GetTokens()
	offsets := encoding.GetOffsets()

	// Convert to our Token structure
	tokens := make([]Token, len(tokenIDs))
	for i := range tokenIDs {
		tokens[i] = Token{
			Text:  tokenStrings[i],
			ID:    int(tokenIDs[i]),
			Start: int(offsets[i][0]),
			End:   int(offsets[i][1]),
		}
	}

	return &TokenizationResult{
		Tokens:     tokens,
		TotalCount: len(tokens),
		Text:       text,
		Model:      t.modelName,
	}, nil
}

// CountTokens returns just the count of tokens without full tokenization details.
func (t *LLaMA3Tokenizer) CountTokens(ctx context.Context, text string) (int, error) {
	// Create input sequence
	input := tokenizer.NewInputSequence(text)
	encodeInput := tokenizer.NewSingleEncodeInput(input)

	// Encode the text
	encoding, err := t.tokenizer.Encode(encodeInput, false)
	if err != nil {
		return 0, fmt.Errorf("failed to count tokens: %w", err)
	}

	return len(encoding.GetIds()), nil
}

// SupportsTokenIDs returns true since LLaMA 3+ tokenizer provides token IDs.
func (t *LLaMA3Tokenizer) SupportsTokenIDs() bool {
	return true
}

// SupportsDecoding returns true since LLaMA 3+ tokenizer supports decoding.
func (t *LLaMA3Tokenizer) SupportsDecoding() bool {
	return true
}
