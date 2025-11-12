# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Token Visualizer is a tool for visualizing and analyzing tokens across multiple LLM tokenizers including GPT-4, GPT-3.5, Claude, LLaMA, and others.

**Language:** Go 1.25 (REQUIRED - do not use older Go versions)

**Current Status:** Early development phase - repository structure is not yet established.

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

## Architecture Considerations

When implementing this project, consider:

1. **Multi-tokenizer support:** The architecture should support pluggable tokenizer implementations for different models (GPT-4, GPT-3.5, Claude, LLaMA, etc.)

2. **Tokenizer interface:** Create a common interface that all tokenizer implementations conform to, enabling easy addition of new tokenizers

3. **Visualization layer:** Separate tokenization logic from visualization/presentation logic

4. **Token libraries:** Research and integrate existing tokenizer libraries:
   - OpenAI's tiktoken (Go port: github.com/pkoukk/tiktoken-go or similar)
   - Anthropic's tokenizer (may need to interface with Python or use API)
   - HuggingFace tokenizers for LLaMA and other models

5. **Performance:** Token visualization may involve processing large texts - consider streaming or chunking for large inputs

## Git Workflow

- **Main branch:** `main`
- All commits should use conventional commit messages where appropriate
