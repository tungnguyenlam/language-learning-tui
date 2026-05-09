package content

import "deutsch-tui/internal/core"

func StandardDecks() []core.Deck {
	return []core.Deck{
		StarterDeck(),
		ScienceTechDeck(),
		GermanExpandedDeck(),
		ProverbsDeck(),
		PhilosophyLiteratureDeck(),
		NewsDeck(),
		EnvironmentDeck(),
	}
}
