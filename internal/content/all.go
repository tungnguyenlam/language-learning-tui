package content

import "deutsch-tui/internal/core"

func StandardDecks() []core.Deck {
	return []core.Deck{
		StarterDeck(),
		GermanExpandedDeck(),
		ScienceTechDeck(),
		ProverbsDeck(),
	}
}
