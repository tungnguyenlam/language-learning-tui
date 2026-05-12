package content

import (
	"path/filepath"

	"deutsch-tui/internal/core"
)

func StandardDecks() []core.Deck {
	decks := []core.Deck{
		StarterDeck(),
		CommonVerbsDeck(),
		ConfusableWordsDeck(),
		GermanExpandedDeck(),
		PrepositionsDeck(),
		B1NatureDeck(),
		AdvancedEmotionsDeck(),
		BusinessDeck(),
		IdiomsDeck(),
		SlangDeck(),
		PhrasalVerbsDeck(),
		MedicalGermanDeck(),
		A1FamilyDeck(),
		A1GreetingsDeck(),
		A1HobbiesDeck(),
		A1FoodDrinkDeck(),
		A1TravelDeck(),
		B1TechnologyDeck(),
		B1FinanceDeck(),
		B1EnvironmentDeck(),
		B1HealthDeck(),
		B1CultureDeck(),
		B1SportsDeck(),
		B1EducationDeck(),
		A2DailyLifeDeck(),
		B2MediaDeck(),
		C1AcademicDeck(),
		A2ShoppingDeck(),
		B2BusinessDeck(),
		C2LegalDeck(),
		B2TravelDeck(),
		B2EnvironmentDeck(),
		C1ScienceDeck(),
		B1ArtDeck(),
		A2HobbiesIIDeck(),
		B2ScienceIIDeck(),
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
