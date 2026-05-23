package tui

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// RenderContext tracks the current rendering position and manages hitboxes.
// It simplifies view rendering by handling Y-offset tracking automatically.
type RenderContext struct {
	model  *Model
	layout viewportLayout
	view   View
	currX  int
	currY  int
	buffer strings.Builder
}

func NewRenderContext(m *Model, layout viewportLayout, view View) *RenderContext {
	return &RenderContext{
		model:  m,
		layout: layout,
		view:   view,
		currX:  layout.X,
		currY:  layout.Y,
	}
}

// Write appends a string to the buffer and updates the current position.
func (c *RenderContext) Write(s string) {
	if s == "" {
		return
	}
	c.buffer.WriteString(s)

	// Split into lines to track Y accurately.
	// Note: lipgloss.Width handles multi-line strings by taking the max width,
	// but we need to track the LAST line's width for currX.
	lines := strings.Split(s, "\n")
	if len(lines) > 1 {
		c.currY += len(lines) - 1
		c.currX = c.layout.X + lipgloss.Width(lines[len(lines)-1])
	} else {
		c.currX += lipgloss.Width(s)
	}
}

// WriteLine appends a string followed by a newline.
func (c *RenderContext) WriteLine(s string) {
	c.Write(s)
	c.NewLine()
}

// NewLine adds a newline and resets currX to layout start.
func (c *RenderContext) NewLine() {
	c.buffer.WriteString("\n")
	c.currY++
	c.currX = c.layout.X
}

// RegisterHitbox adds a hitbox at the current position.
func (c *RenderContext) RegisterHitbox(id string, w, h int) {
	c.model.hitboxes = append(c.model.hitboxes, Hitbox{
		ID:     id,
		View:   c.view,
		X:      c.currX,
		Y:      c.currY,
		Width:  w,
		Height: h,
	})
}

// RegisterHitboxWithAction adds a hitbox with a callback at the current position.
func (c *RenderContext) RegisterHitboxWithAction(id string, w, h int, action func() tea.Cmd) {
	c.model.hitboxes = append(c.model.hitboxes, Hitbox{
		ID:     id,
		View:   c.view,
		X:      c.currX,
		Y:      c.currY,
		Width:  w,
		Height: h,
		Action: action,
	})
}

// RegisterHitboxAt adds a hitbox at a specific offset from the current position.
func (c *RenderContext) RegisterHitboxAt(id string, offsetX, offsetY, w, h int) {
	c.model.hitboxes = append(c.model.hitboxes, Hitbox{
		ID:     id,
		View:   c.view,
		X:      c.currX + offsetX,
		Y:      c.currY + offsetY,
		Width:  w,
		Height: h,
	})
}

// RegisterHitboxAtWithAction adds a hitbox with a callback at a specific offset from the current position.
func (c *RenderContext) RegisterHitboxAtWithAction(id string, offsetX, offsetY, w, h int, action func() tea.Cmd) {
	c.model.hitboxes = append(c.model.hitboxes, Hitbox{
		ID:     id,
		View:   c.view,
		X:      c.currX + offsetX,
		Y:      c.currY + offsetY,
		Width:  w,
		Height: h,
		Action: action,
	})
}

// RegisterAction registers a hitbox for the text about to be written.
func (c *RenderContext) RegisterAction(id string, w, h int, action func() tea.Cmd) {
	c.RegisterHitboxWithAction(id, w, h, action)
}

// WriteAction writes a styled action label and registers its hitbox.
func (c *RenderContext) WriteAction(id string, label string, style lipgloss.Style, action func() tea.Cmd) {
	width := lipgloss.Width(label)
	c.RegisterAction(id, width, 1, action)
	c.Write(style.Render(label))
}

// String returns the accumulated rendered content.
func (c *RenderContext) String() string {
	return c.buffer.String()
}

// Layout returns the current viewport layout.
func (c *RenderContext) Layout() viewportLayout {
	return c.layout
}

// CurrentY returns the current Y coordinate (absolute).
func (c *RenderContext) CurrentY() int {
	return c.currY
}
