package tui

import (
	"strings"
	"testing"

	"charm.land/lipgloss/v2"
)

func TestRenderContextWritePositionTracking(t *testing.T) {
	layout := viewportLayout{X: 5, Y: 10, Width: 80, Height: 24}
	m := NewModel(&mockRepo{}, &mockScheduler{})
	ctx := NewRenderContext(m, layout, ViewBrowser)

	if ctx.currX != 5 || ctx.currY != 10 {
		t.Fatalf("initial position = (%d, %d), want (5, 10)", ctx.currX, ctx.currY)
	}

	// Empty write does nothing
	ctx.Write("")
	if ctx.currX != 5 || ctx.currY != 10 {
		t.Fatalf("after empty write position = (%d, %d), want (5, 10)", ctx.currX, ctx.currY)
	}

	// Single line write advances X by visual width
	ctx.Write("Hello")
	if ctx.currX != 10 || ctx.currY != 10 {
		t.Fatalf("after 'Hello' position = (%d, %d), want (10, 10)", ctx.currX, ctx.currY)
	}

	// Unicode single line write
	ctx.Write(" Straße") // width 7
	if ctx.currX != 17 || ctx.currY != 10 {
		t.Fatalf("after Unicode write position = (%d, %d), want (17, 10)", ctx.currX, ctx.currY)
	}

	// Write newline only
	ctx.Write("\n")
	if ctx.currX != 5 || ctx.currY != 11 {
		t.Fatalf("after newline position = (%d, %d), want (5, 11)", ctx.currX, ctx.currY)
	}

	// Multiline write
	ctx.Write("Line1\nLine2\nLastLine")
	if ctx.currY != 13 {
		t.Fatalf("after multiline Y = %d, want 13", ctx.currY)
	}
	expectedX := layout.X + lipgloss.Width("LastLine") // 5 + 8 = 13
	if ctx.currX != expectedX {
		t.Fatalf("after multiline X = %d, want %d", ctx.currX, expectedX)
	}

	// Trailing newline in multiline write
	ctx.Write("part\n")
	if ctx.currY != 14 {
		t.Fatalf("after trailing newline Y = %d, want 14", ctx.currY)
	}
	if ctx.currX != 5 {
		t.Fatalf("after trailing newline X = %d, want 5", ctx.currX)
	}

	// WriteLine
	ctx.WriteLine("Header")
	if ctx.currY != 15 || ctx.currX != 5 {
		t.Fatalf("after WriteLine position = (%d, %d), want (5, 15)", ctx.currX, ctx.currY)
	}

	// Multiple consecutive newlines
	ctx.Write("\n\n\n")
	if ctx.currY != 18 || ctx.currX != 5 {
		t.Fatalf("after 3 newlines position = (%d, %d), want (5, 18)", ctx.currX, ctx.currY)
	}

	// Styled text with ANSI escape codes
	styled := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205")).Render("Styled")
	ctx.Write(styled)
	expectedX = 5 + lipgloss.Width("Styled") // 5 + 6 = 11
	if ctx.currX != expectedX || ctx.currY != 18 {
		t.Fatalf("after styled write position = (%d, %d), want (%d, 18)", ctx.currX, ctx.currY, expectedX)
	}
}

func TestRenderContextHitboxRegistration(t *testing.T) {
	layout := viewportLayout{X: 2, Y: 4, Width: 80, Height: 24}
	m := NewModel(&mockRepo{}, &mockScheduler{})
	ctx := NewRenderContext(m, layout, ViewDashboard)

	ctx.Write("Prefix: ")
	ctx.RegisterHitbox("btn-1", 10, 1)

	if len(m.hitboxes) != 1 {
		t.Fatalf("hitboxes count = %d, want 1", len(m.hitboxes))
	}
	hb := m.hitboxes[0]
	if hb.ID != "btn-1" || hb.View != ViewDashboard || hb.X != ctx.currX || hb.Y != ctx.currY || hb.Width != 10 || hb.Height != 1 {
		t.Fatalf("unexpected hitbox: %+v", hb)
	}

	// RegisterHitboxAtWithAction
	ctx.RegisterHitboxAtWithAction("btn-offset", 3, -1, 5, 2, nil)
	if len(m.hitboxes) != 2 {
		t.Fatalf("hitboxes count = %d, want 2", len(m.hitboxes))
	}
	hbOffset := m.hitboxes[1]
	if hbOffset.X != ctx.currX+3 || hbOffset.Y != ctx.currY-1 || hbOffset.Width != 5 || hbOffset.Height != 2 {
		t.Fatalf("unexpected offset hitbox: %+v", hbOffset)
	}
}

func TestRenderContextWriteWrapped(t *testing.T) {
	layout := viewportLayout{X: 0, Y: 0, Width: 20, Height: 10}
	m := NewModel(&mockRepo{}, &mockScheduler{})
	ctx := NewRenderContext(m, layout, ViewPractice)

	items := []WrappedItem{
		{ID: "item-1", Label: "[Tag1]"},
		{ID: "item-2", Label: "[Tag2]"},
		{ID: "item-3", Label: "[Tag3-Long]"},
		{ID: "item-4", Label: "[Tag4]"},
	}

	ctx.WriteWrapped(items, 1)
	res := ctx.String()
	lines := strings.Split(res, "\n")
	if len(lines) < 2 {
		t.Fatalf("expected wrapping across multiple lines, got:\n%s", res)
	}
	if len(m.hitboxes) != 0 { // None had actions
		t.Fatalf("expected 0 hitboxes when actions nil, got %d", len(m.hitboxes))
	}
}
