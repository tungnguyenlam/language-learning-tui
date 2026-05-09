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
	"deutsch-tui/internal/content"
	"deutsch-tui/internal/core"

	tea "charm.land/bubbletea/v2"
)

func (m *Model) gradeCard(grade core.ReviewGrade) tea.Cmd {
	if m.activeView != ViewReview || len(m.dueCards) == 0 || m.revealState != RevealRevealed || m.gradingInProgress {
		return nil
	}

	m.gradingInProgress = true
	m.status = fmt.Sprintf("Grade: %s", strings.ToUpper(string(grade[:1]))+string(grade[1:]))
	card := m.dueCards[clampInt(m.cursor, 0, len(m.dueCards)-1)]
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
		if m.bookmarkFilter {
			cards, err2 = m.repo.DueCardsBookmarked(ctx, time.Now(), 50)
		} else {
			cards, err2 = m.repo.DueCards(ctx, time.Now(), 500)
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
		return reviewRecordedMsg{cardID: card.ID, cards: cards, decks: decks, stats: stats, grade: grade}
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
		return reviewUndoneMsg{cardID: cardID, cards: cards, decks: decks, stats: stats, grade: grade}
	}
}

func (m *Model) setDailyGoal(goal int) tea.Cmd {
	if goal < 1 {
		goal = 1
	}
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
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		count := 0
		for _, deck := range content.StandardDecks() {
			if err := m.repo.UpsertDeck(ctx, deck); err != nil {
				return err
			}
			count += len(deck.Notes)
		}
		loadedDecks, err := m.repo.Decks(ctx)
		if err != nil {
			return err
		}
		cards, err := m.repo.DueCards(ctx, time.Now(), 500)
		if err != nil {
			return err
		}
		return importDoneMsg{decks: loadedDecks, cards: cards, count: count, path: "Standard Content"}
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
		cards, err := m.repo.DueCards(ctx, time.Now(), 500)
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
		cards, err := m.repo.DueCards(ctx, time.Now(), 500)
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
		seen := make(map[string]bool)
		notes := make([]core.Note, 0, len(cards))
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
			return fmt.Errorf("failed to create APKG file '%s': %w", filepath.Base(path), err)
		}
		defer file.Close()

		if err := content.ExportAnkiAPKG(file, notes); err != nil {
			return fmt.Errorf("failed to write APKG data to '%s': %w", filepath.Base(path), err)
		}
		return statusMsg{text: fmt.Sprintf("Exported %d notes to %s", len(notes), path)}
	}
}

func (m *Model) approveDraft() tea.Cmd {
	if len(m.drafts) == 0 || m.draftCursor >= len(m.drafts) {
		return nil
	}
	draft := m.drafts[m.draftCursor]
	m.status = "Approving draft..."
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

		cards, err := m.repo.DueCards(ctx, time.Now(), 500)
		if err != nil {
			return err
		}
		return draftApprovedMsg{noteID: draft.Note.ID, cards: cards}
	}
}

func (m *Model) discardDraft() tea.Cmd {
	if len(m.drafts) == 0 || m.draftCursor >= len(m.drafts) {
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

func (m *Model) generateDrafts() tea.Cmd {
	return m.startDrafting()
}

func (m *Model) startDrafting() tea.Cmd {
	if m.drafting {
		return nil
	}
	if m.aiProvider == nil {
		m.status = "AI provider is disabled. Enable it in Settings."
		return nil
	}
	m.drafting = true
	m.status = "AI is drafting flashcards..."
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		deckID := m.deck.ID
		if deckID == "" {
			deckID = "Imported"
		}

		drafts, err := m.aiProvider.GenerateDrafts(ctx, ai.DraftRequest{
			SourceText: m.aiInput,
			DeckID:     deckID,
		})

		if err != nil {
			return err
		}
		return draftsMsg(drafts)
	}
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
		return func() tea.Msg { return m.loadBookmarkedDueCards() }
	}
	return func() tea.Msg { return m.loadDueCards() }
}

func (m *Model) suspendCard() tea.Cmd {
	if len(m.dueCards) == 0 {
		return nil
	}
	card := m.dueCards[clampInt(m.cursor, 0, len(m.dueCards)-1)]
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
	m.cramCursor++
	if m.cramCursor >= len(m.cramCards) {
		m.cramCursor = 0
		m.status = "Cram session finished"
	}
	return nil
}

func (m *Model) handleSettingsEnter() tea.Cmd {
	switch m.settingsCursor {
	case 0:
		switch m.aiProviderName {
		case "offline":
			m.aiProviderName = "template"
			activeSet := ""
			activeSet = m.currentAITemplateSet()
			m.aiProvider = ai.TemplateProvider{
				Templates: m.aiTemplates,
				ActiveSet: activeSet,
			}
		case "template":
			m.aiProviderName = "disabled"
			m.aiProvider = nil
		default:
			m.aiProviderName = "offline"
			m.aiProvider = ai.OfflineProvider{}
		}
		m.status = fmt.Sprintf("Switched to %s AI provider", m.aiProviderName)
		if m.onConfigChange != nil {
			m.onConfigChange(m.aiProviderName, m.aiTemplates, m.autoPlayAudio, m.strictNormalization)
		}

	case 1, 2, 3:
		m.editingTemplate = true
		activeSet := ""
		activeSet = m.currentAITemplateSet()
		key := ""
		switch m.settingsCursor {
		case 1:
			key = "front"
		case 2:
			key = "back"
		case 3:
			key = "example"
		}
		m.originalTemplateValue = m.aiTemplates[activeSet][key]
	case 5:
		m.autoPlayAudio = !m.autoPlayAudio
		status := "disabled"
		if m.autoPlayAudio {
			status = "enabled"
		}
		m.status = fmt.Sprintf("Auto-play audio %s", status)
		if m.onConfigChange != nil {
			m.onConfigChange(m.aiProviderName, m.aiTemplates, m.autoPlayAudio, m.strictNormalization)
		}
	case 6:
		m.strictNormalization = !m.strictNormalization
		status := "disabled"
		if m.strictNormalization {
			status = "enabled"
		}
		m.status = fmt.Sprintf("Strict normalization %s", status)
		if m.onConfigChange != nil {
			m.onConfigChange(m.aiProviderName, m.aiTemplates, m.autoPlayAudio, m.strictNormalization)
		}
	}
	return nil
}

func (m *Model) currentTemplateKey() string {
	return m.templateKeyAtCursor()
}

func (m *Model) templateKeyAtCursor() string {
	switch m.settingsCursor {
	case 1:
		return "front"
	case 2:
		return "back"
	case 3:
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
		m.drafts = nil
		m.draftCursor = 0
		cards, err := m.repo.DueCards(ctx, time.Now(), 500)
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
