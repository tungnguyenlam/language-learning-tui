package tui

import (
	"fmt"
	"strings"
	"testing"

	"charm.land/lipgloss/v2"
	"deutsch-tui/internal/core"
)

func TestRenderDecksBug(t *testing.T) {
	decks := []core.Deck{
		{
			ID:          "a1-travel",
			Name:        "A1 Travel & Transport",
			Description: "Essential A1 vocabulary for traveling and transportation.",
			Tags:        []string{"german", "a1", "travel"},
		},
		{
			ID:          "a2-transport-directions",
			Name:        "A2 Transport & Directions",
			Description: "Imported from Anki TSV.",
		},
	}

	model := NewModel(&mockRepo{}, &mockScheduler{})
	model.decks = decks
	model.activeView = ViewDecks
	model.deckFilter = ""
	model.width = 120
	model.height = 60
	model.breakpoint = BreakpointWide

	layout := model.activeViewContentLayout()
	view := model.renderDecks(layout)
	t.Logf("layout width: %d, height: %d", layout.Width, layout.Height)
	lines := strings.Split(view, "\n")
	for i, l := range lines {
		t.Logf("Line %d: %q", i, l)
	}
}

func TestRenderSettingsScrollbarAlignment(t *testing.T) {
	m := NewModel(&mockRepo{}, &mockScheduler{})
	m.activeView = ViewSettings
	m.width = 60 // Small width to force truncation/wrapping
	m.height = 20
	m.breakpoint = BreakpointCompact

	// Add a very long AI template
	m.aiTemplates = map[string]map[string]string{
		"custom": {
			"front":   "This is a very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very long template",
			"back":    "short",
			"example": "short",
		},
	}
	m.aiTemplateSets = []string{"custom"}
	m.aiTemplateIndex = 0

	// Force scrollbar by having many items or small height
	m.height = 10

	view := m.renderSettings(0, 0)
	lines := strings.Split(view, "\n")

	// Check if scrollbar character '█' or '│' is at the same column for all lines that have it
	scrollbarCol := -1
	for i, l := range lines {
		plain := stripANSI(l)
		if strings.Contains(l, "█") || strings.Contains(l, "│") {
			idx := strings.LastIndex(plain, "█")
			if idx == -1 {
				idx = strings.LastIndex(plain, "│")
			}
			visualCol := lipgloss.Width(plain[:idx])
			if scrollbarCol == -1 {
				scrollbarCol = visualCol
			} else if visualCol != scrollbarCol {
				t.Errorf("Line %d: scrollbar at visual col %d, expected %d. Plain line: %q", i, visualCol, scrollbarCol, plain)
			}
		}
	}
}

func TestRenderStatisticsScrollbarAlignment(t *testing.T) {
	m := NewModel(&mockRepo{}, &mockScheduler{})
	m.activeView = ViewStatistics
	m.width = 60
	m.height = 20
	m.breakpoint = BreakpointCompact

	// Fill with some data
	m.stats.TotalCards = 1000
	m.stats.TotalDecks = 10
	m.decks = make([]core.Deck, 10)
	for i := 0; i < 10; i++ {
		m.decks[i] = core.Deck{
			ID:          fmt.Sprintf("deck-%d", i),
			Name:        fmt.Sprintf("Deck %d with a very very very very very very very very very long name", i),
			SuccessRate: 0.85,
		}
	}

	layout := m.activeViewContentLayout()
	view := m.renderStatistics(layout)
	lines := strings.Split(view, "\n")

	scrollbarCol := -1
	for i, l := range lines {
		plain := stripANSI(l)
		if strings.Contains(l, "█") || strings.Contains(l, "│") {
			idx := strings.LastIndex(plain, "█")
			if idx == -1 {
				idx = strings.LastIndex(plain, "│")
			}
			visualCol := lipgloss.Width(plain[:idx])
			if scrollbarCol == -1 {
				scrollbarCol = visualCol
			} else if visualCol != scrollbarCol {
				t.Errorf("Line %d: scrollbar at visual col %d, expected %d. Plain line: %q", i, visualCol, scrollbarCol, plain)
			}
		}
	}
}
