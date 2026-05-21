package tui

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
)

type scrollableOptions struct {
	hitboxPrefix string
	view         View
	footer       string
	// onLine is called for each visible line.
	// lineIndex is the index in the original split lines.
	// rY is the absolute Y coordinate on the screen.
	// content is the raw content of that line (unpadded).
	onLine func(lineIndex int, rY int, content string)
}

func (m *Model) renderScrollable(layout viewportLayout, content string, scrollOffset int, options scrollableOptions) string {
	lines := strings.Split(content, "\n")
	totalLines := len(lines)
	maxVisible := layout.Height

	// Ensure scroll offset is within bounds
	scrollOffset = clampInt(scrollOffset, 0, maxInt(0, totalLines-maxVisible))

	thumbStart, thumbHeight := scrollbarThumb(totalLines, maxVisible, scrollOffset)
	padWidth := layout.Width

	var visibleContent strings.Builder
	for i := 0; i < maxVisible; i++ {
		lineIdx := i + scrollOffset
		rY := layout.Y + i

		var lineContent string
		if lineIdx < totalLines {
			lineContent = lines[lineIdx]
		}

		if options.onLine != nil && lineIdx < totalLines {
			options.onLine(lineIdx, rY, lineContent)
		}

		// Basic line rendering with truncation to fit viewport
		displayLine := truncateLine(lineContent, padWidth)

		if totalLines > maxVisible {
			// Re-truncate and re-pad to reserve 2 cols for scrollbar
			displayLine = truncateLine(lineContent, padWidth-2)
			displayLine = padLine(displayLine, padWidth-2)

			scrollbarChar := "│"
			if lineIdx >= totalLines {
				// Empty line, still show track
				scrollbarChar = "│"
			} else if i >= thumbStart && i < thumbStart+thumbHeight {
				scrollbarChar = "█"
			}
			displayLine = displayLine + " " + scrollbarChar

			// Add scrollbar hitbox
			m.hitboxes = append(m.hitboxes, Hitbox{
				ID:     fmt.Sprintf("%s-scroll-%d", options.hitboxPrefix, i),
				View:   options.view,
				X:      layout.X + padWidth - 1,
				Y:      rY,
				Width:  1,
				Height: 1,
			})
		} else {
			displayLine = padLine(displayLine, padWidth)
		}
		visibleContent.WriteString(displayLine + "\n")
	}

	res := visibleContent.String()
	if options.footer != "" {
		res += "\n" + lipgloss.NewStyle().Foreground(colorMuted).Render(options.footer)
	}

	return res
}
