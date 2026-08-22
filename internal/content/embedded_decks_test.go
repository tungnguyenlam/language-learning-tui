package content

import (
	"strings"
	"testing"
)

func TestUrbanMobilityDeckIsEmbedded(t *testing.T) {
	deck, err := DeckByID("b2_urban_mobility")
	if err != nil {
		t.Fatalf("DeckByID failed: %v", err)
	}
	if deck == nil {
		t.Fatal("b2_urban_mobility deck was not loaded")
	}
	if deck.Name != "b2-urban-mobility" {
		t.Fatalf("deck.Name = %q, want b2-urban-mobility", deck.Name)
	}
	if len(deck.Notes) < 40 {
		t.Fatalf("len(deck.Notes) = %d, want at least 40", len(deck.Notes))
	}

	var hasMCQ, hasCloze bool
	for _, note := range deck.Notes {
		if note.Type == "MCQ" {
			hasMCQ = true
		}
		if note.Type == "Cloze" {
			hasCloze = true
		}
	}
	if !hasMCQ || !hasCloze {
		t.Fatalf("deck should include both MCQ and Cloze notes; hasMCQ=%v hasCloze=%v", hasMCQ, hasCloze)
	}
}

func TestEmailPhoneCommunicationDeckIsEmbedded(t *testing.T) {
	deck, err := DeckByID("b1_email_phone_communication")
	if err != nil {
		t.Fatalf("DeckByID failed: %v", err)
	}
	if deck == nil {
		t.Fatal("b1_email_phone_communication deck was not loaded")
	}
	if deck.Name != "b1-email-phone-communication" {
		t.Fatalf("deck.Name = %q, want b1-email-phone-communication", deck.Name)
	}
	if len(deck.Notes) < 40 {
		t.Fatalf("len(deck.Notes) = %d, want at least 40", len(deck.Notes))
	}

	var hasEmail, hasPhone bool
	for _, note := range deck.Notes {
		for _, tag := range note.Tags {
			if tag == "email" {
				hasEmail = true
			}
			if tag == "phone" {
				hasPhone = true
			}
		}
		if len(note.Cards) == 0 {
			t.Fatalf("note %q has no generated cards", note.ID)
		}
	}
	if !hasEmail || !hasPhone {
		t.Fatalf("deck should include email and phone tagged notes; hasEmail=%v hasPhone=%v", hasEmail, hasPhone)
	}
}

func TestEmbeddedDecksDoNotLeakHeadersOrLiteralFieldsIntoCards(t *testing.T) {
	decks, err := AllDecks()
	if err != nil {
		t.Fatalf("AllDecks failed: %v", err)
	}
	if len(decks) == 0 {
		t.Fatal("no decks loaded")
	}

	for _, deck := range decks {
		for _, note := range deck.Notes {
			if strings.EqualFold(strings.TrimSpace(note.ID), "front") {
				t.Fatalf("deck %q leaked a TSV header as a note: %+v", deck.ID, note)
			}
			if strings.EqualFold(strings.TrimSpace(note.Front), "back") || strings.EqualFold(strings.TrimSpace(note.Back), "extra") {
				t.Fatalf("deck %q has header fields in note content: %+v", deck.ID, note)
			}
			for _, card := range note.Cards {
				answer := strings.TrimSpace(card.Answer)
				if strings.HasPrefix(answer, "Literal:") || strings.HasPrefix(answer, "Literally:") {
					t.Fatalf("deck %q card %q uses literal explanation as answer: prompt=%q answer=%q", deck.ID, card.ID, card.Prompt, card.Answer)
				}
			}
		}
	}
}

func TestEmbeddedIdiomsRomeReverseCardUsesGermanAnswer(t *testing.T) {
	deck, err := DeckByID("b1_idioms")
	if err != nil {
		t.Fatalf("DeckByID failed: %v", err)
	}
	if deck == nil {
		t.Fatal("b1_idioms deck was not loaded")
	}

	for _, note := range deck.Notes {
		for _, card := range note.Cards {
			if card.Prompt == "All roads lead to Rome" {
				if card.Answer != "Alle Wege führen nach Rom" {
					t.Fatalf("Rome reverse card answer = %q, want Alle Wege führen nach Rom", card.Answer)
				}
				return
			}
		}
	}
	t.Fatal("Rome reverse card was not generated")
}

func TestThinEmbeddedDecksHaveUsefulCoverage(t *testing.T) {
	cases := []struct {
		id        string
		min       int
		wantFront string
	}{
		{"b2_c1_news", 35, "die Eilmeldung"},
		{"a1_emergency", 35, "der Notruf"},
		{"b1_proverbs", 35, "Glück im Unglück."},
		{"a2_grammar_prepositions", 30, ""},
		{"b2_environment", 25, "erneuerbare Energien"},
		{"a1_essential", 40, "Ich heiße …"},
		{"b1_workplace_office", 40, "die Probezeit"},
		{"b2_urban_mobility", 40, "der Pendler"},
		{"a1_numbers_time", 40, "dreizehn"},
	}
	for _, c := range cases {
		deck, err := DeckByID(c.id)
		if err != nil {
			t.Fatalf("DeckByID(%q) failed: %v", c.id, err)
		}
		if deck == nil {
			t.Fatalf("deck %q was not loaded", c.id)
		}
		if len(deck.Notes) < c.min {
			t.Errorf("deck %q notes = %d, want at least %d", c.id, len(deck.Notes), c.min)
		}
		if c.wantFront == "" {
			continue
		}
		found := false
		for _, note := range deck.Notes {
			if note.Front == c.wantFront {
				found = true
				if len(note.Cards) == 0 {
					t.Errorf("deck %q note %q has no cards", c.id, c.wantFront)
				}
				break
			}
		}
		if !found {
			t.Errorf("deck %q missing note with front %q", c.id, c.wantFront)
		}
	}
}

func TestEmbeddedContentFixes(t *testing.T) {
	idioms, err := DeckByID("b1_idioms")
	if err != nil || idioms == nil {
		t.Fatalf("b1_idioms: %v %v", idioms, err)
	}
	var sawOderNie, sawGold bool
	for _, note := range idioms.Notes {
		if note.Front == "Jetzt oder nie" {
			sawOderNie = true
		}
		if note.Front == "Jetzt or nie" {
			t.Error("b1_idioms still has English 'or' in Jetzt or nie")
		}
		if note.Front == "Reden ist Silber, Schweigen ist Gold" {
			sawGold = true
		}
		if note.Front == "Reden ist Silber, Schweigen ist gold" {
			t.Error("b1_idioms still lowercases Gold in Schweigen ist Gold")
		}
	}
	if !sawOderNie {
		t.Error("expected Jetzt oder nie in b1_idioms")
	}
	if !sawGold {
		t.Error("expected capitalized Gold in Schweigen ist Gold")
	}

	env, err := DeckByID("b2_environment")
	if err != nil || env == nil {
		t.Fatalf("b2_environment: %v %v", env, err)
	}
	for _, note := range env.Notes {
		if note.Front == "die Erneuerbare Energien" {
			t.Error("b2_environment still uses ungrammatical die Erneuerbare Energien")
		}
	}

	prep, err := DeckByID("a2_grammar_prepositions")
	if err != nil || prep == nil {
		t.Fatalf("a2_grammar_prepositions: %v %v", prep, err)
	}
	wanted := map[string]bool{"gegen": false, "nach": false, "zu": false, "von": false, "hinter": false, "wegen": false}
	for _, note := range prep.Notes {
		if note.Type != "Cloze" {
			t.Errorf("preposition note %q type = %q, want Cloze", note.ID, note.Type)
		}
		for prepWord := range wanted {
			if strings.Contains(strings.ToLower(note.Front), "{{c1::"+prepWord+"::") {
				wanted[prepWord] = true
			}
		}
	}
	for prepWord, found := range wanted {
		if !found {
			t.Errorf("a2_grammar_prepositions missing cloze for %q", prepWord)
		}
	}

	weather, err := DeckByID("a1_weather_seasons")
	if err != nil || weather == nil {
		t.Fatalf("a1_weather_seasons: %v %v", weather, err)
	}
	var sawSchneien, sawEnglishSnowFront bool
	for _, note := range weather.Notes {
		if note.Front == "schneien" {
			sawSchneien = true
			if note.Back != "to snow" {
				t.Errorf("schneien back = %q, want to snow", note.Back)
			}
		}
		if note.Front == "Schneiden" && note.Back == "to snow" {
			t.Error("a1_weather_seasons still uses schneiden for 'to snow'")
		}
		if note.Front == "It snows" || note.Front == "It rains" {
			sawEnglishSnowFront = true
		}
	}
	if !sawSchneien {
		t.Error("expected schneien (to snow) in a1_weather_seasons")
	}
	if sawEnglishSnowFront {
		t.Error("a1_weather_seasons still has English phrase fronts")
	}

	numbers, err := DeckByID("a1_numbers_time")
	if err != nil || numbers == nil {
		t.Fatalf("a1_numbers_time: %v %v", numbers, err)
	}
	var sawTomorrow, sawMorningNoun, sawAmbiguousMorgen bool
	for _, note := range numbers.Notes {
		if note.Front == "morgen" && strings.EqualFold(note.Back, "tomorrow") {
			sawTomorrow = true
		}
		if note.Front == "der Morgen" {
			sawMorningNoun = true
		}
		if note.Back == "Tomorrow/Morning" {
			sawAmbiguousMorgen = true
		}
	}
	if !sawTomorrow || !sawMorningNoun {
		t.Errorf("a1_numbers_time should split morgen (tomorrow) and der Morgen (morning); tomorrow=%v morning=%v", sawTomorrow, sawMorningNoun)
	}
	if sawAmbiguousMorgen {
		t.Error("a1_numbers_time still has ambiguous Tomorrow/Morning gloss")
	}

	essential, err := DeckByID("a1_essential")
	if err != nil || essential == nil {
		t.Fatalf("a1_essential: %v %v", essential, err)
	}
	var sawSoLala, sawEnglishSoSo bool
	for _, note := range essential.Notes {
		if note.Front == "so lala" {
			sawSoLala = true
		}
		if note.Front == "So-so" {
			sawEnglishSoSo = true
		}
	}
	if !sawSoLala {
		t.Error("expected German 'so lala' in a1_essential")
	}
	if sawEnglishSoSo {
		t.Error("a1_essential still uses English 'So-so' as the front")
	}
}
