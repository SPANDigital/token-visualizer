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
echo "Hello, world!" | ./token-visualizer compare --models gpt4,gpt5

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
- `--model` - Model to use: `gpt4`, `gpt3.5`, `gpt5`, `gpt5-mini`, `gpt5-nano`, `claude:model-name`, `llama:path` (default: `gpt4`)
  - For Claude, use format: `claude:claude-3-5-sonnet-20241022`
  - For LLaMA, use format: `llama:/path/to/tokenizer.model`
- `--format` - Output format: `terminal`, `markdown`, `html` (default: `terminal`)
- `--show-ids`, `-i` - Show token IDs
- `--show-boundaries`, `-b` - Show token boundaries
- `--encoding` - Tiktoken encoding for GPT-4/3.5 models (default: `cl100k_base`)
  - `cl100k_base` - GPT-4, GPT-3.5
  - `o200k_base` - GPT-4o (also used automatically for GPT-5 models)
  - `p50k_base` - Codex
  - `r50k_base` - GPT-3
  - **Note:** GPT-5 models automatically use `o200k_base` encoding regardless of this flag
- `--no-cache`, `-n` - Disable caching for Claude API

### `count`

Show only token counts (no visualization).

```bash
echo "Your text here" | ./token-visualizer count --models gpt4,gpt5,claude:claude-3-5-sonnet-20241022
```

### `compare`

Compare tokenization across multiple models side-by-side.

```bash
echo "Your text here" | ./token-visualizer compare --models gpt4,claude:claude-3-5-sonnet-20241022 [flags]
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
echo "Hello, world!" | ./token-visualizer --model claude:claude-3-5-sonnet-20241022
```

### Compare different Claude models

```bash
export ANTHROPIC_API_KEY="your-api-key"
echo "Hello, world!" | ./token-visualizer compare \
  --models claude:claude-3-5-sonnet-20241022,claude:claude-3-5-haiku-20241022 \
  --show-ids
```

### LLaMA tokenizer

```bash
echo "Hello, world!" | ./token-visualizer \
  --model llama:/path/to/tokenizer.model
```

**Note:** See [Obtaining LLaMA Tokenizer Files](#obtaining-llama-tokenizer-files) below for instructions on how to get the `tokenizer.model` file.

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
# From a file with multiple Claude models
cat article.txt | ./token-visualizer count \
  --models gpt4,gpt5,claude:claude-3-5-sonnet-20241022,claude:claude-3-5-haiku-20241022

# From curl
curl -s https://example.com | ./token-visualizer --model gpt5

# Chain with other tools
cat large-file.txt | head -n 10 | ./token-visualizer --model gpt5 --show-ids

# Compare multiple files across different Claude models
for file in *.txt; do
  echo "=== $file ==="
  cat "$file" | ./token-visualizer count \
    --models gpt4,gpt5,claude:claude-3-5-sonnet-20241022,claude:claude-3-5-haiku-20241022
done
```

## Obtaining LLaMA Tokenizer Files

The token-visualizer tool requires a `tokenizer.model` file to work with LLaMA models. Here's how to obtain it:

### Quick Start (Recommended Method)

1. **Request Access**
   - Visit https://huggingface.co/meta-llama/Llama-2-7b
   - Click "Access repository" and accept Meta's LLaMA Community License
   - Wait for approval (usually processed within an hour)

2. **Install HuggingFace CLI**
   ```bash
   pip install huggingface-hub
   ```

3. **Login to HuggingFace**
   ```bash
   huggingface-cli login
   # Enter your HuggingFace token when prompted
   ```

4. **Download Tokenizer**
   ```bash
   huggingface-cli download meta-llama/Llama-2-7b \
     tokenizer.model \
     --local-dir ./llama-tokenizer
   ```

5. **Use with token-visualizer**
   ```bash
   echo "Hello, world!" | ./token-visualizer \
     --model llama:./llama-tokenizer/tokenizer.model
   ```

### Supported LLaMA Versions

| Version | Tokenizer Type | Vocab Size | Supported | Notes |
|---------|---------------|------------|-----------|-------|
| LLaMA 1 | SentencePiece | 32K | ✅ Yes | Original LLaMA |
| LLaMA 2 | SentencePiece | 32K | ✅ Yes | Recommended |
| Code Llama | SentencePiece | 32K | ✅ Yes | Same as LLaMA 2 |
| LLaMA 3+ | TikToken | 128K | ❌ No | Different tokenizer format |

**Important:** All model sizes within the same version (7B, 13B, 70B) share the same tokenizer file.

### Alternative Method: Official Meta Download

1. Visit https://www.llama.com/llama-downloads/
2. Fill out the request form
3. Wait for the signed download URL via email
4. Clone and run the official download script:
   ```bash
   git clone https://github.com/meta-llama/llama.git
   cd llama
   ./download.sh
   # Paste the signed URL when prompted
   ```
5. The `tokenizer.model` file will be in the downloaded directory

### Direct Download (Advanced)

For users with HuggingFace access:

```bash
wget https://huggingface.co/meta-llama/Llama-2-7b/resolve/main/tokenizer.model
```

### Verifying Your Tokenizer File

Check the file size to ensure it downloaded correctly:

```bash
ls -lh tokenizer.model
# Should be approximately 500KB (499,723 bytes) for LLaMA 2
```

Test with token-visualizer:

```bash
echo "Test tokenization" | ./token-visualizer \
  --model llama:./tokenizer.model \
  --show-ids

# Should display colorized tokens with IDs
```

### License Information

LLaMA tokenizer files are subject to Meta's LLaMA Community License:

- ✅ **Allowed:** Research, development, commercial use (<700M monthly active users)
- ❌ **Not Allowed:** Redistribution without permission, use in apps with >700M MAU without special approval
- **Full License:** https://www.llama.com/llama-downloads/

**Note:** We cannot distribute `tokenizer.model` files with this tool. Users must obtain them directly from Meta or authorized sources.

### Troubleshooting

**Error: "tokenizer model file not found"**
- Use absolute path: `llama:/Users/you/path/to/tokenizer.model`
- Verify file exists: `ls -l /path/to/tokenizer.model`

**Error: "Internal: could not parse ModelProto"**
- You may be using a LLaMA 3 tokenizer (not supported)
- Download LLaMA 2 tokenizer instead

**HuggingFace access denied:**
- Ensure you've accepted the license on the HuggingFace model page
- Wait for approval email (usually takes <1 hour)
- Try logging out and back in: `huggingface-cli logout && huggingface-cli login`

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
