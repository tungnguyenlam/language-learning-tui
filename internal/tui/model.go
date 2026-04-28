package tui

import (
	"fmt"
	"strings"

	"deutsch-tui/internal/content"
	"deutsch-tui/internal/core"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type View string

const (
	ViewDashboard View = "dashboard"
	ViewReview    View = "review"
	ViewImport    View = "import"
	ViewAI        View = "ai"
	ViewSettings  View = "settings"
)

type Breakpoint string

const (
	BreakpointCompact Breakpoint = "compact"
	BreakpointMedium  Breakpoint = "medium"
	BreakpointWide    Breakpoint = "wide"
)

type Hitbox struct {
	ID     string
	View   View
	X      int
	Y      int
	Width  int
	Height int
}

func (h Hitbox) Contains(x, y int) bool {
	return x >= h.X && x < h.X+h.Width && y >= h.Y && y < h.Y+h.Height
}

type Model struct {
	width      int
	height     int
	activeView View
	breakpoint Breakpoint
	deck       core.Deck
	dueCards   []core.Card
	cursor     int
	revealed   bool
	status     string
	mouseX     int
	mouseY     int
	hitboxes   []Hitbox
}

func NewModel() *Model {
	deck := content.StarterDeck()
	var cards []core.Card
	for _, note := range deck.Notes {
		cards = append(cards, note.Cards...)
	}
	return &Model{
		width:      80,
		height:     24,
		activeView: ViewDashboard,
		breakpoint: BreakpointMedium,
		deck:       deck,
		dueCards:   cards,
		status:     "Ready",
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.breakpoint = breakpointForWidth(msg.Width)
	case tea.KeyMsg:
		return m.updateKey(msg)
	case tea.MouseMsg:
		mouse := msg.Mouse()
		m.mouseX = mouse.X
		m.mouseY = mouse.Y
		if hit, ok := m.hitboxAt(mouse.X, mouse.Y); ok {
			m.activateHitbox(hit.ID)
		}
	}
	return m, nil
}

func (m *Model) updateKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit
	case "tab":
		m.nextView()
	case "1":
		m.activeView = ViewDashboard
	case "2":
		m.activeView = ViewReview
	case "3":
		m.activeView = ViewImport
	case "4":
		m.activeView = ViewAI
	case "5":
		m.activeView = ViewSettings
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < len(m.dueCards)-1 {
			m.cursor++
		}
	case "enter", " ":
		if m.activeView == ViewReview {
			m.revealed = !m.revealed
		}
	case "a":
		m.status = "Grade: Again"
	case "h":
		m.status = "Grade: Hard"
	case "g":
		m.status = "Grade: Good"
	case "e":
		m.status = "Grade: Easy"
	}
	return m, nil
}

func (m *Model) View() tea.View {
	view := tea.NewView(m.render())
	view.AltScreen = true
	view.MouseMode = tea.MouseModeAllMotion
	return view
}

func (m *Model) render() string {
	m.hitboxes = nil
	var b strings.Builder

	header := headerStyle.Render("deutsch-tui") + " " + mutedStyle.Render(string(m.breakpoint))
	b.WriteString(header)
	b.WriteString("\n")

	switch m.breakpoint {
	case BreakpointWide:
		b.WriteString(m.renderWide())
	case BreakpointMedium:
		b.WriteString(m.renderMedium())
	default:
		b.WriteString(m.renderCompact())
	}

	footer := fmt.Sprintf("tab switch | 1-5 views | q quit | mouse %d,%d | %s", m.mouseX, m.mouseY, m.status)
	b.WriteString("\n")
	b.WriteString(statusStyle.Width(maxInt(20, m.width)).Render(footer))
	return b.String()
}

func (m *Model) renderWide() string {
	sidebar := m.renderNav(0, 2)
	main := m.renderActiveView()
	detail := panelStyle.Width(28).Render("Deck\n" + m.deck.Name + "\n\nCards due\n" + fmt.Sprint(len(m.dueCards)))
	return lipgloss.JoinHorizontal(lipgloss.Top, sidebar, main, detail)
}

func (m *Model) renderMedium() string {
	return m.renderTabs(0, 2) + "\n" + m.renderActiveView()
}

func (m *Model) renderCompact() string {
	return m.renderTabs(0, 2) + "\n" + compactStyle.Render(m.renderActiveViewPlain())
}

func (m *Model) renderNav(x, y int) string {
	labels := []struct {
		id   string
		view View
		text string
	}{
		{"nav-dashboard", ViewDashboard, "1 Dashboard"},
		{"nav-review", ViewReview, "2 Review"},
		{"nav-import", ViewImport, "3 Import"},
		{"nav-ai", ViewAI, "4 AI Drafts"},
		{"nav-settings", ViewSettings, "5 Settings"},
	}
	lines := make([]string, 0, len(labels))
	for i, label := range labels {
		m.hitboxes = append(m.hitboxes, Hitbox{ID: label.id, View: label.view, X: x, Y: y + i, Width: 18, Height: 1})
		style := navStyle
		if m.activeView == label.view {
			style = navActiveStyle
		}
		lines = append(lines, style.Render(label.text))
	}
	return panelStyle.Width(20).Render(strings.Join(lines, "\n"))
}

func (m *Model) renderTabs(x, y int) string {
	views := []View{ViewDashboard, ViewReview, ViewImport, ViewAI, ViewSettings}
	labels := []string{"Dashboard", "Review", "Import", "AI", "Settings"}
	parts := make([]string, 0, len(views))
	offset := x
	for i, view := range views {
		id := "tab-" + string(view)
		m.hitboxes = append(m.hitboxes, Hitbox{ID: id, View: view, X: offset, Y: y, Width: len(labels[i]) + 2, Height: 1})
		style := tabStyle
		if m.activeView == view {
			style = tabActiveStyle
		}
		rendered := style.Render(labels[i])
		parts = append(parts, rendered)
		offset += len(labels[i]) + 3
	}
	return strings.Join(parts, " ")
}

func (m *Model) renderActiveView() string {
	width := maxInt(30, m.width-54)
	if m.breakpoint == BreakpointMedium {
		width = maxInt(30, m.width-4)
	}
	return panelStyle.Width(width).Render(m.renderActiveViewPlain())
}

func (m *Model) renderActiveViewPlain() string {
	switch m.activeView {
	case ViewReview:
		return m.renderReview()
	case ViewImport:
		return "Import / Export\n\nAnki-friendly TSV roundtrip is available in internal/content.\n.apkg support is a later milestone."
	case ViewAI:
		return "AI Drafts\n\nProvider adapters generate draft notes only. Drafts require validation and approval."
	case ViewSettings:
		return "Settings\n\nLocal-first data directory, keyboard bindings, and provider settings will live here."
	default:
		return fmt.Sprintf("Dashboard\n\nDeck: %s\nDue cards: %d\n\nUse Review to start studying.", m.deck.Name, len(m.dueCards))
	}
}

func (m *Model) renderReview() string {
	if len(m.dueCards) == 0 {
		return "Review\n\nNo cards due."
	}
	card := m.dueCards[m.cursor]
	answer := "Press space or enter to reveal."
	if m.revealed {
		answer = card.Answer + "\n\nGrade: a Again | h Hard | g Good | e Easy"
	}
	return fmt.Sprintf("Review %d/%d\n\n%s\n\n%s", m.cursor+1, len(m.dueCards), card.Prompt, answer)
}

func (m *Model) nextView() {
	views := []View{ViewDashboard, ViewReview, ViewImport, ViewAI, ViewSettings}
	for i, view := range views {
		if m.activeView == view {
			m.activeView = views[(i+1)%len(views)]
			return
		}
	}
	m.activeView = ViewDashboard
}

func (m *Model) hitboxAt(x, y int) (Hitbox, bool) {
	for _, hitbox := range m.hitboxes {
		if hitbox.Contains(x, y) {
			return hitbox, true
		}
	}
	return Hitbox{}, false
}

func (m *Model) activateHitbox(id string) {
	switch {
	case strings.HasPrefix(id, "nav-"):
		m.activeView = View(strings.TrimPrefix(id, "nav-"))
	case strings.HasPrefix(id, "tab-"):
		m.activeView = View(strings.TrimPrefix(id, "tab-"))
	}
}

func breakpointForWidth(width int) Breakpoint {
	if width >= 100 {
		return BreakpointWide
	}
	if width >= 70 {
		return BreakpointMedium
	}
	return BreakpointCompact
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

var (
	headerStyle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("229"))
	mutedStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
	panelStyle     = lipgloss.NewStyle().Padding(1, 2).Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("240"))
	compactStyle   = lipgloss.NewStyle().Padding(1, 0)
	navStyle       = lipgloss.NewStyle().PaddingRight(2)
	navActiveStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("212")).PaddingRight(2)
	tabStyle       = lipgloss.NewStyle().Padding(0, 1).Foreground(lipgloss.Color("250"))
	tabActiveStyle = lipgloss.NewStyle().Padding(0, 1).Bold(true).Foreground(lipgloss.Color("229")).Background(lipgloss.Color("57"))
	statusStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("250")).Background(lipgloss.Color("236")).Padding(0, 1)
)
