# token-visualizer

A modern CLI tool for visualizing and analyzing tokens from various LLM tokenizers (GPT-4, GPT-3.5, GPT-5, Claude, LLaMA).

Built following Unix philosophy: reads from stdin, outputs to stdout, with colorized terminal output or markdown/HTML export.

## Quick Example

```bash
$ echo "Hello, world! 你好世界 🌍" | token-visualizer count --models gpt4,gpt5
╭─────────────────╮
│ 📊 Token Counts │
╰─────────────────╯

cl100k_base: 14 tokens
o200k_base: 10 tokens
```

*GPT-5's new encoding is 40% more efficient for Unicode! ✨*

## Features

- 🎨 **Colorized terminal output** with token boundaries and IDs
- 📊 **Multi-model comparison** side-by-side
- 📄 **Multiple output formats**: terminal, markdown, HTML
- 🔄 **Supports multiple tokenizers**:
  - OpenAI (GPT-4, GPT-3.5, GPT-5, GPT-5-mini, GPT-5-nano) via tiktoken
  - Anthropic Claude via API
  - Meta LLaMA 1/2 via SentencePiece
  - Meta LLaMA 3+ via HuggingFace Tokenizers
- ⚡ **Fast** with local caching for API calls
- 🔌 **Unix-friendly**: pipe text in, get results out

## Installation

### Homebrew (Recommended)

**macOS and Linux:**
```bash
brew install SPANDigital/tap/token-visualizer
```

### Container Image

```bash
# Pull the image
docker pull ghcr.io/spandigital/token-visualizer:latest

# Run with stdin
echo "Your text" | docker run -i ghcr.io/spandigital/token-visualizer:latest --show-ids
```

### Download Binary

Download pre-built binaries from [GitHub Releases](https://github.com/SPANDigital/token-visualizer/releases/latest) for:
- Linux (amd64, arm64)
- macOS (Intel, Apple Silicon)
- Windows (amd64)

### Build from Source

```bash
# Clone the repository
git clone https://github.com/spandigital/token-visualizer.git
cd token-visualizer

# Build
go build -o token-visualizer ./cmd/tokenizer

# Install (optional)
go install ./cmd/tokenizer
```

**Requirements:** Go 1.25+

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
- `--model` - Model to use: `gpt4`, `gpt3.5`, `gpt5`, `gpt5-mini`, `gpt5-nano`, `claude:model-name`, `llama:path`, `llama3:path` (default: `gpt4`)
  - For Claude, use format: `claude:claude-3-5-sonnet-20241022`
  - For LLaMA 1/2, use format: `llama:/path/to/tokenizer.model`
  - For LLaMA 3+, use format: `llama3:/path/to/tokenizer.json`
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

### Basic Visualization with Token IDs

```bash
echo "The quick brown fox jumps over the lazy dog." | ./token-visualizer --show-ids --show-boundaries
```

**Output:**
```
╭────────────────╮
│ 🔤 cl100k_base │
╰────────────────╯

Total tokens: 10

The|[791] quick|[4062] brown|[14198] fox|[39935] jumps|[35308] over|[927] the|[279] lazy|[16053] dog|[5679].[13]
```

### Compare GPT-4 and GPT-5 Encodings

```bash
echo "The quick brown fox jumps over the lazy dog." | ./token-visualizer compare \
  --models gpt4,gpt5 \
  --show-ids
```

**Output:**
```
╭──────────────────────────────────────────────────╮╭──────────────────────────────────────────────────╮
│                                                  ││                                                  │
│ cl100k_base                                      ││ o200k_base                                       │
│                                                  ││                                                  │
│ Tokens: 10                                       ││ Tokens: 10                                       │
│                                                  ││                                                  │
│ The(791)                                         ││ The(976)                                         │
│  quick(4062)                                     ││  quick(4853)                                     │
│  brown(14198)                                    ││  brown(19705)                                    │
│  fox(39935)                                      ││  fox(68347)                                      │
│  jumps(35308)                                    ││  jumps(65613)                                    │
│  over(927)                                       ││  over(1072)                                      │
│  the(279)                                        ││  the(290)                                        │
│  lazy(16053)                                     ││  lazy(29082)                                     │
│  dog(5679)                                       ││  dog(6446)                                       │
│ .(13)                                            ││ .(13)                                            │
╰──────────────────────────────────────────────────╯╰──────────────────────────────────────────────────╯
```

### Token Count Comparison with Unicode

```bash
echo "Hello, world! 你好世界 🌍" | ./token-visualizer count --models gpt4,gpt5
```

**Output:**
```
╭─────────────────╮
│ 📊 Token Counts │
╰─────────────────╯

cl100k_base: 14 tokens
o200k_base: 10 tokens
```

*Note: GPT-5's o200k_base encoding is more efficient for Unicode text!*

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

### LLaMA 1/2 tokenizer

```bash
echo "Hello, world!" | ./token-visualizer \
  --model llama:/path/to/tokenizer.model
```

**Note:** See [Obtaining LLaMA Tokenizer Files](#obtaining-llama-tokenizer-files) below for instructions on how to get the `tokenizer.model` file.

### LLaMA 3+ tokenizer

```bash
echo "Hello, world!" | ./token-visualizer \
  --model llama3:/path/to/tokenizer.json
```

**Note:** See [Obtaining LLaMA 3+ Tokenizer Files](#obtaining-llama-3-tokenizer-files) below for instructions on how to get the `tokenizer.json` file.

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
| LLaMA 1 | SentencePiece | 32K | ✅ Yes | Original LLaMA, use `llama:` |
| LLaMA 2 | SentencePiece | 32K | ✅ Yes | Recommended, use `llama:` |
| Code Llama | SentencePiece | 32K | ✅ Yes | Same as LLaMA 2, use `llama:` |
| LLaMA 3+ | TikToken-based | 128K | ✅ Yes | Use `llama3:` with tokenizer.json |

**Important:** All model sizes within the same version (7B, 13B, 70B, etc.) share the same tokenizer file.

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

## Obtaining LLaMA 3+ Tokenizer Files

The token-visualizer tool requires a `tokenizer.json` file to work with LLaMA 3+ models (LLaMA 3.0, 3.1, 3.2, 3.3). Here's how to obtain it:

### Quick Start (Recommended Method)

1. **Request Access**
   - Visit https://huggingface.co/meta-llama/Meta-Llama-3-8B-Instruct
   - Click "Access repository" and accept Meta's license agreement
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
   # For LLaMA 3/3.1
   huggingface-cli download meta-llama/Meta-Llama-3-8B-Instruct \
     tokenizer.json \
     --local-dir ./llama3-tokenizer

   # For LLaMA 3.2
   huggingface-cli download meta-llama/Llama-3.2-3B-Instruct \
     tokenizer.json \
     --local-dir ./llama3-tokenizer

   # For LLaMA 3.3
   huggingface-cli download meta-llama/Llama-3.3-70B-Instruct \
     tokenizer.json \
     --local-dir ./llama3-tokenizer
   ```

5. **Use with token-visualizer**
   ```bash
   echo "Hello, world!" | ./token-visualizer \
     --model llama3:./llama3-tokenizer/tokenizer.json
   ```

### Available LLaMA 3+ Models on HuggingFace

| Model Family | Model Name | Size | HuggingFace Path |
|-------------|------------|------|------------------|
| LLaMA 3.0 | Base | 8B, 70B | meta-llama/Meta-Llama-3-8B |
| LLaMA 3.0 | Instruct | 8B, 70B | meta-llama/Meta-Llama-3-8B-Instruct |
| LLaMA 3.1 | Base | 8B, 70B, 405B | meta-llama/Meta-Llama-3.1-8B |
| LLaMA 3.1 | Instruct | 8B, 70B, 405B | meta-llama/Meta-Llama-3.1-8B-Instruct |
| LLaMA 3.2 | Base | 1B, 3B | meta-llama/Llama-3.2-1B |
| LLaMA 3.2 | Instruct | 1B, 3B, 11B, 90B | meta-llama/Llama-3.2-3B-Instruct |
| LLaMA 3.3 | Instruct | 70B | meta-llama/Llama-3.3-70B-Instruct |

**Note:** All LLaMA 3.x versions share the same tokenizer, so you can use any `tokenizer.json` from the above models.

### Verifying Your Tokenizer File

Check the file size to ensure it downloaded correctly:

```bash
ls -lh tokenizer.json
# Should be approximately 17MB for LLaMA 3+ tokenizers
```

Test with token-visualizer:

```bash
echo "Test tokenization" | ./token-visualizer \
  --model llama3:./tokenizer.json \
  --show-ids

# Should display colorized tokens with IDs
```

### Comparing LLaMA Versions

You can compare tokenization across different LLaMA versions:

```bash
echo "The quick brown fox" | ./token-visualizer compare \
  --models llama:./llama2/tokenizer.model,llama3:./llama3/tokenizer.json \
  --show-ids
```

### License Information

LLaMA 3+ tokenizer files are subject to Meta's license agreements:

- ✅ **Allowed:** Research, development, commercial use (with restrictions based on model size)
- ❌ **Not Allowed:** Redistribution without permission
- **Full License:** https://www.llama.com/llama-downloads/

**Note:** We cannot distribute `tokenizer.json` files with this tool. Users must obtain them directly from Meta or authorized sources.

### Troubleshooting

**Error: "tokenizer model file not found" (LLaMA 1/2)**
- Use absolute path: `llama:/Users/you/path/to/tokenizer.model`
- Verify file exists: `ls -l /path/to/tokenizer.model`
- Make sure you're using `llama:` for LLaMA 1/2 models

**Error: "tokenizer file not found" (LLaMA 3+)**
- Use absolute path: `llama3:/Users/you/path/to/tokenizer.json`
- Verify file exists: `ls -l /path/to/tokenizer.json`
- Make sure you're using `llama3:` for LLaMA 3+ models

**Error: "Internal: could not parse ModelProto"**
- You're trying to use a LLaMA 3+ tokenizer.json with `llama:` prefix
- Use `llama3:` prefix instead: `llama3:/path/to/tokenizer.json`

**Error: "failed to load LLaMA 3 tokenizer"**
- Ensure you downloaded `tokenizer.json`, not `tokenizer.model`
- Verify the JSON file is valid: `jq . tokenizer.json`
- Try downloading the file again

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
