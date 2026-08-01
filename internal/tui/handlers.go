package tui

import (
	"context"
	"fmt"
	"math/rand/v2"
	"strings"
	"time"
	"unicode"

	"deutsch-tui/internal/ai"
	"deutsch-tui/internal/core"

	tea "charm.land/bubbletea/v2"
)

func (m *Model) nextViewCmd() tea.Cmd {
	views := []View{ViewDashboard, ViewDecks, ViewReview, ViewStatistics, ViewImport, ViewAI, ViewSettings, ViewBrowser, ViewCram, ViewPractice}
	for i, view := range views {
		if m.activeView == view {
			return m.updateView(views[(i+1)%len(views)])
		}
	}
	return m.updateView(ViewDashboard)
}

func (m *Model) previousViewCmd() tea.Cmd {
	views := []View{ViewDashboard, ViewDecks, ViewReview, ViewStatistics, ViewImport, ViewAI, ViewSettings, ViewBrowser, ViewCram, ViewPractice}
	for i, view := range views {
		if m.activeView == view {
			return m.updateView(views[(i-1+len(views))%len(views)])
		}
	}
	return m.updateView(ViewDashboard)
}

func (m *Model) updateView(view View) tea.Cmd {
	var cmd tea.Cmd
	if m.activeView == ViewDictionary && m.dictionarySearch != "" {
		m.recordDictionarySearch(m.dictionarySearch)
	}
	if view == ViewDictionary && m.activeView != ViewDictionary {
		m.dictionaryPreviousView = m.activeView
		m.dictionaryDetailView = false
		m.dictionaryFocusResults = false
		if len(m.dictionaryDiscoverEntries) == 0 && m.dictionarySearch == "" {
			cmd = tea.Batch(cmd, m.loadDictionaryDiscoverEntries())
		}
	}

	// Start view transition animation
	if m.activeView != view && m.activeView != "" {
		m.prevView = m.activeView
		m.viewTransitioning = true
		m.viewTransitionProgress = 0
		m.viewTransitionFrame = 0
	}
	m.activeView = view
	m.dictionaryOverlayActive = false // Close overlay on any view switch
	m.isDragging = false
	m.confirmingDelete = false
	m.clearReviewHistory()
	m.hitboxes = nil // Clear hitboxes for new view
	m.importScreen.resetCursor()
	m.settingsCursor = 0
	m.settingsScroll = 0
	m.editingTemplate = false
	m.typingMode = false
	m.typedAnswer = ""
	m.typingChecked = false
	m.typingCorrect = false
	m.showHint = false
	m.showCardInfo = false
	m.showGrammarHint = false
	m.focusMode = false
	m.fixProposal = nil
	m.fixOldNote = nil
	m.fixingCard = false
	m.fixCardID = ""
	m.fixError = ""

	m.refreshViewStatus()

	if view == ViewImport {
		m.exportDeckID = m.deck.ID
	}
	if view == ViewStatistics {
		return m.loadStatistics()
	}
	if view == ViewBrowser {
		m.browserSearch = ""
		m.browserDeckID = m.deck.ID
		m.browserCursor = 0
		return m.loadBrowserCards()
	}
	if view == ViewReview && m.sessionReviewed == 0 {
		m.sessionStartTime = time.Now()
	}
	if view == ViewCram {
		m.cramType = "bookmarked"
		m.cramCursor = 0
		m.cramReviewed = 0
		m.cramCorrect = 0
		m.cramActive = false
		m.cramRevealed = false
		m.revealState = RevealIdle
		m.revealProgress = 0
		return m.loadCramCards()
	}
	if view == ViewPractice {
		m.practiceSubView = PracticeSubViewHub
		m.practiceHubCursor = 0
		m.practiceScroll = 0
		m.status = "Practice Hub: Choose a trainer"
		return nil // Loaders are called when entering specific modes
	}

	if view == ViewDashboard || view == ViewReview {
		m.applyDeckFilter()
		if view == ViewDashboard {
			return tea.Batch(m.loadStatistics(), m.loadRecentDecks(), m.loadReviewsPerDay())
		}
	}

	if view == ViewDictionary && cmd != nil {
		return cmd
	}

	return nil
}

func (m *Model) refreshViewStatus() {
	m.isErrorStatus = false
	switch m.activeView {
	case ViewReview:
		if len(m.dueCards) > 0 {
			m.status = fmt.Sprintf("%d cards due in %s", len(m.dueCards), m.deckLabel())
		} else {
			m.status = "No cards due"
		}
	case ViewBrowser:
		if len(m.browserCards) > 0 {
			m.status = fmt.Sprintf("%d cards found", len(m.browserCards))
		} else {
			m.status = "No cards found"
		}
	case ViewStatistics:
		if m.stats.TotalCards > 0 {
			m.status = fmt.Sprintf("%d cards in collection", m.stats.TotalCards)
		} else {
			m.status = "Ready"
		}
	case ViewCram:
		if len(m.cramCards) > 0 {
			m.status = fmt.Sprintf("%d cards in cram mode", len(m.cramCards))
		} else {
			m.status = "No cards found for this filter"
		}
	case ViewPractice:
		if m.practiceSubView == PracticeSubViewHub {
			m.status = "Practice Hub: Choose a trainer"
		}
	case ViewDictionary:
		if m.dictionarySearch != "" && len(m.dictionaryResults) > 0 {
			m.status = fmt.Sprintf("Found %d dictionary results", len(m.dictionaryResults))
		} else if m.dictionarySearch != "" {
			m.status = fmt.Sprintf("No dictionary results for %q", m.dictionarySearch)
		} else {
			m.status = "Ready"
		}
	default:
		m.status = "Ready"
	}
}

func (m *Model) enterPracticeMode(mode PracticeSubView) tea.Cmd {
	m.practiceSubView = mode
	if mode == PracticeSubViewGender {
		m.practiceIndex = 0
		m.practiceCorrect = 0
		m.practiceTotal = 0
		m.practiceRevealed = false
		m.status = "Loading nouns for practice..."
		return m.loadPracticeItems()
	}
	if isGenericTrainer(mode) {
		st := m.trainerStateFor(mode)
		st.index = 0
		st.round = 0
		st.correct = 0
		st.total = 0
		st.revealed = false
		st.input = ""
		st.showHint = false
		m.status = "Loading " + st.config.ItemNoun + "..."
		return st.config.Load(m)
	}
	return nil
}

// advanceGenderItem moves the MCQ gender trainer to the next noun, reshuffling
// after a full pass so a second lap is not the same sequence again.
func (m *Model) advanceGenderItem() {
	m.practiceRevealed = false
	if len(m.practiceItems) == 0 {
		m.practiceIndex = 0
		return
	}
	m.practiceIndex++
	if m.practiceIndex >= len(m.practiceItems) {
		m.practiceIndex = 0
		rand.Shuffle(len(m.practiceItems), func(i, j int) {
			m.practiceItems[i], m.practiceItems[j] = m.practiceItems[j], m.practiceItems[i]
		})
	}
}

func (m *Model) toggleBrowserBookmark() tea.Cmd {
	if len(m.browserCards) == 0 {
		return nil
	}
	card := m.browserCards[clampInt(m.browserCursor, 0, len(m.browserCards)-1)]
	next := !card.Bookmarked
	m.status = "Saving bookmark..."
	// Update local state immediately for responsiveness
	m.browserCards[clampInt(m.browserCursor, 0, len(m.browserCards)-1)].Bookmarked = next
	// Also update in allDue if present
	for i := range m.allDue {
		if m.allDue[i].ID == card.ID {
			m.allDue[i].Bookmarked = next
			break
		}
	}
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		if err := m.repo.SetCardBookmark(ctx, card.ID, next); err != nil {
			return err
		}
		return bookmarkToggledMsg{cardID: card.ID, bookmarked: next}
	}
}

func (m *Model) bulkBrowserBookmark(bookmarked bool) tea.Cmd {
	return m.executeBulkAction("Bulk bookmarking...", func(ctx context.Context, id string) error {
		return m.repo.SetCardBookmark(ctx, id, bookmarked)
	})
}

func (m *Model) bulkBrowserSuspend(suspended bool) tea.Cmd {
	return m.executeBulkAction("Bulk suspending...", func(ctx context.Context, id string) error {
		return m.repo.SetCardSuspended(ctx, id, suspended)
	})
}

func (m *Model) handleTagInput() tea.Cmd {
	m.taggingCards = false
	// Split by comma or whitespace
	tags := strings.FieldsFunc(m.tagInput, func(r rune) bool {
		return r == ',' || unicode.IsSpace(r)
	})
	selectedIDs := m.getSelectedCardIDs()

	if len(selectedIDs) == 0 {
		if len(m.browserCards) > 0 {
			selectedIDs = []string{m.browserCards[clampInt(m.browserCursor, 0, len(m.browserCards)-1)].ID}
		}
	}

	if len(selectedIDs) == 0 {
		return nil
	}

	m.status = "Updating tags..."
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := m.repo.SetCardsTags(ctx, selectedIDs, tags); err != nil {
			return err
		}
		return tagsUpdatedMsg{cardIDs: selectedIDs, tags: tags}
	}
}

func (m *Model) cleanupBrowserTags() tea.Cmd {
	m.statusSeq++ // Prevent previous timed clear from clearing this
	m.status = "Cleaning up unused tags..."
	deckID := m.browserDeckID
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := m.repo.CleanupTags(ctx, deckID); err != nil {
			return err
		}
		// Reload decks to see updated tags in UI if we ever switch back to Decks view
		decks, err := m.repo.Decks(ctx)
		if err != nil {
			return err
		}
		return decksMsg(decks)
	}
}

func (m *Model) bulkBrowserDelete() tea.Cmd {
	selectedIDs := m.getSelectedCardIDs()
	if len(selectedIDs) == 0 {
		return nil
	}
	m.confirmingDelete = true
	m.deleteAction = m.executeBulkBrowserDelete
	m.status = fmt.Sprintf("Delete %d selected cards? (y/n)", len(selectedIDs))
	return nil
}

func (m *Model) executeBulkBrowserDelete() tea.Cmd {
	m.browserSelected = make(map[string]bool)
	return m.executeBulkAction("Bulk deleting...", func(ctx context.Context, id string) error {
		return m.repo.DeleteCard(ctx, id)
	})
}

func (m *Model) getSelectedCardIDs() []string {
	visible := make(map[string]bool, len(m.browserCards))
	for _, card := range m.browserCards {
		visible[card.ID] = true
	}
	var ids []string
	for id, selected := range m.browserSelected {
		if selected && visible[id] {
			ids = append(ids, id)
		}
	}
	return ids
}

func (m *Model) executeBulkAction(status string, action func(ctx context.Context, id string) error) tea.Cmd {
	selectedIDs := m.getSelectedCardIDs()
	if len(selectedIDs) == 0 {
		return nil
	}
	m.status = status
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		for _, id := range selectedIDs {
			if err := action(ctx, id); err != nil {
				return err
			}
		}
		return m.loadBrowserCards()()
	}
}

func (m *Model) executeSingleAction(status string, id string, action func(ctx context.Context, id string) error) tea.Cmd {
	m.status = status
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := action(ctx, id); err != nil {
			return err
		}
		return m.loadBrowserCards()()
	}
}

func (m *Model) toggleCardKind() tea.Cmd {
	if len(m.browserCards) == 0 {
		return nil
	}
	card := m.browserCards[clampInt(m.browserCursor, 0, len(m.browserCards)-1)]
	next := core.CardKindFlashcard
	if card.Kind == core.CardKindFlashcard {
		next = core.CardKindMCQ
	}
	return m.executeSingleAction(fmt.Sprintf("Converting to %s...", next), card.ID, func(ctx context.Context, id string) error {
		return m.repo.SetCardKind(ctx, id, next)
	})
}

func (m *Model) bulkBrowserToggleKind() tea.Cmd {
	selectedIDs := m.getSelectedCardIDs()
	if len(selectedIDs) == 0 {
		return nil
	}
	// Snapshot toggles from the already-loaded browser rows. Reloading the
	// whole collection via Cards("",…,…) was capped and silently skipped
	// selected cards outside that window.
	nextByID := make(map[string]core.CardKind, len(selectedIDs))
	selected := make(map[string]bool, len(selectedIDs))
	for _, id := range selectedIDs {
		selected[id] = true
	}
	for _, card := range m.browserCards {
		if !selected[card.ID] {
			continue
		}
		next := core.CardKindFlashcard
		if card.Kind == core.CardKindFlashcard {
			next = core.CardKindMCQ
		}
		nextByID[card.ID] = next
	}
	if len(nextByID) == 0 {
		return nil
	}
	m.status = "Bulk converting kinds..."
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		for id, next := range nextByID {
			if err := m.repo.SetCardKind(ctx, id, next); err != nil {
				return err
			}
		}
		return m.loadBrowserCards()()
	}
}

func (m *Model) toggleBrowserSuspension() tea.Cmd {
	if len(m.browserCards) == 0 {
		return nil
	}
	card := m.browserCards[clampInt(m.browserCursor, 0, len(m.browserCards)-1)]
	next := !card.Suspended
	bookmarkFilter := m.bookmarkFilter
	m.status = "Updating suspension..."
	// Update local state
	m.browserCards[clampInt(m.browserCursor, 0, len(m.browserCards)-1)].Suspended = next
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		if err := m.repo.SetCardSuspended(ctx, card.ID, next); err != nil {
			return err
		}
		// Refresh with the same due-card limit as Review (0 → default).
		// A hard 500 here previously truncated the in-memory review queue.
		var cards []core.Card
		var err error
		if bookmarkFilter {
			cards, err = m.repo.DueCardsBookmarked(ctx, time.Now(), 0)
		} else {
			cards, err = m.repo.DueCards(ctx, time.Now(), 0)
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

func (m *Model) previousDeck() {
	if len(m.decks) == 0 {
		return
	}
	m.deckIndex = (m.deckIndex - 1 + len(m.decks)) % len(m.decks)
	m.selectDeck(m.deckIndex)
}

func (m *Model) handleDeckDelete() tea.Cmd {
	var ids []string
	var names []string

	// Get selected decks
	for _, d := range m.decks {
		if m.deckSelected[d.ID] && d.ID != "" {
			ids = append(ids, d.ID)
			names = append(names, d.Name)
		}
	}

	// If none selected, use cursor
	if len(ids) == 0 {
		filtered := m.filteredDecks()
		if len(filtered) > 0 {
			deck := filtered[clampInt(m.deckCursor, 0, len(filtered)-1)]
			if deck.ID != "" {
				ids = []string{deck.ID}
				names = []string{deck.Name}
			}
		}
	}

	if len(ids) == 0 {
		m.status = "Cannot delete 'All Decks'"
		return m.setStatus(m.status, 2*time.Second)
	}

	m.deleteIDs = ids
	m.confirmingDelete = true
	m.deleteAction = m.executeDeckDelete
	m.status = fmt.Sprintf("Delete %d decks and ALL their cards? (y/n)", len(ids))
	return nil
}

func (m *Model) executeDeckDelete() tea.Cmd {
	ids := m.deleteIDs
	if len(ids) == 0 {
		return nil
	}

	m.status = fmt.Sprintf("Deleting %d decks...", len(ids))
	m.deckSelected = make(map[string]bool)
	m.deleteIDs = nil // Clear for next time
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := m.repo.DeleteDecks(ctx, ids); err != nil {
			return err
		}
		return deckDeletedMsg{}
	}
}

func (m *Model) handleDeckMerge() tea.Cmd {
	var sourceIDs []string
	for id, selected := range m.deckSelected {
		if selected {
			sourceIDs = append(sourceIDs, id)
		}
	}
	if len(sourceIDs) == 0 {
		m.status = "Select source decks first (Space)"
		return nil
	}

	filtered := m.filteredDecks()
	if len(filtered) == 0 {
		return nil
	}
	targetID := filtered[clampInt(m.deckCursor, 0, len(filtered)-1)].ID

	// Remove target from sources if present
	var cleanSources []string
	for _, id := range sourceIDs {
		if id != targetID {
			cleanSources = append(cleanSources, id)
		}
	}
	if len(cleanSources) == 0 {
		m.status = "Cannot merge deck into itself"
		return nil
	}

	m.status = fmt.Sprintf("Merging %d decks into %s...", len(cleanSources), filtered[m.deckCursor].Name)
	m.deckSelected = make(map[string]bool)
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		if err := m.repo.MergeDecks(ctx, cleanSources, targetID); err != nil {
			return err
		}
		return m.loadDecks()
	}
}

func (m *Model) nextExportDeck() {
	if len(m.decks) == 0 {
		return
	}
	currentIndex := -1
	if m.exportDeckID == "" {
		currentIndex = 0 // "All Decks" is virtual at index 0 in my head
	} else {
		for i, d := range m.decks {
			if d.ID == m.exportDeckID {
				currentIndex = i + 1
				break
			}
		}
	}
	nextIndex := (currentIndex + 1) % (len(m.decks) + 1)
	if nextIndex == 0 {
		m.exportDeckID = ""
	} else {
		m.exportDeckID = m.decks[nextIndex-1].ID
	}
}

func (m *Model) previousExportDeck() {
	if len(m.decks) == 0 {
		return
	}
	currentIndex := -1
	if m.exportDeckID == "" {
		currentIndex = 0
	} else {
		for i, d := range m.decks {
			if d.ID == m.exportDeckID {
				currentIndex = i + 1
				break
			}
		}
	}
	prevIndex := (currentIndex - 1 + len(m.decks) + 1) % (len(m.decks) + 1)
	if prevIndex == 0 {
		m.exportDeckID = ""
	} else {
		m.exportDeckID = m.decks[prevIndex-1].ID
	}
}

func (m *Model) cycleExportFilter(forward bool) {
	options := []string{"All", "Mature", "Learning"}
	currentIndex := 0
	for i, opt := range options {
		if m.exportFilter == opt {
			currentIndex = i
			break
		}
	}
	if forward {
		currentIndex = (currentIndex + 1) % len(options)
	} else {
		currentIndex = (currentIndex - 1 + len(options)) % len(options)
	}
	m.exportFilter = options[currentIndex]
}

func (m *Model) nextDeck() {
	if len(m.decks) == 0 {
		return
	}
	m.deckIndex = (m.deckIndex + 1) % len(m.decks)
	m.selectDeck(m.deckIndex)
}

func (m *Model) selectDeckByID(id string) {
	for i, d := range m.decks {
		if d.ID == id {
			m.deckIndex = i
			m.deck = d
			m.deckCursor = i
			m.browserDeckID = id // Sync browser filter
			m.applyDeckFilter()
			return
		}
	}
}

func (m *Model) selectDeck(index int) {
	if len(m.decks) == 0 {
		m.deck = core.Deck{}
		m.deckCursor = 0
		m.browserDeckID = "" // Sync browser filter
		m.applyDeckFilter()
		return
	}
	if index < 0 || index >= len(m.decks) {
		index = 0
	}
	m.deckIndex = index
	m.deck = m.decks[index]
	m.deckCursor = index
	m.browserDeckID = m.deck.ID // Sync browser filter
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
	m.resetMCQState()
	m.clearReviewHistory()
	m.revealState = RevealIdle
	m.revealProgress = 0
	if m.activeView == ViewReview {
		if !m.showHelp {
			if len(m.dueCards) == 0 {
				m.status = fmt.Sprintf("All caught up in %s", m.deckLabel())
			} else {
				m.status = fmt.Sprintf("%d cards due in %s", len(m.dueCards), m.deckLabel())
			}
		} else {
			m.status = "Help overlay shown. Press ? to close."
		}
	} else if m.activeView == ViewDashboard && m.showHelp {
		m.status = "Help overlay shown. Press ? to close."
	}
}

func (m *Model) resetMCQState() {
	m.mcqChoice = -1
	m.mcqAnswered = false
	m.mcqCorrect = false
	m.reviewPredictions = nil
	m.explanation = ""
	m.explainingCard = false
	m.explainCardID = ""
	m.explainError = ""
}

func (m *Model) resetReviewState() {
	m.revealState = RevealIdle
	m.revealProgress = 0
	m.flipProgress = 0
	m.flipFrame = 0
	m.typingMode = false
	m.typedAnswer = ""
	m.typingChecked = false
	m.typingCorrect = false
	m.showHint = false
	m.showCardInfo = false
	m.showGrammarHint = false
	m.focusMode = false
	m.fixProposal = nil
	m.fixOldNote = nil
	m.fixingCard = false
	m.fixCardID = ""
	m.fixError = ""
	m.grammarHint = nil
	m.resetMCQState()
	m.clearReviewHistory()
}

func (m *Model) clearReviewHistory() {
	m.reviewHistory = nil
	m.showReviewHistory = false
	m.reviewHistoryCard = ""
	m.reviewPredictions = nil
	m.explanation = ""
	m.explainingCard = false
	m.explainCardID = ""
	m.explainError = ""
}

func (m *Model) deckLabel() string {
	if strings.TrimSpace(m.deck.Name) != "" {
		return m.deck.Name
	}
	return "No deck"
}

func (m *Model) moveDeckCursor(delta int) {
	filtered := m.filteredDecks()
	if len(filtered) == 0 {
		return
	}
	m.deckCursor = clampInt(m.deckCursor+delta, 0, len(filtered)-1)
}

func (m *Model) moveBrowserCursor(delta int) {
	if len(m.browserCards) == 0 {
		return
	}
	m.browserCursor = clampInt(m.browserCursor+delta, 0, len(m.browserCards)-1)
	m.clearReviewHistory()
}

func (m *Model) moveCramCursor(delta int) {
	if len(m.cramCards) == 0 {
		return
	}
	m.cramCursor = clampInt(m.cramCursor+delta, 0, len(m.cramCards)-1)
}

func (m *Model) scrollStats(delta int) {
	m.setStatsScroll(m.statsScroll + delta)
	if m.statsScroll < 0 {
		m.statsScroll = 0
	}
	maxScroll := m.statsMaxScroll()
	if m.statsScroll > maxScroll {
		m.statsScroll = maxScroll
	}
}

func (m *Model) statsMaxScroll() int {
	visible := m.statisticsVisibleLines(m.activeViewContentLayout().Height)
	return maxInt(0, m.statsTotalLines-visible)
}

func (m *Model) setStatsScroll(offset int) {
	visible := m.statisticsVisibleLines(m.activeViewContentLayout().Height)
	maxScroll := maxInt(0, m.statsTotalLines-visible)
	m.statsScroll = clampInt(offset, 0, maxScroll)
}

func (m *Model) handleMouseDrag(mouseY int) {
	if m.dragVisible <= 1 || m.dragTotal <= m.dragVisible {
		return
	}
	row := mouseY - m.dragTrackStartY
	switch m.dragView {
	case ViewStatistics:
		m.setStatsScroll(scrollOffsetForTrackRow(m.dragTotal, m.dragVisible, row))
	case ViewBrowser:
		next := selectedIndexForTrackRow(m.dragTotal, m.dragVisible, row)
		if next != m.browserCursor {
			m.browserCursor = next
			m.clearReviewHistory()
		}
	case ViewCram:
		m.cramCursor = selectedIndexForTrackRow(m.dragTotal, m.dragVisible, row)
	case ViewSettings:
		maxScroll := maxInt(0, m.settingsTotalLines-m.dragVisible)
		m.settingsScroll = clampInt(scrollOffsetForTrackRow(m.dragTotal, m.dragVisible, row), 0, maxScroll)
	case ViewDictionary:
		if m.dictionaryDetailView {
			m.dictionaryDetailScroll = scrollOffsetForTrackRow(m.dragTotal, m.dragVisible, row)
		} else {
			m.dictionaryScroll = scrollOffsetForTrackRow(m.dragTotal, m.dragVisible, row)
		}
	}
}

// scrollDictionaryWheel moves the dictionary result cursor or detail pane
// with the mouse wheel (full Dictionary tab and Spotlight overlay).
func (m *Model) scrollDictionaryWheel(delta int) {
	if m.dictionaryDetailView {
		if delta < 0 {
			m.dictionaryDetailScroll = maxInt(0, m.dictionaryDetailScroll-1)
			return
		}
		maxScroll := maxInt(0, m.dictionaryDetailTotalLines-m.dictionaryDetailViewportRows(m.activeViewContentLayout()))
		if m.dictionaryDetailScroll < maxScroll {
			m.dictionaryDetailScroll++
		}
		return
	}
	if len(m.dictionaryResults) == 0 {
		return
	}
	m.dictionaryFocusResults = true
	if delta < 0 {
		if m.dictionaryCursor > 0 {
			m.dictionaryCursor--
			m.dictionaryDetailScroll = 0
		}
		return
	}
	if m.dictionaryCursor < len(m.dictionaryResults)-1 {
		m.dictionaryCursor++
		m.dictionaryDetailScroll = 0
	}
}

func (m *Model) scrollAnkiWebWheel(delta int) {
	if m.ankiWebScreen == nil || len(m.ankiWebScreen.results) == 0 {
		return
	}
	s := m.ankiWebScreen
	if delta < 0 {
		if s.cursor > 0 {
			s.cursor--
			s.details = nil
		}
		return
	}
	if s.cursor < len(s.results)-1 {
		s.cursor++
		s.details = nil
	}
}

func (m *Model) statisticsVisibleLines(viewportHeight int) int {
	neededSpace := 2
	if m.statsTotalLines > viewportHeight-2 {
		neededSpace = 3
	}
	return clampInt(viewportHeight-neededSpace, 1, 1000)
}

func (m *Model) listVisibleLines(viewportHeight int) int {
	return clampInt(viewportHeight-8, 3, 10)
}

func (m *Model) nextAITemplate() {
	if len(m.aiTemplateSets) == 0 {
		return
	}
	m.aiTemplateIndex = (m.aiTemplateIndex + 1) % len(m.aiTemplateSets)
	m.updateActiveAITemplate()
}

func (m *Model) previousAITemplate() {
	if len(m.aiTemplateSets) == 0 {
		return
	}
	m.aiTemplateIndex = (m.aiTemplateIndex - 1 + len(m.aiTemplateSets)) % len(m.aiTemplateSets)
	m.updateActiveAITemplate()
}

func (m *Model) updateActiveAITemplate() {
	if len(m.aiTemplateSets) == 0 {
		return
	}
	setName := m.currentAITemplateSet()
	if setName == "" {
		return
	}
	if tp, ok := m.aiProvider.(ai.TemplateProvider); ok {
		tp.ActiveSet = setName
		m.aiProvider = tp
	}
	m.status = fmt.Sprintf("AI Template: %s", setName)
}

func (m *Model) setStatus(text string, d time.Duration) tea.Cmd {
	m.status = text
	m.statusSeq++
	seq := m.statusSeq
	return tea.Tick(d, func(t time.Time) tea.Msg {
		return timedClearStatusMsg{seq: seq}
	})
}

type timedClearStatusMsg struct {
	seq int
}
