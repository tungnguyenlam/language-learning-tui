package tui

import (
	"fmt"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"deutsch-tui/internal/core"
)

func TestStatisticsScreenScrollKeysUseBoundedViewport(t *testing.T) {
	model := NewModel(&mockRepo{}, &mockScheduler{})
	model.activeView = ViewStatistics
	model.width = 100
	model.height = 30
	model.statsTotalLines = 100
	model.statsScroll = 20

	for _, test := range []struct {
		key  tea.KeyPressMsg
		want int
	}{
		{key: tea.KeyPressMsg{Code: tea.KeyDown}, want: 21},
		{key: tea.KeyPressMsg{Code: tea.KeyPgDown}, want: 31},
		{key: tea.KeyPressMsg{Code: tea.KeyPgUp}, want: 21},
		{key: tea.KeyPressMsg{Code: tea.KeyUp}, want: 20},
	} {
		if _, handled := (statisticsScreen{}).HandleKey(model, test.key); !handled {
			t.Fatalf("key %q was not handled", test.key.String())
		}
		if model.statsScroll != test.want {
			t.Fatalf("statsScroll after %q = %d, want %d", test.key.String(), model.statsScroll, test.want)
		}
	}
}

func TestStatisticsRenderHonorsNonzeroScrollOffset(t *testing.T) {
	model := NewModel(&mockRepo{}, &mockScheduler{})
	model.activeView = ViewStatistics
	model.stats.TotalCards = 1_500
	model.stats.TotalDecks = 15
	model.stats.DailyGoal = 10
	model.decks = make([]core.Deck, 15)
	for i := range model.decks {
		model.decks[i] = core.Deck{
			ID:          fmt.Sprintf("deck-%02d", i),
			Name:        fmt.Sprintf("Target Deck %02d", i),
			SuccessRate: 0.85,
		}
	}

	full := strings.Split(stripANSI(model.renderStatistics(viewportLayout{Width: 100, Height: 200})), "\n")
	target := "Target Deck 12"
	targetLine := -1
	for i, line := range full {
		if strings.Contains(line, target) {
			targetLine = i
			break
		}
	}
	if targetLine < 2 {
		t.Fatalf("target row %q missing from full Statistics render", target)
	}

	// The title and blank line live outside the scrollable RenderList content.
	model.statsScroll = targetLine - 2
	scrolled := stripANSI(model.renderStatistics(viewportLayout{Width: 100, Height: 8}))
	if !strings.Contains(scrolled, target) {
		t.Fatalf("Statistics at scroll %d did not render %q:\n%s", model.statsScroll, target, scrolled)
	}
}

func TestStatisticsScreenExportShortcutStartsExport(t *testing.T) {
	model := NewModel(&mockRepo{}, &mockScheduler{})

	cmd, handled := (statisticsScreen{}).HandleKey(model, tea.KeyPressMsg{Text: "x"})
	if !handled {
		t.Fatal("statistics export shortcut was not handled")
	}
	if cmd == nil {
		t.Fatal("statistics export shortcut returned no command")
	}
	if model.status != "Exporting deck statistics..." {
		t.Fatalf("status = %q, want export progress", model.status)
	}
}
