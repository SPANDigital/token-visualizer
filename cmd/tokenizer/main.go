package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/spandigital/token-visualizer/internal/output"
	"github.com/spandigital/token-visualizer/internal/tokenizers"
)

var CLI struct {
	Visualize VisualizeCmd `cmd:"" help:"Visualize tokens with colorized output (default command)" default:"withargs"`
	Count     CountCmd     `cmd:"" help:"Show only token counts"`
	Compare   CompareCmd   `cmd:"" help:"Compare tokenization across multiple models"`
}

type VisualizeCmd struct {
	Model          string `help:"Model to use: gpt4, gpt3.5, gpt5, gpt5-mini, gpt5-nano, claude, llama" default:"gpt4" enum:"gpt4,gpt3.5,gpt5,gpt5-mini,gpt5-nano,claude,llama"`
	Format         string `help:"Output format: terminal, markdown, html" default:"terminal" enum:"terminal,markdown,html"`
	ShowIDs        bool   `help:"Show token IDs" short:"i" name:"show-ids"`
	ShowBoundaries bool   `help:"Show token boundaries" short:"b"`
	Encoding       string `help:"Tiktoken encoding (for GPT models)" default:"cl100k_base"`
	ClaudeModel    string `help:"Claude model name" default:"claude-3-5-sonnet-20241022"`
	LlamaModel     string `help:"Path to LLaMA tokenizer.model file"`
	NoCache        bool   `help:"Disable caching for Claude API" short:"n"`
}

type CountCmd struct {
	Models      []string `help:"Models to count: gpt4, gpt3.5, gpt5, gpt5-mini, gpt5-nano, claude, llama" default:"gpt4"`
	Encoding    string   `help:"Tiktoken encoding (for GPT models)" default:"cl100k_base"`
	ClaudeModel string   `help:"Claude model name" default:"claude-3-5-sonnet-20241022"`
	LlamaModel  string   `help:"Path to LLaMA tokenizer.model file"`
	NoCache     bool     `help:"Disable caching for Claude API" short:"n"`
}

type CompareCmd struct {
	Models         []string `help:"Models to compare: gpt4, gpt3.5, gpt5, gpt5-mini, gpt5-nano, claude, llama" required:""`
	Format         string   `help:"Output format: terminal, markdown, html" default:"terminal" enum:"terminal,markdown,html"`
	ShowIDs        bool     `help:"Show token IDs" short:"i" name:"show-ids"`
	ShowBoundaries bool     `help:"Show token boundaries" short:"b"`
	Encoding       string   `help:"Tiktoken encoding (for GPT models)" default:"cl100k_base"`
	ClaudeModel    string   `help:"Claude model name" default:"claude-3-5-sonnet-20241022"`
	LlamaModel     string   `help:"Path to LLaMA tokenizer.model file"`
	NoCache        bool     `help:"Disable caching for Claude API" short:"n"`
}

func (v *VisualizeCmd) Run() error {
	// Read from stdin
	input, err := readInput()
	if err != nil {
		return fmt.Errorf("failed to read input: %w", err)
	}

	// Create tokenizer
	tokenizer, err := createTokenizer(v.Model, v.Encoding, v.ClaudeModel, v.LlamaModel, !v.NoCache)
	if err != nil {
		return err
	}

	// Tokenize
	ctx := context.Background()
	result, err := tokenizer.Encode(ctx, input)
	if err != nil {
		return fmt.Errorf("tokenization failed: %w", err)
	}

	// Render output
	var outputStr string
	switch v.Format {
	case "terminal":
		renderer := output.NewTerminalRenderer(v.ShowIDs, v.ShowBoundaries)
		outputStr = renderer.RenderSingle(result)
	case "markdown":
		renderer := output.NewMarkdownRenderer(v.ShowIDs)
		outputStr = renderer.RenderSingle(result)
	case "html":
		mdRenderer := output.NewMarkdownRenderer(v.ShowIDs)
		markdown := mdRenderer.RenderSingle(result)
		htmlRenderer := output.NewHTMLRenderer()
		outputStr, err = htmlRenderer.Convert(markdown)
		if err != nil {
			return fmt.Errorf("HTML conversion failed: %w", err)
		}
	}

	fmt.Print(outputStr)
	return nil
}

func (c *CountCmd) Run() error {
	// Read from stdin
	input, err := readInput()
	if err != nil {
		return fmt.Errorf("failed to read input: %w", err)
	}

	ctx := context.Background()
	results := make([]*tokenizers.TokenizationResult, 0, len(c.Models))

	// Process each model
	for _, model := range c.Models {
		tokenizer, err := createTokenizer(model, c.Encoding, c.ClaudeModel, c.LlamaModel, !c.NoCache)
		if err != nil {
			return err
		}

		result, err := tokenizer.Encode(ctx, input)
		if err != nil {
			return fmt.Errorf("tokenization failed for %s: %w", model, err)
		}

		results = append(results, result)
	}

	// Render
	renderer := output.NewTerminalRenderer(false, false)
	fmt.Print(renderer.RenderCountOnly(results))

	return nil
}

func (c *CompareCmd) Run() error {
	// Read from stdin
	input, err := readInput()
	if err != nil {
		return fmt.Errorf("failed to read input: %w", err)
	}

	ctx := context.Background()
	results := make([]*tokenizers.TokenizationResult, 0, len(c.Models))

	// Process each model
	for _, model := range c.Models {
		tokenizer, err := createTokenizer(model, c.Encoding, c.ClaudeModel, c.LlamaModel, !c.NoCache)
		if err != nil {
			return err
		}

		result, err := tokenizer.Encode(ctx, input)
		if err != nil {
			return fmt.Errorf("tokenization failed for %s: %w", model, err)
		}

		results = append(results, result)
	}

	// Render output
	var outputStr string
	switch c.Format {
	case "terminal":
		renderer := output.NewTerminalRenderer(c.ShowIDs, c.ShowBoundaries)
		outputStr = renderer.RenderComparison(results)
	case "markdown":
		renderer := output.NewMarkdownRenderer(c.ShowIDs)
		outputStr = renderer.RenderComparison(results)
	case "html":
		mdRenderer := output.NewMarkdownRenderer(c.ShowIDs)
		markdown := mdRenderer.RenderComparison(results)
		htmlRenderer := output.NewHTMLRenderer()
		outputStr, err = htmlRenderer.Convert(markdown)
		if err != nil {
			return fmt.Errorf("HTML conversion failed: %w", err)
		}
	}

	fmt.Print(outputStr)
	return nil
}

func readInput() (string, error) {
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		return "", err
	}

	input := strings.TrimSpace(string(data))
	if input == "" {
		return "", fmt.Errorf("no input provided (stdin is empty)")
	}

	return input, nil
}

func createTokenizer(model, encoding, claudeModel, llamaModel string, useCache bool) (tokenizers.Tokenizer, error) {
	switch model {
	case "gpt4":
		return tokenizers.NewTikTokenizer(encoding)
	case "gpt3.5":
		return tokenizers.NewTikTokenizer(encoding)
	case "gpt5", "gpt5-mini", "gpt5-nano":
		// GPT-5 models all use o200k_base encoding
		return tokenizers.NewTikTokenizer("o200k_base")
	case "claude":
		return tokenizers.NewClaudeTokenizer(claudeModel, useCache)
	case "llama":
		if llamaModel == "" {
			return nil, fmt.Errorf("LLaMA model requires --llama-model path to tokenizer.model file")
		}
		return tokenizers.NewLLaMATokenizer(llamaModel)
	default:
		return nil, fmt.Errorf("unknown model: %s", model)
	}
}

func main() {
	ctx := kong.Parse(&CLI,
		kong.Name("token-visualizer"),
		kong.Description("Visualize and analyze tokens from various LLM tokenizers"),
		kong.UsageOnError(),
	)

	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}
