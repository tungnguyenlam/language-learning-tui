package tui

import (
	"context"
	"fmt"
	"io"
	"log"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"

	"deutsch-tui/internal/ai"
	"deutsch-tui/internal/app"
	"deutsch-tui/internal/audio"
	"deutsch-tui/internal/content"
	"deutsch-tui/internal/core"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type View string

const (
	ViewDashboard      View = "dashboard"
	ViewDictionary     View = "dictionary"
	ViewDecks          View = "decks"
	ViewReview         View = "review"
	ViewStatistics     View = "statistics"
	ViewImport         View = "import"
	ViewAnkiWeb        View = "ankiweb"
	ViewAI             View = "ai"
	ViewSettings       View = "settings"
	ViewBrowser        View = "browser"
	ViewCram           View = "cram"
	ViewPractice       View = "practice"
	ViewDebug          View = "debug"
	ViewSessionSummary View = "session_summary"
)

type PracticeSubView int

const (
	PracticeSubViewHub PracticeSubView = iota
	PracticeSubViewGender
	PracticeSubViewConjugation
	PracticeSubViewCase
	PracticeSubViewAdjective
	PracticeSubViewPreposition
	PracticeSubViewPlural
	PracticeSubViewSeparable
	PracticeSubViewNumbers
	PracticeSubViewConjunctions
	PracticeSubViewKonjunktiv
	PracticeSubViewPassive
	PracticeSubViewRelative
)

type RevealState int

const (
	RevealIdle RevealState = iota
	RevealFlipping
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

func (l viewportLayout) WithHeight(h int) viewportLayout {
	l.Height = h
	return l
}

func (l viewportLayout) WithY(y int) viewportLayout {
	l.Y = y
	return l
}

type practiceItem struct {
	Word    string
	Article string // "der", "die", "das"
	Meaning string
}

type Model struct {
	repo      core.Repository
	scheduler core.Scheduler
	// screens holds views migrated to the screen interface, keyed by View.
	// Dispatch consults this before the legacy per-view method switches.
	screens map[View]screen
	// importScreen is a typed handle to the registered Import/Export screen, so
	// shared code (the updateView reset, import hitbox actions) can reach its
	// view-local state without a type assertion.
	importScreen *importScreen
	// ankiWebScreen is a typed handle to the AnkiWeb shared-deck browser, so
	// Update can route its async messages to the screen's own state.
	ankiWebScreen *ankiWebScreen
	// ankiWeb is the shared-deck client, created lazily on first use so the app
	// never opens a network stack unless the browser is opened.
	ankiWeb                        ankiWebSearcher
	width                          int
	height                         int
	activeView                     View
	dictionaryPreviousView         View
	breakpoint                     Breakpoint
	decks                          []core.Deck
	deckIndex                      int
	deck                           core.Deck
	deckCursor                     int
	deckScroll                     int
	deckTotalLines                 int
	browserScroll                  int
	cramScroll                     int
	allDue                         []core.Card
	dueCards                       []core.Card
	cursor                         int
	revealState                    RevealState
	revealProgress                 float64
	flipProgress                   float64
	flipFrame                      int
	lastReviewedCardID             string
	lastReviewedGrade              core.ReviewGrade
	status                         string
	mouseX                         int
	mouseY                         int
	hitboxes                       []Hitbox
	prevView                       View
	viewTransitionProgress         float64
	viewTransitionFrame            int
	viewTransitioning              bool
	cardTransitionProgress         float64
	cardTransitionFrame            int
	cardTransitioning              bool
	cardTransitionDir              int // -1 for up, 1 for down
	aiProvider                     ai.Provider
	aiProviderName                 string
	dictCount                      int
	aiTemplates                    map[string]map[string]string
	aiTemplateSets                 []string
	aiTemplateIndex                int
	aiSecrets                      app.Secrets
	editingSecretKey               string // "" or "api_key"/"model"/"base_url"
	editingSecretProvider          string // "openai" or "anthropic"
	originalSecretValue            string
	onSecretsChange                func(app.Secrets)
	autoPlayAudio                  bool
	revealSpeed                    int // 0: instant, 1-10
	speechSynthesizer              audio.Synthesizer
	strictNormalization            bool
	stats                          core.Statistics
	statsLoadID                    int
	settingsCursor                 int
	editingTemplate                bool
	aiInput                        string
	draftSource                    string
	drafts                         []ai.Draft
	draftCursor                    int
	dataDir                        string
	lastBackupPath                 string
	importPath                     string
	exportPath                     string
	exportDeckID                   string
	exportTag                      string
	exportFilter                   string // e.g. "All", "Mature", "Learning"
	onConfigChange                 func(string, map[string]map[string]string, bool, bool, int)
	bookmarkFilter                 bool
	dueLoadID                      int
	originalTemplateValue          string
	mcqChoice                      int
	mcqAnswered                    bool
	mcqCorrect                     bool
	browserCards                   []core.Card
	browserCursor                  int
	browserLoadID                  int
	browserSearch                  string
	browserTag                     string
	browserSearchHistory           []string
	browserDeckID                  string
	browserSelected                map[string]bool
	dictionarySearch               string
	dictionarySearchHistory        []string
	deckSearchHistory              []string
	deckHistorySave                orderedSave
	dictionaryHistorySave          orderedSave
	dictionaryRecentSave           orderedSave
	dictionaryStarredSave          orderedSave
	dictionaryRecentVersion        uint64
	dictionarySearchID             int
	dictionaryRelatedID            int
	dictionarySearchTimerID        int
	browserSearchTimerID           int
	dictionaryResults              []core.DictionaryEntry
	dictionaryRelatedEntries       []core.DictionaryEntry
	dictionaryCursor               int
	dictionaryScroll               int
	dictionaryDetailScroll         int
	dictionaryDetailTotalLines     int
	dictionaryDetailVisibleRows    int
	dictionaryListVisibleRows      int
	dictionaryDetailView           bool
	dictionaryFocusResults         bool
	dictionaryOverlayActive        bool
	dictionaryStarred              map[string]bool
	dictionaryTargetDeckID         string
	dictionaryLastAddAttemptNoteID string
	dictionaryRecentlyViewed       []core.DictionaryEntry
	dictionaryDiscoverEntries      []core.DictionaryEntry
	compoundCache                  map[string][]content.CompoundPart
	isErrorStatus                  bool
	searchingTags                  bool
	sessionReviewed                int
	sessionCorrect                 int
	lastSessionReviewed            int
	lastSessionCorrect             int
	sessionStartTime               time.Time
	lastSessionDuration            time.Duration
	sessionGrades                  map[core.ReviewGrade]int
	showHelp                       bool
	helpScroll                     int
	helpTotalLines                 int
	helpViewportLines              int
	cramCards                      []core.Card
	cramCursor                     int
	cramType                       string
	cramLoadID                     int
	cramReviewed                   int
	cramCorrect                    int
	cramActive                     bool
	cramRevealed                   bool
	reviewsPerDay                  map[string]int
	recentDecks                    []string
	reviewHistory                  []core.ReviewLog
	reviewHistoryCard              string
	showReviewHistory              bool
	reviewPredictions              map[core.ReviewGrade]time.Duration
	reviewPredictionID             int
	spinnerFrame                   int
	deckFilter                     string
	deckSelected                   map[string]bool
	drafting                       bool
	draftCancelled                 bool // discard late draft results after Esc cancel
	statsScroll                    int
	settingsScroll                 int
	practiceScroll                 int
	statsTotalLines                int
	settingsTotalLines             int
	isDragging                     bool
	focusMode                      bool
	dragView                       View
	dragTrackStartY                int
	dragVisible                    int
	dragTotal                      int
	searchingAI                    bool
	searchingBrowser               bool
	searchingDecks                 bool
	taggingCards                   bool
	tagInput                       string
	statusSeq                      int
	confirmingDelete               bool
	deleteAction                   func() tea.Cmd
	deleteIDs                      []string
	editingDeckLimits              bool
	limitCursor                    int                // 0: new limit, 1: review limit
	gradingInProgress              bool               // Prevent double-grading
	logger                         *app.LeveledLogger // Add logger field
	typingMode                     bool               // Typing exercise mode
	typedAnswer                    string             // Current typed answer
	typingChecked                  bool               // Whether typing answer has been checked
	typingCorrect                  bool               // Whether typed answer was correct
	showHint                       bool               // Whether to show hint for the current card
	showCardInfo                   bool               // Whether to show card info overlay

	// Gender Trainer state
	practiceItems      []practiceItem
	practiceLoadID     int
	practiceIndex      int
	practiceCorrect    int
	practiceTotal      int
	practiceRevealed   bool
	practiceLastResult bool

	// trainers holds the live state for every text-input trainer
	// (conjugation, case, adjective, preposition, plural, separable, numbers,
	// conjunctions), keyed by sub-view. See trainer.go / trainer_content.go.
	trainers map[PracticeSubView]*trainerState

	practiceSubView     PracticeSubView
	practiceHubCursor   int
	practiceFilter      string
	practiceFilterFocus bool

	// Card-explanation flow: AI provides a brief pedagogical explanation.
	explainingCard bool
	explainCardID  string // target id for in-flight explain (review card or dict:…)
	explanation    string
	explainError   string

	// Grammar hint overlay
	showGrammarHint bool
	grammarHint     *content.GrammarTip

	// Card-fix flow: user reports the current Review card as wrong; AI
	// proposes a corrected note metadata; the user accepts or discards.
	fixingCard  bool          // AI request in flight
	fixCardID   string        // card the fix targets (for UI clarity)
	fixOldNote  *core.Note    // original note (snapshot before AI)
	fixProposal *ai.FixedNote // AI's proposed correction (nil until ready)
	fixError    string        // last fix error, if any

	testMode bool
}

// buildProvider returns the ai.Provider for the named choice. Unknown
// names fall through to OfflineProvider so the app still drafts cards
// even if the config file gets out of sync.
func buildProvider(name string, secrets app.Secrets, templates map[string]map[string]string, activeSet string) ai.Provider {
	switch name {
	case "disabled":
		return nil
	case "template":
		return ai.TemplateProvider{Templates: templates, ActiveSet: activeSet}
	case "openai":
		return ai.OpenAIProvider{
			APIKey:  secrets.OpenAI.APIKey,
			Model:   secrets.OpenAI.Model,
			BaseURL: secrets.OpenAI.BaseURL,
		}
	case "anthropic":
		return ai.AnthropicProvider{
			APIKey:  secrets.Anthropic.APIKey,
			Model:   secrets.Anthropic.Model,
			BaseURL: secrets.Anthropic.BaseURL,
		}
	case "ollama":
		return ai.OllamaProvider{
			Model:   secrets.Ollama.Model,
			BaseURL: secrets.Ollama.BaseURL,
		}
	default:
		return ai.OfflineProvider{}
	}
}

func NewModel(repo core.Repository, scheduler core.Scheduler) *Model {
	return NewModelWithAI(repo, scheduler, ai.OfflineProvider{})
}

func NewModelWithAI(repo core.Repository, scheduler core.Scheduler, provider ai.Provider) *Model {
	return NewModelWithOptions(repo, scheduler, ModelOptions{
		AIProvider:     provider,
		AIProviderName: "offline",
		RevealSpeed:    5,
	})
}

type ModelOptions struct {
	AIProvider          ai.Provider
	AIProviderName      string
	AITemplates         map[string]map[string]string
	AISecrets           app.Secrets
	TTSProvider         string
	TTSVoice            string
	TTSCacheDir         string
	AutoPlayAudio       bool
	RevealSpeed         int
	StrictNormalization bool
	TestMode            bool
	DataDir             string
	ImportPath          string
	ExportPath          string
	OnConfigChange      func(string, map[string]map[string]string, bool, bool, int)
	OnSecretsChange     func(app.Secrets)
	Logger              *app.LeveledLogger
}

func NewModelWithOptions(repo core.Repository, scheduler core.Scheduler, opts ModelOptions) *Model {
	providerName := opts.AIProviderName
	if providerName == "" {
		providerName = "offline"
	}
	templates := opts.AITemplates
	if len(templates) == 0 {
		templates = app.DefaultConfig().AITemplates
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
	var speechSynthesizer audio.Synthesizer
	switch strings.TrimSpace(opts.TTSProvider) {
	case audio.ProviderEdgeTTS:
		primary := audio.NewEdgeTTS(opts.TTSCacheDir, opts.TTSVoice)
		secondary := audio.NewNativeTTS(opts.TTSCacheDir)
		speechSynthesizer = audio.NewMultiTTS(primary, secondary)
	default:
		speechSynthesizer = nil
	}
	strictNormalization := opts.StrictNormalization
	provider := opts.AIProvider
	if provider == nil {
		activeSet := ""
		if len(sets) > 0 {
			activeSet = sets[aiTemplateIndex]
		}
		provider = buildProvider(providerName, opts.AISecrets, templates, activeSet)
	}
	importPath := opts.ImportPath
	if strings.TrimSpace(importPath) == "" {
		importPath = "import.tsv"
	}
	exportPath := opts.ExportPath
	if strings.TrimSpace(exportPath) == "" {
		exportPath = "export.tsv"
	}

	logger := opts.Logger
	if logger == nil {
		nullLogger := log.New(io.Discard, "", 0)
		logger = app.NewLeveledLogger(nullLogger, app.LogLevelInfo)
	}

	m := &Model{
		repo:                repo,
		scheduler:           scheduler,
		aiProvider:          provider,
		aiProviderName:      providerName,
		aiTemplates:         templates,
		aiTemplateSets:      sets,
		aiTemplateIndex:     aiTemplateIndex,
		autoPlayAudio:       autoPlayAudio,
		revealSpeed:         opts.RevealSpeed,
		speechSynthesizer:   speechSynthesizer,
		strictNormalization: strictNormalization,
		width:               80,
		height:              24,
		activeView:          ViewDashboard,
		breakpoint:          BreakpointMedium,
		status:              "Ready",
		aiInput:             "der Kaffee",
		dataDir:             strings.TrimSpace(opts.DataDir),
		importPath:          filepath.Clean(importPath),
		exportPath:          filepath.Clean(exportPath),
		exportFilter:        "All",
		onConfigChange:      opts.OnConfigChange,
		onSecretsChange:     opts.OnSecretsChange,
		aiSecrets:           opts.AISecrets,
		browserSelected:     make(map[string]bool),
		deckSelected:        make(map[string]bool),
		dictionaryStarred:   make(map[string]bool),
		logger:              logger,
		testMode:            opts.TestMode,
	}
	m.registerScreens()
	return m
}

type dueCardsMsg struct {
	id    int
	cards []core.Card
}
type bookmarkedDueCardsMsg struct {
	id    int
	cards []core.Card
}
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
type reviewRecordedMsg struct {
	cardID         string
	cards          []core.Card
	decks          []core.Deck
	stats          core.Statistics
	grade          core.ReviewGrade
	bookmarkFilter bool
}
type bookmarkToggledMsg struct {
	cardID     string
	bookmarked bool
}
type cardSuspendedMsg struct {
	cardID         string
	cards          []core.Card
	decks          []core.Deck
	stats          core.Statistics
	bookmarkFilter bool
}
type dailyGoalSetMsg core.Statistics
type reviewUndoneMsg struct {
	cardID         string
	cards          []core.Card
	decks          []core.Deck
	stats          core.Statistics
	grade          core.ReviewGrade
	bookmarkFilter bool
}
type cramCardsMsg struct {
	id       int
	cards    []core.Card
	cramType string
	deckID   string
}
type browserCardsResultMsg struct {
	id    int
	cards []core.Card
}
type reviewsPerDayMsg map[string]int
type reviewHistoryMsg struct {
	cardID string
	logs   []core.ReviewLog
}
type reviewPredictionsMsg struct {
	id          int
	cardID      string
	predictions map[core.ReviewGrade]time.Duration
}
type statusMsg struct {
	text string
}
type tagsUpdatedMsg struct {
	cardIDs []string
	tags    []string
}
type explainMsg struct {
	cardID      string
	explanation string
}
type explainErrorMsg struct {
	cardID string
	err    error
}

type debounceSearchMsg struct {
	id   int
	view View
}

type deckDeletedMsg struct{}

type dictHistoryLoadedMsg []string
type dictRecentlyViewedLoadedMsg []core.DictionaryEntry
type deckHistoryLoadedMsg []string
type dictStarredLoadedMsg map[string]bool

func (m *Model) Init() tea.Cmd {
	return tea.Sequence(m.loadDueCards(), m.loadDecks, m.loadStatistics(), m.loadReviewsPerDay(), m.loadRecentDecks(), m.loadDictionaryHistory(), m.loadDictionaryRecentlyViewed(), m.loadDeckHistory(), m.loadDictionaryStarred())
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.logger.Debug("Processing message: %T", msg)

	// The AnkiWeb browser keeps its results on the screen rather than on Model,
	// so its async replies are applied by the screen itself.
	if handled, cmd := m.handleAnkiWebMsg(msg); handled {
		return m, cmd
	}

	// Domain async replies (loads, grading, drafts, animation ticks) are
	// dispatched in messages.go; only resize and input messages remain here.
	if cmd, handled := m.handleAsyncMsg(msg); handled {
		return m, cmd
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.breakpoint = breakpointForWidth(msg.Width)
		m.isDragging = false
		m.logger.Debug("Window resized to %dx%d, breakpoint: %s", msg.Width, msg.Height, m.breakpoint)
		return m, nil

	case tea.KeyPressMsg:
		m.logger.Debug("Key pressed: %s", msg.String())
		return m.updateKey(msg)
	case tea.PasteMsg:
		m.logger.Debug("Pasted text: %s", msg.String())
		return m, m.handlePaste(msg.String())
	case tea.MouseMsg:
		return m, m.handleMouseMsg(msg)
	}
	return m, nil
}

func (m *Model) View() tea.View {
	if m.width < 20 || m.height < 10 {
		return tea.View{Content: "Terminal too small"}
	}

	m.hitboxes = m.hitboxes[:0]
	var b strings.Builder

	header := lipgloss.NewStyle().
		Bold(true).
		Foreground(colorAccent).
		Render(" DEUTSCH-TUI ") +
		lipgloss.NewStyle().Foreground(colorPanel).Render("│") + " " +
		lipgloss.NewStyle().Bold(true).Foreground(colorBlue).Render(strings.ToUpper(string(m.activeView))) + " " +
		lipgloss.NewStyle().Foreground(colorPanel).Render("│") + " " +
		mutedStyle.Render(string(m.breakpoint))
	b.WriteString(header)
	b.WriteString("\n")

	switch m.breakpoint {
	case BreakpointWide:
		b.WriteString(m.renderWide())
	default:
		b.WriteString(m.renderStacked())
	}

	helpHint := ""
	switch m.activeView {
	case ViewDashboard:
		helpHint = "| / search dictionary | [ ] switch decks | 3 Review"
	case ViewDictionary:
		if m.dictionarySearch != "" {
			helpHint = "| Nav: j/k | Shift+↑/↓ details | ctrl+a add | Enter draft | ctrl+f find"
		} else {
			if len(m.dictionarySearchHistory) > 0 {
				helpHint = "| Search: type | ctrl+x: clear history | Esc: back"
			} else {
				helpHint = "| Search: type | Nav: j/k | Esc: back"
			}
		}
	case ViewReview:
		if len(m.dueCards) > 0 {
			if m.revealState == RevealRevealed {
				helpHint = "| Grade: a/h/g/e"
			} else {
				helpHint = "| Reveal: Space/Enter"
			}
		}
	case ViewDecks:
		if m.searchingDecks {
			helpHint = "| Searching... | Stop: Esc/enter"
		} else {
			helpHint = "| Select: enter | Multi-select: m | Search: /"
		}
	case ViewStatistics:
		helpHint = "| Scroll: j/k | Export: x"
	case ViewCram:
		if m.cramActive {
			if m.cramRevealed {
				helpHint = "| Grade: a/h/g/e | Quit: q"
			} else {
				helpHint = "| Reveal: Space/enter"
			}
		} else {
			helpHint = "| Filter: 1-5 | Start: enter"
		}
	case ViewPractice:
		if m.practiceSubView == PracticeSubViewHub {
			helpHint = "| Nav: j/k | Select: enter | Trainer: 1-9,0,-,= | Reset: r | Esc: back"
		} else {
			switch m.practiceSubView {
			case PracticeSubViewGender:
				helpHint = "| der(1/d) die(2/i) das(3/a) | Next click | Esc: back"
			case PracticeSubViewConjugation:
				helpHint = "| Type answer | Enter check | Esc: back"
			case PracticeSubViewCase, PracticeSubViewAdjective, PracticeSubViewPreposition:
				helpHint = "| Type answer | h: hint | Enter check | Esc: back"
			case PracticeSubViewSeparable, PracticeSubViewNumbers:
				helpHint = "| Type answer | Enter check | Esc: back"
			case PracticeSubViewConjunctions, PracticeSubViewKonjunktiv, PracticeSubViewPassive, PracticeSubViewRelative:
				helpHint = "| Type answer | h: hint | Enter check | Esc: back"
			case PracticeSubViewPlural:
				helpHint = "| Type plural | Enter check | Esc: back"
			}
		}
	case ViewBrowser:
		if m.searchingBrowser {
			helpHint = "| Searching... | Stop: Esc/enter"
		} else {
			helpHint = "| Nav: j/k | Select: m | History: enter | Search: /"
		}
	case ViewSettings:
		helpHint = "| Nav: j/k | Edit: enter"
	case ViewImport:
		helpHint = "| Import: enter | Export: x"
	case ViewAI:
		helpHint = "| Generate: enter | Approve: a"
	}

	statusText := singleLine(m.status)
	if m.isErrorStatus {
		statusText = lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true).Render(statusText)
	}
	statusLine := fmt.Sprintf("status: %s", statusText)
	if m.sessionReviewed > 0 {
		accuracy := float64(m.sessionCorrect) / float64(m.sessionReviewed) * 100
		statusLine += fmt.Sprintf(" | session: %d/%d (%.0f%%)", m.sessionCorrect, m.sessionReviewed, accuracy)
	}

	keyInfoStyle := lipgloss.NewStyle().Foreground(colorPink).Bold(true)
	footerParts := []string{
		keyInfoStyle.Render("tab/arrows") + " views",
		keyInfoStyle.Render("0-9") + " views",
		keyInfoStyle.Render("?") + " help",
		keyInfoStyle.Render("q") + " quit",
	}
	if m.mouseX > 0 || m.mouseY > 0 {
		footerParts = append(footerParts, fmt.Sprintf("mouse %d,%d", m.mouseX, m.mouseY))
	}
	if helpHint != "" {
		for _, part := range strings.Split(helpHint, "|") {
			p := strings.TrimSpace(part)
			if p != "" {
				footerParts = append(footerParts, lipgloss.NewStyle().Foreground(colorMuted).Render(p))
			}
		}
	}

	b.WriteString("\n")
	b.WriteString(statusStyle.Render(truncateLine(statusLine, maxInt(20, m.width-2))))
	b.WriteString("\n")
	b.WriteString(formatWrappedFooter(footerParts, maxInt(20, m.width-2), statusStyle))

	finalContent := b.String()
	if m.showHelp {
		// Lipgloss adds the modal's padding and border outside Width. Keep the
		// inner width narrow enough that the complete frame always fits inside
		// the terminal, including on the 80-column test viewport.
		helpLayoutWidth := maxInt(20, minInt(114, m.width-14))
		helpLayoutHeight := maxInt(8, m.height-8)
		helpContent := m.renderHelpViewport(viewportLayout{
			X:      (m.width - helpLayoutWidth) / 2,
			Y:      0,
			Width:  helpLayoutWidth,
			Height: helpLayoutHeight,
		})
		helpBox := lipgloss.NewStyle().
			Width(helpLayoutWidth).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("81")).
			Padding(1, 2).
			Background(lipgloss.Color("233")).
			Render(helpContent)
		helpView := lipgloss.Place(m.width, m.height-3, lipgloss.Center, lipgloss.Center, helpBox)

		statusText := singleLine(m.status)
		if m.isErrorStatus {
			statusText = lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true).Render(statusText)
		}
		statusLine := fmt.Sprintf("status: %s", statusText)
		footer := strings.Join([]string{
			keyInfoStyle.Render("j/k") + " scroll",
			keyInfoStyle.Render("PgUp/PgDn") + " pages",
			keyInfoStyle.Render("Esc/?") + " close",
			keyInfoStyle.Render("Ctrl+c") + " quit",
		}, " │ ")
		helpScreen := helpView + "\n\n" + statusStyle.Render(truncateLine(statusLine, m.width-2)) + "\n" + statusStyle.Render(truncateLine(footer, m.width-2))
		// The underlying view can contain longer text on the same rows. Fill
		// every modal row to the terminal width so the overlay fully erases it
		// in terminals that redraw by writing changed cells.
		screenLines := strings.Split(helpScreen, "\n")
		for i, line := range screenLines {
			screenLines[i] = padLine(truncateLine(line, m.width), m.width)
		}

		return tea.View{
			Content:   strings.Join(screenLines, "\n"),
			AltScreen: true,
			MouseMode: tea.MouseModeAllMotion,
		}
	}

	if m.confirmingDelete {
		finalContent = m.applyOverlay(finalContent, m.renderConfirmation())
	}

	if m.dictionaryOverlayActive {
		finalContent = m.applyOverlay(finalContent, m.renderSpotlightDictionary())
	}

	return tea.View{
		Content:   finalContent,
		AltScreen: true,
		MouseMode: tea.MouseModeAllMotion,
	}
}

func formatWrappedFooter(parts []string, width int, style lipgloss.Style) string {
	if len(parts) == 0 {
		return ""
	}
	maxWidth := maxInt(20, width)
	var lines []string
	var currentParts []string
	currentLen := 0

	for _, part := range parts {
		partLen := lipgloss.Width(part)
		delimLen := 0
		if len(currentParts) > 0 {
			delimLen = 3 // " │ "
		}

		if len(currentParts) > 0 && currentLen+delimLen+partLen > maxWidth {
			lines = append(lines, style.Render(strings.Join(currentParts, " │ ")))
			currentParts = []string{part}
			currentLen = partLen
		} else {
			if len(currentParts) > 0 {
				currentLen += delimLen
			}
			currentParts = append(currentParts, part)
			currentLen += partLen
		}
	}
	if len(currentParts) > 0 {
		lines = append(lines, style.Render(strings.Join(currentParts, " │ ")))
	}
	return strings.Join(lines, "\n")
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
	if m.focusMode && m.activeView == ViewReview {
		return "\n" + m.renderActiveView(0, 2)
	}
	if m.width < 120 {
		nav := m.renderNav(0, 1)
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

	_, height := m.activePanelSize()
	nav := m.renderNav(0, 1)
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

	detail := panelStyle.Width(24).Height(height).Render(sb.String())

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

func (m *Model) renderStacked() string {
	if m.focusMode && m.activeView == ViewReview {
		return "\n" + m.renderActiveView(0, 2)
	}
	return m.renderTabs(0, 1) + "\n" + m.renderActiveView(0, 2)
}

func (m *Model) renderActiveView(x, y int) string {
	width, height := m.activePanelSize()
	style := panelStyle.Width(width).Height(height)
	layout := contentLayoutForStyle(style, x, y)

	return style.Render(m.renderActiveViewPlainAt(layout))
}

func (m *Model) renderActiveViewPlainAt(layout viewportLayout) string {
	if s, ok := m.screens[m.activeView]; ok {
		return s.Render(m, layout)
	}
	return "Unknown View"
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
	if m.focusMode && m.activeView == ViewReview {
		if m.height > 30 {
			return m.width - 4, maxInt(10, m.height-6)
		}
		return m.width - 4, maxInt(10, m.height-8)
	}
	width := maxInt(30, m.width-50)
	height := maxInt(12, m.height-10)
	if m.height > 30 {
		height = maxInt(15, m.height-6)
	}

	if m.breakpoint == BreakpointWide {
		if m.width < 120 {
			width = maxInt(40, m.width-20)
		} else {
			width = maxInt(40, m.width-44)
		}
	} else if m.breakpoint == BreakpointMedium {
		width = maxInt(30, m.width-4)
		if m.height <= 30 {
			height = maxInt(10, m.height-12)
		}
	} else if m.breakpoint == BreakpointCompact {
		width = maxInt(20, m.width-2)
		if m.height <= 30 {
			height = maxInt(10, m.height-12)
		}
	}
	return width, height
}

func (m *Model) applyOverlay(base, overlay string) string {
	baseLines := strings.Split(base, "\n")
	overlayLines := strings.Split(strings.TrimSuffix(overlay, "\n"), "\n")
	if len(overlayLines) == 0 {
		return base
	}

	startY := (len(baseLines) - len(overlayLines)) / 2
	if startY < 0 {
		startY = 0
	}

	overlayWidth := 0
	for _, line := range overlayLines {
		if w := lipgloss.Width(line); w > overlayWidth {
			overlayWidth = w
		}
	}
	startX := (m.width - overlayWidth) / 2
	if startX < 0 {
		startX = 0
	}

	for i, line := range overlayLines {
		if startY+i >= len(baseLines) {
			break
		}
		padded := padLine(line, overlayWidth)
		baseLines[startY+i] = spliceVisual(baseLines[startY+i], startX, overlayWidth, padded)
	}

	return strings.Join(baseLines, "\n")
}

// spliceVisual replaces the visual columns [startCol, startCol+width) in
// line (which may contain ANSI escape sequences) with replacement.
// ANSI escape sequences in the line are preserved: those entirely before
// startCol stay in the prefix, those entirely after startCol+width stay in
// the suffix, and any escape that spans a boundary is kept intact on the
// side that contains more of it.
func spliceVisual(line string, startCol int, width int, replacement string) string {
	prefix, _ := visualPrefix(line, startCol)
	remaining := line[len(prefix):]
	_, afterWidth := visualPrefix(remaining, width)
	return prefix + replacement + afterWidth
}

// visualPrefix splits s at the given visual column boundary.
// It returns (prefix, suffix) where prefix contains all of s up to (but not
// including) visual column col, preserving any ANSI escape sequences that
// start in the prefix. The suffix starts at the visual column boundary.
// ANSI escape sequences that appear immediately after the last visual
// character are included in the prefix to properly close styling.
func visualPrefix(s string, col int) (prefix, suffix string) {
	var buf strings.Builder
	visualCol := 0
	i := 0
	runes := []rune(s)
	consumedVisual := false
	for i < len(runes) && visualCol < col {
		r := runes[i]
		if r == '\x1b' {
			end := findAnsiEnd(runes, i)
			buf.WriteString(string(runes[i:end]))
			i = end
			continue
		}
		w := runeVisualWidth(r)
		if visualCol+w > col {
			break
		}
		buf.WriteRune(r)
		visualCol += w
		consumedVisual = true
		i++
	}
	// Consume any trailing ANSI escape sequences immediately after the boundary,
	// but only if we've consumed at least one visual character
	if consumedVisual {
		for i < len(runes) && runes[i] == '\x1b' {
			end := findAnsiEnd(runes, i)
			buf.WriteString(string(runes[i:end]))
			i = end
		}
	}
	return buf.String(), string(runes[i:])
}

// findAnsiEnd returns the index after the end of the ANSI escape sequence
// starting at runes[start].
func findAnsiEnd(runes []rune, start int) int {
	if start >= len(runes) || runes[start] != '\x1b' {
		return start + 1
	}
	i := start + 1
	if i >= len(runes) {
		return i
	}
	switch runes[i] {
	case '[':
		i++
		for i < len(runes) {
			r := runes[i]
			if r >= '0' && r <= '9' || r == ';' || r == '?' || r == ':' || r == ' ' {
				i++
			} else if (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') || r == '@' {
				i++
				return i
			} else {
				i++
				return i
			}
		}
		return i
	case ']':
		i++
		for i < len(runes) {
			if runes[i] == '\x07' {
				return i + 1
			}
			if runes[i] == '\x1b' && i+1 < len(runes) && runes[i+1] == '\\' {
				return i + 2
			}
			i++
		}
		return i
	case '(', ')', '*', '+':
		if i+1 < len(runes) {
			return i + 2
		}
		return i + 1
	default:
		return i + 1
	}
}

func runeVisualWidth(r rune) int {
	switch {
	case r >= 0x1100 && (r <= 0x115F || r == 0x2329 || r == 0x232A ||
		(r >= 0x2E80 && r <= 0xA4CF) ||
		(r >= 0xAC00 && r <= 0xD7A3) ||
		(r >= 0xF900 && r <= 0xFAFF) ||
		(r >= 0xFE10 && r <= 0xFE19) ||
		(r >= 0xFE30 && r <= 0xFE6F) ||
		(r >= 0xFF01 && r <= 0xFF60) ||
		(r >= 0xFFE0 && r <= 0xFFE6) ||
		(r >= 0x20000 && r <= 0x2FFFD) ||
		(r >= 0x30000 && r <= 0x3FFFD)):
		return 2
	default:
		return 1
	}
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
type flipTickMsg struct{}
type viewTransitionTickMsg struct{}
type cardTransitionTickMsg struct{}

func (m *Model) startRevealAnimation(card core.Card) tea.Cmd {
	if m.revealSpeed == 0 {
		m.revealState = RevealRevealed
		m.revealProgress = 100
		if m.activeView == ViewCram {
			m.cramRevealed = true
		}
		if m.autoPlayAudio {
			return m.playCardAudio(card)
		}
		return nil
	}
	m.revealState = RevealFlipping
	m.flipProgress = 0
	m.flipFrame = 0
	cmds := []tea.Cmd{m.tickFlip()}
	if m.autoPlayAudio {
		cmds = append(cmds, m.playCardAudio(card))
	}
	return tea.Batch(cmds...)
}

func (m *Model) tickReveal() tea.Cmd {
	return tea.Tick(time.Millisecond*60, func(t time.Time) tea.Msg {
		return revealTickMsg{}
	})
}

func (m *Model) tickFlip() tea.Cmd {
	return tea.Tick(time.Millisecond*30, func(t time.Time) tea.Msg {
		return flipTickMsg{}
	})
}

func (m *Model) tickViewTransition() tea.Cmd {
	return tea.Tick(time.Millisecond*30, func(t time.Time) tea.Msg {
		return viewTransitionTickMsg{}
	})
}

func (m *Model) tickCardTransition() tea.Cmd {
	return tea.Tick(time.Millisecond*30, func(t time.Time) tea.Msg {
		return cardTransitionTickMsg{}
	})
}

func (m *Model) playAudio(audioPath string) tea.Cmd {
	if audioPath == "" {
		return nil
	}
	return func() tea.Msg {
		if err := startAudioPlayer(audioPath); err != nil {
			return err
		}
		return nil
	}
}

func (m *Model) playCardAudio(card core.Card) tea.Cmd {
	if card.Audio != "" {
		m.statusSeq++
		m.status = "Playing audio..."
		return m.playAudio(card.Audio)
	}
	if m.speechSynthesizer == nil {
		m.status = "No audio for this card"
		return nil
	}
	text := speechTextForCard(card)
	if strings.TrimSpace(text) == "" {
		m.status = "No text available for TTS"
		return nil
	}
	m.statusSeq++
	m.status = fmt.Sprintf("Generating %s TTS audio...", m.speechSynthesizer.VoiceName())
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
		defer cancel()
		path, err := m.speechSynthesizer.Synthesize(ctx, text)
		if err != nil {
			return err
		}
		if err := startAudioPlayer(path); err != nil {
			return err
		}
		return nil
	}
}

func (m *Model) playDictionaryAudio(word string) tea.Cmd {
	if m.speechSynthesizer == nil {
		m.status = "Audio not available (no TTS configured)"
		return nil
	}
	if strings.TrimSpace(word) == "" {
		m.status = "No word selected"
		return nil
	}
	m.statusSeq++
	m.status = fmt.Sprintf("Generating %s TTS audio...", m.speechSynthesizer.VoiceName())
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
		defer cancel()
		path, err := m.speechSynthesizer.Synthesize(ctx, word)
		if err != nil {
			return err
		}
		if err := startAudioPlayer(path); err != nil {
			return err
		}
		return nil
	}
}

func startAudioPlayer(audioPath string) error {
	if audioPath == "" {
		return nil
	}
	spec, err := selectAudioPlayer(audioPath, runtime.GOOS, exec.LookPath)
	if err != nil {
		return err
	}
	cmd := exec.Command(spec.name, spec.args...)
	if err := cmd.Start(); err != nil {
		return err
	}
	go cmd.Wait()
	return nil
}

type audioPlayerSpec struct {
	name string
	args []string
}

func selectAudioPlayer(audioPath, goos string, lookPath func(string) (string, error)) (audioPlayerSpec, error) {
	if strings.TrimSpace(audioPath) == "" {
		return audioPlayerSpec{}, nil
	}
	if lookPath == nil {
		lookPath = exec.LookPath
	}
	switch goos {
	case "darwin":
		return firstAvailablePlayer(lookPath, []audioPlayerSpec{
			{name: "afplay", args: []string{audioPath}},
			{name: "mpv", args: []string{"--no-terminal", "--really-quiet", audioPath}},
			{name: "ffplay", args: []string{"-nodisp", "-autoexit", "-loglevel", "quiet", audioPath}},
			{name: "play", args: []string{audioPath}},
		})
	case "windows":
		return firstAvailablePlayer(lookPath, []audioPlayerSpec{
			{name: "mpv", args: []string{"--no-terminal", "--really-quiet", audioPath}},
			{name: "ffplay", args: []string{"-nodisp", "-autoexit", "-loglevel", "quiet", audioPath}},
			windowsPowerShellPlayer("powershell.exe", audioPath),
			windowsPowerShellPlayer("powershell", audioPath),
			windowsPowerShellPlayer("pwsh", audioPath),
		})
	default:
		return firstAvailablePlayer(lookPath, []audioPlayerSpec{
			{name: "mpv", args: []string{"--no-terminal", "--really-quiet", audioPath}},
			{name: "ffplay", args: []string{"-nodisp", "-autoexit", "-loglevel", "quiet", audioPath}},
			{name: "play", args: []string{audioPath}},
			{name: "paplay", args: []string{audioPath}},
			{name: "aplay", args: []string{audioPath}},
		})
	}
}

func firstAvailablePlayer(lookPath func(string) (string, error), candidates []audioPlayerSpec) (audioPlayerSpec, error) {
	checked := make([]string, 0, len(candidates))
	for _, candidate := range candidates {
		if candidate.name == "" {
			continue
		}
		checked = append(checked, candidate.name)
		path, err := lookPath(candidate.name)
		if err == nil {
			candidate.name = path
			return candidate, nil
		}
	}
	return audioPlayerSpec{}, fmt.Errorf("no audio player found; install one of: %s", strings.Join(checked, ", "))
}

func windowsPowerShellPlayer(name, audioPath string) audioPlayerSpec {
	script := `$ErrorActionPreference = 'Stop'; ` +
		`Add-Type -AssemblyName PresentationCore; ` +
		`$resolved = (Resolve-Path -LiteralPath $args[0]).ProviderPath; ` +
		`$player = New-Object System.Windows.Media.MediaPlayer; ` +
		`$player.Open([Uri]::new($resolved)); ` +
		`while (-not $player.NaturalDuration.HasTimeSpan) { Start-Sleep -Milliseconds 50 }; ` +
		`$player.Play(); ` +
		`while ($player.Position -lt $player.NaturalDuration.TimeSpan) { Start-Sleep -Milliseconds 100 }`
	return audioPlayerSpec{
		name: name,
		args: []string{"-NoProfile", "-ExecutionPolicy", "Bypass", "-Command", script, audioPath},
	}
}

func speechTextForCard(card core.Card) string {
	if card.Kind == core.CardKindCloze && card.Answer != "" {
		return card.Answer
	}
	if strings.HasSuffix(card.ID, ":back") && card.Answer != "" {
		return card.Answer
	}
	if card.Prompt != "" {
		return card.Prompt
	}
	return card.Answer
}

func (m *Model) ttsAvailable() bool {
	return m.speechSynthesizer != nil
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
