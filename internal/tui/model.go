package tui

import (
	"fmt"
	"io"
	"log"
	"net/url"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	"deutsch-tui/internal/ai"
	"deutsch-tui/internal/app"
	"deutsch-tui/internal/core"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type View string

const (
	ViewDashboard      View = "dashboard"
	ViewDecks          View = "decks"
	ViewReview         View = "review"
	ViewStatistics     View = "statistics"
	ViewImport         View = "import"
	ViewAI             View = "ai"
	ViewSettings       View = "settings"
	ViewBrowser        View = "browser"
	ViewCram           View = "cram"
	ViewSessionSummary View = "session_summary"
)

type RevealState int

const (
	RevealIdle RevealState = iota
	RevealRevealing
	RevealRevealed
)

type Breakpoint string

const (
	BreakpointCompact Breakpoint = "compact"
	BreakpointMedium  Breakpoint = "medium"
	BreakpointWide    Breakpoint = "wide"
)

type viewportLayout struct {
	X      int
	Y      int
	Width  int
	Height int
}

type Model struct {
	repo                  core.Repository
	scheduler             core.Scheduler
	width                 int
	height                int
	activeView            View
	breakpoint            Breakpoint
	decks                 []core.Deck
	deckIndex             int
	deck                  core.Deck
	deckCursor            int
	allDue                []core.Card
	dueCards              []core.Card
	cursor                int
	revealState           RevealState
	revealProgress        float64
	lastReviewedCardID    string
	lastReviewedGrade     core.ReviewGrade
	status                string
	mouseX                int
	mouseY                int
	hitboxes              []Hitbox
	aiProvider            ai.Provider
	aiProviderName        string
	aiTemplates           map[string]map[string]string
	aiTemplateSets        []string
	aiTemplateIndex       int
	autoPlayAudio         bool
	strictNormalization   bool
	stats                 core.Statistics
	settingsCursor        int
	editingTemplate       bool
	aiInput               string
	draftSource           string
	drafts                []ai.Draft
	draftCursor           int
	importPath            string
	exportPath            string
	exportDeckID          string
	exportTag             string
	exportFilter          string // e.g. "All", "Mature", "Learning"
	importCursor          int    // 0: import path, 1: export path, 2: export deck, 3: export tag, 4: export filter
	editingImportPath     bool
	editingExportTag      bool
	theme                 string
	onConfigChange        func(string, string, map[string]map[string]string, bool, bool)
	bookmarkFilter        bool
	originalTemplateValue string
	mcqChoice             int
	mcqAnswered           bool
	mcqCorrect            bool
	browserCards          []core.Card
	browserCursor         int
	browserSearch         string
	browserTag            string
	browserDeckID         string
	searchingTags         bool
	browserSelected       map[string]bool
	sessionReviewed       int
	sessionCorrect        int
	lastSessionReviewed   int
	lastSessionCorrect    int
	sessionStartTime      time.Time
	showHelp              bool
	cramCards             []core.Card
	cramCursor            int
	cramType              string
	cramReviewed          int
	cramCorrect           int
	cramActive            bool
	cramRevealed          bool
	reviewsPerDay         map[string]int
	reviewHistory         []core.ReviewLog
	reviewHistoryCard     string
	showReviewHistory     bool
	reviewPredictions     map[core.ReviewGrade]time.Duration
	spinnerFrame          int
	deckFilter            string
	deckSelected          map[string]bool
	drafting              bool
	statsScroll           int
	statsTotalLines       int
	isDragging            bool
	dragView              View
	dragTrackStartY       int
	dragVisible           int
	dragTotal             int
	searchingAI           bool
	searchingBrowser      bool
	searchingDecks        bool
	taggingCards          bool
	tagInput              string
	statusSeq             int
	confirmingDelete      bool
	deleteAction          func() tea.Cmd
	deleteIDs             []string
	editingDeckLimits     bool
	limitCursor           int                // 0: new limit, 1: review limit
	gradingInProgress     bool               // Prevent double-grading
	logger                *app.LeveledLogger // Add logger field
	typingMode            bool               // Typing exercise mode
	typedAnswer           string             // Current typed answer
	typingChecked         bool               // Whether typing answer has been checked
	typingCorrect         bool               // Whether typed answer was correct
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
	Theme               string
	AIProvider          ai.Provider
	AIProviderName      string
	AITemplates         map[string]map[string]string
	AutoPlayAudio       bool
	StrictNormalization bool
	ImportPath          string
	ExportPath          string
	OnConfigChange      func(string, string, map[string]map[string]string, bool, bool)
	Logger              *app.LeveledLogger // Add logger option
}

func NewModelWithOptions(repo core.Repository, scheduler core.Scheduler, opts ModelOptions) *Model {
	providerName := opts.AIProviderName
	if providerName == "" {
		providerName = "offline"
	}
	templates := opts.AITemplates
	if templates == nil {
		templates = map[string]map[string]string{
			"vocabulary": {
				"front":   "{{.Topic}}",
				"back":    "German prompt for {{.Topic}}",
				"example": "Practice sentence using {{.Topic}}.",
			},
			"phrases": {
				"front":   "Common German phrase for {{.Topic}}",
				"back":    "English translation",
				"example": "Context sentence using the phrase.",
			},
			"grammar": {
				"front":   "Grammar rule for {{.Topic}}",
				"back":    "Explanation with examples",
				"example": "Example sentence demonstrating the rule.",
			},
		}
	}
	if len(templates) == 0 {
		templates = map[string]map[string]string{
			"vocabulary": {
				"front":   "{{.Topic}}",
				"back":    "German prompt for {{.Topic}}",
				"example": "Practice sentence using {{.Topic}}.",
			},
			"phrases": {
				"front":   "Common German phrase for {{.Topic}}",
				"back":    "English translation",
				"example": "Context sentence using the phrase.",
			},
			"grammar": {
				"front":   "Grammar rule for {{.Topic}}",
				"back":    "Explanation with examples",
				"example": "Example sentence demonstrating the rule.",
			},
		}
	}

	var sets []string
	for k := range templates {
		sets = append(sets, k)
	}
	sort.Strings(sets)

	aiTemplateIndex := 0
	for i, s := range sets {
		if s == "vocabulary" {
			aiTemplateIndex = i
			break
		}
	}

	autoPlayAudio := opts.AutoPlayAudio
	strictNormalization := opts.StrictNormalization
	provider := opts.AIProvider
	if provider == nil {
		switch providerName {
		case "template":
			activeSet := ""
			if len(sets) > 0 {
				activeSet = sets[aiTemplateIndex]
			}
			provider = ai.TemplateProvider{
				Templates: templates,
				ActiveSet: activeSet,
			}
		case "disabled":
			provider = nil
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

	// Create a default logger if none provided
	logger := opts.Logger
	if logger == nil {
		// Create a minimal logger that discards output
		nullLogger := log.New(io.Discard, "", 0)
		logger = app.NewLeveledLogger(nullLogger, app.LogLevelInfo)
	}

	return &Model{
		repo:                repo,
		scheduler:           scheduler,
		theme:               opts.Theme,
		aiProvider:          provider,
		aiProviderName:      providerName,
		aiTemplates:         templates,
		aiTemplateSets:      sets,
		aiTemplateIndex:     aiTemplateIndex,
		autoPlayAudio:       autoPlayAudio,
		strictNormalization: strictNormalization,
		width:               80,
		height:              24,
		activeView:          ViewDashboard,
		breakpoint:          BreakpointMedium,
		status:              "Ready",
		aiInput:             "der Kaffee",
		importPath:          filepath.Clean(importPath),
		exportPath:          filepath.Clean(exportPath),
		exportFilter:        "All",
		onConfigChange:      opts.OnConfigChange,
		browserSelected:     make(map[string]bool),
		deckSelected:        make(map[string]bool),
		logger:              logger, // Set the logger
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
type reviewUndoneMsg struct {
	cardID string
	cards  []core.Card
	decks  []core.Deck
	stats  core.Statistics
	grade  core.ReviewGrade
}
type cramCardsMsg []core.Card
type browserCardsMsg []core.Card
type reviewsPerDayMsg map[string]int
type reviewHistoryMsg struct {
	cardID string
	logs   []core.ReviewLog
}
type reviewPredictionsMsg map[core.ReviewGrade]time.Duration
type statusMsg struct {
	text string
}
type tagsUpdatedMsg struct {
	cardIDs []string
	tags    []string
}
type deckDeletedMsg struct{}

func (m *Model) Init() tea.Cmd {
	return tea.Sequence(m.loadDueCards, m.loadDecks, m.loadStatistics(), m.loadReviewsPerDay())
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.logger.Debug("Processing message: %T", msg)

	switch msg := msg.(type) {
	case spinnerTickMsg:
		if m.drafting {
			m.spinnerFrame++
			return m, m.tickSpinner()
		}
		return m, nil
	case revealTickMsg:
		if m.revealState == RevealRevealing {
			m.revealProgress += 10
			if m.revealProgress >= 100 {
				m.revealProgress = 100
				m.revealState = RevealRevealed
				m.logger.Debug("Card fully revealed")
			} else {
				return m, m.tickReveal()
			}
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.breakpoint = breakpointForWidth(msg.Width)
		m.isDragging = false
		m.logger.Debug("Window resized to %dx%d, breakpoint: %s", msg.Width, msg.Height, m.breakpoint)
		return m, nil
	case error:
		m.drafting = false
		m.gradingInProgress = false
		m.status = friendlyError(msg)
		m.logger.Error("Error occurred: %v", msg)
	case decksMsg:
		m.logger.Debug("Received %d decks", len(msg))
		m.syncDecks([]core.Deck(msg))
	case dueCardsMsg:
		m.logger.Debug("Received %d due cards", len(msg))
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
		m.logger.Debug("Received %d bookmarked due cards", len(msg))

	case statsMsg:
		m.stats = core.Statistics(msg)
		m.logger.Debug("Received statistics update")
	case reviewsPerDayMsg:
		m.reviewsPerDay = map[string]int(msg)
		m.logger.Debug("Received reviews per day data")
	case reviewHistoryMsg:
		if msg.cardID == m.reviewHistoryCard {
			m.reviewHistory = msg.logs
			m.showReviewHistory = true
			if len(msg.logs) == 0 {
				m.status = "No review history for card"
			} else {
				m.status = fmt.Sprintf("%d review history entries", len(msg.logs))
			}
			m.logger.Debug("Received %d review history entries for card %s", len(msg.logs), msg.cardID)
		}
	case reviewPredictionsMsg:
		m.reviewPredictions = map[core.ReviewGrade]time.Duration(msg)
		m.logger.Debug("Received review predictions")
	case draftsMsg:
		m.drafting = false
		m.drafts = []ai.Draft(msg)
		m.draftCursor = 0
		if len(m.drafts) == 0 {
			m.status = "No drafts generated"
		} else {
			m.status = fmt.Sprintf("%d drafts ready", len(m.drafts))
		}
		m.logger.Info("Generated %d AI drafts", len(msg))
	case draftApprovedMsg:
		m.removeDraft(msg.noteID)
		m.allDue = msg.cards
		m.applyDeckFilter()
		m.status = "Draft approved"
		m.logger.Info("Approved AI draft %s", msg.noteID)
	case importDoneMsg:
		m.syncDecks(msg.decks)
		m.allDue = msg.cards
		m.applyDeckFilter()
		m.logger.Info("Import completed: %d notes from %s", msg.count, filepath.Base(msg.path))
		return m, m.setStatus(fmt.Sprintf("Imported %d notes from %s", msg.count, filepath.Base(msg.path)), 3*time.Second)
	case statusMsg:
		m.logger.Debug("Setting status: %s", msg.text)
		return m, m.setStatus(msg.text, 3*time.Second)
	case tagsUpdatedMsg:
		m.taggingCards = false
		m.tagInput = ""
		// Update local browser cards if visible
		for _, id := range msg.cardIDs {
			for i := range m.browserCards {
				if m.browserCards[i].ID == id {
					m.browserCards[i].Tags = msg.tags
				}
			}
		}
		m.logger.Info("Updated tags for %d cards", len(msg.cardIDs))
		return m, m.setStatus(fmt.Sprintf("Updated tags for %d cards", len(msg.cardIDs)), 3*time.Second)
	case deckDeletedMsg:
		m.logger.Info("Deck deleted")
		return m, tea.Batch(m.loadDecks, m.loadDueCards)
	case timedClearStatusMsg:
		if msg.seq == m.statusSeq {
			m.status = "Ready"
		}
		return m, nil

	case reviewRecordedMsg:
		m.gradingInProgress = false
		m.lastReviewedCardID = msg.cardID
		m.lastReviewedGrade = msg.grade
		m.syncDecks(msg.decks)
		m.stats = msg.stats
		m.allDue = msg.cards
		m.applyDeckFilter()
		m.sessionReviewed++
		if msg.grade != core.GradeAgain {
			m.sessionCorrect++
		}
		m.logger.Info("Recorded review for card %s with grade %s", msg.cardID, msg.grade)

		// Reset typing state for next card
		m.typedAnswer = ""
		m.typingChecked = false
		m.typingCorrect = false

		if len(m.dueCards) == 0 && m.sessionReviewed > 0 {
			m.activeView = ViewSessionSummary
		}
		return m, m.loadReviewsPerDay()
	case bookmarkToggledMsg:
		m.setCardBookmarkLocal(msg.cardID, msg.bookmarked)
		if msg.bookmarked {
			m.status = "Card bookmarked"
		} else {
			m.status = "Card bookmark removed"
		}
		m.logger.Debug("Toggled bookmark for card %s: %t", msg.cardID, msg.bookmarked)
	case cardSuspendedMsg:
		m.syncDecks(msg.decks)
		m.stats = msg.stats
		m.allDue = msg.cards
		m.applyDeckFilter()
		m.status = "Card suspended"
		m.logger.Info("Suspended card %s", msg.cardID)
		return m, m.loadReviewsPerDay()
	case dailyGoalSetMsg:
		m.stats = core.Statistics(msg)
		m.status = fmt.Sprintf("Daily goal set to %d", m.stats.DailyGoal)
		m.logger.Info("Daily goal set to %d", m.stats.DailyGoal)
	case reviewUndoneMsg:
		m.lastReviewedCardID = ""
		if m.sessionReviewed > 0 {
			m.sessionReviewed--
			if msg.grade != core.GradeAgain && m.sessionCorrect > 0 {
				m.sessionCorrect--
			}
		}
		m.syncDecks(msg.decks)
		m.stats = msg.stats
		m.allDue = msg.cards
		m.applyDeckFilter()
		m.status = "Last review undone"
		m.logger.Info("Undid last review for card %s", msg.cardID)
		return m, m.loadReviewsPerDay()
	case browserCardsMsg:
		m.browserCards = []core.Card(msg)
		if m.browserCursor >= len(m.browserCards) {
			m.browserCursor = maxInt(0, len(m.browserCards)-1)
		}
		if len(m.browserCards) == 0 {
			m.status = "No cards found"
		} else {
			m.status = fmt.Sprintf("%d cards found", len(m.browserCards))
		}
		m.logger.Debug("Loaded %d browser cards", len(msg))
	case cramCardsMsg:
		allCards := []core.Card(msg)
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
		if m.cramCursor >= len(m.cramCards) {
			m.cramCursor = maxInt(0, len(m.cramCards)-1)
		}
		if len(m.cramCards) == 0 {
			m.status = "No cards found for this filter"
		} else {
			m.status = fmt.Sprintf("%d cards in cram mode", len(m.cramCards))
		}
		m.logger.Debug("Loaded %d cram cards with filter %s", len(m.cramCards), m.cramType)
	case tea.KeyPressMsg:
		m.logger.Debug("Key pressed: %s", msg.String())
		return m.updateKey(msg)
	case tea.MouseMsg:
		if m.confirmingDelete {
			return m, nil
		}
		mouse := msg.Mouse()
		m.mouseX = mouse.X
		m.mouseY = mouse.Y

		switch msg.(type) {
		case tea.MouseClickMsg, *tea.MouseClickMsg:
			if mouse.Button == tea.MouseLeft {
				if hit, ok := m.hitboxAt(mouse.X, mouse.Y); ok {
					m.logger.Debug("Mouse click at (%d, %d) on hitbox %s", mouse.X, mouse.Y, hit.ID)
					if strings.Contains(hit.ID, "-scroll-") {
						m.isDragging = true
						m.dragView = hit.View
						parts := strings.Split(hit.ID, "-")
						if len(parts) > 0 {
							if row, err := strconv.Atoi(parts[len(parts)-1]); err == nil {
								m.dragTrackStartY = mouse.Y - row
							}
						}
						layout := m.activeViewContentLayout()
						switch hit.View {
						case ViewStatistics:
							m.dragVisible = m.statisticsVisibleLines(layout.Height)
							m.dragTotal = m.statsTotalLines
						case ViewBrowser:
							m.dragVisible = m.listVisibleLines(layout.Height)
							m.dragTotal = len(m.browserCards)
						case ViewCram:
							m.dragVisible = m.listVisibleLines(layout.Height)
							m.dragTotal = len(m.cramCards)
						}
					}
					cmd := m.activateHitbox(hit.ID)
					return m, cmd
				}
			}
		case tea.MouseMotionMsg, *tea.MouseMotionMsg:
			if m.isDragging {
				m.handleMouseDrag(mouse.Y)
			}
		case tea.MouseReleaseMsg, *tea.MouseReleaseMsg:
			if m.isDragging {
				m.logger.Debug("Mouse drag released")
				m.isDragging = false
			}
		case tea.MouseWheelMsg, *tea.MouseWheelMsg:
			m.logger.Debug("Mouse wheel event: button=%v", mouse.Button)
			if mouse.Button == tea.MouseWheelUp {
				switch m.activeView {
				case ViewStatistics:
					m.scrollStats(-1)
				case ViewBrowser:
					m.moveBrowserCursor(-1)
				case ViewCram:
					m.moveCramCursor(-1)
				}
			} else if mouse.Button == tea.MouseWheelDown {
				switch m.activeView {
				case ViewStatistics:
					m.scrollStats(1)
				case ViewBrowser:
					m.moveBrowserCursor(1)
				case ViewCram:
					m.moveCramCursor(1)
				}
			}
		}
	}
	return m, nil
}

func (m *Model) View() tea.View {
	if m.width < 20 || m.height < 10 {
		return tea.View{Content: "Terminal too small"}
	}

	m.hitboxes = m.hitboxes[:0]
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

	helpHint := ""
	switch m.activeView {
	case ViewDashboard:
		helpHint = "| [ / ] switch decks | 3 Review"
	case ViewReview:
		if len(m.dueCards) > 0 {
			if m.revealState == RevealRevealed {
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
		if m.searchingBrowser {
			helpHint = "| Searching... | Stop: Esc/Enter"
		} else {
			helpHint = "| Search: / | Bookmark: b | Suspend: x"
		}
	case ViewSettings:
		helpHint = "| Nav: j/k | Edit: Enter"
	case ViewImport:
		helpHint = "| Import: Enter | Export: x"
	case ViewAI:
		helpHint = "| Generate: Enter | Approve: a"
	}

	statusLine := fmt.Sprintf("status: %s", singleLine(m.status))
	if m.sessionReviewed > 0 {
		accuracy := float64(m.sessionCorrect) / float64(m.sessionReviewed) * 100
		statusLine += fmt.Sprintf(" | session: %d/%d (%.0f%%)", m.sessionCorrect, m.sessionReviewed, accuracy)
	}
	footer := fmt.Sprintf("tab/arrows views | 1-9 views | ? help | q quit | mouse %d,%d %s", m.mouseX, m.mouseY, helpHint)
	b.WriteString("\n")
	b.WriteString(statusStyle.Render(truncateLine(statusLine, maxInt(20, m.width-2))))
	b.WriteString("\n")
	b.WriteString(statusStyle.Render(truncateLine(footer, maxInt(20, m.width-2))))

	if m.showHelp {
		b.WriteString("\n\n")
		b.WriteString(m.renderHelp())
	}

	finalContent := b.String()
	if m.confirmingDelete {
		finalContent = m.applyOverlay(finalContent, m.renderConfirmation())
	}

	return tea.View{
		Content:   finalContent,
		AltScreen: true,
		MouseMode: tea.MouseModeAllMotion,
	}
}

func (m *Model) renderConfirmation() string {
	style := lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		BorderForeground(lipgloss.Color("196")).
		Padding(1, 4).
		Align(lipgloss.Center)

	title := lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true).Render(" CONFIRM DELETION ")
	content := m.status
	buttons := lipgloss.NewStyle().Bold(true).Render(" [y] Yes / [n] No ")

	return style.Render(fmt.Sprintf("%s\n\n%s\n\n%s", title, content, buttons))
}

func (m *Model) renderWide() string {
	if m.width < 120 {
		nav := m.renderNav(0, 2)
		content := m.renderActiveView(20, 1)

		navLines := strings.Split(nav, "\n")
		contentLines := strings.Split(content, "\n")

		maxLines := maxInt(len(navLines), len(contentLines))
		var b strings.Builder
		for i := 0; i < maxLines; i++ {
			navLine := ""
			if i < len(navLines) {
				navLine = navLines[i]
			}
			contentLine := ""
			if i < len(contentLines) {
				contentLine = contentLines[i]
			}
			b.WriteString(padLine(navLine, 20))
			b.WriteString(contentLine)
			if i < maxLines-1 {
				b.WriteString("\n")
			}
		}
		return b.String()
	}

	nav := m.renderNav(0, 2)
	content := m.renderActiveView(20, 1)

	// Create a more informative sidebar
	var sb strings.Builder
	sb.WriteString(lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("212")).Render("Session Info") + "\n")
	sb.WriteString(fmt.Sprintf("Reviews: %d\n", m.sessionReviewed))
	if m.sessionReviewed > 0 {
		accuracy := float64(m.sessionCorrect) / float64(m.sessionReviewed) * 100
		sb.WriteString(fmt.Sprintf("Accuracy: %.0f%%\n", accuracy))
	}
	sb.WriteString("\n" + lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("212")).Render("Today") + "\n")
	sb.WriteString(fmt.Sprintf("Progress: %d/%d\n", m.stats.ReviewsToday, m.stats.DailyGoal))
	sb.WriteString(fmt.Sprintf("Streak: %d days\n", m.stats.CurrentStreak))

	detail := panelStyle.Width(24).Render(sb.String())

	navLines := strings.Split(nav, "\n")
	contentLines := strings.Split(content, "\n")
	detailLines := strings.Split(detail, "\n")

	maxLines := maxInt(len(navLines), maxInt(len(contentLines), len(detailLines)))
	var b strings.Builder
	for i := 0; i < maxLines; i++ {
		navLine := ""
		if i < len(navLines) {
			navLine = navLines[i]
		}

		contentLine := ""
		if i < len(contentLines) {
			contentLine = contentLines[i]
		}

		detailLine := ""
		if i < len(detailLines) {
			detailLine = detailLines[i]
		}

		b.WriteString(padLine(navLine, 20))
		b.WriteString(contentLine)
		b.WriteString(detailLine)
		if i < maxLines-1 {
			b.WriteString("\n")
		}
	}

	return b.String()
}

func (m *Model) renderMedium() string {
	return m.renderTabs(0, 2) + "\n" + m.renderActiveView(0, 2)
}

func (m *Model) renderCompact() string {
	return m.renderTabs(0, 2) + "\n" + m.renderActiveView(0, 2)
}

func (m *Model) renderNav(x, y int) string {
	labels := []struct {
		id   string
		view View
		text string
	}{
		{"nav-dashboard", ViewDashboard, "Dashboard"},
		{"nav-decks", ViewDecks, "Decks"},
		{"nav-review", ViewReview, "Review"},
		{"nav-statistics", ViewStatistics, "Statistics"},
		{"nav-import", ViewImport, "Import"},
		{"nav-ai", ViewAI, "AI Drafts"},
		{"nav-settings", ViewSettings, "Settings"},
		{"nav-browser", ViewBrowser, "Browser"},
		{"nav-cram", ViewCram, "Cram"},
	}

	var b strings.Builder
	b.WriteString(headerStyle.Render("deutsch-tui") + "\n\n")
	for i, l := range labels {
		style := navStyle
		if m.activeView == l.view {
			style = navActiveStyle
		}
		item := style.Render(l.text)
		b.WriteString(item + "\n")
		m.hitboxes = append(m.hitboxes, Hitbox{
			ID:     l.id,
			View:   l.view,
			X:      x,
			Y:      y + 2 + i,
			Width:  lipgloss.Width(item),
			Height: 1,
		})
	}
	return b.String()
}

func (m *Model) renderTabs(x, y int) string {
	tabs := []struct {
		id   string
		view View
		text string
	}{
		{"tab-dashboard", ViewDashboard, "Dashboard"},
		{"tab-decks", ViewDecks, "Decks"},
		{"tab-review", ViewReview, "Review"},
		{"tab-statistics", ViewStatistics, "Statistics"},
		{"tab-import", ViewImport, "Import"},
		{"tab-ai", ViewAI, "AI"},
		{"tab-settings", ViewSettings, "Settings"},
		{"tab-browser", ViewBrowser, "Browser"},
		{"tab-cram", ViewCram, "Cram"},
	}

	var renderedTabs []string
	currentX := x
	for _, t := range tabs {
		style := tabStyle
		if m.activeView == t.view {
			style = tabActiveStyle
		}
		item := style.Render(t.text)
		renderedTabs = append(renderedTabs, item)
		m.hitboxes = append(m.hitboxes, Hitbox{
			ID:     t.id,
			View:   t.view,
			X:      currentX,
			Y:      y,
			Width:  lipgloss.Width(item),
			Height: 1,
		})
		currentX += lipgloss.Width(item) + 1
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
}

func (m *Model) renderActiveView(x, y int) string {
	width, height := m.activePanelSize()
	style := panelStyle.Width(width).Height(height)
	layout := contentLayoutForStyle(style, x, y)

	return style.Render(m.renderActiveViewPlainAt(layout))
}

func (m *Model) renderActiveViewPlainAt(layout viewportLayout) string {
	var content string
	switch m.activeView {
	case ViewDashboard:
		content = m.renderDashboard(layout)
	case ViewDecks:
		content = m.renderDecks(layout)
	case ViewReview:
		content = m.renderReview(layout.X, layout.Y)
	case ViewStatistics:
		content = m.renderStatisticsAt(layout)
	case ViewImport:
		content = m.renderImport(layout.X, layout.Y)
	case ViewAI:
		content = m.renderAI(layout.X, layout.Y)
	case ViewSettings:
		content = m.renderSettings(layout.X, layout.Y)
	case ViewBrowser:
		content = m.renderBrowserAt(layout)
	case ViewCram:
		content = m.renderCramAt(layout)
	case ViewSessionSummary:
		content = m.renderSessionSummary(layout)
	default:
		content = "Unknown View"
	}

	return content
}

func (m *Model) activeViewContentLayout() viewportLayout {
	width, height := m.activePanelSize()
	return contentLayoutForStyle(panelStyle.Width(width).Height(height), 0, 0)
}

func contentLayoutForStyle(style lipgloss.Style, x, y int) viewportLayout {
	return viewportLayout{
		X:      x + style.GetMarginLeft() + style.GetBorderLeftSize() + style.GetPaddingLeft(),
		Y:      y + style.GetMarginTop() + style.GetBorderTopSize() + style.GetPaddingTop(),
		Width:  maxInt(1, style.GetWidth()-style.GetHorizontalFrameSize()),
		Height: maxInt(1, style.GetHeight()-style.GetVerticalFrameSize()),
	}
}

func (m *Model) activePanelSize() (int, int) {
	width := maxInt(30, m.width-50)
	height := maxInt(15, m.height-10)
	if m.breakpoint == BreakpointWide {
		if m.width < 120 {
			width = maxInt(40, m.width-20)
		} else {
			width = maxInt(40, m.width-44)
		}
	} else if m.breakpoint == BreakpointMedium {
		width = maxInt(30, m.width-4)
		height = maxInt(15, m.height-12)
	} else if m.breakpoint == BreakpointCompact {
		width = maxInt(20, m.width-2)
		height = maxInt(10, m.height-12)
	}
	return width, height
}

func (m *Model) renderStatus(x, y int) string {
	w := m.width
	status := m.status
	if m.drafting {
		frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
		status = frames[m.spinnerFrame%len(frames)] + " " + status
	}

	sessionStats := ""
	if m.sessionReviewed > 0 {
		accuracy := float64(m.sessionCorrect) / float64(m.sessionReviewed) * 100
		sessionStats = fmt.Sprintf(" | Session: %d rev, %.0f%% acc", m.sessionReviewed, accuracy)
	}

	res := statusStyle.Width(w).Render(fmt.Sprintf(" %s%s", status, sessionStats))
	return res
}

func (m *Model) applyOverlay(base, overlay string) string {
	baseLines := strings.Split(base, "\n")
	overlayLines := strings.Split(overlay, "\n")

	startY := (len(baseLines) - len(overlayLines)) / 2
	if startY < 0 {
		startY = 0
	}
	startX := (m.width - lipgloss.Width(overlayLines[0])) / 2
	if startX < 0 {
		startX = 0
	}

	for i, line := range overlayLines {
		if startY+i >= len(baseLines) {
			break
		}
		baseLines[startY+i] = baseLines[startY+i][:startX] + line + baseLines[startY+i][minInt(len(baseLines[startY+i]), startX+lipgloss.Width(line)):]
	}

	return strings.Join(baseLines, "\n")
}

func (m *Model) syncDecks(newDecks []core.Deck) {
	m.decks = newDecks
	if len(m.decks) > 0 {
		found := false
		for i, d := range m.decks {
			if d.ID == m.deck.ID {
				m.deckIndex = i
				m.deck = d
				found = true
				break
			}
		}
		if !found {
			if m.deckIndex >= len(m.decks) {
				m.deckIndex = maxInt(0, len(m.decks)-1)
			}
			m.deck = m.decks[m.deckIndex]
		}
		// Clamp deckCursor to prevent race conditions
		filteredLen := len(m.filteredDecks())
		if filteredLen > 0 {
			m.deckCursor = clampInt(m.deckCursor, 0, filteredLen-1)
		} else {
			m.deckCursor = 0
		}
	} else {
		m.deck = core.Deck{}
		m.deckIndex = 0
		m.deckCursor = 0
	}
}

func (m *Model) tickSpinner() tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
		return spinnerTickMsg{}
	})
}

type spinnerTickMsg struct{}
type revealTickMsg struct{}

func (m *Model) startRevealAnimation(audioPath string) tea.Cmd {
	m.revealState = RevealRevealing
	m.revealProgress = 0
	cmds := []tea.Cmd{m.tickReveal()}
	if audioPath != "" && m.autoPlayAudio {
		m.statusSeq++
		m.status = "Auto-playing audio..."
		cmds = append(cmds, m.playAudio(audioPath))
	}
	return tea.Batch(cmds...)
}

func (m *Model) tickReveal() tea.Cmd {
	return tea.Tick(time.Millisecond*60, func(t time.Time) tea.Msg {
		return revealTickMsg{}
	})
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
		if err := cmd.Run(); err != nil {
			return err
		}
		return nil
	}
}

func (m *Model) openDictionary(word string) tea.Cmd {
	if word == "" {
		return nil
	}
	return func() tea.Msg {
		url := "https://www.dict.cc/?s=" + url.QueryEscape(word)
		var cmd *exec.Cmd
		switch runtime.GOOS {
		case "darwin":
			cmd = exec.Command("open", url)
		case "windows":
			cmd = exec.Command("cmd", "/c", "start", url)
		default:
			cmd = exec.Command("xdg-open", url)
		}
		if err := cmd.Run(); err != nil {
			return err
		}
		return nil
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
