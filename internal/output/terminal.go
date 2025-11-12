package output

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/spandigital/token-visualizer/internal/tokenizers"
)

var (
	// Color palette for tokens
	tokenColors = []lipgloss.Color{
		lipgloss.Color("205"), // Pink
		lipgloss.Color("141"), // Purple
		lipgloss.Color("87"),  // Cyan
		lipgloss.Color("228"), // Yellow
		lipgloss.Color("118"), // Green
		lipgloss.Color("213"), // Magenta
		lipgloss.Color("117"), // Light Blue
		lipgloss.Color("223"), // Peach
	}

	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("86")).
			Border(lipgloss.RoundedBorder()).
			Padding(0, 1).
			MarginBottom(1)

	statsStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Italic(true)

	tokenIDStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("243")).
			Faint(true)

	columnStyle = lipgloss.NewStyle().
			Width(50).
			Padding(1).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240"))
)

// TerminalRenderer renders tokenization results to colorized terminal output
type TerminalRenderer struct {
	showIDs        bool
	showBoundaries bool
}

// NewTerminalRenderer creates a new terminal renderer
func NewTerminalRenderer(showIDs, showBoundaries bool) *TerminalRenderer {
	return &TerminalRenderer{
		showIDs:        showIDs,
		showBoundaries: showBoundaries,
	}
}

// RenderSingle renders a single tokenization result
func (r *TerminalRenderer) RenderSingle(result *tokenizers.TokenizationResult) string {
	var output strings.Builder

	// Header with model name
	header := headerStyle.Render(fmt.Sprintf("🔤 %s", result.Model))
	output.WriteString(header)
	output.WriteString("\n\n")

	// Stats
	stats := statsStyle.Render(fmt.Sprintf("Total tokens: %d", result.TotalCount))
	output.WriteString(stats)
	output.WriteString("\n\n")

	// Render tokens
	for i, token := range result.Tokens {
		colorIdx := i % len(tokenColors)
		tokenStyle := lipgloss.NewStyle().
			Foreground(tokenColors[colorIdx]).
			Bold(true)

		// Token text
		output.WriteString(tokenStyle.Render(token.Text))

		// Optional: Show boundaries
		if r.showBoundaries && i < len(result.Tokens)-1 {
			output.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("238")).Render("|"))
		}

		// Optional: Show token IDs
		if r.showIDs && token.ID >= 0 {
			idStr := tokenIDStyle.Render(fmt.Sprintf("[%d]", token.ID))
			output.WriteString(idStr)
		}
	}

	output.WriteString("\n")
	return output.String()
}

// RenderComparison renders multiple tokenization results side-by-side
func (r *TerminalRenderer) RenderComparison(results []*tokenizers.TokenizationResult) string {
	if len(results) == 0 {
		return ""
	}

	if len(results) == 1 {
		return r.RenderSingle(results[0])
	}

	columns := make([]string, len(results))

	for i, result := range results {
		columns[i] = r.renderColumn(result)
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, columns...)
}

// renderColumn renders a single column for comparison view
func (r *TerminalRenderer) renderColumn(result *tokenizers.TokenizationResult) string {
	var content strings.Builder

	// Model name header
	header := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("86")).
		Underline(true).
		Render(result.Model)
	content.WriteString(header)
	content.WriteString("\n\n")

	// Token count
	stats := statsStyle.Render(fmt.Sprintf("Tokens: %d", result.TotalCount))
	content.WriteString(stats)
	content.WriteString("\n\n")

	// Render tokens
	for i, token := range result.Tokens {
		colorIdx := i % len(tokenColors)
		tokenStyle := lipgloss.NewStyle().
			Foreground(tokenColors[colorIdx]).
			Bold(true)

		// Token text with optional boundary
		if r.showBoundaries {
			content.WriteString(tokenStyle.Render(fmt.Sprintf("[%s]", token.Text)))
		} else {
			content.WriteString(tokenStyle.Render(token.Text))
		}

		// Optional: Show token IDs
		if r.showIDs && token.ID >= 0 {
			idStr := tokenIDStyle.Render(fmt.Sprintf("(%d)", token.ID))
			content.WriteString(idStr)
		}

		content.WriteString("\n")
	}

	return columnStyle.Render(content.String())
}

// RenderCountOnly renders just the token count for multiple models
func (r *TerminalRenderer) RenderCountOnly(results []*tokenizers.TokenizationResult) string {
	var output strings.Builder

	output.WriteString(headerStyle.Render("📊 Token Counts"))
	output.WriteString("\n\n")

	for _, result := range results {
		modelStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("141")).
			Bold(true)

		countStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("228")).
			Bold(true)

		line := fmt.Sprintf("%s: %s",
			modelStyle.Render(result.Model),
			countStyle.Render(fmt.Sprintf("%d tokens", result.TotalCount)))

		output.WriteString(line)
		output.WriteString("\n")
	}

	return output.String()
}
