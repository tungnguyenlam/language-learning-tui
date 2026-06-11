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
		A1ColorsShapesDeck(),
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
		C1AcademicDeck(),
		B2TravelDeck(),
		B2BusinessDeck(),
		C2LegalDeck(),
		B2EnvironmentDeck(),
		C1ScienceDeck(),
		B1ArtDeck(),
		A2HobbiesIIDeck(),
		B2ScienceIIDeck(),
		C2LiteratureDeck(),
		C1PoliticsDeck(),
		B1JobsDeck(),
		B2MusicDeck(),
		B1TransportDeck(),
		A1OfficeDeck(),
		B2ClimateDeck(),
		A1HouseFurnitureDeck(),
		A1AnimalsDeck(),
		A2BodyHealthDeck(),
		B1CookingDeck(),
		B1WeatherDeck(),
		A2ShoppingDeck(),
		B2CultureLeisureDeck(),
		B2TravelAdventureDeck(),
		C1PsychologyMindDeck(),
		B1BureaucracyAppointmentsDeck(),
		B2DigitalPrivacyDeck(),
		A2TravelBookingDeck(),
		B1HousingApartmentDeck(),
		C1EnvironmentSustainabilityDeck(),
		B2JobApplicationDeck(),
		C1PhilosophyEthicsDeck(),
		A1CityDirectionsDeck(),
		B2ProgrammingDeck(),
		B2LogisticsDeck(),
		B2BusinessMeetingsDeck(),
		A2MedicalAppointmentDeck(),
		A2WorkOfficeDeck(),
		C1SocialIssuesDeck(),
		FalseFriendsDeck(),
		C2FinanceDeck(),
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
