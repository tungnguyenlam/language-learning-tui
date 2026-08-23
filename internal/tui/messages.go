package tui

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"deutsch-tui/internal/ai"
	"deutsch-tui/internal/content"
	"deutsch-tui/internal/core"

	tea "charm.land/bubbletea/v2"
)

// handleAsyncMsg applies domain async replies (load results, grading, drafts,
// animation ticks) to the model. It reports whether the message was consumed;
// unconsumed messages (input, resize) fall through to Update's own switch.
func (m *Model) handleAsyncMsg(msg tea.Msg) (tea.Cmd, bool) {
	switch msg := msg.(type) {
	case debounceSearchMsg:
		switch msg.view {
		case ViewDictionary:
			if msg.id == m.dictionarySearchTimerID {
				return m.searchDictionary(), true
			}
		case ViewBrowser:
			if msg.id == m.browserSearchTimerID {
				return m.loadBrowserCards(), true
			}
		}
		return nil, true

	case dictionarySearchResultsMsg:
		if msg.id != m.dictionarySearchID {
			return nil, true
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
			cmd = tea.Batch(m.findRelatedEntries(msg.results[0].Word), m.loadCompoundBreakdown(msg.results[0].Word))
		}
		return cmd, true
	case compoundBreakdownMsg:
		if msg.searchID != m.dictionarySearchID {
			return nil, true
		}
		if msg.err != nil {
			m.isErrorStatus = true
			m.status = friendlyError(msg.err)
			m.logger.Error("Error occurred: %v", msg.err)
			return nil, true
		}
		if m.compoundCache == nil {
			m.compoundCache = make(map[string][]content.CompoundPart)
		}
		m.compoundCache[msg.word] = msg.parts
		return nil, true
	case dictHistoryLoadedMsg:
		m.dictionarySearchHistory = msg
		return nil, true
	case dictRecentlyViewedLoadedMsg:
		m.dictionaryRecentlyViewed = msg
		return nil, true
	case dictRelatedEntriesMsg:
		if msg.id != m.dictionaryRelatedID || (m.dictionaryCursor >= 0 && m.dictionaryCursor < len(m.dictionaryResults) && m.dictionaryResults[m.dictionaryCursor].Word != msg.word) {
			return nil, true
		}
		m.dictionaryRelatedEntries = msg.entries
		return nil, true
	case dictDiscoverEntriesMsg:
		m.dictionaryDiscoverEntries = msg.entries
		return nil, true
	case dictStarredLoadedMsg:
		if msg != nil {
			m.dictionaryStarred = msg
		}
		return nil, true
	case deckHistoryLoadedMsg:
		m.deckSearchHistory = msg
		return nil, true
	case latestBackupPathMsg:
		return m.handleLatestBackupPathMsg(msg), true
	case spinnerTickMsg:
		if m.drafting {
			m.spinnerFrame++
			return m.tickSpinner(), true
		}
		return nil, true
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
				return m.tickReveal(), true
			}
		}
		return nil, true
	case flipTickMsg:
		if m.revealState == RevealFlipping {
			m.flipFrame++
			m.flipProgress += 10
			if m.flipProgress >= 100 {
				m.flipProgress = 100
				m.revealState = RevealRevealing
				m.revealProgress = 0
				m.logger.Debug("Flip complete, starting reveal")
				return m.tickReveal(), true
			}
			return m.tickFlip(), true
		}
		return nil, true
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
				return m.tickViewTransition(), true
			}
		}
		return nil, true
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
				return m.tickCardTransition(), true
			}
		}
		return nil, true
	case error:
		m.drafting = false
		m.draftCancelled = false
		m.gradingInProgress = false
		m.isErrorStatus = true
		m.status = friendlyError(msg)
		m.logger.Error("Error occurred: %v", msg)
		return nil, true
	case decksMsg:
		m.logger.Debug("Received %d decks", len(msg))
		m.syncDecks([]core.Deck(msg))
		return nil, true
	case dueCardsMsg:
		if msg.id != m.dueLoadID || m.bookmarkFilter {
			return nil, true
		}
		m.logger.Debug("Received %d due cards", len(msg.cards))
		m.allDue = msg.cards
		m.applyDeckFilter()
		m.showHint = false
		return nil, true
	case bookmarkedDueCardsMsg:
		if msg.id != m.dueLoadID || !m.bookmarkFilter {
			return nil, true
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
		return nil, true

	case statsMsg:
		// Statistics are deck-scoped and loaded asynchronously. Ignore a
		// response for an older request or a deck the user has already left.
		if msg.id != 0 && (msg.id != m.statsLoadID || msg.deckID != m.deck.ID) {
			return nil, true
		}
		m.stats = msg.stats
		m.dictCount = msg.dictCount
		m.logger.Debug("Received statistics update")
		return nil, true
	case reviewsPerDayMsg:
		m.reviewsPerDay = map[string]int(msg)
		m.logger.Debug("Received reviews per day data")
		return nil, true
	case recentDecksMsg:
		m.recentDecks = []string(msg)
		m.logger.Debug("Received %d recent decks", len(msg))
		return nil, true
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
		return nil, true
	case reviewPredictionsMsg:
		// Prediction requests can finish after navigation or grading. Only
		// install a response for the still-current request and card.
		if msg.id != m.reviewPredictionID || m.activeView != ViewReview || len(m.dueCards) == 0 ||
			m.dueCards[clampInt(m.cursor, 0, len(m.dueCards)-1)].ID != msg.cardID {
			return nil, true
		}
		m.reviewPredictions = msg.predictions
		m.logger.Debug("Received review predictions")
		return nil, true
	case draftsMsg:
		m.drafting = false
		if m.draftCancelled {
			m.draftCancelled = false
			return nil, true
		}
		m.drafts = []ai.Draft(msg)
		m.draftCursor = clampInt(m.draftCursor, 0, maxInt(0, len(m.drafts)-1))
		if len(m.drafts) == 0 {
			m.status = "No drafts generated"
		} else {
			m.status = fmt.Sprintf("%d drafts ready", len(m.drafts))
		}
		m.logger.Info("Generated %d AI drafts", len(msg))
		return nil, true
	case draftApprovedMsg:
		m.removeDraft(msg.noteID)
		m.allDue = msg.cards
		m.applyDeckFilter()
		m.status = "Draft saved"
		m.logger.Info("Approved AI draft %s", msg.noteID)
		return nil, true
	case importDoneMsg:
		if msg.path == "AI Drafts" {
			m.drafts = nil
			m.draftCursor = 0
		}
		m.syncDecks(msg.decks)
		m.allDue = msg.cards
		m.applyDeckFilter()
		m.logger.Info("Import completed: %d notes from %s", msg.count, filepath.Base(msg.path))
		return tea.Batch(
			m.setStatus(fmt.Sprintf("Imported %d notes from %s", msg.count, filepath.Base(msg.path)), 3*time.Second),
			m.loadStatistics(),
			m.loadReviewsPerDay(),
			m.loadRecentDecks(),
		), true
	case backupDoneMsg:
		if msg.restore {
			m.syncDecks(msg.decks)
			m.allDue = msg.cards
			m.applyDeckFilter()
			m.logger.Info("Restored %d rows from %s", msg.info.TotalRows, filepath.Base(msg.info.Path))
			return tea.Batch(
				m.setStatus(fmt.Sprintf("Restored %d rows from %s", msg.info.TotalRows, filepath.Base(msg.info.Path)), 3*time.Second),
				m.loadStatistics(),
				m.loadReviewsPerDay(),
				m.loadRecentDecks(),
			), true
		}
		m.lastBackupPath = msg.info.Path
		m.logger.Info("Backed up %d rows to %s", msg.info.TotalRows, filepath.Base(msg.info.Path))
		return m.setStatus(fmt.Sprintf("Backed up %d rows to %s", msg.info.TotalRows, filepath.Base(msg.info.Path)), 3*time.Second), true
	case statusMsg:
		m.logger.Debug("Setting status: %s", msg.text)
		m.isErrorStatus = false
		return m.setStatus(msg.text, 3*time.Second), true
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
		return m.setStatus(fmt.Sprintf("Updated tags for %d cards", len(msg.cardIDs)), 3*time.Second), true
	case fixProposalMsg:
		if !m.fixingCard || msg.cardID != m.fixCardID {
			return nil, true
		}
		m.fixingCard = false
		note := msg.oldNote
		m.fixOldNote = &note
		prop := msg.proposal
		m.fixProposal = &prop
		m.fixError = ""
		m.logger.Info("AI returned fix proposal for card %s", msg.cardID)
		m.status = "Review the proposed fix: y to apply, n to discard"
		return nil, true
	case fixErrorMsg:
		if !m.fixingCard || msg.cardID != m.fixCardID {
			return nil, true
		}
		m.fixingCard = false
		m.fixOldNote = nil
		m.fixProposal = nil
		m.fixError = msg.err.Error()
		m.logger.Error("Fix failed for card %s: %v", msg.cardID, msg.err)
		return m.setStatus("Fix failed: "+msg.err.Error(), 5*time.Second), true
	case fixAppliedMsg:
		m.allDue = msg.cards
		m.applyDeckFilter()
		m.logger.Info("Applied fix for card %s", msg.cardID)
		return m.setStatus("Card updated by AI", 3*time.Second), true
	case explainMsg:
		if !m.explainingCard || msg.cardID != m.explainCardID {
			return nil, true
		}
		m.explainingCard = false
		m.explanation = msg.explanation
		m.explainError = ""
		m.logger.Info("AI returned explanation for card %s", msg.cardID)
		m.status = "Pedagogical explanation received"
		return nil, true
	case explainErrorMsg:
		if !m.explainingCard || msg.cardID != m.explainCardID {
			return nil, true
		}
		m.explainingCard = false
		m.explanation = ""
		m.explainError = msg.err.Error()
		m.logger.Error("Explanation failed for card %s: %v", msg.cardID, msg.err)
		return m.setStatus("Explanation failed: "+msg.err.Error(), 5*time.Second), true
	case deckDeletedMsg:
		m.logger.Info("Deck deleted")
		return tea.Batch(m.loadDecks, m.loadDueCards()), true
	case practiceItemsMsg:
		// Ignore Gender loads that finished after leaving the trainer or
		// after a newer Gender load was started.
		if msg.id != m.practiceLoadID || m.practiceSubView != PracticeSubViewGender {
			return nil, true
		}
		m.practiceItems = msg.items
		if len(m.practiceItems) == 0 {
			m.status = "No nouns found for practice"
		} else {
			m.status = fmt.Sprintf("Loaded %d nouns for practice", len(m.practiceItems))
		}
		return nil, true
	case trainerItemsMsg:
		st := m.trainerStateFor(msg.kind)
		if m.practiceSubView != msg.kind || msg.id != st.loadID {
			return nil, true
		}
		st.items = msg.items
		if len(msg.items) == 0 {
			m.status = "No " + st.config.ItemNoun + " found"
		} else {
			m.status = fmt.Sprintf("Loaded %d %s", len(msg.items), st.config.ItemNoun)
		}
		return nil, true
	case timedClearStatusMsg:
		if msg.seq == m.statusSeq {
			m.status = "Ready"
			m.isErrorStatus = false
		}
		return nil, true

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
			return tea.Batch(append([]tea.Cmd{m.updateView(ViewSessionSummary)}, followUps...)...), true
		}
		return tea.Batch(followUps...), true
	case bookmarkToggledMsg:
		m.setCardBookmarkLocal(msg.cardID, msg.bookmarked)
		if msg.bookmarked {
			m.status = "Card bookmarked"
		} else {
			m.status = "Card bookmark removed"
		}
		m.logger.Debug("Toggled bookmark for card %s: %t", msg.cardID, msg.bookmarked)
		return nil, true
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
			return tea.Batch(dueReload, m.loadReviewsPerDay()), true
		}
		return m.loadReviewsPerDay(), true
	case dailyGoalSetMsg:
		newStats := core.Statistics(msg)
		// Preserve optimistic goal to avoid race conditions when spamming keys
		goal := m.stats.DailyGoal
		m.stats = newStats
		m.stats.DailyGoal = goal
		m.status = fmt.Sprintf("Daily goal set to %d", m.stats.DailyGoal)
		m.logger.Info("Daily goal set to %d (optimistic preserved)", m.stats.DailyGoal)
		return nil, true
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
			return tea.Batch(dueReload, m.loadReviewsPerDay()), true
		}
		return m.loadReviewsPerDay(), true
	case browserCardsResultMsg:
		if msg.id != m.browserLoadID {
			return nil, true
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
		return nil, true
	case cramCardsMsg:
		// Drop loads that finished after the user changed filter or deck.
		if msg.id != 0 && msg.id != m.cramLoadID {
			return nil, true
		}
		if msg.cramType != "" && msg.cramType != m.cramType {
			return nil, true
		}
		if msg.deckID != m.deck.ID {
			return nil, true
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
		return nil, true
	}
	return nil, false
}
