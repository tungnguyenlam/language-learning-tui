package tui

import (
	"errors"
	"fmt"
	"image/color"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode/utf8"

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

func formatReviewInterval(interval time.Duration) string {
	if interval <= 0 {
		return "same day"
	}
	hours := int(interval.Hours())
	if hours < 24 {
		if hours <= 1 {
			return "1 hour"
		}
		return fmt.Sprintf("%d hours", hours)
	}
	days := hours / 24
	if days == 1 {
		return "1 day"
	}
	if days >= 365 {
		years := days / 365
		remainMonths := (days % 365) / 30
		if remainMonths > 0 {
			if years == 1 {
				return fmt.Sprintf("1 year %d mo", remainMonths)
			}
			return fmt.Sprintf("%d years %d mo", years, remainMonths)
		}
		if years == 1 {
			return "1 year"
		}
		return fmt.Sprintf("%d years", years)
	}
	if days >= 30 {
		months := days / 30
		remainDays := days % 30
		if remainDays > 0 {
			if months == 1 {
				return fmt.Sprintf("1 month %dd", remainDays)
			}
			return fmt.Sprintf("%d months %dd", months, remainDays)
		}
		if months == 1 {
			return "1 month"
		}
		return fmt.Sprintf("%d months", months)
	}
	return fmt.Sprintf("%d days", days)
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

func (m *Model) deckNameByID(deckID string) string {
	if strings.TrimSpace(deckID) == "" {
		return "Unknown Deck"
	}
	for _, deck := range m.decks {
		if deck.ID == deckID {
			return deck.Name
		}
	}
	if m.deck.ID == deckID && m.deck.Name != "" {
		return m.deck.Name
	}
	return deckID
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

	var b strings.Builder
	currW := 0
	inEsc := false
	runes := []rune(s)

	for i := 0; i < len(runes); i++ {
		r := runes[i]

		if r == '\x1b' {
			inEsc = true
			b.WriteRune(r)
			continue
		}

		if inEsc {
			b.WriteRune(r)
			// ANSI sequences end with a letter (for CSI sequences like colors)
			// or other characters for other sequences. This is a common heuristic.
			if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
				inEsc = false
			}
			continue
		}

		rw := lipgloss.Width(string(r))
		if currW+rw+3 > maxWidth {
			b.WriteString("...")
			// Make sure we close any open ANSI sequences if possible?
			// For simplicity, just append a reset sequence if we were highlighting.
			b.WriteString("\x1b[0m")
			break
		}
		b.WriteRune(r)
		currW += rw
	}
	return b.String()
}

func progressBar(width int, percentage float64, filledColor, emptyColor color.Color) string {
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

	filled := lipgloss.NewStyle().Foreground(filledColor).Render(strings.Repeat("█", filledWidth))
	empty := lipgloss.NewStyle().Foreground(emptyColor).Render(strings.Repeat("░", emptyWidth))

	return filled + empty
}

func sparkline(data []int, width int) string {
	if width <= 0 {
		return ""
	}
	if len(data) == 0 {
		return strings.Repeat(" ", width)
	}
	if len(data) > width {
		data = data[len(data)-width:]
	} else if len(data) < width {
		padded := make([]int, width)
		copy(padded[width-len(data):], data)
		data = padded
	}
	maxVal := 1
	for _, v := range data {
		if v > maxVal {
			maxVal = v
		}
	}
	blocks := []string{" ", "▂", "▃", "▄", "▅", "▆", "▇", "█"}
	var res strings.Builder
	for _, v := range data {
		if v == 0 {
			res.WriteString(" ")
			continue
		}
		idx := (v * (len(blocks) - 1)) / maxVal
		res.WriteString(blocks[idx])
	}
	return res.String()
}

func padLine(line string, width int) string {
	padding := width - lipgloss.Width(line)
	if padding <= 0 {
		return line
	}
	return line + strings.Repeat(" ", padding)
}

// fillViewportContent pads content to the viewport height and width when the
// rendered view is shorter than the panel interior, preventing stale
// characters from prior frames bleeding through. Content taller than the
// viewport is left intact so scrollable views keep list layout and hitboxes.
func fillViewportContent(content string, layout viewportLayout) string {
	lines := strings.Split(strings.TrimSuffix(content, "\n"), "\n")
	if len(lines) >= layout.Height {
		return content
	}
	filled := make([]string, layout.Height)
	for i := 0; i < layout.Height; i++ {
		line := ""
		if i < len(lines) {
			line = lines[i]
		}
		filled[i] = padLine(truncateLine(line, layout.Width), layout.Width)
	}
	return strings.Join(filled, "\n")
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
	if s == "space" {
		return " ", true
	}
	r, size := utf8.DecodeRuneInString(s)
	if size <= 0 || size != len(s) || r < ' ' {
		return "", false
	}
	return string(r), true
}

func (m *Model) normalizeAnswer(s string) string {
	s = strings.TrimSpace(strings.ToLower(s))
	if !m.strictNormalization {
		s = strings.ReplaceAll(s, "ß", "ss")
		s = strings.ReplaceAll(s, "ä", "ae")
		s = strings.ReplaceAll(s, "ö", "oe")
		s = strings.ReplaceAll(s, "ü", "ue")
	}
	s = strings.ReplaceAll(s, "  ", " ")
	// Strip common punctuation from end
	s = strings.TrimRight(s, ".!?,;:")
	return s
}

func normalizeUmlauts(s string) string {
	s = strings.ReplaceAll(s, "ß", "ss")
	s = strings.ReplaceAll(s, "ä", "ae")
	s = strings.ReplaceAll(s, "ö", "oe")
	s = strings.ReplaceAll(s, "ü", "ue")
	return s
}

func highlightQuery(text, query string, style lipgloss.Style) string {
	query = strings.TrimSpace(query)
	if query == "" {
		return text
	}
	textRunes := []rune(text)
	lowerTextRunes := []rune(strings.ToLower(text))
	lowerQueryRunes := []rune(strings.ToLower(query))
	if len(lowerQueryRunes) == 0 || len(lowerQueryRunes) > len(lowerTextRunes) {
		return text
	}

	var result strings.Builder
	pos := 0
	for pos < len(lowerTextRunes) {
		matchStart := -1
		for i := pos; i <= len(lowerTextRunes)-len(lowerQueryRunes); i++ {
			matched := true
			for j := range lowerQueryRunes {
				if lowerTextRunes[i+j] != lowerQueryRunes[j] {
					matched = false
					break
				}
			}
			if matched {
				matchStart = i
				break
			}
		}
		if matchStart == -1 {
			result.WriteString(string(textRunes[pos:]))
			break
		}

		result.WriteString(string(textRunes[pos:matchStart]))
		matchEnd := matchStart + len(lowerQueryRunes)
		result.WriteString(style.Render(string(textRunes[matchStart:matchEnd])))
		pos = matchEnd
	}
	return result.String()
}

func renderTypingDiff(typed, expected string) string {
	tRunes := []rune(typed)
	eRunes := []rune(expected)

	tLen, eLen := len(tRunes), len(eRunes)
	m := make([][]int, tLen+1)
	for i := range m {
		m[i] = make([]int, eLen+1)
	}
	for i := 1; i <= tLen; i++ {
		for j := 1; j <= eLen; j++ {
			if tRunes[i-1] == eRunes[j-1] {
				m[i][j] = m[i-1][j-1] + 1
			} else {
				m[i][j] = maxInt(m[i-1][j], m[i][j-1])
			}
		}
	}

	correctStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("46")).Bold(true)
	wrongStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true).Strikethrough(true)
	missingStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Bold(true)

	var diff []string
	i, j := tLen, eLen
	for i > 0 || j > 0 {
		if i > 0 && j > 0 && tRunes[i-1] == eRunes[j-1] {
			diff = append([]string{correctStyle.Render(string(tRunes[i-1]))}, diff...)
			i--
			j--
		} else if j > 0 && (i == 0 || m[i][j-1] >= m[i-1][j]) {
			char := string(eRunes[j-1])
			if char == " " {
				char = "_"
			}
			diff = append([]string{missingStyle.Render(char)}, diff...)
			j--
		} else if i > 0 && (j == 0 || m[i][j-1] < m[i-1][j]) {
			diff = append([]string{wrongStyle.Render(string(tRunes[i-1]))}, diff...)
			i--
		}
	}

	return strings.Join(diff, "")
}
