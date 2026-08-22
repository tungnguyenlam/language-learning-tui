package content

import (
	"embed"
	"path/filepath"
	"strings"

	"deutsch-tui/internal/core"
)

func AllDecks() ([]core.Deck, error) {
	embedded, err := loadEmbeddedTSVDecks()
	if err != nil {
		return nil, err
	}
	seen := make(map[string]bool, len(embedded))
	decks := make([]core.Deck, 0, len(embedded))
	for _, deck := range embedded {
		if seen[deck.ID] {
			continue
		}
		seen[deck.ID] = true
		decks = append(decks, deck)
	}
	for _, deck := range StandardDecks() {
		if seen[deck.ID] {
			continue
		}
		seen[deck.ID] = true
		decks = append(decks, deck)
	}
	return decks, nil
}

func DeckByID(id string) (*core.Deck, error) {
	decks, err := AllDecks()
	if err != nil {
		return nil, err
	}
	for _, deck := range decks {
		if deck.ID == id {
			d := deck
			return &d, nil
		}
	}
	return nil, nil
}

//go:embed testdata/german-decks/*.tsv
var EmbeddedDecks embed.FS

func EmbeddedDeckPaths() []string {
	entries, err := EmbeddedDecks.ReadDir("testdata/german-decks")
	if err != nil {
		return nil
	}
	var paths []string
	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), ".tsv") {
			paths = append(paths, filepath.Join("testdata/german-decks", entry.Name()))
		}
	}
	return paths
}

func loadEmbeddedTSVDecks() ([]core.Deck, error) {
	entries, err := EmbeddedDecks.ReadDir("testdata/german-decks")
	if err != nil {
		return nil, err
	}

	deckMap := make(map[string]*core.Deck)
	for _, entry := range entries {
		if !strings.HasSuffix(entry.Name(), ".tsv") {
			continue
		}
		path := filepath.Join("testdata/german-decks", entry.Name())
		data, err := EmbeddedDecks.ReadFile(path)
		if err != nil {
			continue
		}
		notes, err := ImportAnkiTSV(strings.NewReader(string(data)), ImportOptions{})
		if err != nil {
			continue
		}

		deckName := strings.TrimSuffix(entry.Name(), ".tsv")
		deckID := ToDeckID(deckName)
		for i := range notes {
			notes[i].DeckID = deckID
			notes[i].Cards = CardsForNote(notes[i])
			for j := range notes[i].Cards {
				notes[i].Cards[j].DeckID = deckID
			}
		}
		if deck, ok := deckMap[deckID]; ok {
			deck.Notes = append(deck.Notes, notes...)
		} else {
			deckMap[deckID] = &core.Deck{
				ID:    deckID,
				Name:  deckName,
				Notes: notes,
			}
		}
	}

	decks := make([]core.Deck, 0, len(deckMap))
	for _, deck := range deckMap {
		decks = append(decks, *deck)
	}
	return decks, nil
}

func ToDeckID(name string) string {
	id := strings.ReplaceAll(name, "-", "_")
	id = strings.ReplaceAll(id, " ", "_")
	return strings.ToLower(id)
}
