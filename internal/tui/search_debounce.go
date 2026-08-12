package tui

import (
	"time"

	tea "charm.land/bubbletea/v2"
)

const searchDebounceDelay = 250 * time.Millisecond

// debounceSearch invalidates the previous pending search for a view and
// schedules a new one. Browser text/tag filters intentionally share a
// generation because either field changes the same result set.
func (m *Model) debounceSearch(view View) tea.Cmd {
	var id int
	switch view {
	case ViewBrowser:
		m.browserSearchTimerID++
		id = m.browserSearchTimerID
	case ViewDictionary:
		m.dictionarySearchTimerID++
		id = m.dictionarySearchTimerID
	default:
		return nil
	}

	return tea.Tick(searchDebounceDelay, func(time.Time) tea.Msg {
		return debounceSearchMsg{id: id, view: view}
	})
}
