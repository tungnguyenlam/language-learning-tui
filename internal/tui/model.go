package tui

import (
	"context"
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
	dictionaryProvider             string
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
	settingsCursor                 int
	editingTemplate                bool
	aiInput                        string
	draftSource                    string
	drafts                         []ai.Draft
	draftCursor                    int
	importPath                     string
	exportPath                     string
	exportDeckID                   string
	exportTag                      string
	exportFilter                   string // e.g. "All", "Mature", "Learning"
	theme                          string
	onConfigChange                 func(string, string, string, map[string]map[string]string, bool, bool, int)
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
	Theme               string
	AIProvider          ai.Provider
	AIProviderName      string
	DictionaryProvider  string
	AITemplates         map[string]map[string]string
	AISecrets           app.Secrets
	TTSProvider         string
	TTSVoice            string
	TTSCacheDir         string
	AutoPlayAudio       bool
	RevealSpeed         int
	StrictNormalization bool
	TestMode            bool
	ImportPath          string
	ExportPath          string
	OnConfigChange      func(string, string, string, map[string]map[string]string, bool, bool, int)
	OnSecretsChange     func(app.Secrets)
	Logger              *app.LeveledLogger
}

func NewModelWithOptions(repo core.Repository, scheduler core.Scheduler, opts ModelOptions) *Model {
	providerName := opts.AIProviderName
	if providerName == "" {
		providerName = "offline"
	}
	templates := opts.AITemplates
	if templates == nil || len(templates) == 0 {
		templates = map[string]map[string]string{
			"vocabulary": {
				"front":   "{{.Topic}}",
				"back":    "Translation: German prompt for {{.Topic}}.\nPlural: die {{.Topic}}e (example)\nGender: der/die/das",
				"example": "Ich lerne {{.Topic}}.",
			},
			"phrases": {
				"front":   "Common German phrase for {{.Topic}}",
				"back":    "English translation",
				"example": "Context sentence using the phrase.",
			},
			"grammar": {
				"front":   "Ich {{c1::...}} {{.Topic}}.",
				"back":    "Grammar: {{.Topic}}\nRule: Explanation of the rule for {{.Topic}}.",
				"example": "Ich {{c1::bin}} {{.Topic}}.",
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

	// Create a default logger if none provided
	logger := opts.Logger
	if logger == nil {
		// Create a minimal logger that discards output
		nullLogger := log.New(io.Discard, "", 0)
		logger = app.NewLeveledLogger(nullLogger, app.LogLevelInfo)
	}

	m := &Model{
		repo:                repo,
		scheduler:           scheduler,
		theme:               opts.Theme,
		aiProvider:          provider,
		aiProviderName:      providerName,
		dictionaryProvider:  opts.DictionaryProvider,
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
		importPath:          filepath.Clean(importPath),
		exportPath:          filepath.Clean(exportPath),
		exportFilter:        "All",
		onConfigChange:      opts.OnConfigChange,
		onSecretsChange:     opts.OnSecretsChange,
		aiSecrets:           opts.AISecrets,
		browserSelected:     make(map[string]bool),
		deckSelected:        make(map[string]bool),
		dictionaryStarred:   make(map[string]bool),
		logger:              logger, // Set the logger
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
type exportDoneMsg struct {
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
type browserCardsMsg []core.Card
type browserCardsResultMsg struct {
	id    int
	cards []core.Card
}
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

	switch msg := msg.(type) {
	case debounceSearchMsg:
		switch msg.view {
		case ViewDictionary:
			if msg.id == m.dictionarySearchTimerID {
				return m, m.searchDictionary()
			}
		case ViewBrowser:
			if msg.id == m.browserSearchTimerID {
				return m, m.loadBrowserCards()
			}
		}
		return m, nil

	case dictionarySearchResultsMsg:
		if msg.id != m.dictionarySearchID {
			return m, nil
		}
		m.dictionaryResults = msg.results
		m.dictionaryRelatedID++
		m.dictionaryRelatedEntries = nil
		m.dictionaryCursor = 0
		m.dictionaryScroll = 0
		m.dictionaryDetailScroll = 0
		query := strings.TrimSpace(m.dictionarySearch)
		if query == "" {
			m.status = "Dictionary search cleared"
		} else if len(msg.results) > 0 {
			m.status = fmt.Sprintf("Found %d dictionary results", len(msg.results))
		} else {
			m.status = fmt.Sprintf("No dictionary results for %q", query)
		}

		var cmd tea.Cmd
		wide := m.width > 80
		if (m.dictionaryDetailView || wide) && len(msg.results) > 0 {
			cmd = m.findRelatedEntries(msg.results[0].Word)
		}
		return m, cmd
	case dictHistoryLoadedMsg:
		m.dictionarySearchHistory = msg
	case dictRecentlyViewedLoadedMsg:
		m.dictionaryRecentlyViewed = msg
	case dictRelatedEntriesMsg:
		if msg.id != m.dictionaryRelatedID || (m.dictionaryCursor >= 0 && m.dictionaryCursor < len(m.dictionaryResults) && m.dictionaryResults[m.dictionaryCursor].Word != msg.word) {
			return m, nil
		}
		m.dictionaryRelatedEntries = msg.entries
		return m, nil
	case dictDiscoverEntriesMsg:
		m.dictionaryDiscoverEntries = msg.entries
		return m, nil
	case dictStarredLoadedMsg:
		if msg != nil {
			m.dictionaryStarred = msg
		}
	case deckHistoryLoadedMsg:
		m.deckSearchHistory = msg
	case spinnerTickMsg:
		if m.drafting {
			m.spinnerFrame++
			return m, m.tickSpinner()
		}
		return m, nil
	case revealTickMsg:
		if m.revealState == RevealRevealing {
			step := 10
			if m.revealSpeed > 0 {
				step = m.revealSpeed * 2
			}
			m.revealProgress += float64(step)
			if m.revealProgress >= 100 {
				m.revealProgress = 100
				m.revealState = RevealRevealed
				if m.activeView == ViewCram {
					m.cramRevealed = true
				}
				m.logger.Debug("Card fully revealed")
			} else {
				return m, m.tickReveal()
			}
		}
	case flipTickMsg:
		if m.revealState == RevealFlipping {
			m.flipFrame++
			m.flipProgress += 10
			if m.flipProgress >= 100 {
				m.flipProgress = 100
				m.revealState = RevealRevealing
				m.revealProgress = 0
				m.logger.Debug("Flip complete, starting reveal")
				return m, m.tickReveal()
			}
			return m, m.tickFlip()
		}
	case viewTransitionTickMsg:
		if m.viewTransitioning {
			m.viewTransitionFrame++
			m.viewTransitionProgress += 10
			if m.viewTransitionProgress >= 100 {
				m.viewTransitionProgress = 100
				m.viewTransitioning = false
				m.prevView = ""
				m.logger.Debug("View transition complete")
			} else {
				return m, m.tickViewTransition()
			}
		}
	case cardTransitionTickMsg:
		if m.cardTransitioning {
			m.cardTransitionFrame++
			m.cardTransitionProgress += 15
			if m.cardTransitionProgress >= 100 {
				m.cardTransitionProgress = 100
				m.cardTransitioning = false
				m.cardTransitionDir = 0
				m.logger.Debug("Card transition complete")
			} else {
				return m, m.tickCardTransition()
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
		m.draftCancelled = false
		m.gradingInProgress = false
		m.isErrorStatus = true
		m.status = friendlyError(msg)
		m.logger.Error("Error occurred: %v", msg)
	case decksMsg:
		m.logger.Debug("Received %d decks", len(msg))
		m.syncDecks([]core.Deck(msg))
	case dueCardsMsg:
		if msg.id != m.dueLoadID || m.bookmarkFilter {
			return m, nil
		}
		m.logger.Debug("Received %d due cards", len(msg.cards))
		m.allDue = msg.cards
		m.applyDeckFilter()
		m.showHint = false
	case bookmarkedDueCardsMsg:
		if msg.id != m.dueLoadID || !m.bookmarkFilter {
			return m, nil
		}
		m.allDue = msg.cards
		m.applyDeckFilter()
		m.showHint = false
		if len(m.allDue) == 0 {
			m.status = "No bookmarked cards due"
		} else {
			m.status = fmt.Sprintf("%d bookmarked cards due", len(m.allDue))
		}
		m.logger.Debug("Received %d bookmarked due cards", len(msg.cards))

	case statsMsg:
		m.stats = msg.stats
		m.dictCount = msg.dictCount
		m.logger.Debug("Received statistics update")
	case reviewsPerDayMsg:
		m.reviewsPerDay = map[string]int(msg)
		m.logger.Debug("Received reviews per day data")
	case recentDecksMsg:
		m.recentDecks = []string(msg)
		m.logger.Debug("Received %d recent decks", len(msg))
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
		if m.draftCancelled {
			m.draftCancelled = false
			return m, nil
		}
		m.drafts = []ai.Draft(msg)
		m.draftCursor = clampInt(m.draftCursor, 0, maxInt(0, len(m.drafts)-1))
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
		m.status = "Draft saved"
		m.logger.Info("Approved AI draft %s", msg.noteID)
	case importDoneMsg:
		if msg.path == "AI Drafts" {
			m.drafts = nil
			m.draftCursor = 0
		}
		m.syncDecks(msg.decks)
		m.allDue = msg.cards
		m.applyDeckFilter()
		m.logger.Info("Import completed: %d notes from %s", msg.count, filepath.Base(msg.path))
		return m, m.setStatus(fmt.Sprintf("Imported %d notes from %s", msg.count, filepath.Base(msg.path)), 3*time.Second)
	case statusMsg:
		m.logger.Debug("Setting status: %s", msg.text)
		m.isErrorStatus = false
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
	case fixProposalMsg:
		if !m.fixingCard || msg.cardID != m.fixCardID {
			return m, nil
		}
		m.fixingCard = false
		note := msg.oldNote
		m.fixOldNote = &note
		prop := msg.proposal
		m.fixProposal = &prop
		m.fixError = ""
		m.logger.Info("AI returned fix proposal for card %s", msg.cardID)
		m.status = "Review the proposed fix: y to apply, n to discard"
		return m, nil
	case fixErrorMsg:
		if !m.fixingCard || msg.cardID != m.fixCardID {
			return m, nil
		}
		m.fixingCard = false
		m.fixOldNote = nil
		m.fixProposal = nil
		m.fixError = msg.err.Error()
		m.logger.Error("Fix failed for card %s: %v", msg.cardID, msg.err)
		return m, m.setStatus("Fix failed: "+msg.err.Error(), 5*time.Second)
	case fixAppliedMsg:
		m.allDue = msg.cards
		m.applyDeckFilter()
		m.logger.Info("Applied fix for card %s", msg.cardID)
		return m, m.setStatus("Card updated by AI", 3*time.Second)
	case explainMsg:
		if !m.explainingCard || msg.cardID != m.explainCardID {
			return m, nil
		}
		m.explainingCard = false
		m.explanation = msg.explanation
		m.explainError = ""
		m.logger.Info("AI returned explanation for card %s", msg.cardID)
		m.status = "Pedagogical explanation received"
		return m, nil
	case explainErrorMsg:
		if !m.explainingCard || msg.cardID != m.explainCardID {
			return m, nil
		}
		m.explainingCard = false
		m.explanation = ""
		m.explainError = msg.err.Error()
		m.logger.Error("Explanation failed for card %s: %v", msg.cardID, msg.err)
		return m, m.setStatus("Explanation failed: "+msg.err.Error(), 5*time.Second)
	case deckDeletedMsg:
		m.logger.Info("Deck deleted")
		return m, tea.Batch(m.loadDecks, m.loadDueCards())
	case practiceItemsMsg:
		// Ignore Gender loads that finished after leaving the trainer or
		// after a newer Gender load was started.
		if msg.id != m.practiceLoadID || m.practiceSubView != PracticeSubViewGender {
			return m, nil
		}
		m.practiceItems = msg.items
		if len(m.practiceItems) == 0 {
			m.status = "No nouns found for practice"
		} else {
			m.status = fmt.Sprintf("Loaded %d nouns for practice", len(m.practiceItems))
		}
		return m, nil
	case trainerItemsMsg:
		st := m.trainerStateFor(msg.kind)
		st.items = msg.items
		if len(msg.items) == 0 {
			m.status = "No " + st.config.ItemNoun + " found"
		} else {
			m.status = fmt.Sprintf("Loaded %d %s", len(msg.items), st.config.ItemNoun)
		}
		return m, nil
	case timedClearStatusMsg:
		if msg.seq == m.statusSeq {
			m.status = "Ready"
			m.isErrorStatus = false
		}
		return m, nil

	case reviewRecordedMsg:
		m.gradingInProgress = false
		m.lastReviewedCardID = msg.cardID
		m.lastReviewedGrade = msg.grade
		m.syncDecks(msg.decks)
		m.stats = msg.stats
		var dueReload tea.Cmd
		if msg.bookmarkFilter != m.bookmarkFilter {
			// Filter changed while grading; keep session stats but refresh
			// the queue with the filter the user currently has selected.
			dueReload = m.reloadDueForCurrentFilter()
		} else {
			m.allDue = msg.cards
			m.applyDeckFilter()
		}
		m.showHint = false

		// Only update session stats if a valid grade was provided
		if msg.grade != "" {
			m.sessionReviewed++
			if msg.grade != core.GradeAgain {
				m.sessionCorrect++
			}
			if m.sessionGrades == nil {
				m.sessionGrades = make(map[core.ReviewGrade]int)
			}
			m.sessionGrades[msg.grade]++

			gradeIcon := map[core.ReviewGrade]string{
				core.GradeAgain: "✗",
				core.GradeHard:  "~",
				core.GradeGood:  "✓",
				core.GradeEasy:  "★",
			}
			gradeText := map[core.ReviewGrade]string{
				core.GradeAgain: "Again",
				core.GradeHard:  "Hard",
				core.GradeGood:  "Good",
				core.GradeEasy:  "Easy",
			}
			icon := gradeIcon[msg.grade]
			gradeName := gradeText[msg.grade]
			accuracy := 0
			if m.sessionReviewed > 0 {
				accuracy = int(float64(m.sessionCorrect) / float64(m.sessionReviewed) * 100)
			}
			remaining := len(m.dueCards)
			m.status = fmt.Sprintf("%s %s | %d cards due | %d%% accuracy", icon, gradeName, remaining, accuracy)
			m.logger.Info("Recorded review for card %s with grade %s", msg.cardID, msg.grade)
		} else {
			m.status = "Ready"
		}

		// Reset typing state for next card
		m.typedAnswer = ""
		m.typingChecked = false
		m.typingCorrect = false

		followUps := []tea.Cmd{m.loadReviewsPerDay(), m.loadRecentDecks(), m.loadStatistics()}
		if dueReload != nil {
			followUps = append([]tea.Cmd{dueReload}, followUps...)
		}
		if len(m.dueCards) == 0 && m.sessionReviewed > 0 && dueReload == nil {
			return m, tea.Batch(append([]tea.Cmd{m.updateView(ViewSessionSummary)}, followUps...)...)
		}
		return m, tea.Batch(followUps...)
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
		var dueReload tea.Cmd
		if msg.bookmarkFilter != m.bookmarkFilter {
			dueReload = m.reloadDueForCurrentFilter()
		} else {
			m.allDue = msg.cards
			m.applyDeckFilter()
		}
		m.status = "Card suspended"
		m.logger.Info("Suspended card %s", msg.cardID)
		if dueReload != nil {
			return m, tea.Batch(dueReload, m.loadReviewsPerDay())
		}
		return m, m.loadReviewsPerDay()
	case dailyGoalSetMsg:
		newStats := core.Statistics(msg)
		// Preserve optimistic goal to avoid race conditions when spamming keys
		goal := m.stats.DailyGoal
		m.stats = newStats
		m.stats.DailyGoal = goal
		m.status = fmt.Sprintf("Daily goal set to %d", m.stats.DailyGoal)
		m.logger.Info("Daily goal set to %d (optimistic preserved)", m.stats.DailyGoal)
		return m, nil
	case reviewUndoneMsg:
		m.lastReviewedCardID = ""
		if m.sessionReviewed > 0 {
			m.sessionReviewed--
			if msg.grade != core.GradeAgain && m.sessionCorrect > 0 {
				m.sessionCorrect--
			}
		}
		if msg.grade != "" && m.sessionGrades != nil && m.sessionGrades[msg.grade] > 0 {
			m.sessionGrades[msg.grade]--
		}
		m.syncDecks(msg.decks)
		m.stats = msg.stats
		var dueReload tea.Cmd
		if msg.bookmarkFilter != m.bookmarkFilter {
			dueReload = m.reloadDueForCurrentFilter()
		} else {
			m.allDue = msg.cards
			m.applyDeckFilter()
		}
		m.status = "Last review undone"
		m.logger.Info("Undid last review for card %s", msg.cardID)
		if dueReload != nil {
			return m, tea.Batch(dueReload, m.loadReviewsPerDay())
		}
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
	case browserCardsResultMsg:
		if msg.id != m.browserLoadID {
			return m, nil
		}
		m.browserCards = msg.cards
		if m.browserCursor >= len(m.browserCards) {
			m.browserCursor = maxInt(0, len(m.browserCards)-1)
		}
		if len(m.browserCards) == 0 {
			m.status = "No cards found"
		} else {
			m.status = fmt.Sprintf("%d cards found", len(m.browserCards))
		}
		m.logger.Debug("Loaded %d browser cards for request %d", len(msg.cards), msg.id)
	case cramCardsMsg:
		// Drop loads that finished after the user changed filter or deck.
		if msg.id != 0 && msg.id != m.cramLoadID {
			return m, nil
		}
		if msg.cramType != "" && msg.cramType != m.cramType {
			return m, nil
		}
		if msg.deckID != m.deck.ID {
			return m, nil
		}
		// Flag modes are filtered in SQL; keep a defensive in-memory filter
		// using the request's snapshotted type so a mid-flight filter change
		// cannot mis-attribute results.
		filter := msg.cramType
		if filter == "" {
			filter = m.cramType
		}
		m.cramCards = m.cramCards[:0]
		for _, card := range msg.cards {
			switch filter {
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
			default:
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
		m.logger.Debug("Loaded %d cram cards with filter %s", len(m.cramCards), filter)
	case tea.KeyPressMsg:
		m.logger.Debug("Key pressed: %s", msg.String())
		return m.updateKey(msg)
	case tea.PasteMsg:
		m.logger.Debug("Pasted text: %s", msg.String())
		return m, m.handlePaste(msg.String())
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
						switch {
						case hit.View == ViewStatistics:
							m.dragVisible = m.statisticsVisibleLines(layout.Height)
							m.dragTotal = m.statsTotalLines
						case hit.View == ViewBrowser:
							m.dragVisible = m.listVisibleLines(layout.Height)
							m.dragTotal = len(m.browserCards)
						case hit.View == ViewCram:
							m.dragVisible = m.listVisibleLines(layout.Height)
							m.dragTotal = len(m.cramCards)
						case hit.View == ViewSettings:
							m.dragVisible = layout.Height
							m.dragTotal = m.settingsTotalLines
						case strings.HasPrefix(hit.ID, "dict-detail-scroll-"):
							m.dragView = ViewDictionary
							m.dragVisible = m.dictionaryDetailViewportRows(layout)
							m.dragTotal = m.dictionaryDetailTotalLines
							m.isDragging = m.dragTotal > m.dragVisible
						case strings.HasPrefix(hit.ID, "dict-scroll-"):
							m.dragView = ViewDictionary
							m.dragVisible = m.dictionaryListViewportRows(layout)
							m.dragTotal = len(m.dictionaryResults)
							m.isDragging = m.dragTotal > m.dragVisible
						}
					}
					cmd := m.activateHitbox(hit)
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
			delta := 0
			if mouse.Button == tea.MouseWheelUp {
				delta = -1
			} else if mouse.Button == tea.MouseWheelDown {
				delta = 1
			}
			if delta != 0 {
				if m.dictionaryOverlayActive || m.activeView == ViewDictionary {
					m.scrollDictionaryWheel(delta)
				} else {
					switch m.activeView {
					case ViewDecks:
						m.moveDeckCursor(delta)
					case ViewStatistics:
						m.scrollStats(delta)
					case ViewBrowser:
						m.moveBrowserCursor(delta)
					case ViewCram:
						m.moveCramCursor(delta)
					case ViewSettings:
						if delta < 0 {
							m.settingsScroll = maxInt(0, m.settingsScroll-1)
						} else {
							maxScroll := maxInt(0, m.settingsTotalLines-1)
							if m.settingsScroll < maxScroll {
								m.settingsScroll++
							}
						}
					case ViewPractice:
						if m.practiceSubView == PracticeSubViewHub {
							if delta < 0 {
								m.practiceScroll = maxInt(0, m.practiceScroll-1)
							} else {
								m.practiceScroll++
							}
						}
					case ViewAnkiWeb:
						m.scrollAnkiWebWheel(delta)
					}
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
	case BreakpointMedium:
		b.WriteString(m.renderMedium())
	default:
		b.WriteString(m.renderCompact())
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
		helpBox := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("81")).
			Padding(1, 2).
			Background(lipgloss.Color("233")).
			Render(m.renderHelp(viewportLayout{
				X:      m.width / 4,
				Y:      m.height / 4,
				Width:  m.width / 2,
				Height: m.height / 2,
			}))
		helpView := lipgloss.Place(m.width, m.height-3, lipgloss.Center, lipgloss.Center, helpBox)

		statusText := singleLine(m.status)
		if m.isErrorStatus {
			statusText = lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true).Render(statusText)
		}
		statusLine := fmt.Sprintf("status: %s", statusText)
		footer := strings.Join([]string{
			keyInfoStyle.Render("tab/arrows") + " views",
			keyInfoStyle.Render("0-9") + " views",
			keyInfoStyle.Render("?") + " help",
			keyInfoStyle.Render("q") + " quit",
		}, " │ ")

		return tea.View{
			Content:   helpView + "\n\n" + statusStyle.Render(truncateLine(statusLine, m.width-2)) + "\n" + statusStyle.Render(truncateLine(footer, m.width-2)),
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

func (m *Model) renderMedium() string {
	if m.focusMode && m.activeView == ViewReview {
		return "\n" + m.renderActiveView(0, 2)
	}
	return m.renderTabs(0, 1) + "\n" + m.renderActiveView(0, 2)
}

func (m *Model) renderCompact() string {
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

func (m *Model) openDictionary(word string) tea.Cmd {
	if word == "" {
		return nil
	}

	if strings.EqualFold(m.dictionaryProvider, "Local TUI") {
		m.dictionaryPreviousView = m.activeView
		m.activeView = ViewDictionary
		m.dictionarySearch = word
		m.dictionaryCursor = 0
		return m.searchDictionary()
	}

	return func() tea.Msg {
		var urlStr string
		switch strings.ToLower(m.dictionaryProvider) {
		case "linguee":
			urlStr = "https://www.linguee.com/german-english/search?source=auto&query=" + url.QueryEscape(word)
		case "leo":
			urlStr = "https://dict.leo.org/german-english/" + url.QueryEscape(word)
		case "duden":
			urlStr = "https://www.duden.de/suchen/dudenonline/" + url.QueryEscape(word)
		case "pons":
			urlStr = "https://en.pons.com/translate/german-english/" + url.QueryEscape(word)
		case "cambridge":
			urlStr = "https://dictionary.cambridge.org/dictionary/german-english/" + url.QueryEscape(word)
		case "google translate":
			urlStr = "https://translate.google.com/?sl=de&tl=en&text=" + url.QueryEscape(word) + "&op=translate"
		default: // dict.cc
			urlStr = "https://www.dict.cc/?s=" + url.QueryEscape(word)
		}

		var cmd *exec.Cmd
		switch runtime.GOOS {
		case "darwin":
			cmd = exec.Command("open", urlStr)
		case "windows":
			cmd = exec.Command("cmd", "/c", "start", urlStr)
		default:
			cmd = exec.Command("xdg-open", urlStr)
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
