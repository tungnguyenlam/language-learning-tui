package tui

import (
	"testing"

	tea "charm.land/bubbletea/v2"
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
