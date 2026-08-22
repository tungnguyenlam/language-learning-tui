package tui

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"deutsch-tui/internal/ai"
	"deutsch-tui/internal/app"
	"deutsch-tui/internal/content"
	"deutsch-tui/internal/core"

	tea "charm.land/bubbletea/v2"
)

func (m *Model) gradeCard(grade core.ReviewGrade) tea.Cmd {
	if m.activeView != ViewReview || len(m.dueCards) == 0 || m.revealState != RevealRevealed || m.gradingInProgress || len(grade) == 0 {
		return nil
	}

	m.gradingInProgress = true
	m.status = fmt.Sprintf("Grade: %s", strings.ToUpper(string(grade[:1]))+string(grade[1:]))
	card := m.dueCards[clampInt(m.cursor, 0, len(m.dueCards)-1)]
	bookmarkFilter := m.bookmarkFilter
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
		var err2 error
		if bookmarkFilter {
			cards, err2 = m.repo.DueCardsBookmarked(ctx, time.Now(), 0)
		} else {
			cards, err2 = m.repo.DueCards(ctx, time.Now(), 0)
		}
		if err2 != nil {
			return err2
		}
		decks, err2 := m.repo.Decks(ctx)
		if err2 != nil {
			return err2
		}
		stats, err2 := m.repo.Statistics(ctx)
		if err2 != nil {
			return err2
		}
		return reviewRecordedMsg{cardID: card.ID, cards: cards, decks: decks, stats: stats, grade: grade, bookmarkFilter: bookmarkFilter}
	}
}

func (m *Model) undoLastReview() tea.Cmd {
	if m.lastReviewedCardID == "" {
		m.status = "Nothing to undo"
		return nil
	}
	m.status = "Undoing last review..."
	cardID := m.lastReviewedCardID
	grade := m.lastReviewedGrade
	bookmarkFilter := m.bookmarkFilter
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		if err := m.repo.UndoLastReview(ctx, cardID); err != nil {
			return err
		}
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
		return reviewUndoneMsg{cardID: cardID, cards: cards, decks: decks, stats: stats, grade: grade, bookmarkFilter: bookmarkFilter}
	}
}

func (m *Model) setDailyGoal(goal int) tea.Cmd {
	if goal < 1 {
		goal = 1
	}
	m.stats.DailyGoal = goal
	m.status = "Saving daily goal..."
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
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

func (m *Model) setRevealSpeed(speed int) tea.Cmd {
	if speed < 0 {
		speed = 0
	}
	if speed > 10 {
		speed = 10
	}
	m.revealSpeed = speed
	m.status = fmt.Sprintf("Reveal speed set to %d", speed)
	if speed == 0 {
		m.status = "Reveal animation disabled (instant)"
	}
	m.persistConfig()
	return nil
}

func (m *Model) persistConfig() {
	if m.onConfigChange != nil {
		m.onConfigChange(m.aiProviderName, m.dictionaryProvider, m.aiTemplates, m.autoPlayAudio, m.strictNormalization, m.revealSpeed)
	}
}

func (m *Model) setDeckLimits(deckID string, newLimit, reviewLimit int) tea.Cmd {
	if deckID == "" {
		return nil
	}
	m.status = "Saving deck limits..."
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		if err := m.repo.SetDeckLimits(ctx, deckID, newLimit, reviewLimit); err != nil {
			return err
		}
		decks, err := m.repo.Decks(ctx)
		if err != nil {
			return err
		}
		return decksMsg(decks)
	}
}

func (m *Model) seedStandardContent() tea.Cmd {
	m.status = "Seeding standard content..."
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second) // Increased timeout for dictionary
		defer cancel()
		count := 0
		var dictEntries []core.DictionaryEntry

		for _, deck := range content.StandardDecks() {
			if err := m.repo.UpsertDeck(ctx, deck); err != nil {
				return err
			}
			count += len(deck.Notes)

			// Populate dictionary entries from notes
			for _, note := range deck.Notes {
				// We use the front as word and back as translation for basic seeding
				// Standard decks usually have clean word - translation pairs in notes
				entry := core.DictionaryEntry{
					ID:          note.ID,
					Word:        note.Front,
					Translation: note.Back,
					Tags:        note.Tags,
				}
				dictEntries = append(dictEntries, entry)
			}
		}

		if dictRepo, ok := m.repo.(core.DictionaryRepository); ok {
			if err := dictRepo.ImportEntries(ctx, dictEntries); err != nil {
				// Don't fail the whole seed if dictionary fails
				m.logger.Error("Failed to seed dictionary: %v", err)
			}
		}

		loadedDecks, err := m.repo.Decks(ctx)
		if err != nil {
			return err
		}
		// Use a small buffer to ensure we catch newly created cards
		cards, err := m.repo.DueCards(ctx, time.Now().Add(time.Minute), 0)
		if err != nil {
			return err
		}
		return importDoneMsg{decks: loadedDecks, cards: cards, count: count, path: "Standard Content"}
	}
}

func (m *Model) handleResetDatabase() tea.Cmd {
	m.confirmingDelete = true
	m.deleteAction = m.executeResetDatabase
	m.status = "WIPE ALL DATA? (y/n)"
	return nil
}

func (m *Model) executeResetDatabase() tea.Cmd {
	m.status = "Resetting database..."
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := m.repo.Reset(ctx); err != nil {
			return err
		}
		// Reset everything
		decks, _ := m.repo.Decks(ctx)
		cards, _ := m.repo.DueCards(ctx, time.Now(), 0)

		// Create a compound message or just trigger reloads
		return importDoneMsg{decks: decks, cards: cards, count: 0, path: "Database Reset"}
	}
}

func (m *Model) importTSV() tea.Cmd {
	path := strings.TrimSpace(m.importPath)
	if path == "" {
		m.status = "Import path is empty"
		return nil
	}
	m.status = "Importing TSV..."
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		file, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("failed to open file '%s': %w", filepath.Base(path), err)
		}
		defer file.Close()

		defaultDeck := m.deckLabel()
		if m.deck.ID == "" {
			defaultDeck = "Imported"
		}

		notes, err := content.ImportAnkiTSV(file, content.ImportOptions{
			DefaultDeck: defaultDeck,
		})
		if err != nil {
			return fmt.Errorf("failed to parse TSV file '%s': %w", filepath.Base(path), err)
		}
		decks := content.DecksFromNotes(
			notes)
		for _, deck := range decks {
			if err := m.repo.UpsertDeck(ctx, deck); err != nil {
				return fmt.Errorf("failed to save deck '%s': %w", deck.Name, err)
			}
		}
		loadedDecks, err := m.repo.Decks(ctx)
		if err != nil {
			return fmt.Errorf("failed to load decks: %w", err)
		}
		cards, err := m.repo.DueCards(ctx, time.Now(), 0)
		if err != nil {
			return fmt.Errorf("failed to load due cards: %w", err)
		}
		return importDoneMsg{decks: loadedDecks, cards: cards, count: len(notes), path: path}
	}
}

func (m *Model) exportTSV() tea.Cmd {
	path := strings.TrimSpace(m.exportPath)
	if path == "" {
		m.status = "Export path is empty"
		return nil
	}
	deckID := m.exportDeckID
	tag := m.exportTag
	filter := m.exportFilter
	m.status = "Exporting TSV..."
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		var notes []core.Note
		cards, err := m.repo.Cards(ctx, deckID, "", tag)
		if err != nil {
			return fmt.Errorf("failed to load cards for export: %w", err)
		}
		seen := make(map[string]bool)
		for _, c := range cards {
			if filter == "Mature" && !c.Mature {
				continue
			}
			if filter == "Learning" && c.Mature {
				continue
			}
			if seen[c.NoteID] {
				continue
			}
			seen[c.NoteID] = true
			notes = append(notes, core.Note{
				ID:     c.NoteID,
				DeckID: c.DeckID,
				Front:  c.Prompt,
				Back:   c.Answer,
				Tags:   c.Tags,
				Audio:  c.Audio,
			})
		}

		file, err := os.Create(path)
		if err != nil {
			return fmt.Errorf("failed to create export file '%s': %w", filepath.Base(path), err)
		}
		defer file.Close()

		if err := content.ExportAnkiTSV(file, notes); err != nil {
			return fmt.Errorf("failed to write TSV data to '%s': %w", filepath.Base(path), err)
		}
		return statusMsg{text: fmt.Sprintf("Exported %d notes to %s", len(notes), path)}
	}
}

func (m *Model) importAPKG() tea.Cmd {
	path := strings.TrimSpace(m.importPath)
	if path == "" {
		m.status = "Import path is empty"
		return nil
	}
	m.status = "Importing APKG..."
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		file, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("failed to open APKG file '%s': %w", filepath.Base(path), err)
		}
		defer file.Close()

		notes, err := content.ImportAnkiAPKG(file)
		if err != nil {
			return fmt.Errorf("failed to parse APKG file '%s': %w", filepath.Base(path), err)
		}
		decks := content.DecksFromNotes(
			notes)
		for _, deck := range decks {
			if err := m.repo.UpsertDeck(ctx, deck); err != nil {
				return fmt.Errorf("failed to save deck '%s': %w", deck.Name, err)
			}
		}
		loadedDecks, err := m.repo.Decks(ctx)
		if err != nil {
			return fmt.Errorf("failed to load decks: %w", err)
		}
		cards, err := m.repo.DueCards(ctx, time.Now(), 0)
		if err != nil {
			return fmt.Errorf("failed to load due cards: %w", err)
		}
		return importDoneMsg{decks: loadedDecks, cards: cards, count: len(notes), path: path}
	}
}

func (m *Model) exportAPKG() tea.Cmd {
	path := strings.TrimSpace(m.exportPath)
	if path == "" {
		m.status = "Export path is empty"
		return nil
	}
	deckID := m.exportDeckID
	tag := m.exportTag
	filter := m.exportFilter
	m.status = "Exporting APKG..."
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		cards, err := m.repo.Cards(ctx, deckID, "", tag)
		if err != nil {
			return fmt.Errorf("failed to load cards for export: %w", err)
		}

		// Anki shows the deck name, so map our deck IDs onto their titles
		// rather than exporting decks called "a1_food_drink".
		deckNames := make(map[string]string)
		if decks, err := m.repo.Decks(ctx); err == nil {
			for _, d := range decks {
				deckNames[d.ID] = d.Name
			}
		}

		index := make(map[string]int)
		notes := make([]core.Note, 0, len(cards))
		for _, c := range cards {
			if filter == "Mature" && !c.Mature {
				continue
			}
			if filter == "Learning" && c.Mature {
				continue
			}
			if at, ok := index[c.NoteID]; ok {
				// A second card whose sides are swapped is the reverse
				// direction of the same note; Anki models that as a note type,
				// not as an extra note.
				if notes[at].Front == c.Answer && notes[at].Back == c.Prompt {
					notes[at].Type = "Reverse"
				}
				continue
			}
			index[c.NoteID] = len(notes)
			notes = append(notes, core.Note{
				ID:     c.NoteID,
				DeckID: c.DeckID,
				Front:  c.Prompt,
				Back:   c.Answer,
				Extra:  c.Extra,
				Tags:   c.Tags,
				Audio:  c.Audio,
			})
		}

		file, err := os.Create(path)
		if err != nil {
			return fmt.Errorf("failed to create APKG file '%s': %w", filepath.Base(path), err)
		}
		defer file.Close()

		if err := content.ExportAnkiAPKGWithDeckNames(file, notes, deckNames); err != nil {
			return fmt.Errorf("failed to write APKG data to '%s': %w", filepath.Base(path), err)
		}
		return statusMsg{text: fmt.Sprintf("Exported %d notes to %s", len(notes), path)}
	}
}

func (m *Model) approveDraft() tea.Cmd {
	if len(m.drafts) == 0 || m.draftCursor < 0 || m.draftCursor >= len(m.drafts) {
		m.status = "No draft selected"
		return nil
	}
	draft := m.drafts[m.draftCursor]
	m.status = "Draft saved"
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		decks := content.DecksFromNotes(
			[]core.Note{draft.Note})
		for _, deck := range decks {
			if err := m.repo.UpsertDeck(ctx, deck); err != nil {
				return err
			}
		}

		cards, err := m.repo.DueCards(ctx, time.Now(), 0)
		if err != nil {
			return err
		}
		return draftApprovedMsg{noteID: draft.Note.ID, cards: cards}
	}
}

func (m *Model) discardDraft() tea.Cmd {
	if len(m.drafts) == 0 || m.draftCursor < 0 || m.draftCursor >= len(m.drafts) {
		m.status = "No draft selected"
		return nil
	}
	m.drafts = append(m.drafts[:m.draftCursor], m.drafts[m.draftCursor+1:]...)
	if m.draftCursor >= len(m.drafts) {
		m.draftCursor = maxInt(0, len(m.drafts)-1)
	}
	m.status = "Draft discarded"
	return nil
}

func (m *Model) discardAllDrafts() tea.Cmd {
	if len(m.drafts) == 0 {
		m.status = "No drafts to discard"
		return nil
	}
	count := len(m.drafts)
	m.drafts = nil
	m.draftCursor = 0
	m.status = fmt.Sprintf("Discarded %d drafts", count)
	return nil
}

func (m *Model) generateDrafts() tea.Cmd {
	return m.startDrafting()
}

func (m *Model) startDrafting() tea.Cmd {
	if m.drafting {
		return nil
	}
	if strings.TrimSpace(m.aiInput) == "" {
		m.status = "Enter a topic before generating AI drafts"
		return nil
	}
	if m.aiProvider == nil {
		m.status = "AI provider is disabled. Enable it in Settings."
		return nil
	}

	// Extract tags from input (e.g. "doctor visit #medical #b1")
	var tags []string
	cleanTopic := ""
	words := strings.Fields(m.aiInput)
	for _, w := range words {
		if strings.HasPrefix(w, "#") && len(w) > 1 {
			tags = append(tags, w[1:])
		} else {
			if cleanTopic != "" {
				cleanTopic += " "
			}
			cleanTopic += w
		}
	}

	sourceText := cleanTopic
	if m.draftSource != "" {
		sourceText = fmt.Sprintf("Topic/User Input: %s\n\nDictionary Context:\n%s", cleanTopic, m.draftSource)
	}

	m.drafting = true
	m.draftCancelled = false
	m.status = "AI is drafting flashcards..."
	return tea.Batch(
		m.tickSpinner(),
		func() tea.Msg {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			deckID := m.deck.ID
			if deckID == "" {
				deckID = "Imported"
			}

			drafts, err := m.aiProvider.GenerateDrafts(ctx, ai.DraftRequest{
				SourceText: sourceText,
				DeckID:     deckID,
				Tags:       tags,
			})

			if err != nil {
				return err
			}
			return draftsMsg(drafts)
		},
	)
}

func (m *Model) setCardBookmarkLocal(cardID string, bookmarked bool) {
	for i := range m.allDue {
		if m.allDue[i].ID == cardID {
			m.allDue[i].Bookmarked = bookmarked
			break
		}
	}
	m.applyDeckFilter()
}

func (m *Model) toggleBookmark() tea.Cmd {
	if len(m.dueCards) == 0 {
		return nil
	}
	card := m.dueCards[clampInt(m.cursor, 0, len(m.dueCards)-1)]
	next := !card.Bookmarked
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		if err := m.repo.SetCardBookmark(ctx, card.ID, next); err != nil {
			return err
		}
		return bookmarkToggledMsg{cardID: card.ID, bookmarked: next}
	}
}

func (m *Model) toggleBookmarkFilter() tea.Cmd {
	m.bookmarkFilter = !m.bookmarkFilter
	if m.bookmarkFilter {
		return m.loadBookmarkedDueCards()
	}
	return m.loadDueCards()
}

// reloadDueForCurrentFilter starts a fresh due-queue load that matches the
// live bookmark filter. Used when an async grade/undo/suspend finishes after
// the user flipped the filter mid-flight.
func (m *Model) reloadDueForCurrentFilter() tea.Cmd {
	if m.bookmarkFilter {
		return m.loadBookmarkedDueCards()
	}
	return m.loadDueCards()
}

func (m *Model) suspendCard() tea.Cmd {
	if len(m.dueCards) == 0 {
		return nil
	}
	card := m.dueCards[clampInt(m.cursor, 0, len(m.dueCards)-1)]
	bookmarkFilter := m.bookmarkFilter
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		if err := m.repo.SetCardSuspended(ctx, card.ID, true); err != nil {
			return err
		}
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
		return cardSuspendedMsg{cardID: card.ID, cards: cards, decks: decks, stats: stats, bookmarkFilter: bookmarkFilter}
	}
}

func (m *Model) toggleReviewHistory() tea.Cmd {
	if len(m.dueCards) == 0 {
		return nil
	}
	cardID := m.dueCards[clampInt(m.cursor, 0, len(m.dueCards)-1)].ID
	if m.showReviewHistory && m.reviewHistoryCard == cardID {
		m.showReviewHistory = false
		m.reviewHistory = nil
		m.reviewHistoryCard = ""
		m.status = "Review history hidden"
		return nil
	}
	m.reviewHistoryCard = cardID
	return m.loadReviewHistory(cardID)
}

func (m *Model) selectMCQChoice(key string) {
	if len(m.dueCards) == 0 || m.cursor >= len(m.dueCards) {
		return
	}
	idx, _ := strconv.Atoi(key)
	m.mcqChoice = idx - 1
	m.mcqAnswered = true
	card := m.dueCards[clampInt(m.cursor, 0, len(m.dueCards)-1)]
	if m.mcqChoice >= 0 && m.mcqChoice < len(card.Choices) {
		m.mcqCorrect = card.Answer == card.Choices[m.mcqChoice]
	} else {
		m.mcqCorrect = false
	}
}

func (m *Model) deleteSelectedCard() tea.Cmd {
	if len(m.browserCards) == 0 {
		return nil
	}
	card := m.browserCards[clampInt(m.browserCursor, 0, len(m.browserCards)-1)]
	m.confirmingDelete = true
	m.deleteAction = m.executeDeleteSelectedCard
	m.status = fmt.Sprintf("Delete card '%s'? (y/n)", card.Prompt)
	return nil
}

func (m *Model) executeDeleteSelectedCard() tea.Cmd {
	if len(m.browserCards) == 0 {
		return nil
	}
	card := m.browserCards[clampInt(m.browserCursor, 0, len(m.browserCards)-1)]
	m.status = "Deleting card..."
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		if err := m.repo.DeleteCard(ctx, card.ID); err != nil {
			return err
		}
		return m.loadBrowserCards()()
	}
}

func (m *Model) deleteReviewCard() tea.Cmd {
	if len(m.dueCards) == 0 {
		return nil
	}
	card := m.dueCards[clampInt(m.cursor, 0, len(m.dueCards)-1)]
	m.confirmingDelete = true
	m.deleteAction = m.executeDeleteReviewCard
	m.status = fmt.Sprintf("Delete card '%s'? (y/n)", strings.Split(card.Prompt, "\n")[0])
	return nil
}

func (m *Model) executeDeleteReviewCard() tea.Cmd {
	if len(m.dueCards) == 0 {
		return nil
	}
	card := m.dueCards[clampInt(m.cursor, 0, len(m.dueCards)-1)]
	m.status = "Deleting card..."
	m.dueLoadID++
	id := m.dueLoadID
	bookmarkFilter := m.bookmarkFilter
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		if err := m.repo.DeleteCard(ctx, card.ID); err != nil {
			return err
		}
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
		if bookmarkFilter {
			return bookmarkedDueCardsMsg{id: id, cards: cards}
		}
		return dueCardsMsg{id: id, cards: cards}
	}
}

func (m *Model) setCramFilter(idx int) tea.Cmd {
	types := []string{"bookmarked", "suspended", "leech", "flagged", "all"}
	if idx < 1 || idx > len(types) {
		return nil
	}
	m.cramType = types[idx-1]
	m.cramActive = false
	return m.loadCramCards()
}

func (m *Model) gradeCramCard(grade core.ReviewGrade) tea.Cmd {
	m.cramReviewed++
	if grade != core.GradeAgain {
		m.cramCorrect++
	}
	m.cramActive = false
	m.cramRevealed = false
	m.revealState = RevealIdle
	m.revealProgress = 0
	m.flipProgress = 0
	m.flipFrame = 0
	m.cramCursor++
	if m.cramCursor >= len(m.cramCards) {
		m.cramCursor = 0
		m.status = fmt.Sprintf("Cram session finished (%d/%d correct)", m.cramCorrect, m.cramReviewed)
	}
	return nil
}

func (m *Model) handleSettingsEnter() tea.Cmd {
	switch m.settingsCursor {
	case 0:
		// Provider cycle: disabled → offline → template → openai → anthropic → ollama → disabled
		cycle := []string{"disabled", "offline", "template", "openai", "anthropic", "ollama"}
		next := "disabled"
		for i, name := range cycle {
			if name == m.aiProviderName {
				next = cycle[(i+1)%len(cycle)]
				break
			}
		}
		m.aiProviderName = next
		m.aiProvider = buildProvider(next, m.aiSecrets, m.aiTemplates, m.currentAITemplateSet())
		m.status = fmt.Sprintf("Switched to %s AI provider", m.aiProviderName)
		m.persistConfig()

	case 1:
		// Dictionary cycle
		cycle := []string{"Local TUI", "dict.cc", "Linguee", "Leo", "Duden", "Pons", "Cambridge", "Google Translate"}
		next := "Local TUI"
		for i, name := range cycle {
			if strings.EqualFold(name, m.dictionaryProvider) {
				next = cycle[(i+1)%len(cycle)]
				break
			}
		}
		m.dictionaryProvider = next
		m.status = fmt.Sprintf("Switched to %s dictionary", m.dictionaryProvider)
		m.persistConfig()

	case 2, 3, 4:
		activeSet := m.currentAITemplateSet()
		if activeSet == "" {
			m.status = "No template set available to edit"
			return nil
		}
		if m.aiTemplates == nil {
			m.aiTemplates = make(map[string]map[string]string)
		}
		if m.aiTemplates[activeSet] == nil {
			m.aiTemplates[activeSet] = make(map[string]string)
		}
		m.editingTemplate = true
		key := ""
		switch m.settingsCursor {
		case 2:
			key = "front"
		case 3:
			key = "back"
		case 4:
			key = "example"
		}
		m.originalTemplateValue = m.aiTemplates[activeSet][key]
	case 6:
		m.autoPlayAudio = !m.autoPlayAudio
		status := "disabled"
		if m.autoPlayAudio {
			status = "enabled"
		}
		m.status = fmt.Sprintf("Auto-play audio %s", status)
		m.persistConfig()
	case 7:
		m.strictNormalization = !m.strictNormalization
		status := "disabled"
		if m.strictNormalization {
			status = "enabled"
		}
		m.status = fmt.Sprintf("Strict normalization %s", status)
		m.persistConfig()
	case 8, 9, 10, 11, 12, 13, 14, 15:
		var provider string
		var key string
		switch m.settingsCursor {
		case 8:
			provider = "openai"
			key = "api_key"
		case 9:
			provider = "openai"
			key = "model"
		case 10:
			provider = "openai"
			key = "base_url"
		case 11:
			provider = "anthropic"
			key = "api_key"
		case 12:
			provider = "anthropic"
			key = "model"
		case 13:
			provider = "anthropic"
			key = "base_url"
		case 14:
			provider = "ollama"
			key = "model"
		case 15:
			provider = "ollama"
			key = "base_url"
		}
		m.editingSecretProvider = provider
		m.editingSecretKey = key
		m.originalSecretValue = m.getCredValue(provider, key)
		m.status = fmt.Sprintf("Editing %s %s — Enter to save, Esc to cancel", provider, key)
	case 16:
		if m.revealSpeed < 10 {
			m.revealSpeed++
		} else {
			m.revealSpeed = 0
		}
		m.status = fmt.Sprintf("Reveal speed set to %d", m.revealSpeed)
		if m.revealSpeed == 0 {
			m.status = "Reveal animation disabled (instant)"
		}
		m.persistConfig()
	}
	return nil
}

func credKeyForCursor(cursor int) string {
	switch cursor {
	case 8, 11:
		return "api_key"
	case 9, 12, 14:
		return "model"
	case 10, 13, 15:
		return "base_url"
	}
	return ""
}

// getCredValue returns the current credential value for the named provider and key.
func (m *Model) getCredValue(provider, key string) string {
	c := m.credsFor(provider)
	switch key {
	case "api_key":
		return c.APIKey
	case "model":
		return c.Model
	case "base_url":
		return c.BaseURL
	}
	return ""
}

func (m *Model) credsFor(provider string) app.ProviderCreds {
	switch provider {
	case "openai":
		return m.aiSecrets.OpenAI
	case "anthropic":
		return m.aiSecrets.Anthropic
	case "ollama":
		return m.aiSecrets.Ollama
	}
	return app.ProviderCreds{}
}

func (m *Model) setCredValue(provider, key, value string) {
	c := m.credsFor(provider)
	switch key {
	case "api_key":
		c.APIKey = value
	case "model":
		c.Model = value
	case "base_url":
		c.BaseURL = value
	}
	switch provider {
	case "openai":
		m.aiSecrets.OpenAI = c
	case "anthropic":
		m.aiSecrets.Anthropic = c
	case "ollama":
		m.aiSecrets.Ollama = c
	}
}

func (m *Model) templateKeyAtCursor() string {
	switch m.settingsCursor {
	case 2:
		return "front"
	case 3:
		return "back"
	case 4:
		return "example"
	}
	return "front"
}

func (m *Model) approveAllDrafts() tea.Cmd {
	if len(m.drafts) == 0 {
		return nil
	}
	m.status = "Approving all drafts..."
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		var notes []core.Note
		for _, d := range m.drafts {
			notes = append(notes, d.Note)
		}
		decks := content.DecksFromNotes(
			notes)
		for _, deck := range decks {
			if err := m.repo.UpsertDeck(ctx, deck); err != nil {
				return err
			}
		}
		cards, err := m.repo.DueCards(ctx, time.Now(), 0)
		if err != nil {
			return err
		}
		decks, err = m.repo.Decks(ctx)
		if err != nil {
			return err
		}
		return importDoneMsg{decks: decks, cards: cards, count: len(notes), path: "AI Drafts"}
	}
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

func (m *Model) exportStatsCSV() tea.Cmd {
	m.status = "Exporting deck statistics..."
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		decks, err := m.repo.Decks(ctx)
		if err != nil {
			return err
		}

		path := filepath.Join(filepath.Dir(m.exportPath), "deck_statistics.csv")
		file, err := os.Create(path)
		if err != nil {
			return err
		}
		defer file.Close()

		writer := csv.NewWriter(file)
		defer writer.Flush()

		headers := []string{"DeckID", "DeckName", "TotalCards", "NewCards", "YoungCards", "MatureCards", "TotalReviews", "SuccessRate"}
		if err := writer.Write(headers); err != nil {
			return err
		}

		for _, deck := range decks {
			stats, err := m.repo.DeckStatistics(ctx, deck.ID)
			if err != nil {
				return err
			}
			row := []string{
				deck.ID,
				deck.Name,
				fmt.Sprintf("%d", stats.TotalCards),
				fmt.Sprintf("%d", stats.NewCards),
				fmt.Sprintf("%d", stats.YoungCards),
				fmt.Sprintf("%d", stats.MatureCards),
				fmt.Sprintf("%d", stats.TotalReviews),
				fmt.Sprintf("%.2f", stats.SuccessRate),
			}
			if err := writer.Write(row); err != nil {
				return err
			}
		}

		return statusMsg{text: fmt.Sprintf("Exported deck stats to %s", path)}
	}
}

// fixProposalMsg carries the AI's proposed correction back to the model.
type fixProposalMsg struct {
	cardID   string
	oldNote  core.Note
	proposal ai.FixedNote
}

// fixErrorMsg carries an error from the fix flow back to the model.
type fixErrorMsg struct {
	cardID string
	err    error
}

// fixAppliedMsg signals that a fix has been persisted; carries refreshed
// due cards so the visible card is updated immediately.
type fixAppliedMsg struct {
	cardID string
	cards  []core.Card
}

// explainCard starts the AI-powered pedagogical explanation flow for the
// currently focused review card.
func (m *Model) explainCard() tea.Cmd {
	if m.activeView != ViewReview || len(m.dueCards) == 0 {
		return nil
	}
	if m.explanation != "" {
		m.explanation = ""
		m.status = "Explanation hidden"
		return nil
	}
	if m.explainingCard {
		return nil
	}
	card := m.dueCards[clampInt(m.cursor, 0, len(m.dueCards)-1)]
	provider := m.aiProvider
	if provider == nil {
		m.status = "Enable an AI provider in Settings to get explanations"
		return nil
	}
	m.explainingCard = true
	m.explainCardID = card.ID
	m.explanation = ""
	m.explainError = ""
	m.status = "Asking AI for explanation…"
	cardID := card.ID
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()
		explanation, err := ai.ExplainCard(ctx, provider, card)
		if err != nil {
			return explainErrorMsg{cardID: cardID, err: err}
		}
		return explainMsg{cardID: cardID, explanation: explanation}
	}
}

// reportCardWrong starts the fix-card flow for the currently focused
// review card. It fetches the underlying note, then sends it to the
// active AI provider for correction. The model state is updated via
// fixProposalMsg / fixErrorMsg.
func (m *Model) reportCardWrong() tea.Cmd {
	if m.activeView != ViewReview || len(m.dueCards) == 0 || m.fixingCard {
		return nil
	}
	card := m.dueCards[clampInt(m.cursor, 0, len(m.dueCards)-1)]
	if card.NoteID == "" {
		m.status = "Cannot fix this card (missing note id)"
		return nil
	}
	provider := m.aiProvider
	if provider == nil {
		m.status = "Enable an AI provider in Settings to fix cards"
		return nil
	}
	m.fixingCard = true
	m.fixCardID = card.ID
	m.fixOldNote = nil
	m.fixProposal = nil
	m.fixError = ""
	m.status = "Reporting card to AI for review…"
	noteID := card.NoteID
	cardID := card.ID
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()
		note, err := m.repo.GetNote(ctx, noteID)
		if err != nil || note.ID == "" {
			// Fall back to synthesizing a note from the visible card so
			// users on older data without notes can still get help.
			note = core.Note{
				ID:     noteID,
				DeckID: card.DeckID,
				Front:  card.Prompt,
				Back:   card.Answer,
				Extra:  card.Extra,
				Tags:   card.Tags,
			}
		}
		proposal, err := ai.FixCard(ctx, provider, ai.FixRequest{Note: note})
		if err != nil {
			return fixErrorMsg{cardID: cardID, err: err}
		}
		return fixProposalMsg{cardID: cardID, oldNote: note, proposal: proposal}
	}
}

// applyFixProposal persists the AI's proposed correction to the active
// note. SRS state is preserved because card IDs are deterministic from
// the note ID and CardsForNote regenerates them with the same IDs.
func (m *Model) applyFixProposal() tea.Cmd {
	if m.fixOldNote == nil || m.fixProposal == nil {
		return nil
	}
	old := *m.fixOldNote
	fix := *m.fixProposal
	cardID := m.fixCardID
	m.fixingCard = false
	m.fixOldNote = nil
	m.fixProposal = nil
	m.fixCardID = ""
	m.status = "Applying fix…"
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		updated := old
		updated.Front = fix.Front
		updated.Back = fix.Back
		updated.Extra = fix.Extra
		if fix.Example != "" {
			updated.Examples = []string{fix.Example}
		}
		updated.Cards = content.CardsForNote(updated)
		if err := m.repo.UpsertNote(ctx, updated); err != nil {
			return fixErrorMsg{cardID: cardID, err: err}
		}
		cards, err := m.repo.DueCards(ctx, time.Now(), 0)
		if err != nil {
			return fixErrorMsg{cardID: cardID, err: err}
		}
		return fixAppliedMsg{cardID: cardID, cards: cards}
	}
}

// discardFixProposal cancels the fix preview without persisting anything.
func (m *Model) discardFixProposal() {
	m.fixingCard = false
	m.fixOldNote = nil
	m.fixProposal = nil
	m.fixCardID = ""
	m.fixError = ""
	m.status = "Fix discarded"
}

func (m *Model) addWordOfTheDayToCollection() tea.Cmd {
	word := content.GetWordOfTheDay()
	m.status = fmt.Sprintf("Adding %s to collection...", word.German)
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		note := core.Note{
			ID:     "wotd-" + time.Now().Format("2006-01-02"),
			DeckID: "quick-add",
			Front:  word.German,
			Back:   word.English,
			Extra:  fmt.Sprintf("Plural: %s\nExample: %s", word.Plural, word.Example),
			Tags:   []string{"word-of-the-day"},
		}
		note.Cards = content.CardsForNote(note)

		if err := m.repo.UpsertNote(ctx, note); err != nil {
			return err
		}
		return statusMsg{text: "Added " + word.German + " to Quick Add deck"}
	}
}

func (m *Model) practiceVerbOfTheDay() tea.Cmd {
	verb := content.GetVerbOfTheDay()
	m.status = fmt.Sprintf("Practicing %s...", verb.German)

	return tea.Sequence(
		m.updateView(ViewPractice),
		func() tea.Msg {
			// Find the verb in the loaded conjugation items and jump to it.
			st := m.trainerStateFor(PracticeSubViewConjugation)
			for i, item := range st.items {
				if item.Title == verb.German {
					m.practiceSubView = PracticeSubViewConjugation
					st.index = i
					st.revealed = false
					st.input = ""
					m.status = fmt.Sprintf("Conjugate: %s", verb.German)
					return nil
				}
			}
			// Fallback: load conjugation items if not present yet.
			return m.enterPracticeMode(PracticeSubViewConjugation)()
		},
	)
}

func (m *Model) searchGrammarTipInBrowser() tea.Cmd {
	tip := content.GetDailyGrammarTip()
	m.status = fmt.Sprintf("Searching for %s in browser...", tip.Title)

	m.activeView = ViewBrowser
	m.browserSearch = tip.Title
	m.browserDeckID = ""
	m.browserTag = ""
	return m.loadBrowserCards()
}
