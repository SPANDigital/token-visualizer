# token-visualizer

A modern CLI tool for visualizing and analyzing tokens from various LLM tokenizers (GPT-4, GPT-3.5, GPT-5, Claude, LLaMA).

Built following Unix philosophy: reads from stdin, outputs to stdout, with colorized terminal output or markdown/HTML export.

## Features

- 🎨 **Colorized terminal output** with token boundaries and IDs
- 📊 **Multi-model comparison** side-by-side
- 📄 **Multiple output formats**: terminal, markdown, HTML
- 🔄 **Supports multiple tokenizers**:
  - OpenAI (GPT-4, GPT-3.5, GPT-5, GPT-5-mini, GPT-5-nano) via tiktoken
  - Anthropic Claude via API
  - Meta LLaMA via SentencePiece
- ⚡ **Fast** with local caching for API calls
- 🔌 **Unix-friendly**: pipe text in, get results out

## Installation

```bash
# Clone the repository
git clone https://github.com/spandigital/token-visualizer.git
cd token-visualizer

# Build
go build -o token-visualizer ./cmd/tokenizer

# Install (optional)
go install ./cmd/tokenizer
```

## Quick Start

```bash
# Basic usage - pipe text to visualize tokens
echo "Hello, world!" | ./token-visualizer

# Show token IDs and boundaries
echo "Hello, world!" | ./token-visualizer --show-ids --show-boundaries

# Compare tokenization across models
echo "Hello, world!" | ./token-visualizer compare --models gpt4,gpt3.5

# Export to markdown
echo "Hello, world!" | ./token-visualizer --format markdown > output.md

# Export to HTML
echo "Hello, world!" | ./token-visualizer --format html > output.html
```

## Commands

### `visualize` (default)

Visualize tokens with colorized output.

```bash
echo "Your text here" | ./token-visualizer visualize [flags]
```

**Flags:**
- `--model` - Model to use: `gpt4`, `gpt3.5`, `gpt5`, `gpt5-mini`, `gpt5-nano`, `claude`, `llama` (default: `gpt4`)
- `--format` - Output format: `terminal`, `markdown`, `html` (default: `terminal`)
- `--show-ids`, `-i` - Show token IDs
- `--show-boundaries`, `-b` - Show token boundaries
- `--encoding` - Tiktoken encoding for GPT-4/3.5 models (default: `cl100k_base`)
  - `cl100k_base` - GPT-4, GPT-3.5
  - `o200k_base` - GPT-4o (also used automatically for GPT-5 models)
  - `p50k_base` - Codex
  - `r50k_base` - GPT-3
  - **Note:** GPT-5 models automatically use `o200k_base` encoding regardless of this flag
- `--claude-model` - Claude model name (default: `claude-3-5-sonnet-20241022`)
- `--llama-model` - Path to LLaMA `tokenizer.model` file
- `--no-cache`, `-n` - Disable caching for Claude API

### `count`

Show only token counts (no visualization).

```bash
echo "Your text here" | ./token-visualizer count --models gpt4,gpt5,claude
```

### `compare`

Compare tokenization across multiple models side-by-side.

```bash
echo "Your text here" | ./token-visualizer compare --models gpt4,claude [flags]
```

## Examples

### GPT-4 with token IDs

```bash
echo "The quick brown fox" | ./token-visualizer --show-ids
```

### Compare GPT-4 and GPT-5 encodings

```bash
echo "Hello, world!" | ./token-visualizer compare \
  --models gpt4,gpt5 \
  --show-ids
```

### Compare GPT-5 variants

```bash
# All GPT-5 variants use the same o200k_base encoding
echo "The quick brown fox" | ./token-visualizer count \
  --models gpt5,gpt5-mini,gpt5-nano
```

### Use Claude tokenizer (requires API key)

```bash
export ANTHROPIC_API_KEY="your-api-key"
echo "Hello, world!" | ./token-visualizer --model claude
```

### LLaMA tokenizer

```bash
echo "Hello, world!" | ./token-visualizer \
  --model llama \
  --llama-model /path/to/tokenizer.model
```

### Export comparison to HTML

```bash
echo "The quick brown fox jumps over the lazy dog." | \
  ./token-visualizer compare \
  --models gpt4,gpt5 \
  --format html \
  --show-ids > comparison.html
```

## Configuration

### Environment Variables

- `ANTHROPIC_API_KEY` - Required for Claude tokenizer
- `TIKTOKEN_CACHE_DIR` - Cache directory for tiktoken encodings (default: `~/.cache/tiktoken`)

### Cache

Claude API responses are cached locally at `~/.cache/token-visualizer/` to speed up repeated queries and reduce API calls.

Use `--no-cache` to disable caching.

## Requirements

- Go 1.25+
- For Claude: `ANTHROPIC_API_KEY` environment variable
- For LLaMA: `tokenizer.model` file from LLaMA distribution

## Architecture

```
token-visualizer/
├── cmd/tokenizer/        # CLI entry point
├── internal/
│   ├── tokenizers/       # Tokenizer implementations
│   ├── output/           # Output renderers (terminal, markdown, HTML)
│   └── cache/            # Caching layer
└── go.mod
```

## Unix Philosophy

This tool follows Unix philosophy:

- **Do one thing well**: Tokenize and visualize text
- **Text streams**: Reads from stdin, writes to stdout
- **Composable**: Pipe with other tools
- **Silent on success**: Only outputs results

## Examples with Pipes

```bash
# From a file
cat article.txt | ./token-visualizer count --models gpt4,gpt5,claude

# From curl
curl -s https://example.com | ./token-visualizer --model gpt5

# Chain with other tools
cat large-file.txt | head -n 10 | ./token-visualizer --model gpt5 --show-ids

# Compare multiple files
for file in *.txt; do
  echo "=== $file ==="
  cat "$file" | ./token-visualizer count --models gpt4,gpt5,claude
done
```

## CI/CD

This project uses a modern CI/CD pipeline with:
- **GitHub Actions** for automated testing and deployment
- **GoReleaser** for building and publishing releases
- **ko** for building container images

See [CI_CD.md](CI_CD.md) for detailed documentation.

### Quick Release

```bash
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0
```

This automatically:
- Builds binaries for all platforms
- Creates container images (multi-arch)
- Publishes to GitHub Releases
- Updates Homebrew tap
- Generates changelog

## License

See [LICENSE](LICENSE) file.

## Contributing

Contributions welcome! Please ensure:
- Code follows Go 1.25 standards
- Uses Kong for CLI handling
- Maintains Unix philosophy principles
- Tests pass locally before submitting PR
