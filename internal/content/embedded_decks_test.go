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
	if len(deck.Notes) < 25 {
		t.Fatalf("len(deck.Notes) = %d, want at least 25", len(deck.Notes))
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
