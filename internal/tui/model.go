package tui

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
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
	ViewDashboard  View = "dashboard"
	ViewDecks      View = "decks"
	ViewReview     View = "review"
	ViewStatistics View = "statistics"
	ViewImport     View = "import"
	ViewAI         View = "ai"
	ViewSettings   View = "settings"
	ViewBrowser    View = "browser"
	ViewCram       View = "cram"
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
	repo               core.Repository
	scheduler          core.Scheduler
	width              int
	height             int
	activeView         View
	breakpoint         Breakpoint
	decks              []core.Deck
	deckIndex          int
	deck               core.Deck
	deckCursor         int
	allDue             []core.Card
	dueCards           []core.Card
	cursor             int
	revealed           bool
	lastReviewedCardID string
	status             string
	mouseX             int
	mouseY             int
	hitboxes           []Hitbox
	aiProvider         ai.Provider
	aiProviderName     string
	aiTemplates        map[string]string
	autoPlayAudio      bool
	stats              core.Statistics
	settingsCursor     int
	editingTemplate    bool
	aiInput            string
	drafts             []ai.Draft
	draftCursor        int
	importPath         string
	exportPath         string
	importCursor       int // 0 for importPath, 1 for exportPath
	editingImportPath  bool
	onConfigChange     func(string, map[string]string, bool)
	bookmarkFilter     bool
	mcqChoice          int
	mcqAnswered        bool
	mcqCorrect         bool
	browserCards       []core.Card
	browserCursor      int
	browserSearch      string
	browserDeckID      string
	sessionReviewed    int
	sessionCorrect     int
	showHelp           bool
	cramCards          []core.Card
	cramCursor         int
	cramType           string // "bookmarked", "suspended", "leech", "flagged", or "all"
	cramReviewed       int
	cramCorrect        int
	cramActive         bool
	cramRevealed       bool
	reviewsPerDay      map[string]int
	reviewHistory      []core.ReviewLog
	reviewHistoryCard  string
	showReviewHistory  bool
	spinnerFrame       int
	deckFilter         string // New field for deck filtering
	drafting           bool
}

func NewModel(repo core.Repository, scheduler core.Scheduler) *Model {
	return NewModelWithAI(repo, scheduler, ai.OfflineProvider{})
}

func NewModelWithAI(repo core.Repository, scheduler core.Scheduler, provider ai.Provider) *Model {
	return NewModelWithOptions(repo, scheduler, ModelOptions{
		AIProvider:     provider,
		AIProviderName: "offline",
	})
}

type ModelOptions struct {
	AIProvider     ai.Provider
	AIProviderName string
	AITemplates    map[string]string
	AutoPlayAudio  bool
	ImportPath     string
	ExportPath     string
	OnConfigChange func(string, map[string]string, bool)
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
	autoPlayAudio := opts.AutoPlayAudio
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
		autoPlayAudio:  autoPlayAudio,
		width:          80,
		height:         24,
		activeView:     ViewDashboard,
		breakpoint:     BreakpointMedium,
		status:         "Ready",
		aiInput:        "der Kaffee",
		importPath:     filepath.Clean(importPath),
		exportPath:     filepath.Clean(exportPath),
		onConfigChange: opts.OnConfigChange,
	}
}

type dueCardsMsg []core.Card
type bookmarkedDueCardsMsg []core.Card
type decksMsg []core.Deck
type statsMsg core.Statistics
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
type reviewRecordedMsg struct {
	cardID string
	cards  []core.Card
	decks  []core.Deck
	stats  core.Statistics
	grade  core.ReviewGrade
}
type bookmarkToggledMsg struct {
	cardID     string
	bookmarked bool
}
type cardSuspendedMsg struct {
	cardID string
	cards  []core.Card
	decks  []core.Deck
	stats  core.Statistics
}
type dailyGoalSetMsg core.Statistics
type reviewUndoMsg struct {
	cardID string
	cards  []core.Card
	decks  []core.Deck
	stats  core.Statistics
}
type cramCardsMsg []core.Card
type browserCardsMsg []core.Card
type reviewsPerDayMsg map[string]int
type reviewHistoryMsg struct {
	cardID string
	logs   []core.ReviewLog
}

func (m *Model) Init() tea.Cmd {
	return tea.Batch(m.loadDueCards, m.loadDecks, m.loadStatistics(), m.loadReviewsPerDay())
}

func (m *Model) loadBrowserCards() tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		cards, err := m.repo.Cards(ctx, m.browserDeckID, m.browserSearch)
		if err != nil {
			return err
		}
		return browserCardsMsg(cards)
	}
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

func (m *Model) loadBookmarkedDueCards() tea.Msg {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	cards, err := m.repo.DueCardsBookmarked(ctx, time.Now(), 50)
	if err != nil {
		return err
	}
	return bookmarkedDueCardsMsg(cards)
}

func (m *Model) loadCramCards() tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		cards, err := m.repo.Cards(ctx, "", "")
		if err != nil {
			return err
		}
		return cramCardsMsg(cards)
	}
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

func (m *Model) loadStatistics() tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		stats, err := m.repo.Statistics(ctx)
		if err != nil {
			return err
		}
		return statsMsg(stats)
	}
}

func (m *Model) loadReviewsPerDay() tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		data, err := m.repo.ReviewsPerDay(ctx, 30)
		if err != nil {
			return err
		}
		return reviewsPerDayMsg(data)
	}
}

func (m *Model) loadReviewHistory(cardID string) tea.Cmd {
	if strings.TrimSpace(cardID) == "" {
		return nil
	}
	m.status = "Loading review history..."
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		logs, err := m.repo.ReviewHistory(ctx, cardID, 5)
		if err != nil {
			return err
		}
		return reviewHistoryMsg{cardID: cardID, logs: logs}
	}
}

type spinnerTickMsg struct{}

func (m *Model) tickSpinner() tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
		return spinnerTickMsg{}
	})
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case spinnerTickMsg:
		if m.drafting {
			m.spinnerFrame++
			return m, m.tickSpinner()
		}
		return m, nil
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.breakpoint = breakpointForWidth(msg.Width)
	case error:
		m.drafting = false
		m.status = friendlyError(msg)
	case decksMsg:
		m.decks = []core.Deck(msg)
		if m.deckIndex >= len(m.decks) {
			m.deckIndex = maxInt(0, len(m.decks)-1)
		}
		m.selectDeck(m.deckIndex)
	case dueCardsMsg:
		m.allDue = []core.Card(msg)
		m.applyDeckFilter()
	case bookmarkedDueCardsMsg:
		m.allDue = []core.Card(msg)
		m.applyDeckFilter()
		if len(m.allDue) == 0 {
			m.status = "No bookmarked cards due"
		} else {
			m.status = fmt.Sprintf("%d bookmarked cards due", len(m.allDue))
		}
	case statsMsg:
		m.stats = core.Statistics(msg)
	case reviewsPerDayMsg:
		m.reviewsPerDay = map[string]int(msg)
	case reviewHistoryMsg:
		if msg.cardID == m.reviewHistoryCard {
			m.reviewHistory = msg.logs
			m.showReviewHistory = true
			if len(msg.logs) == 0 {
				m.status = "No review history for card"
			} else {
				m.status = fmt.Sprintf("%d review history entries", len(msg.logs))
			}
		}
	case draftsMsg:
		m.drafting = false
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
			m.deckCursor = len(m.decks) - 1
		}
		m.allDue = msg.cards
		m.applyDeckFilter()
		m.status = fmt.Sprintf("Imported %d notes from %s", msg.count, filepath.Base(msg.path))
	case exportDoneMsg:
		m.status = fmt.Sprintf("Exported %d notes to %s", msg.count, filepath.Base(msg.path))
	case reviewRecordedMsg:
		m.lastReviewedCardID = msg.cardID
		m.decks = msg.decks
		m.stats = msg.stats
		m.allDue = msg.cards
		m.applyDeckFilter()
		// Track session stats
		m.sessionReviewed++
		if msg.grade != core.GradeAgain {
			m.sessionCorrect++
		}
		// Track cram stats if in cram view
		if m.activeView == ViewCram {
			m.cramReviewed++
			if msg.grade != core.GradeAgain {
				m.cramCorrect++
			}
		}
		return m, m.loadReviewsPerDay()
	case bookmarkToggledMsg:
		m.setCardBookmarkLocal(msg.cardID, msg.bookmarked)
		if msg.bookmarked {
			m.status = "Card bookmarked"
		} else {
			m.status = "Card bookmark removed"
		}
	case cardSuspendedMsg:
		m.decks = msg.decks
		m.stats = msg.stats
		m.allDue = msg.cards
		m.applyDeckFilter()
		m.status = "Card suspended"
		return m, m.loadReviewsPerDay()
	case dailyGoalSetMsg:
		m.stats = core.Statistics(msg)
		m.status = fmt.Sprintf("Daily goal set to %d", m.stats.DailyGoal)
	case reviewUndoMsg:
		m.lastReviewedCardID = ""
		m.decks = msg.decks
		m.stats = msg.stats
		m.allDue = msg.cards
		m.applyDeckFilter()
		m.status = "Last review undone"
		return m, m.loadReviewsPerDay()
	case browserCardsMsg:
		m.browserCards = []core.Card(msg)
		if len(m.browserCards) == 0 {
			m.status = "No cards found"
		} else {
			m.status = fmt.Sprintf("%d cards found", len(m.browserCards))
		}
	case cramCardsMsg:
		allCards := []core.Card(msg)
		// Filter cards based on cram type
		m.cramCards = m.cramCards[:0]
		for _, card := range allCards {
			switch m.cramType {
			case "bookmarked":
				if card.Bookmarked {
					m.cramCards = append(m.cramCards, card)
				}
			case "suspended":
				if card.Suspended {
					m.cramCards = append(m.cramCards, card)
				}
			case "leech":
				if card.Leech {
					m.cramCards = append(m.cramCards, card)
				}
			case "flagged":
				if card.Bookmarked || card.Suspended || card.Leech {
					m.cramCards = append(m.cramCards, card)
				}
			case "all":
				m.cramCards = append(m.cramCards, card)
			}
		}
		if len(m.cramCards) == 0 {
			m.status = "No cards found for this filter"
		} else {
			m.status = fmt.Sprintf("%d cards in cram mode", len(m.cramCards))
		}
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

func (m *Model) playAudio(audioPath string) tea.Cmd {
	if audioPath == "" {
		return nil
	}
	return func() tea.Msg {
		var cmd *exec.Cmd
		if runtime.GOOS == "darwin" {
			cmd = exec.Command("afplay", audioPath)
		} else {
			cmd = exec.Command("play", audioPath)
		}
		_ = cmd.Run()
		return nil
	}
}

func (m *Model) updateKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "q":
		if m.textInputActive() {
			break
		}
		if m.activeView == ViewCram && m.cramActive {
			m.cramActive = false
			return m, nil
		}
		return m, tea.Quit
	case "?":
		if m.textInputActive() {
			break
		}
		m.showHelp = !m.showHelp
		if m.showHelp {
			m.status = "Help overlay shown. Press ? to close."
		} else {
			m.status = "Help overlay closed."
		}
		return m, nil
	case "tab":
		if m.textInputActive() {
			break
		}
		return m, m.nextViewCmd()
	case "left", "shift+tab":
		if m.textInputActive() {
			break
		}
		return m, m.previousViewCmd()
	case "right":
		if m.textInputActive() {
			break
		}
		return m, m.nextViewCmd()
	case "w":
		if !m.textInputActive() {
			return m, m.previousViewCmd()
		}
	case "s":
		if !m.textInputActive() {
			return m, m.nextViewCmd()
		}
	case "up":
		if !m.activeViewHandlesVerticalNavigation() {
			return m, m.previousViewCmd()
		}
	case "down":
		if !m.activeViewHandlesVerticalNavigation() {
			return m, m.nextViewCmd()
		}
	case "[":
		if m.textInputActive() && m.activeView != ViewBrowser {
			break
		}
		if m.activeView == ViewBrowser {
			m.previousDeck()
			return m, m.reloadBrowserForSelectedDeck()
		}
		m.previousDeck()
		return m, nil
	case "]":
		if m.textInputActive() && m.activeView != ViewBrowser {
			break
		}
		if m.activeView == ViewBrowser {
			m.nextDeck()
			return m, m.reloadBrowserForSelectedDeck()
		}
		m.nextDeck()
		return m, nil
	}

	if cmd, handled := m.updateNumberKey(msg); handled {
		return m, cmd
	}

	// Handle view-specific keys after global navigation.
	if cmd, handled := m.updateActiveViewKey(msg); handled {
		return m, cmd
	}

	switch msg.String() {
	case "p":
		if m.activeView == ViewReview && len(m.dueCards) > 0 {
			return m, m.playAudio(m.dueCards[m.cursor].Audio)
		}
		if m.activeView == ViewCram && m.cramActive && len(m.cramCards) > 0 {
			return m, m.playAudio(m.cramCards[m.cramCursor].Audio)
		}
	default:
		// Handle remaining keys based on active view
		switch m.activeView {
		case ViewReview:
			switch msg.String() {
			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
					m.resetMCQState()
					m.clearReviewHistory()
				}
			case "down", "j":
				if m.cursor < len(m.dueCards)-1 {
					m.cursor++
					m.resetMCQState()
					m.clearReviewHistory()
				}
			case "enter", "space":
				card := m.dueCards[m.cursor]
				m.revealed = !m.revealed
				m.mcqChoice = -1
				m.mcqAnswered = false
				if m.revealed && m.autoPlayAudio && card.Audio != "" {
					m.status = "Auto-playing audio: " + card.Audio
					return m, m.playAudio(card.Audio)
				}
			case "b":
				return m, m.toggleBookmark()
			case "B":
				return m, m.toggleBookmarkFilter()
			case "x":
				return m, m.suspendCard()
			case "u":
				return m, m.undoLastReview()
			case "r":
				return m, m.toggleReviewHistory()
			case "a":
				return m, m.gradeCard(core.GradeAgain)
			case "h":
				return m, m.gradeCard(core.GradeHard)
			case "g":
				return m, m.gradeCard(core.GradeGood)
			case "e":
				return m, m.gradeCard(core.GradeEasy)
			}
		}
	}

	return m, nil
}

func (m *Model) textInputActive() bool {
	return (m.activeView == ViewImport && m.editingImportPath) ||
		(m.activeView == ViewSettings && m.editingTemplate)
}

func (m *Model) activeViewHandlesVerticalNavigation() bool {
	switch m.activeView {
	case ViewAI, ViewBrowser, ViewCram, ViewDecks, ViewReview, ViewSettings, ViewImport:
		return true
	default:
		return false
	}
}

func (m *Model) updateNumberKey(msg tea.KeyMsg) (tea.Cmd, bool) {
	key := msg.String()
	if m.activeView == ViewCram {
		return nil, false
	}
	if m.textInputActive() {
		return nil, false
	}
	if m.activeView == ViewReview && len(m.dueCards) > 0 && m.dueCards[m.cursor].Kind == core.CardKindMCQ && m.revealed && !m.mcqAnswered {
		switch key {
		case "1", "2", "3", "4":
			m.selectMCQChoice(key)
			return nil, true
		}
	}

	switch key {
	case "1":
		return m.updateView(ViewDashboard), true
	case "2":
		return m.updateView(ViewDecks), true
	case "3":
		return m.updateView(ViewReview), true
	case "4":
		return m.updateView(ViewStatistics), true
	case "5":
		return m.updateView(ViewImport), true
	case "6":
		return m.updateView(ViewAI), true
	case "7":
		return m.updateView(ViewSettings), true
	case "8":
		return m.updateView(ViewBrowser), true
	case "9":
		return m.updateView(ViewCram), true
	default:
		return nil, false
	}
}

func (m *Model) updateActiveViewKey(msg tea.KeyMsg) (tea.Cmd, bool) {
	switch m.activeView {
	case ViewAI:
		return m.updateAIKey(msg)
	case ViewImport:
		return m.updateImportKey(msg)
	case ViewBrowser:
		return m.updateBrowserKey(msg)
	case ViewSettings:
		return m.updateSettingsKey(msg)
	case ViewCram:
		return m.updateCramKey(msg)
	case ViewDecks:
		return m.updateDecksKey(msg)
	default:
		return nil, false
	}
}

func (m *Model) updateDecksKey(msg tea.KeyMsg) (tea.Cmd, bool) {
	switch msg.String() {
	case "up", "k":
		if m.deckCursor > 0 {
			m.deckCursor--
		}
		return nil, true
	case "down", "j":
		if m.deckCursor < len(m.filteredDecks())-1 {
			m.deckCursor++
		}
		return nil, true
	case "enter":
		if len(m.filteredDecks()) > 0 {
			m.selectDeck(m.deckCursor)
			m.activeView = ViewDashboard
		}
		return nil, true
	case "backspace":
		if len(m.deckFilter) > 0 {
			m.deckFilter = m.deckFilter[:len(m.deckFilter)-1]
			m.deckCursor = 0 // Reset cursor when filter changes
		}
		return nil, true
	case "esc":
		m.deckFilter = ""
		m.deckCursor = 0 // Reset cursor when filter is cleared
		return nil, true
	}

	// Handle text input for filtering
	if len(msg.String()) == 1 {
		m.deckFilter += msg.String()
		m.deckCursor = 0 // Reset cursor when filter changes
		return nil, true
	}

	return nil, false
}

func (m *Model) updateImportKey(msg tea.KeyMsg) (tea.Cmd, bool) {
	if m.editingImportPath {
		switch msg.String() {
		case "enter", "esc":
			m.editingImportPath = false
			return nil, true
		case "backspace":
			if m.importCursor == 0 {
				if len(m.importPath) > 0 {
					m.importPath = m.importPath[:len(m.importPath)-1]
				}
			} else {
				if len(m.exportPath) > 0 {
					m.exportPath = m.exportPath[:len(m.exportPath)-1]
				}
			}
			return nil, true
		case "ctrl+u":
			if m.importCursor == 0 {
				m.importPath = ""
			} else {
				m.exportPath = ""
			}
			return nil, true
		}
		if len(msg.String()) == 1 {
			if m.importCursor == 0 {
				m.importPath += msg.String()
			} else {
				m.exportPath += msg.String()
			}
			return nil, true
		}
		return nil, true
	}

	switch msg.String() {
	case "up", "k":
		if m.importCursor > 0 {
			m.importCursor--
		}
		return nil, true
	case "down", "j":
		if m.importCursor < 1 {
			m.importCursor++
		}
		return nil, true
	case "enter":
		m.editingImportPath = true
		return nil, true
	case "i":
		return m.importTSV(), true
	case "I": // Shift+i for APKG import
		return m.importAPKG(), true
	case "x":
		return m.exportTSV(), true
	case "X": // Shift+x for APKG export
		return m.exportAPKG(), true
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
			if m.onConfigChange != nil {
				m.onConfigChange(m.aiProviderName, m.aiTemplates, m.autoPlayAudio)
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
		if m.settingsCursor < 5 {
			m.settingsCursor++
		}
		return nil, true
	case "+", "=":
		return m.setDailyGoal(m.stats.DailyGoal + 1), true
	case "-":
		return m.setDailyGoal(m.stats.DailyGoal - 1), true
	case "enter":
		if m.settingsCursor == 0 {
			switch m.aiProviderName {
			case "offline":
				m.aiProviderName = "template"
				m.aiProvider = ai.TemplateProvider{Templates: m.aiTemplates}
			case "template":
				m.aiProviderName = "disabled"
				m.aiProvider = nil
			default:
				m.aiProviderName = "offline"
				m.aiProvider = ai.OfflineProvider{}
			}
			m.status = fmt.Sprintf("Switched to %s AI provider", m.aiProviderName)
			if m.onConfigChange != nil {
				m.onConfigChange(m.aiProviderName, m.aiTemplates, m.autoPlayAudio)
			}
			return nil, true
		}
		if m.settingsCursor == 5 {
			m.autoPlayAudio = !m.autoPlayAudio
			status := "disabled"
			if m.autoPlayAudio {
				status = "enabled"
			}
			m.status = fmt.Sprintf("Auto-play audio %s", status)
			if m.onConfigChange != nil {
				m.onConfigChange(m.aiProviderName, m.aiTemplates, m.autoPlayAudio)
			}
			return nil, true
		} else if m.settingsCursor > 0 && m.settingsCursor < 4 {
			m.editingTemplate = true
			return nil, true
		}
	}
	return nil, false
}

func (m *Model) updateBrowserKey(msg tea.KeyMsg) (tea.Cmd, bool) {
	switch msg.String() {
	case "up", "k":
		if m.browserCursor > 0 {
			m.browserCursor--
			m.clearReviewHistory()
		}
		return nil, true
	case "down", "j":
		if m.browserCursor < len(m.browserCards)-1 {
			m.browserCursor++
			m.clearReviewHistory()
		}
		return nil, true
	case "enter":
		return m.toggleBrowserHistory(), true
	case "backspace":
		if len(m.browserSearch) > 0 {
			m.browserSearch = m.browserSearch[:len(m.browserSearch)-1]
			m.clearReviewHistory()
			return m.loadBrowserCards(), true
		}
	}
	if len(msg.String()) == 1 && msg.String() >= " " && msg.String() <= "~" {
		m.browserSearch += msg.String()
		m.clearReviewHistory()
		return m.loadBrowserCards(), true
	}
	return nil, true
}

func (m *Model) updateCramKey(msg tea.KeyMsg) (tea.Cmd, bool) {
	if m.cramActive {
		switch msg.String() {
		case "esc":
			m.cramActive = false
			return nil, true
		case "space", "enter":
			if !m.cramRevealed {
				m.cramRevealed = true
				if m.autoPlayAudio {
					card := m.cramCards[m.cramCursor]
					if card.Audio != "" {
						return m.playAudio(card.Audio), true
					}
				}
				return nil, true
			}
		case "a", "1":
			if m.cramRevealed {
				m.cramReviewed++
				m.cramActive = false
				m.nextCramCard()
				return nil, true
			}
		case "h", "g", "e", "2", "3", "4":
			if m.cramRevealed {
				m.cramReviewed++
				m.cramCorrect++
				m.cramActive = false
				m.nextCramCard()
				return nil, true
			}
		}
		return nil, true
	}

	switch msg.String() {
	case "up", "k":
		if m.cramCursor > 0 {
			m.cramCursor--
		}
		return nil, true
	case "down", "j":
		if m.cramCursor < len(m.cramCards)-1 {
			m.cramCursor++
		}
		return nil, true
	case "enter":
		if len(m.cramCards) > 0 {
			m.cramActive = true
			m.cramRevealed = false
			return nil, true
		}
	case "1":
		m.cramType = "bookmarked"
		m.cramCursor = 0
		m.cramCards = nil
		return m.loadCramCards(), true
	case "2":
		m.cramType = "suspended"
		m.cramCursor = 0
		m.cramCards = nil
		return m.loadCramCards(), true
	case "3":
		m.cramType = "leech"
		m.cramCursor = 0
		m.cramCards = nil
		return m.loadCramCards(), true
	case "4":
		m.cramType = "flagged"
		m.cramCursor = 0
		m.cramCards = nil
		return m.loadCramCards(), true
	case "5":
		m.cramType = "all"
		m.cramCursor = 0
		m.cramCards = nil
		return m.loadCramCards(), true
	}
	return nil, false
}

func (m *Model) nextCramCard() {
	if len(m.cramCards) > 0 {
		m.cramCursor = (m.cramCursor + 1) % len(m.cramCards)
	}
}

func (m *Model) clearReviewHistory() {
	m.reviewHistory = nil
	m.reviewHistoryCard = ""
	m.showReviewHistory = false
}

func (m *Model) toggleReviewHistory() tea.Cmd {
	if m.activeView != ViewReview || len(m.dueCards) == 0 {
		return nil
	}
	cardID := m.dueCards[m.cursor].ID
	if m.showReviewHistory && m.reviewHistoryCard == cardID {
		m.clearReviewHistory()
		m.status = "Review history hidden"
		return nil
	}
	m.reviewHistoryCard = cardID
	m.reviewHistory = nil
	m.showReviewHistory = true
	return m.loadReviewHistory(cardID)
}

func (m *Model) toggleBrowserHistory() tea.Cmd {
	if m.activeView != ViewBrowser || len(m.browserCards) == 0 {
		return nil
	}
	cardID := m.browserCards[m.browserCursor].ID
	if m.showReviewHistory && m.reviewHistoryCard == cardID {
		m.clearReviewHistory()
		m.status = "Review history hidden"
		return nil
	}
	m.reviewHistoryCard = cardID
	m.reviewHistory = nil
	m.showReviewHistory = true
	return m.loadReviewHistory(cardID)
}

func (m *Model) filteredDecks() []core.Deck {
	if m.deckFilter == "" {
		return m.decks
	}

	filtered := make([]core.Deck, 0)
	filterLower := strings.ToLower(m.deckFilter)

	for _, deck := range m.decks {
		// Check if deck name matches filter
		if strings.Contains(strings.ToLower(deck.Name), filterLower) {
			filtered = append(filtered, deck)
			continue
		}

		// Check if deck description matches filter
		if strings.Contains(strings.ToLower(deck.Description), filterLower) {
			filtered = append(filtered, deck)
			continue
		}

		// Check if any deck tags match filter
		for _, tag := range deck.Tags {
			if strings.Contains(strings.ToLower(tag), filterLower) {
				filtered = append(filtered, deck)
				break
			}
		}
	}

	return filtered
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

	autoPlayStatus := "off"
	if m.autoPlayAudio {
		autoPlayStatus = "on"
	}
	options := []string{
		fmt.Sprintf("AI Provider: %s", m.aiProviderName),
		fmt.Sprintf("Front Template: %s", m.aiTemplates["front"]),
		fmt.Sprintf("Back Template: %s", m.aiTemplates["back"]),
		fmt.Sprintf("Example Template: %s", m.aiTemplates["example"]),
		fmt.Sprintf("Daily Goal: %d", m.stats.DailyGoal),
		fmt.Sprintf("Auto-play audio: %s", autoPlayStatus),
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
		b.WriteString("\nUse j/k to move, +/- to adjust daily goal, Enter to toggle provider, audio, or edit template.")
	}
	return b.String()
}

func (m *Model) setDailyGoal(goal int) tea.Cmd {
	if goal < 1 {
		goal = 1
	}
	m.status = "Saving daily goal..."
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		if err := m.repo.SetDailyGoal(ctx, goal); err != nil {
			return err
		}
		stats, err := m.repo.Statistics(ctx)
		if err != nil {
			return err
		}
		return dailyGoalSetMsg(stats)
	}
}

func (m *Model) updateAIKey(msg tea.KeyMsg) (tea.Cmd, bool) {
	switch msg.String() {
	case "enter":
		return tea.Batch(m.generateDrafts(), m.tickSpinner()), true
	case "backspace":
		if len(m.aiInput) > 0 {
			m.aiInput = m.aiInput[:len(m.aiInput)-1]
		}
		return nil, true
	case "up", "k":
		if m.draftCursor > 0 {
			m.draftCursor--
		}
		return nil, true
	case "down", "j":
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
	return nil, true
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
	m.drafting = true
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

func (m *Model) exportAPKG() tea.Cmd {
	deckID := m.deck.ID
	path := m.exportPath
	if strings.TrimSpace(deckID) == "" {
		m.status = "Select a deck before exporting"
		return nil
	}
	m.status = "Exporting APKG..."
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
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
		if err := content.ExportAnkiAPKG(file, deck.Notes); err != nil {
			return err
		}
		return exportDoneMsg{count: len(deck.Notes), path: path}
	}
}

func (m *Model) importAPKG() tea.Cmd {
	path := m.importPath
	m.status = "Importing APKG..."
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		notes, err := content.ImportAnkiAPKG(file)
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

		var cards []core.Card
		if m.bookmarkFilter {
			cards, err = m.repo.DueCardsBookmarked(ctx, time.Now(), 50)
		} else {
			cards, err = m.repo.DueCards(ctx, time.Now(), 50)
		}
		if err != nil {
			return err
		}
		decks, err := m.repo.Decks(ctx)
		if err != nil {
			return err
		}
		stats, err := m.repo.Statistics(ctx)
		if err != nil {
			return err
		}
		return reviewRecordedMsg{cardID: card.ID, cards: cards, decks: decks, stats: stats, grade: grade}
	}
}

func (m *Model) toggleBookmark() tea.Cmd {
	if m.activeView != ViewReview || len(m.dueCards) == 0 {
		return nil
	}
	card := m.dueCards[m.cursor]
	next := !card.Bookmarked
	m.status = "Saving bookmark..."
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		if err := m.repo.SetCardBookmark(ctx, card.ID, next); err != nil {
			return err
		}
		return bookmarkToggledMsg{cardID: card.ID, bookmarked: next}
	}
}

func (m *Model) suspendCard() tea.Cmd {
	if m.activeView != ViewReview || len(m.dueCards) == 0 {
		return nil
	}
	card := m.dueCards[m.cursor]
	m.status = "Suspending card..."
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		if err := m.repo.SetCardSuspended(ctx, card.ID, true); err != nil {
			return err
		}
		var cards []core.Card
		var err error
		if m.bookmarkFilter {
			cards, err = m.repo.DueCardsBookmarked(ctx, time.Now(), 50)
		} else {
			cards, err = m.repo.DueCards(ctx, time.Now(), 50)
		}
		if err != nil {
			return err
		}
		decks, err := m.repo.Decks(ctx)
		if err != nil {
			return err
		}
		stats, err := m.repo.Statistics(ctx)
		if err != nil {
			return err
		}
		return cardSuspendedMsg{cardID: card.ID, cards: cards, decks: decks, stats: stats}
	}
}

func (m *Model) undoLastReview() tea.Cmd {
	if strings.TrimSpace(m.lastReviewedCardID) == "" {
		m.status = "No review to undo"
		return nil
	}
	cardID := m.lastReviewedCardID
	m.status = "Undoing last review..."
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		if err := m.repo.UndoLastReview(ctx, cardID); err != nil {
			return err
		}
		var cards []core.Card
		var err error
		if m.bookmarkFilter {
			cards, err = m.repo.DueCardsBookmarked(ctx, time.Now(), 50)
		} else {
			cards, err = m.repo.DueCards(ctx, time.Now(), 50)
		}
		if err != nil {
			return err
		}
		decks, err := m.repo.Decks(ctx)
		if err != nil {
			return err
		}
		stats, err := m.repo.Statistics(ctx)
		if err != nil {
			return err
		}
		return reviewUndoMsg{cardID: cardID, cards: cards, decks: decks, stats: stats}
	}
}

func (m *Model) toggleBookmarkFilter() tea.Cmd {
	m.bookmarkFilter = !m.bookmarkFilter
	m.revealed = false
	m.cursor = 0
	m.clearReviewHistory()
	if m.bookmarkFilter {
		m.status = "Loading bookmarked cards..."
		return m.loadBookmarkedDueCards
	}
	m.status = "Loading all cards..."
	return m.loadDueCards
}

func (m *Model) selectMCQChoice(key string) {
	choiceIdx := -1
	switch key {
	case "1":
		choiceIdx = 0
	case "2":
		choiceIdx = 1
	case "3":
		choiceIdx = 2
	case "4":
		choiceIdx = 3
	}
	if choiceIdx < 0 || choiceIdx >= len(m.dueCards[m.cursor].Choices) {
		return
	}
	m.mcqChoice = choiceIdx
	m.mcqAnswered = true
	m.revealed = true
	chosen := m.dueCards[m.cursor].Choices[choiceIdx]
	m.mcqCorrect = chosen == m.dueCards[m.cursor].Answer
}

func (m *Model) resetMCQState() {
	m.mcqChoice = -1
	m.mcqAnswered = false
	m.mcqCorrect = false
}

func (m *Model) setCardBookmarkLocal(cardID string, bookmarked bool) {
	for i := range m.allDue {
		if m.allDue[i].ID == cardID {
			m.allDue[i].Bookmarked = bookmarked
		}
	}
	for i := range m.dueCards {
		if m.dueCards[i].ID == cardID {
			m.dueCards[i].Bookmarked = bookmarked
		}
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

	// Add contextual help hints based on active view
	helpHint := ""
	switch m.activeView {
	case ViewReview:
		if len(m.dueCards) > 0 {
			if m.revealed {
				helpHint = "| Grade: a/h/g/e"
			} else {
				helpHint = "| Reveal: Space/Enter"
			}
		}
	case ViewCram:
		if m.cramActive {
			if m.cramRevealed {
				helpHint = "| Grade: a/h/g/e"
			} else {
				helpHint = "| Reveal: Space/Enter"
			}
		} else {
			helpHint = "| Filter: 1-5"
		}
	case ViewBrowser:
		helpHint = "| Search: type text"
	case ViewSettings:
		helpHint = "| Nav: j/k | Edit: Enter"
	case ViewImport:
		helpHint = "| Import: Enter | Export: x"
	case ViewAI:
		helpHint = "| Generate: Enter | Approve: a"
	}

	statusLine := fmt.Sprintf("status: %s", singleLine(m.status))
	footer := fmt.Sprintf("tab/arrows views | 1-9 views | ? help | q quit | mouse %d,%d %s", m.mouseX, m.mouseY, helpHint)
	b.WriteString("\n")
	b.WriteString(statusStyle.Render(truncateLine(statusLine, maxInt(20, m.width-2))))
	b.WriteString("\n")
	b.WriteString(statusStyle.Render(truncateLine(footer, maxInt(20, m.width-2))))

	if m.showHelp {
		b.WriteString("\n\n")
		b.WriteString(m.renderHelp())
	}
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
		{"nav-statistics", ViewStatistics, "4 Statistics"},
		{"nav-import", ViewImport, "5 Import"},
		{"nav-ai", ViewAI, "6 AI Drafts"},
		{"nav-settings", ViewSettings, "7 Settings"},
		{"nav-browser", ViewBrowser, "8 Browser"},
		{"nav-cram", ViewCram, "9 Cram"},
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
	views := []View{ViewDashboard, ViewDecks, ViewReview, ViewStatistics, ViewImport, ViewAI, ViewSettings, ViewBrowser, ViewCram}
	labels := []string{"Dashboard", "Decks", "Review", "Statistics", "Import", "AI", "Settings", "Browser", "Cram"}
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
	case ViewStatistics:
		return m.renderStatistics()
	case ViewImport:
		return m.renderImport()
	case ViewAI:
		return m.renderAI(x, y)
	case ViewSettings:
		return m.renderSettings(x, y)
	case ViewBrowser:
		return m.renderBrowser()
	case ViewCram:
		return m.renderCram()
	default:
		streakIndicator := ""
		if m.stats.CurrentStreak > 0 {
			streakIndicator = " 🔥"
		}

		titleStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("205")).
			MarginBottom(1)

		headerBox := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")).
			Padding(0, 1).
			Width(maxInt(40, m.width-60)).
			Render(fmt.Sprintf("Active Deck: %s", lipgloss.NewStyle().Foreground(lipgloss.Color("81")).Render(m.deckLabel())))

		statsStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("226")).Bold(true)

		reviewQueue := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("63")).
			Padding(0, 1).
			Width(maxInt(30, (m.width-64)/2)).
			Render(statsStyle.Render("Review Queue") + "\n" +
				fmt.Sprintf("  Due cards:   %d\n", len(m.dueCards)) +
				fmt.Sprintf("  Bookmarked:  %d (%d due)", m.stats.BookmarkedCards, m.stats.BookmarkedDue))

		collectionStats := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("64")).
			Padding(0, 1).
			Width(maxInt(30, (m.width-64)/2)).
			Render(statsStyle.Render("Collection") + "\n" +
				fmt.Sprintf("  Leech:       %d\n", m.stats.LeechCards) +
				fmt.Sprintf("  Suspended:   %d", m.stats.SuspendedCards))

		goalStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("46")).Bold(true)
		progressBox := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("65")).
			Padding(0, 1).
			Width(maxInt(40, m.width-60)).
			Render(goalStyle.Render("Today's Progress") + "\n" +
				fmt.Sprintf("  Reviews:     %d/%d\n", m.stats.ReviewsToday, m.stats.DailyGoal) +
				fmt.Sprintf("  Streak:      %d days%s", m.stats.CurrentStreak, streakIndicator))

		var db strings.Builder
		db.WriteString(titleStyle.Render("DASHBOARD") + "\n")
		db.WriteString(headerBox + "\n")
		db.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, reviewQueue, " ", collectionStats) + "\n")
		db.WriteString(progressBox + "\n\n")
		db.WriteString(mutedStyle.Render("Use [ and ] to switch decks.\nUse Review (3) to start studying."))

		return db.String()
	}
}

func (m *Model) renderStatistics() string {
	var b strings.Builder
	b.WriteString("Statistics\n\n")

	b.WriteString(fmt.Sprintf("Total Cards:   %d\n", m.stats.TotalCards))
	b.WriteString(fmt.Sprintf("  New:         %d\n", m.stats.NewCards))
	b.WriteString(fmt.Sprintf("  Young:       %d\n", m.stats.YoungCards))
	b.WriteString(fmt.Sprintf("  Mature:      %d\n\n", m.stats.MatureCards))
	b.WriteString(fmt.Sprintf("  Bookmarked:  %d (%d due)\n", m.stats.BookmarkedCards, m.stats.BookmarkedDue))
	b.WriteString(fmt.Sprintf("  Leech:       %d\n", m.stats.LeechCards))
	b.WriteString(fmt.Sprintf("  Suspended:   %d\n\n", m.stats.SuspendedCards))

	b.WriteString(fmt.Sprintf("Total Reviews: %d\n", m.stats.TotalReviews))
	b.WriteString(fmt.Sprintf("Reviews Today: %d/%d\n", m.stats.ReviewsToday, m.stats.DailyGoal))
	// Colored progress bar for daily goal
	progressWidth := 30
	progress := 0
	if m.stats.DailyGoal > 0 {
		progress = (m.stats.ReviewsToday * progressWidth) / m.stats.DailyGoal
		if progress > progressWidth {
			progress = progressWidth
		}
	}
	// Color: green if complete, yellow if halfway, red otherwise
	barColor := "196" // red
	if progress >= progressWidth {
		barColor = "46" // green
	} else if progress >= progressWidth/2 {
		barColor = "226" // yellow
	}
	barStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(barColor))
	bar := strings.Repeat("█", progress) + strings.Repeat("░", progressWidth-progress)
	b.WriteString(fmt.Sprintf("  %s %d%%\n", barStyle.Render(bar), (m.stats.ReviewsToday*100)/maxInt(m.stats.DailyGoal, 1)))
	// Streak with fire emoji if > 0
	streakIndicator := ""
	if m.stats.CurrentStreak > 0 {
		streakIndicator = " 🔥"
	}
	b.WriteString(fmt.Sprintf("Current Streak: %d days%s\n", m.stats.CurrentStreak, streakIndicator))
	b.WriteString(fmt.Sprintf("Success Rate:  %.1f%%\n\n", m.stats.SuccessRate*100))

	b.WriteString("Session Stats:\n")
	b.WriteString(fmt.Sprintf("  Reviewed:    %d\n", m.sessionReviewed))
	if m.sessionReviewed > 0 {
		b.WriteString(fmt.Sprintf("  Correct:     %d\n", m.sessionCorrect))
		b.WriteString(fmt.Sprintf("  Accuracy:     %.1f%%\n\n", float64(m.sessionCorrect)/float64(m.sessionReviewed)*100))
	} else {
		b.WriteString("  (no reviews yet)\n\n")
	}

	b.WriteString("Reviews by Grade:\n")
	grades := []core.ReviewGrade{core.GradeAgain, core.GradeHard, core.GradeGood, core.GradeEasy}
	for _, g := range grades {
		count := m.stats.Grades[g]
		b.WriteString(fmt.Sprintf("  %-5s: %d\n", g, count))
	}

	// Review Activity (last 14 days) with colored bars
	b.WriteString("\nReview Activity (last 14 days):\n")
	if len(m.reviewsPerDay) == 0 {
		b.WriteString("  (no review data yet)\n")
	} else {
		now := time.Now().UTC()
		maxPerDay := 0
		for i := 13; i >= 0; i-- {
			day := now.AddDate(0, 0, -i)
			dayStr := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, time.UTC).Format("2006-01-02")
			count := m.reviewsPerDay[dayStr]
			if count > maxPerDay {
				maxPerDay = count
			}
		}
		if maxPerDay == 0 {
			maxPerDay = 1
		}
		barWidth := 30
		for i := 13; i >= 0; i-- {
			day := now.AddDate(0, 0, -i)
			dayStr := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, time.UTC).Format("2006-01-02")
			count := m.reviewsPerDay[dayStr]
			barLen := (count * barWidth) / maxPerDay
			if barLen == 0 && count > 0 {
				barLen = 1
			}
			// Color based on count: green for high, yellow for medium, red for low
			barColor := "244" // muted/default
			if count > 0 {
				if count >= maxPerDay*3/4 {
					barColor = "46" // green
				} else if count >= maxPerDay/2 {
					barColor = "226" // yellow
				} else {
					barColor = "196" // red
				}
			}
			barStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(barColor))
			bar := strings.Repeat("█", barLen)
			b.WriteString(fmt.Sprintf("  %s %2d %s\n", day.Format("01-02"), count, barStyle.Render(bar)))
		}
	}

	return b.String()
}

func (m *Model) renderDecks(x, y int) string {
	var b strings.Builder
	b.WriteString("Decks\n\n")

	// Show filter if active
	if m.deckFilter != "" {
		b.WriteString(fmt.Sprintf("Filter: %s_\n\n", m.deckFilter))
	}

	filteredDecks := m.filteredDecks()
	if len(filteredDecks) == 0 {
		if m.deckFilter != "" {
			b.WriteString("No decks match filter. Press Esc to clear filter.\n")
		} else {
			b.WriteString("No decks found. Use Import to add notes.\n")
		}
		b.WriteString("\nPress Esc to clear filter.")
		return b.String()
	}

	for i, deck := range filteredDecks {
		prefix := "  "
		style := lipgloss.NewStyle()
		if i == m.deckCursor {
			prefix = "> "
			style = style.Bold(true).Foreground(lipgloss.Color("212"))
		}
		label := fmt.Sprintf("%s%s (%d total, %d due, today %d, %.0f%% success)", prefix, deck.Name, deck.TotalCards, deck.DueCards, deck.ReviewsToday, deck.SuccessRate*100)
		b.WriteString(style.Render(label))
		b.WriteString("\n")
		if deck.Description != "" {
			b.WriteString(fmt.Sprintf("     %s\n", mutedStyle.Render(deck.Description)))
		}
		// Show tags if they exist
		if len(deck.Tags) > 0 {
			tags := strings.Join(deck.Tags, ", ")
			b.WriteString(fmt.Sprintf("     Tags: %s\n", mutedStyle.Render(tags)))
		}
	}
	b.WriteString("\nPress enter to select deck. Type to filter. Esc to clear filter.")
	return b.String()
}

func (m *Model) renderImport() string {
	var b strings.Builder
	b.WriteString("Import / Export\n\n")

	importLabel := "Import file: " + m.importPath
	exportLabel := "Export file: " + m.exportPath

	style := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("212"))
	editStyle := lipgloss.NewStyle().Bold(true).Background(lipgloss.Color("62"))

	if m.importCursor == 0 {
		importLabel = "> " + importLabel
		exportLabel = "  " + exportLabel
		if m.editingImportPath {
			importLabel = editStyle.Render(importLabel)
		} else {
			importLabel = style.Render(importLabel)
		}
	} else {
		importLabel = "  " + importLabel
		exportLabel = "> " + exportLabel
		if m.editingImportPath {
			exportLabel = editStyle.Render(exportLabel)
		} else {
			exportLabel = style.Render(exportLabel)
		}
	}

	b.WriteString(importLabel + "\n")
	b.WriteString(exportLabel + "\n\n")

	b.WriteString("Actions:\n")
	b.WriteString("  i         : Import TSV\n")
	b.WriteString("  I         : Import APKG\n")
	b.WriteString("  x         : Export TSV\n")
	b.WriteString("  X         : Export APKG\n\n")

	b.WriteString(fmt.Sprintf("Current Deck: %s\n\n", m.deckLabel()))

	if m.editingImportPath {
		b.WriteString("EDITING - Enter to save, Esc to cancel.")
	} else {
		b.WriteString("Use j/k to select path, Enter to edit, i/I/x/X to execute.")
	}

	return b.String()
}

func (m *Model) renderAI(x, y int) string {
	var b strings.Builder
	spinner := ""
	if m.drafting {
		frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
		spinner = " " + frames[m.spinnerFrame%len(frames)]
	}
	fmt.Fprintf(&b, "AI Drafts%s\n\nDeck: %s\nTopic: %s\n\nEnter generate | a approve | d discard\n", spinner, m.deckLabel(), m.aiInput)
	if len(m.drafts) == 0 {
		if m.drafting {
			b.WriteString("\nDrafting in progress...")
		} else {
			b.WriteString("\nNo drafts yet.")
		}
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

func (m *Model) renderBrowser() string {
	var b strings.Builder
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205")).MarginBottom(1)
	b.WriteString(titleStyle.Render("Card Browser") + "\n")

	searchStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("62")).
		Padding(0, 1).
		Width(maxInt(30, m.width-60))

	b.WriteString(searchStyle.Render(fmt.Sprintf("Search: %s_", m.browserSearch)) + "\n\n")

	if len(m.browserCards) == 0 {
		b.WriteString("No cards found. Type to search.\n\n")
		b.WriteString(mutedStyle.Render("Use left/right/[ to change deck filter.\n"))
		return b.String()
	}
	for i, card := range m.browserCards {
		prefix := "  "
		style := lipgloss.NewStyle()
		if i == m.browserCursor {
			prefix = "> "
			style = style.Bold(true).Foreground(lipgloss.Color("212"))
		}
		kind := "FC"
		if card.Kind == core.CardKindMCQ {
			kind = "MCQ"
		}
		bookmark := ""
		if card.Bookmarked {
			bookmark = " [B]"
		}
		leech := ""
		if card.Leech {
			leech = " [L]"
		}
		suspended := ""
		if card.Suspended {
			suspended = " [S]"
		}
		label := fmt.Sprintf("%s[%s] %s%s%s%s", prefix, kind, card.Prompt, bookmark, leech, suspended)
		b.WriteString(style.Render(label))
		b.WriteString("\n")
	}
	if m.showReviewHistory && m.browserCursor < len(m.browserCards) && m.reviewHistoryCard == m.browserCards[m.browserCursor].ID {
		b.WriteString("\n")
		b.WriteString(m.renderReviewHistory(m.browserCards[m.browserCursor].Prompt))
		b.WriteString("\n")
	}
	b.WriteString("\nUse j/k to navigate, type to search, Enter for history, backspace to delete.\n")
	return b.String()
}

func (m *Model) renderReviewHistory(label string) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("Review History: %s\n", label))
	if len(m.reviewHistory) == 0 {
		b.WriteString("  No reviews yet.")
		return b.String()
	}
	for i, log := range m.reviewHistory {
		fmt.Fprintf(&b, "  %d. %s at %s -> next %s (%s, reviews %d, lapses %d)\n",
			i+1,
			log.Grade,
			log.Reviewed.Local().Format("Jan 02 15:04"),
			log.Due.Local().Format("Jan 02"),
			formatReviewInterval(log.Interval),
			log.Reviews,
			log.Lapses,
		)
	}
	return strings.TrimRight(b.String(), "\n")
}

func formatReviewInterval(interval time.Duration) string {
	if interval <= 0 {
		return "same day"
	}
	hours := int(interval.Hours())
	if hours < 24 {
		if hours <= 1 {
			return "1 hour"
		}
		return fmt.Sprintf("%d hours", hours)
	}
	days := hours / 24
	if days == 1 {
		return "1 day"
	}
	return fmt.Sprintf("%d days", days)
}

func (m *Model) renderCram() string {
	if m.cramActive {
		var b strings.Builder
		card := m.cramCards[m.cramCursor]
		audioIndicator := ""
		if card.Audio != "" {
			audioIndicator = " [Audio]"
		}
		b.WriteString("Cram Review\n\n")
		b.WriteString(fmt.Sprintf("Prompt: %s%s\n\n", card.Prompt, audioIndicator))
		if m.cramRevealed {
			b.WriteString(fmt.Sprintf("Answer: %s\n\n", card.Answer))
			b.WriteString("Grade: a Again | h Hard | g Good | e Easy\n")
		} else {
			b.WriteString("Press space or enter to reveal.\n")
		}
		b.WriteString("\np play audio | q to exit cram review.")
		return b.String()
	}

	var b strings.Builder
	b.WriteString("Cram Mode\n\n")
	b.WriteString(fmt.Sprintf("Filter: %s\n\n", m.cramType))
	if len(m.cramCards) == 0 {
		b.WriteString("No cards found for this filter.\n\n")
		b.WriteString("Press 1-5 to change filter:\n")
		b.WriteString("  1: Bookmarked\n")
		b.WriteString("  2: Suspended\n")
		b.WriteString("  3: Leeches\n")
		b.WriteString("  4: All flagged\n")
		b.WriteString("  5: All cards\n")
		return b.String()
	}
	for i, card := range m.cramCards {
		prefix := "  "
		style := lipgloss.NewStyle()
		if i == m.cramCursor {
			prefix = "> "
			style = style.Bold(true).Foreground(lipgloss.Color("212"))
		}
		kind := "FC"
		if card.Kind == core.CardKindMCQ {
			kind = "MCQ"
		}
		bookmark := ""
		if card.Bookmarked {
			bookmark = " [B]"
		}
		leech := ""
		if card.Leech {
			leech = " [L]"
		}
		suspended := ""
		if card.Suspended {
			suspended = " [S]"
		}
		label := fmt.Sprintf("%s[%s] %s%s%s%s", prefix, kind, card.Prompt, bookmark, leech, suspended)
		b.WriteString(style.Render(label))
		b.WriteString("\n")
	}
	if m.cramReviewed > 0 {
		accuracy := 0.0
		if m.cramReviewed > 0 {
			accuracy = float64(m.cramCorrect) / float64(m.cramReviewed) * 100
		}
		b.WriteString(fmt.Sprintf("\nCram Stats: %d reviewed, %d correct (%.1f%%)\n", m.cramReviewed, m.cramCorrect, accuracy))
	}
	b.WriteString("\nUse j/k to navigate. Type 1-5 for filter. q to quit.\n")
	return b.String()
}

func (m *Model) renderHelp() string {
	var b strings.Builder
	b.WriteString("Keyboard Shortcuts\n\n")
	b.WriteString("Global:\n")
	b.WriteString("  1-9          Switch to view\n")
	b.WriteString("  Tab/arrows   Cycle views\n")
	b.WriteString("  w/s          Previous/next view\n")
	b.WriteString("  ?            Toggle this help\n")
	b.WriteString("  q/Ctrl+c     Quit\n\n")

	b.WriteString("Dashboard/Decks:\n")
	b.WriteString("  [ ]          Previous/next deck\n")
	b.WriteString("  Enter        Select deck (Decks view)\n\n")

	b.WriteString("Review:\n")
	b.WriteString("  Space/Enter  Reveal answer\n")
	b.WriteString("  a/h/g/e      Grade Again/Hard/Good/Easy\n")
	b.WriteString("  b            Toggle bookmark\n")
	b.WriteString("  B            Toggle bookmarked-only mode\n")
	b.WriteString("  x            Suspend card\n")
	b.WriteString("  u            Undo last review\n")
	b.WriteString("  r            Toggle card review history\n")
	b.WriteString("  p            Play audio\n")
	b.WriteString("  1-4          Select MCQ choice\n\n")

	b.WriteString("Browser:\n")
	b.WriteString("  j/k          Navigate cards\n")
	b.WriteString("  Enter        Toggle card review history\n")
	b.WriteString("  Type         Search cards\n")
	b.WriteString("  Backspace    Delete search char\n\n")

	b.WriteString("Cram:\n")
	b.WriteString("  j/k          Navigate cards\n")
	b.WriteString("  Enter        Start cram review\n")
	b.WriteString("  p            Play audio (in review)\n")
	b.WriteString("  1-5          Filter: bookmarked/suspended/leech/flagged/all\n\n")

	b.WriteString("Import:\n")
	b.WriteString("  j/k          Select field\n")
	b.WriteString("  Enter        Start/stop editing path\n")
	b.WriteString("  i/I          Import TSV/APKG\n")
	b.WriteString("  x/X          Export TSV/APKG\n\n")

	b.WriteString("Settings:\n")
	b.WriteString("  j/k          Navigate options\n")
	b.WriteString("  +/-          Adjust daily goal\n")
	b.WriteString("  Enter        Toggle AI provider / edit template")

	return panelStyle.Width(60).Render(b.String())
}

func (m *Model) renderReview(x, y int) string {
	if len(m.dueCards) == 0 {
		if m.bookmarkFilter {
			return "Review (Bookmarked)\n\nNo bookmarked cards due."
		}
		return "Review\n\nNo cards due."
	}
	card := m.dueCards[m.cursor]
	bookmark := "Bookmark: off"
	if card.Bookmarked {
		bookmark = "Bookmark: on"
	}
	leech := ""
	if card.Leech {
		leech = " | LEECH"
	}
	suspended := ""
	if card.Suspended {
		suspended = " | SUSPENDED"
	}
	filterBanner := ""
	if m.bookmarkFilter {
		filterBanner = " (Bookmarked)"
	}
	keys := "b toggle | x suspend | B filter | u undo | r history | p audio"
	if m.bookmarkFilter {
		keys = "b toggle | x suspend | B all cards | u undo | r history | p audio"
	}
	audioIndicator := ""
	if card.Audio != "" {
		audioIndicator = " [Audio]"
	}

	var answer string
	if card.Kind == core.CardKindMCQ && len(card.Choices) > 0 {
		if m.revealed {
			if m.mcqAnswered {
				feedback := "Incorrect"
				if m.mcqCorrect {
					feedback = "Correct"
				}
				answer = fmt.Sprintf("%s: %s\n\n%s\n\nGrade: a Again | h Hard | g Good | e Easy", feedback, card.Answer, renderMCQChoices(card.Choices, m.mcqChoice))
				m.hitboxes = append(m.hitboxes, Hitbox{ID: "grade-again", View: ViewReview, X: x + 7, Y: y + 10, Width: 5, Height: 1})
				m.hitboxes = append(m.hitboxes, Hitbox{ID: "grade-hard", View: ViewReview, X: x + 17, Y: y + 10, Width: 4, Height: 1})
				m.hitboxes = append(m.hitboxes, Hitbox{ID: "grade-good", View: ViewReview, X: x + 26, Y: y + 10, Width: 4, Height: 1})
				m.hitboxes = append(m.hitboxes, Hitbox{ID: "grade-easy", View: ViewReview, X: x + 35, Y: y + 10, Width: 4, Height: 1})
			} else {
				answer = fmt.Sprintf("1-4 select answer\n\n%s\n\nGrade: a Again | h Hard | g Good | e Easy", renderMCQChoices(card.Choices, m.mcqChoice))
				m.hitboxes = append(m.hitboxes, Hitbox{ID: "grade-again", View: ViewReview, X: x + 7, Y: y + 8, Width: 5, Height: 1})
				m.hitboxes = append(m.hitboxes, Hitbox{ID: "grade-hard", View: ViewReview, X: x + 17, Y: y + 8, Width: 4, Height: 1})
				m.hitboxes = append(m.hitboxes, Hitbox{ID: "grade-good", View: ViewReview, X: x + 26, Y: y + 8, Width: 4, Height: 1})
				m.hitboxes = append(m.hitboxes, Hitbox{ID: "grade-easy", View: ViewReview, X: x + 35, Y: y + 8, Width: 4, Height: 1})
			}
		} else {
			answer = "Press space or enter to reveal choices."
		}
	} else if m.revealed {
		answer = card.Answer + "\n\nGrade: a Again | h Hard | g Good | e Easy"
		m.hitboxes = append(m.hitboxes, Hitbox{ID: "grade-again", View: ViewReview, X: x + 7, Y: y + 6, Width: 5, Height: 1})
		m.hitboxes = append(m.hitboxes, Hitbox{ID: "grade-hard", View: ViewReview, X: x + 17, Y: y + 6, Width: 4, Height: 1})
		m.hitboxes = append(m.hitboxes, Hitbox{ID: "grade-good", View: ViewReview, X: x + 26, Y: y + 6, Width: 4, Height: 1})
		m.hitboxes = append(m.hitboxes, Hitbox{ID: "grade-easy", View: ViewReview, X: x + 35, Y: y + 6, Width: 4, Height: 1})
	} else {
		answer = "Press space or enter to reveal."
	}

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
	header := fmt.Sprintf("%s%s %d/%d", titleStyle.Render("Review"), filterBanner, m.cursor+1, len(m.dueCards))

	view := fmt.Sprintf("%s\n%s | %s\n%s%s%s\n\n%s\n\n%s", header, bookmark, keys, leech, suspended, audioIndicator, card.Prompt, answer)
	if m.showReviewHistory && m.reviewHistoryCard == card.ID {
		view += "\n\n" + m.renderReviewHistory(card.Prompt)
	}
	return view
}

func renderMCQChoices(choices []string, selected int) string {
	var b strings.Builder
	for i, choice := range choices {
		prefix := "  "
		mark := " "
		if i == selected {
			mark = ">"
		}
		b.WriteString(fmt.Sprintf("%s%d: %s%s\n", prefix, i+1, mark, choice))
	}
	return strings.TrimRight(b.String(), "\n")
}

func (m *Model) nextViewCmd() tea.Cmd {
	views := []View{ViewDashboard, ViewDecks, ViewReview, ViewStatistics, ViewImport, ViewAI, ViewSettings, ViewBrowser, ViewCram}
	for i, view := range views {
		if m.activeView == view {
			return m.updateView(views[(i+1)%len(views)])
		}
	}
	return m.updateView(ViewDashboard)
}

func (m *Model) previousViewCmd() tea.Cmd {
	views := []View{ViewDashboard, ViewDecks, ViewReview, ViewStatistics, ViewImport, ViewAI, ViewSettings, ViewBrowser, ViewCram}
	for i, view := range views {
		if m.activeView == view {
			return m.updateView(views[(i-1+len(views))%len(views)])
		}
	}
	return m.updateView(ViewDashboard)
}

func (m *Model) updateView(view View) tea.Cmd {
	m.activeView = view
	m.clearReviewHistory()
	if view == ViewStatistics {
		return m.loadStatistics()
	}
	if view == ViewBrowser {
		m.browserSearch = ""
		m.browserDeckID = m.deck.ID
		m.browserCursor = 0
		return m.loadBrowserCards()
	}
	if view == ViewCram {
		m.cramType = "bookmarked"
		m.cramCursor = 0
		m.cramReviewed = 0
		m.cramCorrect = 0
		return m.loadCramCards()
	}
	return nil
}

func (m *Model) reloadBrowserForSelectedDeck() tea.Cmd {
	m.browserDeckID = m.deck.ID
	m.browserSearch = ""
	m.browserCards = nil
	m.browserCursor = 0
	m.clearReviewHistory()
	m.status = fmt.Sprintf("Browsing %s", m.deckLabel())
	return m.loadBrowserCards()
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
		return m.updateView(View(strings.TrimPrefix(id, "nav-")))
	case strings.HasPrefix(id, "tab-"):
		return m.updateView(View(strings.TrimPrefix(id, "tab-")))
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

func friendlyError(err error) string {
	if errors.Is(err, os.ErrNotExist) {
		var pathErr *os.PathError
		if errors.As(err, &pathErr) {
			return fmt.Sprintf("Error: no such file or directory: %s", filepath.Base(pathErr.Path))
		}
		return "Error: no such file or directory"
	}
	return fmt.Sprintf("Error: %v", err)
}

func singleLine(s string) string {
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "\r", " ")
	return strings.Join(strings.Fields(s), " ")
}

func truncateLine(s string, maxWidth int) string {
	if maxWidth <= 0 || len(s) <= maxWidth {
		return s
	}
	if maxWidth <= 1 {
		return s[:maxWidth]
	}
	if maxWidth <= 3 {
		return s[:maxWidth]
	}
	return s[:maxWidth-3] + "..."
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
