package tui

import (
	tea "charm.land/bubbletea/v2"
)

// A screen is a self-contained view: it renders itself and handles its own key
// presses. It receives *Model for shared services (repo, navigation, the
// hitbox list, the status line, session counters) while its own view-local
// state lives on the concrete screen type rather than on Model.
//
// Migrating views to this interface is how the oversized Model struct is
// unwound incrementally: each migrated view moves its render + key logic here
// and its private fields off Model, one view per change. renderActiveViewPlainAt
// and updateActiveViewKey consult the registry first and fall back to the
// legacy per-view *Model methods for views not yet migrated.
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
		ViewSessionSummary: summaryScreen{},
		ViewImport:         m.importScreen,
		ViewAnkiWeb:        m.ankiWebScreen,
	}
}
