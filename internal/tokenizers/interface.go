package tokenizers

import "context"

// Token represents a single token with its text, ID, and position information
type Token struct {
	Text  string // The decoded text of this token
	ID    int    // The numeric token ID
	Start int    // Start position in original text (bytes)
	End   int    // End position in original text (bytes)
}

// TokenizationResult contains the full result of tokenizing text
type TokenizationResult struct {
	Tokens     []Token // List of tokens
	TotalCount int     // Total number of tokens
	Text       string  // Original text
	Model      string  // Model/encoding used
}

// Tokenizer is the interface that all tokenizer implementations must satisfy
type Tokenizer interface {
	// Name returns the human-readable name of this tokenizer
	Name() string

	// Encode converts text into a list of tokens
	Encode(ctx context.Context, text string) (*TokenizationResult, error)

	// CountTokens returns just the count without full tokenization (optimization)
	CountTokens(ctx context.Context, text string) (int, error)

	// SupportsTokenIDs returns true if this tokenizer can provide token IDs
	SupportsTokenIDs() bool

	// SupportsDecoding returns true if this tokenizer can decode tokens back to text
	SupportsDecoding() bool
}
