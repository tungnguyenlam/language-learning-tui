package tui

import (
	"context"
	"fmt"
	"math/rand/v2"
	"strings"
	"time"

	"deutsch-tui/internal/content"
	"deutsch-tui/internal/core"

	tea "charm.land/bubbletea/v2"
)

func (m *Model) loadBrowserCards() tea.Cmd {
	m.browserLoadID++
	id := m.browserLoadID
	deckID := m.browserDeckID
	search := m.browserSearch
	tag := m.browserTag
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		cards, err := m.repo.Cards(ctx, deckID, search, tag)
		if err != nil {
			return err
		}
		return browserCardsResultMsg{id: id, cards: cards}
	}
}

func (m *Model) loadDueCards() tea.Msg {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	cards, err := m.repo.DueCards(ctx, time.Now(), 0)
	if err != nil {
		return err
	}
	return dueCardsMsg(cards)
}

func (m *Model) loadBookmarkedDueCards() tea.Msg {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	// Use default due limit (0) so bookmark filter matches Review's full queue.
	cards, err := m.repo.DueCardsBookmarked(ctx, time.Now(), 0)
	if err != nil {
		return err
	}
	return bookmarkedDueCardsMsg(cards)
}

func (m *Model) loadCramCards() tea.Cmd {
	deckID := m.deck.ID
	cramType := m.cramType
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		var cards []core.Card
		var err error
		switch cramType {
		case "bookmarked", "suspended", "leech", "flagged":
			cards, err = m.repo.CardsWithFlag(ctx, deckID, cramType)
		default:
			cards, err = m.repo.Cards(ctx, deckID, "", "")
		}
		if err != nil {
			return err
		}
		return cramCardsMsg{cards: cards, cramType: cramType}
	}
}

func (m *Model) loadDecks() tea.Msg {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	decks, err := m.repo.Decks(ctx)
	if err != nil {
		return err
	}
	allDecks := append([]core.Deck{{ID: "", Name: "All Decks", Description: "Study cards from all decks combined."}}, decks...)
	return decksMsg(allDecks)
}

type statsMsg struct {
	stats     core.Statistics
	dictCount int
}

func (m *Model) loadStatistics() tea.Cmd {
	deckID := m.deck.ID
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		var stats core.Statistics
		var err error

		if deckID != "" {
			stats, err = m.repo.DeckStatistics(ctx, deckID)
		} else {
			stats, err = m.repo.Statistics(ctx)
		}
		if err != nil {
			return err
		}

		var count int
		if dictRepo, ok := m.repo.(core.DictionaryRepository); ok {
			count, _ = dictRepo.DictionaryCount(ctx)
		}

		return statsMsg{stats: stats, dictCount: count}
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

type recentDecksMsg []string

func (m *Model) loadRecentDecks() tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		data, err := m.repo.RecentDecks(ctx, 5)
		if err != nil {
			return err
		}
		return recentDecksMsg(data)
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

func (m *Model) loadReviewPredictions(cardID string) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		state, err := m.repo.GetReviewState(ctx, cardID)
		if err != nil {
			return err
		}
		predictions := m.scheduler.Predict(state, time.Now())
		return reviewPredictionsMsg(predictions)
	}
}

func (m *Model) reloadBrowserForSelectedDeck() tea.Cmd {
	m.browserDeckID = m.deck.ID
	m.browserSearch = ""
	m.browserCards = nil
	m.browserCursor = 0
	m.browserSelected = make(map[string]bool)
	m.clearReviewHistory()
	m.status = fmt.Sprintf("Browsing %s", m.deckLabel())
	return m.loadBrowserCards()
}

func (m *Model) loadPracticeItems() tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		cards, err := m.repo.Cards(ctx, "", "", "")
		if err != nil {
			return err
		}

		var items []practiceItem
		for _, card := range cards {
			info := content.AnalyzeCard(card.Prompt, card.Answer)
			if info.Kind == content.KindNoun && info.Article != "" {
				// Meaning is usually the side that isn't German
				meaning := card.Answer
				if content.Analyze(card.Answer).Article != "" {
					meaning = card.Prompt
				}

				items = append(items, practiceItem{
					Word:    info.Base,
					Article: info.Article,
					Meaning: meaning,
				})
			}
		}

		// Cards arrive in repository order, which would otherwise present the
		// same nouns in the same sequence every session.
		rand.Shuffle(len(items), func(i, j int) { items[i], items[j] = items[j], items[i] })

		return practiceItemsMsg(items)
	}
}

func extractPlural(extra string) string {
	lower := strings.ToLower(extra)
	idx := strings.Index(lower, "plural:")
	if idx == -1 {
		return ""
	}
	sub := extra[idx+7:]
	if end := strings.IndexAny(sub, ";\n\r"); end != -1 {
		sub = sub[:end]
	}
	return strings.TrimSpace(sub)
}

type practiceItemsMsg []practiceItem

func (m *Model) filteredDecks() []core.Deck {
	if m.deckFilter == "" {
		return m.decks
	}
	var filtered []core.Deck
	filter := strings.ToLower(m.deckFilter)
	for _, d := range m.decks {
		if strings.Contains(strings.ToLower(d.Name), filter) {
			filtered = append(filtered, d)
			continue
		}
		for _, tag := range d.Tags {
			if strings.Contains(strings.ToLower(tag), filter) {
				filtered = append(filtered, d)
				break
			}
		}
	}
	return filtered
}
