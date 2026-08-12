package tui

import (
	tea "charm.land/bubbletea/v2"
)

// A screen owns a view's rendering and key handling. It receives *Model for
// shared services and for legacy view state that has not yet moved onto a
// concrete screen type.
//
// Migrating views to this interface is how the oversized Model struct is
// unwound incrementally: first co-locate render + key logic, then move private
// state when its cross-view dependencies are understood. renderActiveViewPlainAt
// and updateActiveViewKey dispatch through this registry.
type screen interface {
	Render(m *Model, layout viewportLayout) string
	HandleKey(m *Model, msg tea.KeyPressMsg) (tea.Cmd, bool)
}

// registerScreens populates m.screens with the views migrated to the screen
// interface. Stateful screens are stored both in the map and via a typed handle
// on Model so shared code can reach their view-local state. Add migrated views
// here.
func (m *Model) registerScreens() {
	m.importScreen = &importScreen{}
	m.ankiWebScreen = &ankiWebScreen{}
	m.screens = map[View]screen{
		ViewDashboard:      dashboardScreen{},
		ViewDecks:          decksScreen{},
		ViewReview:         reviewScreen{},
		ViewBrowser:        browserScreen{},
		ViewAI:             aiScreen{},
		ViewSettings:       settingsScreen{},
		ViewCram:           cramScreen{},
		ViewPractice:       practiceScreen{},
		ViewDebug:          debugScreen{},
		ViewStatistics:     statisticsScreen{},
		ViewDictionary:     dictionaryScreen{},
		ViewSessionSummary: summaryScreen{},
		ViewImport:         m.importScreen,
		ViewAnkiWeb:        m.ankiWebScreen,
	}
}
