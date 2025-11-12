# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Token Visualizer is a tool for visualizing and analyzing tokens across multiple LLM tokenizers including GPT-4, GPT-3.5, GPT-5 (including gpt5-mini and gpt5-nano), Claude, LLaMA, and others.

**Language:** Go 1.25 (REQUIRED - do not use older Go versions)

**Current Status:** Active development with working CLI, multi-tokenizer support, and CI/CD pipeline.

## Development Commands

This is a new Go project. Standard Go commands will be used once the codebase is established:

- **Build:** `go build ./...`
- **Test:** `go test ./...`
- **Test with coverage:** `go test -coverprofile=coverage.out ./...`
- **View coverage:** `go tool cover -html=coverage.out`
- **Run single test:** `go test -run TestName ./path/to/package`
- **Lint:** Use `golangci-lint run` (once configured)
- **Format code:** `go fmt ./...`
- **Vet code:** `go vet ./...`

## Go Version Requirement

**CRITICAL:** This project requires Go 1.25. Always ensure:
- `go.mod` specifies `go 1.25`
- Use Go 1.25-specific features and improvements where appropriate
- Do not write code compatible with older Go versions if it compromises using 1.25 features

## CLI Framework Requirement

**REQUIRED:** This project must use Kong (`github.com/alecthomas/kong`) for all CLI handling.

Kong is a command-line parser for Go that supports complex command-line structures with minimal developer effort by expressing CLIs as Go types.

Key Kong patterns to follow:
- Define CLI structure using Go structs with Kong tags
- Use `cmd:""` tag for commands
- Use `arg:""` tag for positional arguments
- Use `help:""` tag for documentation
- Implement `Run()` methods on command structs for command execution
- Use `kong.Parse()` to parse arguments and `ctx.Run()` to execute commands
- Leverage Kong's built-in help generation and validation

Example CLI structure:
```go
var CLI struct {
  Command struct {
    Flag string `help:"Description of flag."`
    Args []string `arg:"" name:"arg" help:"Positional arguments."`
  } `cmd:"" help:"Command description."`
}

func main() {
  ctx := kong.Parse(&CLI)
  err := ctx.Run()
  kong.FatalIfErrorf(err)
}
```

## Supported Models

The tool supports the following LLM tokenizers:

### OpenAI Models
- **GPT-4**: Uses `cl100k_base` encoding (default) or configurable via `--encoding` flag
- **GPT-3.5**: Uses `cl100k_base` encoding (default) or configurable via `--encoding` flag
- **GPT-5**: Uses `o200k_base` encoding (hardcoded, cannot be overridden)
- **GPT-5-mini**: Uses `o200k_base` encoding (hardcoded, cannot be overridden)
- **GPT-5-nano**: Uses `o200k_base` encoding (hardcoded, cannot be overridden)

**Note:** All GPT-5 variants use the same tokenization (`o200k_base`), so token counts and boundaries are identical across gpt5, gpt5-mini, and gpt5-nano.

### Claude Models
- **Claude**: Uses Anthropic's Token Counting API with local caching
- **Required format**: `claude:model-name`
  - Example: `claude:claude-3-5-sonnet-20241022`
  - Example: `claude:claude-3-5-haiku-20241022`
  - This allows comparing multiple Claude models side-by-side

**Available Claude Models (as of January 2025):**
- `claude-3-5-sonnet-20241022`
- `claude-3-5-haiku-20241022`
- `claude-3-opus-20240229`
- And other Anthropic models available via the Token Counting API

### LLaMA Models
- **LLaMA**: Uses SentencePiece tokenizer
- **Required format**: `llama:/path/to/tokenizer.model`
  - Example: `llama:/models/llama2/tokenizer.model`
  - Example: `llama:/models/llama3/tokenizer.model`
  - This allows comparing multiple LLaMA models with different tokenizer files

## Architecture Considerations

When implementing this project, consider:

1. **Multi-tokenizer support:** The architecture supports pluggable tokenizer implementations via the `Tokenizer` interface

2. **Tokenizer interface:** All tokenizers implement a common interface in `internal/tokenizers/interface.go`

3. **Visualization layer:** Tokenization logic is separated from presentation logic in `internal/output/`

4. **Token libraries:** Currently integrated:
   - OpenAI's tiktoken: `github.com/pkoukk/tiktoken-go`
   - Anthropic's Token Counting API: Direct HTTP calls with caching
   - SentencePiece for LLaMA: `github.com/lwch/sentencepiece`

5. **Performance:** Implements caching for Claude API calls to reduce latency and API costs

## Git Workflow

- **Main branch:** `main`
- All commits should use conventional commit messages where appropriate
