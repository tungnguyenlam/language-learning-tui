package content

import (
	"path/filepath"

	"deutsch-tui/internal/core"
)

func StandardDecks() []core.Deck {
	decks := []core.Deck{
		StarterDeck(),
		CommonVerbsDeck(),
		GermanExpandedDeck(),
		PrepositionsDeck(),
		B1NatureDeck(),
		AdvancedEmotionsDeck(),
		BusinessDeck(),
	}

	// Load all embedded TSV decks
	embeddedPaths := EmbeddedDeckPaths()
	for _, path := range embeddedPaths {
		file, err := EmbeddedDecks.Open(path)
		if err != nil {
			continue
		}
		notes, err := ImportAnkiTSV(file, ImportOptions{
			DefaultDeck: filepath.Base(path)[:len(filepath.Base(path))-4],
		})
		file.Close()
		if err != nil {
			continue
		}
		decks = append(decks, DecksFromNotes(notes)...)
	}

	return decks
}
