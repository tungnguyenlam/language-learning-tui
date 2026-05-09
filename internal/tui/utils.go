package tui

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode/utf8"

	"deutsch-tui/internal/core"

	"charm.land/lipgloss/v2"
)

func formatDuration(d time.Duration) string {
	days := int(d.Hours() / 24)
	if days > 0 {
		return fmt.Sprintf("%dd", days)
	}
	hours := int(d.Hours())
	if hours > 0 {
		return fmt.Sprintf("%dh", hours)
	}
	minutes := int(d.Minutes())
	if minutes > 0 {
		return fmt.Sprintf("%dm", minutes)
	}
	return "now"
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func clampInt(v, low, high int) int {
	if high < low {
		return low
	}
	if v < low {
		return low
	}
	if v > high {
		return high
	}
	return v
}

func (m *Model) currentAITemplateSet() string {
	if len(m.aiTemplateSets) == 0 {
		return ""
	}
	return m.aiTemplateSets[clampInt(m.aiTemplateIndex, 0, len(m.aiTemplateSets)-1)]
}

func friendlyError(err error) string {
	if err == nil {
		return ""
	}
	if errors.Is(err, os.ErrNotExist) {
		var pathErr *os.PathError
		if errors.As(err, &pathErr) {
			return fmt.Sprintf("Error: no such file or directory: %s", filepath.Base(pathErr.Path))
		}
		return "Error: no such file or directory"
	}
	return fmt.Sprintf("Error: %v", err)
}

func singleLine(s string) string {
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "\r", " ")
	return strings.Join(strings.Fields(s), " ")
}

func truncateLine(s string, maxWidth int) string {
	if maxWidth <= 0 {
		return ""
	}
	if lipgloss.Width(s) <= maxWidth {
		return s
	}
	if maxWidth <= 3 {
		return strings.Repeat(".", maxWidth)
	}

	runes := []rune(s)
	res := ""
	for _, r := range runes {
		if lipgloss.Width(res+string(r)+"...") > maxWidth {
			break
		}
		res += string(r)
	}
	return res + "..."
}

func decksFromNotes(notes []core.Note) []core.Deck {
	byID := make(map[string]int)
	decks := make([]core.Deck, 0, 1)
	for _, note := range notes {
		deckID := strings.TrimSpace(note.DeckID)
		if deckID == "" {
			deckID = "Imported"
			note.DeckID = deckID
		}
		index, ok := byID[deckID]
		if !ok {
			index = len(decks)
			byID[deckID] = index
			decks = append(decks, core.Deck{
				ID:          deckID,
				Name:        deckID,
				Description: "Imported from Anki TSV.",
			})
		}
		decks[index].Notes = append(decks[index].Notes, note)
	}
	return decks
}

func progressBar(width int, percentage float64, filledColor, emptyColor string) string {
	if width <= 0 {
		return ""
	}
	filledWidth := int(float64(width) * percentage)
	if filledWidth > width {
		filledWidth = width
	}
	if filledWidth < 0 {
		filledWidth = 0
	}
	emptyWidth := width - filledWidth

	filled := lipgloss.NewStyle().Foreground(lipgloss.Color(filledColor)).Render(strings.Repeat("█", filledWidth))
	empty := lipgloss.NewStyle().Foreground(lipgloss.Color(emptyColor)).Render(strings.Repeat("░", emptyWidth))

	return filled + empty
}

var (
	headerStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("229"))
	mutedStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
	panelStyle  = lipgloss.NewStyle().Padding(1, 2).Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("240"))

	compactStyle   = lipgloss.NewStyle().Padding(1, 0)
	navStyle       = lipgloss.NewStyle().PaddingRight(2)
	navActiveStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("212")).PaddingRight(2)
	tabStyle       = lipgloss.NewStyle().Padding(0, 1).Foreground(lipgloss.Color("250"))
	tabActiveStyle = lipgloss.NewStyle().Padding(0, 1).Bold(true).Foreground(lipgloss.Color("229")).Background(lipgloss.Color("57"))
	statusStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("250")).Background(lipgloss.Color("236")).Padding(0, 1)
)

func padLine(line string, width int) string {
	padding := width - lipgloss.Width(line)
	if padding <= 0 {
		return line
	}
	return line + strings.Repeat(" ", padding)
}

func scrollbarThumb(totalLines, visibleLines, offset int) (int, int) {
	if totalLines <= 0 || visibleLines <= 0 || totalLines <= visibleLines {
		return 0, 0
	}
	thumbHeight := maxInt(1, (visibleLines*visibleLines)/totalLines)
	thumbHeight = minInt(visibleLines, thumbHeight)
	maxScroll := maxInt(1, totalLines-visibleLines)
	maxThumbStart := maxInt(0, visibleLines-thumbHeight)
	thumbStart := (clampInt(offset, 0, maxScroll) * maxThumbStart) / maxScroll
	return thumbStart, thumbHeight
}

func scrollOffsetForTrackRow(totalLines, visibleLines, row int) int {
	maxScroll := maxInt(0, totalLines-visibleLines)
	if maxScroll == 0 || visibleLines <= 1 {
		return 0
	}
	return clampInt((clampInt(row, 0, visibleLines-1)*maxScroll)/(visibleLines-1), 0, maxScroll)
}

func selectedIndexForTrackRow(totalItems, visibleLines, row int) int {
	if totalItems <= 0 || visibleLines <= 1 {
		return 0
	}
	return clampInt((clampInt(row, 0, visibleLines-1)*(totalItems-1))/(visibleLines-1), 0, totalItems-1)
}

func scrollbarLineWidth(viewportWidth int) int {
	return maxInt(1, viewportWidth-2)
}

func (m *Model) renderActiveViewPlain(x, y int) string {
	layout := m.activeViewContentLayout()
	layout.X = x
	layout.Y = y
	return m.renderActiveViewPlainAt(layout)
}

func stripANSI(s string) string {
	var b strings.Builder
	inEsc := false
	for i := 0; i < len(s); i++ {
		if s[i] == '\x1b' {
			inEsc = true
			continue
		}
		if inEsc {
			if (s[i] >= 'a' && s[i] <= 'z') || (s[i] >= 'A' && s[i] <= 'Z') {
				inEsc = false
			}
			continue
		}
		b.WriteByte(s[i])
	}
	return b.String()
}

func trimLastRune(s string) string {
	if s == "" {
		return s
	}
	_, size := utf8.DecodeLastRuneInString(s)
	if size <= 0 {
		return ""
	}
	return s[:len(s)-size]
}

func singlePrintableInput(s string) (string, bool) {
	if s == "" {
		return "", false
	}
	r, size := utf8.DecodeRuneInString(s)
	if size <= 0 || size != len(s) || r < ' ' {
		return "", false
	}
	return string(r), true
}

func normalizeAnswer(s string) string {
	s = strings.TrimSpace(strings.ToLower(s))
	s = strings.ReplaceAll(s, "ß", "ss")
	s = strings.ReplaceAll(s, "  ", " ")
	// Strip common punctuation from end
	s = strings.TrimRight(s, ".!?,;:")
	return s
}
