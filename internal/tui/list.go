package tui

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
)

type ListOptions struct {
	HitboxPrefix string
	View         View
	Footer       string
	ScrollOffset *int
	TotalLines   *int
	Cursor       int
	// OnLine is called for each visible line to register hitboxes.
	OnLine func(lineIdx int, ctx *RenderContext, content string)
}

// RenderList provides a unified way to render scrollable lists with hitboxes and auto-scrolling.
func (m *Model) RenderList(layout viewportLayout, content string, opts ListOptions) string {
	lines := strings.Split(strings.TrimSuffix(content, "\n"), "\n")
	totalLines := len(lines)
	if opts.TotalLines != nil {
		*opts.TotalLines = totalLines
	}

	maxVisible := layout.Height
	if maxVisible <= 0 {
		return ""
	}

	// Default scroll offset if pointer provided
	var scroll int
	if opts.ScrollOffset != nil {
		scroll = *opts.ScrollOffset
	}

	// Ensure scroll offset is within bounds
	scroll = clampInt(scroll, 0, maxInt(0, totalLines-maxVisible))

	thumbStart, thumbHeight := scrollbarThumb(totalLines, maxVisible, scroll)
	padWidth := layout.Width

	ctx := NewRenderContext(m, layout, opts.View)

	for i := 0; i < maxVisible; i++ {
		lineIdx := i + scroll

		var lineContent string
		if lineIdx < totalLines {
			lineContent = lines[lineIdx]
		}

		if opts.OnLine != nil && lineIdx < totalLines {
			// Provide a context for this specific line
			lineCtx := NewRenderContext(m, viewportLayout{X: layout.X, Y: ctx.currY, Width: layout.Width, Height: 1}, opts.View)
			opts.OnLine(lineIdx, lineCtx, lineContent)
		}

		if totalLines > maxVisible {
			// Re-truncate and re-pad to reserve 2 cols for scrollbar
			displayLine := truncateLine(lineContent, padWidth-2)
			displayLine = padLine(displayLine, padWidth-2)
			ctx.Write(displayLine)
			ctx.Write(" ")

			scrollbarChar := "│"
			style := lipgloss.NewStyle().Foreground(colorPanel)
			if lineIdx < totalLines {
				if i >= thumbStart && i < thumbStart+thumbHeight {
					scrollbarChar = "█"
					style = lipgloss.NewStyle().Foreground(colorAccent)
				}
			}

			// Register scrollbar hitbox precisely where it is written
			ctx.RegisterHitbox(fmt.Sprintf("%s-scroll-%d", opts.HitboxPrefix, i), 1, 1)
			ctx.Write(style.Render(scrollbarChar))
		} else {
			displayLine := padLine(truncateLine(lineContent, padWidth), padWidth)
			ctx.Write(displayLine)
		}

		ctx.NewLine()
	}

	if opts.ScrollOffset != nil {
		*opts.ScrollOffset = scroll
	}

	res := ctx.String()
	if opts.Footer != "" {
		res += "\n" + lipgloss.NewStyle().Foreground(colorMuted).Render(opts.Footer)
	}

	return res
}

// AutoScroll calculates the new scroll offset to keep the cursor visible.
func AutoScroll(cursorLine, scroll, visibleHeight, totalLines int) int {
	if cursorLine < scroll {
		return cursorLine
	}
	if cursorLine >= scroll+visibleHeight {
		return cursorLine - visibleHeight + 1
	}
	return clampInt(scroll, 0, maxInt(0, totalLines-visibleHeight))
}
