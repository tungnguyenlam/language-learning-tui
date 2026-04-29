package tui

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"deutsch-tui/internal/ai"
	"deutsch-tui/internal/content"
	"deutsch-tui/internal/core"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type View string

const (
	ViewDashboard View = "dashboard"
	ViewDecks     View = "decks"
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
	repo        core.Repository
	scheduler   core.Scheduler
	width       int
	height      int
	activeView  View
	breakpoint  Breakpoint
	decks       []core.Deck
	deckIndex   int
	deck        core.Deck
	deckCursor  int
	allDue      []core.Card
	dueCards    []core.Card
	cursor      int
	revealed    bool
	status      string
	mouseX      int
	mouseY      int
	hitboxes    []Hitbox
	aiProvider  ai.Provider
	aiProviderName string
	aiTemplates    map[string]string
	settingsCursor int
	editingTemplate bool
	aiInput     string
	drafts      []ai.Draft
	draftCursor int
	importPath  string
	exportPath  string
}

func NewModel(repo core.Repository, scheduler core.Scheduler) *Model {
	return NewModelWithAI(repo, scheduler, ai.OfflineProvider{})
}

func NewModelWithAI(repo core.Repository, scheduler core.Scheduler, provider ai.Provider) *Model {
	return NewModelWithOptions(repo, scheduler, ModelOptions{
		AIProvider: provider,
		AIProviderName: "offline",
	})
}

type ModelOptions struct {
	AIProvider     ai.Provider
	AIProviderName string
	AITemplates    map[string]string
	ImportPath     string
	ExportPath     string
}

func NewModelWithOptions(repo core.Repository, scheduler core.Scheduler, opts ModelOptions) *Model {
	providerName := opts.AIProviderName
	if providerName == "" {
		providerName = "offline"
	}
	templates := opts.AITemplates
	if templates == nil {
		templates = map[string]string{
			"front":   "{{.Topic}}",
			"back":    "German prompt for {{.Topic}}",
			"example": "Practice sentence using {{.Topic}}.",
		}
	}
	provider := opts.AIProvider
	if provider == nil {
		switch providerName {
		case "template":
			provider = ai.TemplateProvider{Templates: templates}
		default:
			provider = ai.OfflineProvider{}
		}
	}
	importPath := opts.ImportPath
	if strings.TrimSpace(importPath) == "" {
		importPath = "import.tsv"
	}
	exportPath := opts.ExportPath
	if strings.TrimSpace(exportPath) == "" {
		exportPath = "export.tsv"
	}
	return &Model{
		repo:           repo,
		scheduler:      scheduler,
		aiProvider:     provider,
		aiProviderName: providerName,
		aiTemplates:    templates,
		width:          80,
		height:         24,
		activeView:     ViewDashboard,
		breakpoint:     BreakpointMedium,
		status:         "Ready",
		aiInput:        "der Kaffee",
		importPath:     filepath.Clean(importPath),
		exportPath:     filepath.Clean(exportPath),
	}
}

type dueCardsMsg []core.Card
type decksMsg []core.Deck
type draftsMsg []ai.Draft
type draftApprovedMsg struct {
	noteID string
	cards  []core.Card
}
type importDoneMsg struct {
	decks []core.Deck
	cards []core.Card
	count int
	path  string
}
type exportDoneMsg struct {
	count int
	path  string
}

func (m *Model) Init() tea.Cmd {
	return tea.Batch(m.loadDueCards, m.loadDecks)
}

func (m *Model) loadDueCards() tea.Msg {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	cards, err := m.repo.DueCards(ctx, time.Now(), 50)
	if err != nil {
		return err
	}
	return dueCardsMsg(cards)
}

func (m *Model) loadDecks() tea.Msg {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	decks, err := m.repo.Decks(ctx)
	if err != nil {
		return err
	}
	return decksMsg(decks)
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.breakpoint = breakpointForWidth(msg.Width)
	case error:
		m.status = fmt.Sprintf("Error: %v", msg)
	case decksMsg:
		m.decks = []core.Deck(msg)
		if m.deckIndex >= len(m.decks) {
			m.deckIndex = maxInt(0, len(m.decks)-1)
		}
		m.selectDeck(m.deckIndex)
	case dueCardsMsg:
		m.allDue = []core.Card(msg)
		m.applyDeckFilter()
	case draftsMsg:
		m.drafts = []ai.Draft(msg)
		m.draftCursor = 0
		if len(m.drafts) == 0 {
			m.status = "No drafts generated"
		} else {
			m.status = fmt.Sprintf("%d draft ready", len(m.drafts))
		}
	case draftApprovedMsg:
		m.removeDraft(msg.noteID)
		m.allDue = msg.cards
		m.applyDeckFilter()
		m.status = "Draft approved"
	case importDoneMsg:
		m.decks = msg.decks
		if len(m.decks) > 0 {
			m.deckIndex = maxInt(0, m.deckIndex)
			if m.deckIndex >= len(m.decks) {
				m.deckIndex = len(m.decks) - 1
			}
			m.deck = m.decks[m.deckIndex]
		}
		m.allDue = msg.cards
		m.applyDeckFilter()
		m.status = fmt.Sprintf("Imported %d notes from %s", msg.count, filepath.Base(msg.path))
	case exportDoneMsg:
		m.status = fmt.Sprintf("Exported %d notes to %s", msg.count, filepath.Base(msg.path))
	case tea.KeyMsg:
		return m.updateKey(msg)
	case tea.MouseMsg:
		mouse := msg.Mouse()
		m.mouseX = mouse.X
		m.mouseY = mouse.Y
		if hit, ok := m.hitboxAt(mouse.X, mouse.Y); ok {
			cmd := m.activateHitbox(hit.ID)
			return m, cmd
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
		m.activeView = ViewDecks
	case "3":
		m.activeView = ViewReview
	case "4":
		m.activeView = ViewImport
	case "5":
		m.activeView = ViewAI
	case "6":
		m.activeView = ViewSettings
	case "left", "[":
		m.previousDeck()
	case "right", "]":
		m.nextDeck()
	default:
		if m.activeView == ViewAI {
			cmd, handled := m.updateAIKey(msg)
			if handled {
				return m, cmd
			}
		}
		if m.activeView == ViewImport {
			cmd, handled := m.updateImportKey(msg)
			if handled {
				return m, cmd
			}
		}
		if m.activeView == ViewDecks {
			cmd, handled := m.updateDecksKey(msg)
			if handled {
				return m, cmd
			}
		}
		if m.activeView == ViewSettings {
			cmd, handled := m.updateSettingsKey(msg)
			if handled {
				return m, cmd
			}
		}
	}

	if m.activeView == ViewAI || m.activeView == ViewImport || m.activeView == ViewDecks || m.activeView == ViewSettings {
		return m, nil
	}

	switch msg.String() {
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < len(m.dueCards)-1 {
			m.cursor++
		}
	case "enter", "space":
		if m.activeView == ViewReview {
			m.revealed = !m.revealed
		}
	case "a":
		return m, m.gradeCard(core.GradeAgain)
	case "h":
		return m, m.gradeCard(core.GradeHard)
	case "g":
		return m, m.gradeCard(core.GradeGood)
	case "e":
		return m, m.gradeCard(core.GradeEasy)
	}
	return m, nil
}

func (m *Model) updateDecksKey(msg tea.KeyMsg) (tea.Cmd, bool) {
	switch msg.String() {
	case "up", "k":
		if m.deckCursor > 0 {
			m.deckCursor--
		}
		return nil, true
	case "down", "j":
		if m.deckCursor < len(m.decks)-1 {
			m.deckCursor++
		}
		return nil, true
	case "enter":
		m.selectDeck(m.deckCursor)
		m.activeView = ViewDashboard
		return nil, true
	}
	return nil, false
}

func (m *Model) updateImportKey(msg tea.KeyMsg) (tea.Cmd, bool) {
	switch msg.String() {
	case "i":
		return m.importTSV(), true
	case "x":
		return m.exportTSV(), true
	}
	return nil, false
}

func (m *Model) updateSettingsKey(msg tea.KeyMsg) (tea.Cmd, bool) {
	if m.editingTemplate {
		switch msg.String() {
		case "enter", "esc":
			m.editingTemplate = false
			if m.aiProviderName == "template" {
				m.aiProvider = ai.TemplateProvider{Templates: m.aiTemplates}
			}
			return nil, true
		case "backspace":
			key := m.templateKeyAtCursor()
			if val := m.aiTemplates[key]; len(val) > 0 {
				m.aiTemplates[key] = val[:len(val)-1]
			}
			return nil, true
		}
		if len(msg.String()) == 1 {
			key := m.templateKeyAtCursor()
			m.aiTemplates[key] += msg.String()
			return nil, true
		}
		return nil, true
	}

	switch msg.String() {
	case "up", "k":
		if m.settingsCursor > 0 {
			m.settingsCursor--
		}
		return nil, true
	case "down", "j":
		if m.settingsCursor < 3 {
			m.settingsCursor++
		}
		return nil, true
	case "enter":
		if m.settingsCursor == 0 {
			if m.aiProviderName == "offline" {
				m.aiProviderName = "template"
				m.aiProvider = ai.TemplateProvider{Templates: m.aiTemplates}
			} else {
				m.aiProviderName = "offline"
				m.aiProvider = ai.OfflineProvider{}
			}
			m.status = fmt.Sprintf("Switched to %s AI provider", m.aiProviderName)
			return nil, true
		} else {
			m.editingTemplate = true
			return nil, true
		}
	}
	return nil, false
}

func (m *Model) templateKeyAtCursor() string {
	switch m.settingsCursor {
	case 1:
		return "front"
	case 2:
		return "back"
	case 3:
		return "example"
	default:
		return ""
	}
}

func (m *Model) renderSettings(x, y int) string {
	var b strings.Builder
	b.WriteString("Settings\n\n")

	options := []string{
		fmt.Sprintf("AI Provider: %s", m.aiProviderName),
		fmt.Sprintf("Front Template: %s", m.aiTemplates["front"]),
		fmt.Sprintf("Back Template: %s", m.aiTemplates["back"]),
		fmt.Sprintf("Example Template: %s", m.aiTemplates["example"]),
	}

	for i, opt := range options {
		prefix := "  "
		style := lipgloss.NewStyle()
		if i == m.settingsCursor {
			prefix = "> "
			if m.editingTemplate {
				style = style.Bold(true).Background(lipgloss.Color("62"))
			} else {
				style = style.Bold(true).Foreground(lipgloss.Color("212"))
			}
		}
		b.WriteString(style.Render(prefix+opt) + "\n")
	}

	if m.editingTemplate {
		b.WriteString("\nEDITING - Enter to save, Esc to cancel.")
	} else {
		b.WriteString("\nUse j/k to move, Enter to toggle provider or edit template.")
	}
	return b.String()
}

func (m *Model) updateAIKey(msg tea.KeyMsg) (tea.Cmd, bool) {
	switch msg.String() {
	case "enter":
		return m.generateDrafts(), true
	case "backspace":
		if len(m.aiInput) > 0 {
			m.aiInput = m.aiInput[:len(m.aiInput)-1]
		}
		return nil, true
	case "up":
		if m.draftCursor > 0 {
			m.draftCursor--
		}
		return nil, true
	case "down":
		if m.draftCursor < len(m.drafts)-1 {
			m.draftCursor++
		}
		return nil, true
	case "a":
		if len(m.drafts) > 0 {
			return m.approveDraft(), true
		}
	case "d":
		if len(m.drafts) > 0 {
			m.discardDraft()
			return nil, true
		}
	}
	if len(msg.String()) == 1 {
		m.aiInput += msg.String()
		return nil, true
	}
	return nil, false
}

func (m *Model) generateDrafts() tea.Cmd {
	if m.aiProvider == nil {
		m.status = "AI provider unavailable"
		return nil
	}
	if strings.TrimSpace(m.deck.ID) == "" {
		m.status = "Select a deck before drafting"
		return nil
	}
	request := ai.DraftRequest{
		SourceText: m.aiInput,
		DeckID:     m.deck.ID,
		Tags:       []string{"reviewed"},
	}
	m.status = "Generating draft..."
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		drafts, err := m.aiProvider.GenerateDrafts(ctx, request)
		if err != nil {
			return err
		}
		if err := ai.ValidateDrafts(drafts); err != nil {
			return err
		}
		return draftsMsg(drafts)
	}
}

func (m *Model) approveDraft() tea.Cmd {
	if len(m.drafts) == 0 || m.draftCursor >= len(m.drafts) {
		m.status = "No draft selected"
		return nil
	}
	draft := m.drafts[m.draftCursor]
	if err := ai.ValidateDraft(draft); err != nil {
		m.status = fmt.Sprintf("Error: %v", err)
		return nil
	}
	m.status = "Approving draft..."
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		deck := m.deck
		deck.Notes = []core.Note{draft.Note}
		if err := m.repo.UpsertDeck(ctx, deck); err != nil {
			return err
		}
		cards, err := m.repo.DueCards(ctx, time.Now(), 50)
		if err != nil {
			return err
		}
		return draftApprovedMsg{noteID: draft.Note.ID, cards: cards}
	}
}

func (m *Model) importTSV() tea.Cmd {
	path := m.importPath
	m.status = "Importing TSV..."
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		notes, err := content.ImportAnkiTSV(file, content.ImportOptions{DefaultDeck: m.deck.ID})
		if err != nil {
			return err
		}
		decks := decksFromNotes(notes)
		for _, deck := range decks {
			if err := m.repo.UpsertDeck(ctx, deck); err != nil {
				return err
			}
		}
		loadedDecks, err := m.repo.Decks(ctx)
		if err != nil {
			return err
		}
		cards, err := m.repo.DueCards(ctx, time.Now(), 50)
		if err != nil {
			return err
		}
		return importDoneMsg{decks: loadedDecks, cards: cards, count: len(notes), path: path}
	}
}

func (m *Model) exportTSV() tea.Cmd {
	deckID := m.deck.ID
	path := m.exportPath
	if strings.TrimSpace(deckID) == "" {
		m.status = "Select a deck before exporting"
		return nil
	}
	m.status = "Exporting TSV..."
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		deck, err := m.repo.GetDeck(ctx, deckID)
		if err != nil {
			return err
		}
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return err
		}
		file, err := os.Create(path)
		if err != nil {
			return err
		}
		defer file.Close()
		if err := content.ExportAnkiTSV(file, deck.Notes); err != nil {
			return err
		}
		return exportDoneMsg{count: len(deck.Notes), path: path}
	}
}

func (m *Model) discardDraft() {
	if len(m.drafts) == 0 || m.draftCursor >= len(m.drafts) {
		m.status = "No draft selected"
		return
	}
	m.drafts = append(m.drafts[:m.draftCursor], m.drafts[m.draftCursor+1:]...)
	if m.draftCursor >= len(m.drafts) {
		m.draftCursor = maxInt(0, len(m.drafts)-1)
	}
	m.status = "Draft discarded"
}

func (m *Model) removeDraft(noteID string) {
	for i, draft := range m.drafts {
		if draft.Note.ID == noteID {
			m.drafts = append(m.drafts[:i], m.drafts[i+1:]...)
			break
		}
	}
	if m.draftCursor >= len(m.drafts) {
		m.draftCursor = maxInt(0, len(m.drafts)-1)
	}
}

func (m *Model) gradeCard(grade core.ReviewGrade) tea.Cmd {
	if m.activeView != ViewReview || len(m.dueCards) == 0 || !m.revealed {
		return nil
	}

	m.status = fmt.Sprintf("Grade: %s", strings.ToUpper(string(grade[:1]))+string(grade[1:]))
	card := m.dueCards[m.cursor]
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		now := time.Now()
		state, err := m.repo.GetReviewState(ctx, card.ID)
		if err != nil {
			return err
		}

		next, err := m.scheduler.Review(state, grade, now)
		if err != nil {
			return err
		}

		result := core.ReviewResult{
			CardID:   card.ID,
			Grade:    grade,
			Reviewed: now,
			Next:     next,
		}

		if err := m.repo.RecordReview(ctx, result); err != nil {
			return err
		}

		return m.loadDueCards()
	}
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
	main := m.renderActiveView(20, 2)
	detail := panelStyle.Width(28).Render("Deck\n" + m.deckLabel() + "\n\nCards due\n" + fmt.Sprint(len(m.dueCards)) + "\n\n[ ] switch deck")
	return lipgloss.JoinHorizontal(lipgloss.Top, sidebar, main, detail)
}

func (m *Model) renderMedium() string {
	return m.renderTabs(0, 2) + "\n" + m.renderActiveView(0, 3)
}

func (m *Model) renderCompact() string {
	return m.renderTabs(0, 2) + "\n" + compactStyle.Render(m.renderActiveViewPlain(0, 3))
}

func (m *Model) renderNav(x, y int) string {
	labels := []struct {
		id   string
		view View
		text string
	}{
		{"nav-dashboard", ViewDashboard, "1 Dashboard"},
		{"nav-decks", ViewDecks, "2 Decks"},
		{"nav-review", ViewReview, "3 Review"},
		{"nav-import", ViewImport, "4 Import"},
		{"nav-ai", ViewAI, "5 AI Drafts"},
		{"nav-settings", ViewSettings, "6 Settings"},
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
	views := []View{ViewDashboard, ViewDecks, ViewReview, ViewImport, ViewAI, ViewSettings}
	labels := []string{"Dashboard", "Decks", "Review", "Import", "AI", "Settings"}
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

func (m *Model) renderActiveView(x, y int) string {
	width := maxInt(30, m.width-54)
	if m.breakpoint == BreakpointMedium {
		width = maxInt(30, m.width-4)
	}
	return panelStyle.Width(width).Render(m.renderActiveViewPlain(x+2, y+1))
}

func (m *Model) renderActiveViewPlain(x, y int) string {
	switch m.activeView {
	case ViewDecks:
		return m.renderDecks(x, y)
	case ViewReview:
		return m.renderReview(x, y)
	case ViewImport:
		return m.renderImport()
	case ViewAI:
		return m.renderAI(x, y)
	case ViewSettings:
		return m.renderSettings(x, y)
	default:
		return fmt.Sprintf("Dashboard\n\nDeck: %s\nDue cards: %d\n\nUse [ and ] to switch decks.\nUse Review to start studying.", m.deckLabel(), len(m.dueCards))
	}
}

func (m *Model) renderDecks(x, y int) string {
	var b strings.Builder
	b.WriteString("Decks\n\n")
	if len(m.decks) == 0 {
		b.WriteString("No decks found. Use Import to add notes.")
		return b.String()
	}
	for i, deck := range m.decks {
		prefix := "  "
		style := lipgloss.NewStyle()
		if i == m.deckCursor {
			prefix = "> "
			style = style.Bold(true).Foreground(lipgloss.Color("212"))
		}
		label := fmt.Sprintf("%s%s (%d total, %d due)", prefix, deck.Name, deck.TotalCards, deck.DueCards)
		b.WriteString(style.Render(label))
		b.WriteString("\n")
		if deck.Description != "" {
			b.WriteString(fmt.Sprintf("     %s\n", mutedStyle.Render(deck.Description)))
		}
	}
	b.WriteString("\nPress enter to select deck.")
	return b.String()
}

func (m *Model) renderImport() string {
	return fmt.Sprintf("Import / Export\n\nImport file: %s\nExport file: %s\n\nPress i to import TSV.\nPress x to export selected deck.\n\nDeck: %s\n.apkg support is a later milestone.", m.importPath, m.exportPath, m.deckLabel())
}

func (m *Model) renderAI(x, y int) string {
	var b strings.Builder
	fmt.Fprintf(&b, "AI Drafts\n\nDeck: %s\nTopic: %s\n\nEnter generate | a approve | d discard\n", m.deckLabel(), m.aiInput)
	if len(m.drafts) == 0 {
		b.WriteString("\nNo drafts yet.")
		return b.String()
	}
	for i, draft := range m.drafts {
		prefix := "  "
		if i == m.draftCursor {
			prefix = "> "
			m.hitboxes = append(m.hitboxes, Hitbox{ID: "draft-approve", View: ViewAI, X: x, Y: y + 7 + i, Width: 32, Height: 1})
		}
		fmt.Fprintf(&b, "\n%s%s -> %s", prefix, draft.Note.Front, draft.Note.Back)
	}
	return b.String()
}

func (m *Model) renderReview(x, y int) string {
	if len(m.dueCards) == 0 {
		return "Review\n\nNo cards due."
	}
	card := m.dueCards[m.cursor]
	answer := "Press space or enter to reveal."
	if m.revealed {
		answer = card.Answer + "\n\nGrade: a Again | h Hard | g Good | e Easy"
		// Add hitboxes for grades. These are approximate but should work for basic interaction.
		// "Again" is around offset 7, "Hard" around 17, "Good" around 26, "Easy" around 35
		m.hitboxes = append(m.hitboxes, Hitbox{ID: "grade-again", View: ViewReview, X: x + 7, Y: y + 6, Width: 5, Height: 1})
		m.hitboxes = append(m.hitboxes, Hitbox{ID: "grade-hard", View: ViewReview, X: x + 17, Y: y + 6, Width: 4, Height: 1})
		m.hitboxes = append(m.hitboxes, Hitbox{ID: "grade-good", View: ViewReview, X: x + 26, Y: y + 6, Width: 4, Height: 1})
		m.hitboxes = append(m.hitboxes, Hitbox{ID: "grade-easy", View: ViewReview, X: x + 35, Y: y + 6, Width: 4, Height: 1})
	}
	return fmt.Sprintf("Review %d/%d\n\n%s\n\n%s", m.cursor+1, len(m.dueCards), card.Prompt, answer)
}

func (m *Model) nextView() {
	views := []View{ViewDashboard, ViewDecks, ViewReview, ViewImport, ViewAI, ViewSettings}
	for i, view := range views {
		if m.activeView == view {
			m.activeView = views[(i+1)%len(views)]
			return
		}
	}
	m.activeView = ViewDashboard
}

func (m *Model) previousDeck() {
	if len(m.decks) == 0 {
		return
	}
	m.deckIndex = (m.deckIndex - 1 + len(m.decks)) % len(m.decks)
	m.selectDeck(m.deckIndex)
}

func (m *Model) nextDeck() {
	if len(m.decks) == 0 {
		return
	}
	m.deckIndex = (m.deckIndex + 1) % len(m.decks)
	m.selectDeck(m.deckIndex)
}

func (m *Model) selectDeck(index int) {
	if len(m.decks) == 0 {
		m.deck = core.Deck{}
		m.deckCursor = 0
		m.applyDeckFilter()
		return
	}
	if index < 0 || index >= len(m.decks) {
		index = 0
	}
	m.deckIndex = index
	m.deck = m.decks[index]
	m.deckCursor = index
	m.applyDeckFilter()
}

func (m *Model) applyDeckFilter() {
	m.dueCards = m.dueCards[:0]
	for _, card := range m.allDue {
		if m.deck.ID == "" || card.DeckID == m.deck.ID {
			m.dueCards = append(m.dueCards, card)
		}
	}
	if m.cursor >= len(m.dueCards) {
		m.cursor = maxInt(0, len(m.dueCards)-1)
	}
	m.revealed = false
	if len(m.dueCards) == 0 {
		m.status = "All caught up!"
	} else {
		m.status = fmt.Sprintf("%d cards due", len(m.dueCards))
	}
}

func (m *Model) deckLabel() string {
	if strings.TrimSpace(m.deck.Name) != "" {
		return m.deck.Name
	}
	return "No deck"
}

func (m *Model) hitboxAt(x, y int) (Hitbox, bool) {
	for _, hitbox := range m.hitboxes {
		if hitbox.Contains(x, y) {
			return hitbox, true
		}
	}
	return Hitbox{}, false
}

func (m *Model) activateHitbox(id string) tea.Cmd {
	switch {
	case strings.HasPrefix(id, "nav-"):
		m.activeView = View(strings.TrimPrefix(id, "nav-"))
	case strings.HasPrefix(id, "tab-"):
		m.activeView = View(strings.TrimPrefix(id, "tab-"))
	case id == "grade-again":
		return m.gradeCard(core.GradeAgain)
	case id == "grade-hard":
		return m.gradeCard(core.GradeHard)
	case id == "grade-good":
		return m.gradeCard(core.GradeGood)
	case id == "grade-easy":
		return m.gradeCard(core.GradeEasy)
	case id == "draft-approve":
		return m.approveDraft()
	}
	return nil
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

func decksFromNotes(notes []core.Note) []core.Deck {
	byID := make(map[string]int)
	decks := make([]core.Deck, 0, 1)
	for _, note := range notes {
		deckID := strings.TrimSpace(note.DeckID)
		if deckID == "" {
			deckID = "Imported"
			note.DeckID = deckID
		}
		index, ok := byID[deckID]
		if !ok {
			index = len(decks)
			byID[deckID] = index
			decks = append(decks, core.Deck{
				ID:          deckID,
				Name:        deckID,
				Description: "Imported from Anki TSV.",
			})
		}
		decks[index].Notes = append(decks[index].Notes, note)
	}
	return decks
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
