package output

import (
	"fmt"
	"html"
	"strings"

	"github.com/spandigital/token-visualizer/internal/tokenizers"
)

// HTMLInlineRenderer generates inline HTML output with colored tokens similar to terminal output
type HTMLInlineRenderer struct {
	showIDs        bool
	showBoundaries bool
}

// NewHTMLInlineRenderer creates a new HTML inline renderer
func NewHTMLInlineRenderer(showIDs, showBoundaries bool) *HTMLInlineRenderer {
	return &HTMLInlineRenderer{
		showIDs:        showIDs,
		showBoundaries: showBoundaries,
	}
}

// tokenColors matches the terminal output color palette
var htmlTokenColors = []string{
	"#ff5faf", // Pink
	"#af87ff", // Purple
	"#5fffff", // Cyan
	"#ffff87", // Yellow
	"#87ff00", // Green
	"#ff87ff", // Magenta
	"#87d7ff", // Light Blue
	"#ffd7af", // Peach
}

// generateCSS generates the CSS styles for the HTML output
func (r *HTMLInlineRenderer) generateCSS() string {
	var css strings.Builder
	css.WriteString(`<style>
body {
    font-family: 'SF Mono', 'Monaco', 'Inconsolata', 'Fira Code', 'Consolas', monospace;
    padding: 20px;
    background-color: #1e1e1e;
    color: #d4d4d4;
    line-height: 1.6;
}
.container {
    max-width: 1200px;
    margin: 0 auto;
}
.model-header {
    font-size: 1.2em;
    font-weight: bold;
    margin-bottom: 10px;
    color: #569cd6;
}
.token-count {
    color: #9cdcfe;
    margin-bottom: 15px;
}
.tokens {
    background-color: #252526;
    padding: 15px;
    border-radius: 5px;
    margin-bottom: 20px;
    word-wrap: break-word;
}
.token {
    font-weight: bold;
    white-space: pre-wrap;
}
.token-id {
    font-size: 0.8em;
    color: #6a6a6a;
    font-weight: normal;
}
.boundary {
    color: #6a6a6a;
    font-weight: normal;
}
.comparison-container {
    display: flex;
    gap: 20px;
    flex-wrap: wrap;
}
.comparison-model {
    flex: 1;
    min-width: 300px;
}
`)

	// Generate color classes for each token rotation
	for i, color := range htmlTokenColors {
		css.WriteString(fmt.Sprintf(".token-%d { color: %s; }\n", i, color))
	}

	css.WriteString("</style>\n")
	return css.String()
}

// RenderSingle renders a single tokenization result as inline HTML
func (r *HTMLInlineRenderer) RenderSingle(result *tokenizers.TokenizationResult) string {
	var html strings.Builder

	// HTML header with CSS
	html.WriteString("<!DOCTYPE html>\n<html>\n<head>\n")
	html.WriteString("<meta charset=\"UTF-8\">\n")
	html.WriteString("<meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">\n")
	html.WriteString("<title>Token Visualization</title>\n")
	html.WriteString(r.generateCSS())
	html.WriteString("</head>\n<body>\n<div class=\"container\">\n")

	// Model header
	html.WriteString(fmt.Sprintf("<div class=\"model-header\">%s</div>\n", escapeHTML(result.Model)))
	html.WriteString(fmt.Sprintf("<div class=\"token-count\">Total tokens: %d</div>\n", result.TotalCount))

	// Tokens
	html.WriteString("<div class=\"tokens\">\n")
	for i, token := range result.Tokens {
		colorIdx := i % len(htmlTokenColors)

		if r.showBoundaries && i > 0 {
			html.WriteString("<span class=\"boundary\">|</span>")
		}

		html.WriteString(fmt.Sprintf("<span class=\"token token-%d\">%s</span>", colorIdx, escapeHTML(token.Text)))

		if r.showIDs {
			html.WriteString(fmt.Sprintf("<span class=\"token-id\">[%d]</span>", token.ID))
		}
	}
	html.WriteString("\n</div>\n")

	// Footer
	html.WriteString("</div>\n</body>\n</html>\n")

	return html.String()
}

// RenderComparison renders multiple tokenization results side by side as inline HTML
func (r *HTMLInlineRenderer) RenderComparison(results []*tokenizers.TokenizationResult) string {
	var html strings.Builder

	// HTML header with CSS
	html.WriteString("<!DOCTYPE html>\n<html>\n<head>\n")
	html.WriteString("<meta charset=\"UTF-8\">\n")
	html.WriteString("<meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">\n")
	html.WriteString("<title>Token Comparison</title>\n")
	html.WriteString(r.generateCSS())
	html.WriteString("</head>\n<body>\n<div class=\"container\">\n")

	// Comparison container with side-by-side models
	html.WriteString("<div class=\"comparison-container\">\n")

	for _, result := range results {
		html.WriteString("<div class=\"comparison-model\">\n")

		// Model header
		html.WriteString(fmt.Sprintf("<div class=\"model-header\">%s</div>\n", escapeHTML(result.Model)))
		html.WriteString(fmt.Sprintf("<div class=\"token-count\">Total tokens: %d</div>\n", result.TotalCount))

		// Tokens
		html.WriteString("<div class=\"tokens\">\n")
		for i, token := range result.Tokens {
			colorIdx := i % len(htmlTokenColors)

			if r.showBoundaries && i > 0 {
				html.WriteString("<span class=\"boundary\">|</span>")
			}

			html.WriteString(fmt.Sprintf("<span class=\"token token-%d\">%s</span>", colorIdx, escapeHTML(token.Text)))

			if r.showIDs {
				html.WriteString(fmt.Sprintf("<span class=\"token-id\">[%d]</span>", token.ID))
			}
		}
		html.WriteString("\n</div>\n")

		html.WriteString("</div>\n")
	}

	html.WriteString("</div>\n")

	// Footer
	html.WriteString("</div>\n</body>\n</html>\n")

	return html.String()
}

// RenderCountOnly renders only the token counts as HTML
func (r *HTMLInlineRenderer) RenderCountOnly(results []*tokenizers.TokenizationResult) string {
	var html strings.Builder

	// HTML header with minimal CSS
	html.WriteString("<!DOCTYPE html>\n<html>\n<head>\n")
	html.WriteString("<meta charset=\"UTF-8\">\n")
	html.WriteString("<meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">\n")
	html.WriteString("<title>Token Counts</title>\n")
	html.WriteString(`<style>
body {
    font-family: 'SF Mono', 'Monaco', 'Inconsolata', 'Fira Code', 'Consolas', monospace;
    padding: 20px;
    background-color: #1e1e1e;
    color: #d4d4d4;
}
.container {
    max-width: 800px;
    margin: 0 auto;
}
.count-item {
    background-color: #252526;
    padding: 15px;
    border-radius: 5px;
    margin-bottom: 10px;
}
.model-name {
    color: #569cd6;
    font-weight: bold;
}
.token-count {
    color: #9cdcfe;
    margin-left: 10px;
}
</style>
`)
	html.WriteString("</head>\n<body>\n<div class=\"container\">\n")

	for _, result := range results {
		html.WriteString("<div class=\"count-item\">")
		html.WriteString(fmt.Sprintf("<span class=\"model-name\">%s</span>", escapeHTML(result.Model)))
		html.WriteString(fmt.Sprintf("<span class=\"token-count\">%d tokens</span>", result.TotalCount))
		html.WriteString("</div>\n")
	}

	html.WriteString("</div>\n</body>\n</html>\n")

	return html.String()
}

// escapeHTML escapes HTML special characters
func escapeHTML(s string) string {
	return html.EscapeString(s)
}
