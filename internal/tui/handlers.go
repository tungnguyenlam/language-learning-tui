package tui

import (
	"context"
	"fmt"
	"strings"
	"time"

	"deutsch-tui/internal/ai"
	"deutsch-tui/internal/core"

	tea "charm.land/bubbletea/v2"
)

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
	m.isDragging = false
	m.clearReviewHistory()
	m.hitboxes = nil // Clear hitboxes for new view
	m.importCursor = 0
	m.settingsCursor = 0
	m.editingImportPath = false
	m.editingExportTag = false
	m.editingTemplate = false

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
	if view == ViewCram {
		m.cramType = "bookmarked"
		m.cramCursor = 0
		m.cramReviewed = 0
		m.cramCorrect = 0
		return m.loadCramCards()
	}
	return nil
}

func (m *Model) toggleBrowserBookmark() tea.Cmd {
	if len(m.browserCards) == 0 {
		return nil
	}
	card := m.browserCards[m.browserCursor]
	next := !card.Bookmarked
	m.status = "Saving bookmark..."
	// Update local state immediately for responsiveness
	m.browserCards[m.browserCursor].Bookmarked = next
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
	selectedIDs := m.getSelectedCardIDs()
	if len(selectedIDs) == 0 {
		return nil
	}
	m.status = "Bulk bookmarking..."
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		for _, id := range selectedIDs {
			_ = m.repo.SetCardBookmark(ctx, id, bookmarked)
		}
		return m.loadBrowserCards()()
	}
}

func (m *Model) bulkBrowserSuspend(suspended bool) tea.Cmd {
	selectedIDs := m.getSelectedCardIDs()
	if len(selectedIDs) == 0 {
		return nil
	}
	m.status = "Bulk suspending..."
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		for _, id := range selectedIDs {
			_ = m.repo.SetCardSuspended(ctx, id, suspended)
		}
		return m.loadBrowserCards()()
	}
}

func (m *Model) handleTagInput() tea.Cmd {
	m.taggingCards = false
	tags := strings.Fields(m.tagInput)
	selectedIDs := m.getSelectedCardIDs()

	if len(selectedIDs) == 0 {
		if len(m.browserCards) > 0 {
			selectedIDs = []string{m.browserCards[m.browserCursor].ID}
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
		decks, _ := m.repo.Decks(ctx)
		return decksMsg(decks)
	}
}

func (m *Model) bulkBrowserDelete() tea.Cmd {
	selectedIDs := m.getSelectedCardIDs()
	if len(selectedIDs) == 0 {
		return nil
	}
	m.status = "Bulk deleting..."
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		for _, id := range selectedIDs {
			_ = m.repo.DeleteCard(ctx, id)
		}
		m.browserSelected = make(map[string]bool)
		return m.loadBrowserCards()()
	}
}

func (m *Model) getSelectedCardIDs() []string {
	var ids []string
	for id, selected := range m.browserSelected {
		if selected {
			ids = append(ids, id)
		}
	}
	return ids
}

func (m *Model) toggleCardKind() tea.Cmd {
	if len(m.browserCards) == 0 {
		return nil
	}
	card := m.browserCards[m.browserCursor]
	next := core.CardKindFlashcard
	if card.Kind == core.CardKindFlashcard {
		next = core.CardKindMCQ
	}
	m.status = fmt.Sprintf("Converting to %s...", next)
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		if err := m.repo.SetCardKind(ctx, card.ID, next); err != nil {
			return err
		}
		return m.loadBrowserCards()()
	}
}

func (m *Model) bulkBrowserToggleKind() tea.Cmd {
	selectedIDs := m.getSelectedCardIDs()
	if len(selectedIDs) == 0 {
		return nil
	}
	m.status = "Bulk converting kinds..."
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		for _, id := range selectedIDs {
			// For bulk, we'll just toggle to Flashcard if it was MCQ, or MCQ if it was Flashcard
			// but to keep it safe and predictable, let's just convert all to the opposite of the first selected
			var card core.Card
			cards, _ := m.repo.Cards(ctx, "", "")
			for _, c := range cards {
				if c.ID == id {
					card = c
					break
				}
			}
			next := core.CardKindFlashcard
			if card.Kind == core.CardKindFlashcard {
				next = core.CardKindMCQ
			}
			_ = m.repo.SetCardKind(ctx, id, next)
		}
		return m.loadBrowserCards()()
	}
}

func (m *Model) toggleBrowserSuspension() tea.Cmd {
	if len(m.browserCards) == 0 {
		return nil
	}
	card := m.browserCards[m.browserCursor]
	next := !card.Suspended
	m.status = "Updating suspension..."
	// Update local state
	m.browserCards[m.browserCursor].Suspended = next
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		if err := m.repo.SetCardSuspended(ctx, card.ID, next); err != nil {
			return err
		}
		// Refresh everything to ensure stats and lists are consistent
		var cards []core.Card
		var err error
		if m.bookmarkFilter {
			cards, err = m.repo.DueCardsBookmarked(ctx, time.Now(), 50)
		} else {
			cards, err = m.repo.DueCards(ctx, time.Now(), 500)
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
	for id, selected := range m.deckSelected {
		if selected {
			ids = append(ids, id)
		}
	}
	if len(ids) == 0 {
		filtered := m.filteredDecks()
		if len(filtered) > 0 {
			ids = []string{filtered[m.deckCursor].ID}
		}
	}
	if len(ids) == 0 {
		return nil
	}

	m.status = fmt.Sprintf("Deleting %d decks...", len(ids))
	m.deckSelected = make(map[string]bool)
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := m.repo.DeleteDecks(ctx, ids); err != nil {
			return err
		}
		return m.loadDecks()
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
	targetID := filtered[m.deckCursor].ID

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
			m.applyDeckFilter()
			return
		}
	}
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
	m.resetMCQState()
	m.clearReviewHistory()
	m.revealState = RevealIdle
	m.revealProgress = 0
	if m.activeView == ViewReview || m.activeView == ViewDashboard {
		if len(m.dueCards) == 0 {
			m.status = "All caught up!"
		} else {
			m.status = fmt.Sprintf("%d cards due", len(m.dueCards))
		}
	}
}

func (m *Model) resetMCQState() {
	m.mcqChoice = -1
	m.mcqAnswered = false
	m.mcqCorrect = false
	m.reviewPredictions = nil
}

func (m *Model) clearReviewHistory() {
	m.reviewHistory = nil
	m.showReviewHistory = false
	m.reviewHistoryCard = ""
	m.reviewPredictions = nil
}

func (m *Model) deckLabel() string {
	if strings.TrimSpace(m.deck.Name) != "" {
		return m.deck.Name
	}
	return "No deck"
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
	}
}

func (m *Model) statisticsVisibleLines(viewportHeight int) int {
	return clampInt(viewportHeight-2, 5, 40)
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
	setName := m.aiTemplateSets[m.aiTemplateIndex]
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
