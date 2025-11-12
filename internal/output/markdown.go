package output

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/spandigital/token-visualizer/internal/tokenizers"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
)

// MarkdownRenderer renders tokenization results as markdown
type MarkdownRenderer struct {
	showIDs bool
}

// NewMarkdownRenderer creates a new markdown renderer
func NewMarkdownRenderer(showIDs bool) *MarkdownRenderer {
	return &MarkdownRenderer{
		showIDs: showIDs,
	}
}

// RenderSingle renders a single tokenization result as markdown
func (r *MarkdownRenderer) RenderSingle(result *tokenizers.TokenizationResult) string {
	var md strings.Builder

	md.WriteString(fmt.Sprintf("# %s\n\n", result.Model))
	md.WriteString(fmt.Sprintf("**Total tokens:** %d\n\n", result.TotalCount))

	if len(result.Tokens) > 0 && result.Tokens[0].ID >= 0 {
		md.WriteString("## Tokens\n\n")
		md.WriteString("| # | Text | ID |\n")
		md.WriteString("|---|------|----|\n")

		for i, token := range result.Tokens {
			text := strings.ReplaceAll(token.Text, "|", "\\|")
			text = strings.ReplaceAll(text, "\n", "\\n")

			if r.showIDs {
				md.WriteString(fmt.Sprintf("| %d | `%s` | %d |\n", i+1, text, token.ID))
			} else {
				md.WriteString(fmt.Sprintf("| %d | `%s` | |\n", i+1, text))
			}
		}
	} else {
		md.WriteString("## Text\n\n")
		md.WriteString(fmt.Sprintf("```\n%s\n```\n", result.Text))
	}

	return md.String()
}

// RenderComparison renders multiple tokenization results as markdown
func (r *MarkdownRenderer) RenderComparison(results []*tokenizers.TokenizationResult) string {
	if len(results) == 0 {
		return ""
	}

	if len(results) == 1 {
		return r.RenderSingle(results[0])
	}

	var md strings.Builder

	md.WriteString("# Token Comparison\n\n")

	// Summary table
	md.WriteString("## Token Counts\n\n")
	md.WriteString("| Model | Token Count |\n")
	md.WriteString("|-------|-------------|\n")

	for _, result := range results {
		md.WriteString(fmt.Sprintf("| %s | %d |\n", result.Model, result.TotalCount))
	}

	md.WriteString("\n")

	// Individual results
	for _, result := range results {
		md.WriteString(fmt.Sprintf("## %s\n\n", result.Model))

		if len(result.Tokens) > 0 && result.Tokens[0].ID >= 0 {
			md.WriteString("| # | Text | ID |\n")
			md.WriteString("|---|------|----|\n")

			for i, token := range result.Tokens {
				text := strings.ReplaceAll(token.Text, "|", "\\|")
				text = strings.ReplaceAll(text, "\n", "\\n")

				if r.showIDs {
					md.WriteString(fmt.Sprintf("| %d | `%s` | %d |\n", i+1, text, token.ID))
				} else {
					md.WriteString(fmt.Sprintf("| %d | `%s` | |\n", i+1, text))
				}
			}
		} else {
			md.WriteString(fmt.Sprintf("**Total tokens:** %d\n", result.TotalCount))
		}

		md.WriteString("\n")
	}

	return md.String()
}

// RenderCountOnly renders just token counts as markdown
func (r *MarkdownRenderer) RenderCountOnly(results []*tokenizers.TokenizationResult) string {
	var md strings.Builder

	md.WriteString("# Token Counts\n\n")
	md.WriteString("| Model | Token Count |\n")
	md.WriteString("|-------|-------------|\n")

	for _, result := range results {
		md.WriteString(fmt.Sprintf("| %s | %d |\n", result.Model, result.TotalCount))
	}

	return md.String()
}

// HTMLRenderer converts markdown to HTML
type HTMLRenderer struct {
	md goldmark.Markdown
}

// NewHTMLRenderer creates a new HTML renderer
func NewHTMLRenderer() *HTMLRenderer {
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.Table,
		),
		goldmark.WithRendererOptions(
			html.WithUnsafe(), // Allow raw HTML
			html.WithXHTML(),
		),
	)

	return &HTMLRenderer{md: md}
}

// Convert converts markdown to HTML
func (h *HTMLRenderer) Convert(markdown string) (string, error) {
	var buf bytes.Buffer

	// Add some basic styling
	buf.WriteString(`<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>Token Visualization</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Helvetica, Arial, sans-serif;
            max-width: 1200px;
            margin: 0 auto;
            padding: 2rem;
            line-height: 1.6;
        }
        table {
            border-collapse: collapse;
            width: 100%;
            margin: 1rem 0;
        }
        th, td {
            border: 1px solid #ddd;
            padding: 0.75rem;
            text-align: left;
        }
        th {
            background-color: #f6f8fa;
            font-weight: 600;
        }
        tr:nth-child(even) {
            background-color: #f6f8fa;
        }
        code {
            background-color: #f6f8fa;
            padding: 0.2em 0.4em;
            border-radius: 3px;
            font-family: "SF Mono", Monaco, Menlo, Consolas, monospace;
        }
        h1 {
            border-bottom: 2px solid #e1e4e8;
            padding-bottom: 0.3rem;
        }
        h2 {
            border-bottom: 1px solid #e1e4e8;
            padding-bottom: 0.3rem;
            margin-top: 2rem;
        }
    </style>
</head>
<body>
`)

	// Convert markdown to HTML
	if err := h.md.Convert([]byte(markdown), &buf); err != nil {
		return "", fmt.Errorf("failed to convert markdown to HTML: %w", err)
	}

	buf.WriteString(`
</body>
</html>
`)

	return buf.String(), nil
}
