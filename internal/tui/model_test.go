package tui

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
)

func TestBreakpointForWidth(t *testing.T) {
	tests := []struct {
		width int
		want  Breakpoint
	}{
		{50, BreakpointCompact},
		{80, BreakpointMedium},
		{120, BreakpointWide},
	}
	for _, tt := range tests {
		if got := breakpointForWidth(tt.width); got != tt.want {
			t.Fatalf("breakpointForWidth(%d) = %s, want %s", tt.width, got, tt.want)
		}
	}
}

func TestWindowResizeChangesViewShape(t *testing.T) {
	model := NewModel()
	updated, _ := model.Update(tea.WindowSizeMsg{Width: 50, Height: 20})
	model = updated.(*Model)
	if model.breakpoint != BreakpointCompact {
		t.Fatalf("breakpoint = %s, want compact", model.breakpoint)
	}
	if !strings.Contains(model.View().Content, "Dashboard") {
		t.Fatal("view should render dashboard")
	}
}

func TestHitboxActivation(t *testing.T) {
	model := NewModel()
	model.View()
	model.activateHitbox("tab-review")
	if model.activeView != ViewReview {
		t.Fatalf("active view = %s, want review", model.activeView)
	}
}
