package content

import (
	"deutsch-tui/internal/core"
	"strings"
	"testing"
)

func TestImportAnkiTSV(t *testing.T) {
	input := "#separator:tab\n#deck:A1\nid-1\tder Apfel\tapple\tPlural: die Aepfel\ta1 noun\n"
	notes, err := ImportAnkiTSV(strings.NewReader(input), ImportOptions{})
	if err != nil {
		t.Fatalf("import failed: %v", err)
	}
	if len(notes) != 1 {
		t.Fatalf("notes = %d, want 1", len(notes))
	}
	if notes[0].DeckID != "A1" || notes[0].Front != "der Apfel" || notes[0].Back != "apple" {
		t.Fatalf("unexpected note: %+v", notes[0])
	}
}

func TestImportAnkiTSVWithNewlines(t *testing.T) {
	// Quoted field with a newline followed by something that looks like a comment
	input := "#separator:tab\nid-1\t\"Front\n#not-a-comment\"\tBack\t\t\t\t\n"
	notes, err := ImportAnkiTSV(strings.NewReader(input), ImportOptions{})
	if err != nil {
		t.Fatalf("import failed: %v", err)
	}
	if len(notes) != 1 {
		t.Fatalf("notes = %d, want 1", len(notes))
	}
	if notes[0].Front != "Front\n#not-a-comment" {
		t.Fatalf("front = %q, want \"Front\\n#not-a-comment\"", notes[0].Front)
	}
}

func TestExportAnkiTSVRoundTrip(t *testing.T) {
	deck := StarterDeck()
	out, err := ExportAnkiTSVString(deck.Notes[:1])
	if err != nil {
		t.Fatalf("export failed: %v", err)
	}
	notes, err := ImportAnkiTSV(strings.NewReader(out), ImportOptions{DefaultDeck: "fallback"})
	if err != nil {
		t.Fatalf("reimport failed: %v", err)
	}
	if notes[0].ID != deck.Notes[0].ID {
		t.Fatalf("id = %q, want %q", notes[0].ID, deck.Notes[0].ID)
	}
	if notes[0].DeckID != deck.Notes[0].DeckID {
		t.Fatalf("deck id = %q, want %q", notes[0].DeckID, deck.Notes[0].DeckID)
	}
}

func TestImportExportMCQ(t *testing.T) {
	input := "#separator:tab\n#deck:Grammar\nid-1\tIch gebe ___ das Buch.\tdem Mann\tDativ\tgrammar\tGrammar\tMCQ:den Mann|||der Mann|||dem Mann\n"
	notes, err := ImportAnkiTSV(strings.NewReader(input), ImportOptions{})
	if err != nil {
		t.Fatalf("import failed: %v", err)
	}
	if len(notes) != 1 {
		t.Fatalf("notes = %d, want 1", len(notes))
	}
	note := notes[0]
	if len(note.Choices) != 3 || note.Choices[2] != "dem Mann" {
		t.Fatalf("unexpected choices: %v", note.Choices)
	}

	mcqCard := false
	for _, card := range note.Cards {
		if card.Kind == core.CardKindMCQ {
			mcqCard = true
			if len(card.Choices) != 3 || card.Choices[2] != "dem Mann" {
				t.Fatalf("unexpected card choices: %v", card.Choices)
			}
		}
	}
	if !mcqCard {
		t.Fatal("no MCQ card generated")
	}

	// Test round-trip export
	out, err := ExportAnkiTSVString(notes)
	if err != nil {
		t.Fatalf("export failed: %v", err)
	}
	if !strings.Contains(out, "MCQ:den Mann|||der Mann|||dem Mann") {
		t.Fatalf("exported TSV missing MCQ choices: %q", out)
	}
}

func TestImportCloze(t *testing.T) {
	input := "#separator:tab\n#deck:Cloze\nid-1\tIch {{c1::gebe}} dem Mann das Buch.\tgive\t\t\t\t\n"
	notes, err := ImportAnkiTSV(strings.NewReader(input), ImportOptions{})
	if err != nil {
		t.Fatalf("import failed: %v", err)
	}
	if len(notes) != 1 {
		t.Fatalf("notes = %d, want 1", len(notes))
	}
	note := notes[0]

	clozeCard := false
	for _, card := range note.Cards {
		if card.Kind == core.CardKindCloze {
			clozeCard = true
			if card.Prompt != "Ich [...] dem Mann das Buch." {
				t.Fatalf("unexpected prompt: %q", card.Prompt)
			}
			if card.Answer != "Ich gebe dem Mann das Buch." {
				t.Fatalf("unexpected answer: %q", card.Answer)
			}
			if len(card.Choices) != 1 || card.Choices[0] != "gebe" {
				t.Fatalf("unexpected choices: %v", card.Choices)
			}
		}
	}
	if !clozeCard {
		t.Fatal("no Cloze card generated")
	}
}

func TestImportMultiCloze(t *testing.T) {
	input := "#separator:tab\n#deck:Cloze\nid-1\t{{c1::Ich}} {{c2::gebe}} dem Mann das Buch.\tgive\t\t\t\t\n"
	notes, err := ImportAnkiTSV(strings.NewReader(input), ImportOptions{})
	if err != nil {
		t.Fatalf("import failed: %v", err)
	}
	note := notes[0]

	clozeCount := 0
	for _, card := range note.Cards {
		if card.Kind == core.CardKindCloze {
			clozeCount++
			if clozeCount == 1 {
				if card.Prompt != "[...] gebe dem Mann das Buch." {
					t.Fatalf("unexpected prompt 1: %q", card.Prompt)
				}
			} else if clozeCount == 2 {
				if card.Prompt != "Ich [...] dem Mann das Buch." {
					t.Fatalf("unexpected prompt 2: %q", card.Prompt)
				}
			}
		}
	}
	if clozeCount != 2 {
		t.Fatalf("clozeCount = %d, want 2", clozeCount)
	}
}

func TestImportGroupedAndOutOfOrderClozeOrdinals(t *testing.T) {
	note := core.Note{
		ID:    "grouped-cloze",
		Front: "Ich {{c2::gebe}} {{c1::dem}} {{c1::Mann}}.",
	}
	cards := CardsForNote(note)
	if len(cards) != 2 {
		t.Fatalf("len(cards) = %d, want 2 cards (one per ordinal)", len(cards))
	}

	first, second := cards[0], cards[1]
	if first.ID != "grouped-cloze:cloze-1" || second.ID != "grouped-cloze:cloze-2" {
		t.Fatalf("cloze IDs = %q, %q; want ordinal order", first.ID, second.ID)
	}
	if first.Prompt != "Ich gebe [...] [...]." {
		t.Fatalf("ordinal 1 prompt = %q", first.Prompt)
	}
	if got, want := strings.Join(first.Choices, "|"), "dem|Mann"; got != want {
		t.Fatalf("ordinal 1 choices = %q, want %q", got, want)
	}
	if second.Prompt != "Ich [...] dem Mann." {
		t.Fatalf("ordinal 2 prompt = %q", second.Prompt)
	}
	if len(second.Choices) != 1 || second.Choices[0] != "gebe" {
		t.Fatalf("ordinal 2 choices = %v, want [gebe]", second.Choices)
	}
}

func TestImportReverse(t *testing.T) {
	input := "#separator:tab\n#deck:Reverse\nid-1\tder Apfel\tapple\t\t\t\tReverse\n"
	notes, err := ImportAnkiTSV(strings.NewReader(input), ImportOptions{})
	if err != nil {
		t.Fatalf("import failed: %v", err)
	}
	if len(notes) != 1 {
		t.Fatalf("notes = %d, want 1", len(notes))
	}
	note := notes[0]
	if note.Type != "Reverse" {
		t.Fatalf("note.Type = %q, want Reverse", note.Type)
	}

	if len(note.Cards) != 2 {
		t.Fatalf("len(note.Cards) = %d, want 2", len(note.Cards))
	}

	frontCard := note.Cards[0]
	if frontCard.Prompt != "der Apfel" || frontCard.Answer != "apple" {
		t.Fatalf("unexpected front card: %+v", frontCard)
	}

	backCard := note.Cards[1]
	if backCard.Prompt != "apple" || backCard.Answer != "der Apfel" {
		t.Fatalf("unexpected back card: %+v", backCard)
	}
}

func TestImportHeaderedReverseDeckWithoutIDs(t *testing.T) {
	input := strings.Join([]string{
		"#separator:tab",
		"#html:false",
		"#deck:B1 German Idioms",
		"front\tback\textra\ttags\tnotetype",
		"Alle Wege führen nach Rom\tAll roads lead to Rome\tLiteral: all ways lead to Rome\tidiom b1\tB1 German Idioms\tReverse",
		"",
	}, "\n")

	notes, err := ImportAnkiTSV(strings.NewReader(input), ImportOptions{})
	if err != nil {
		t.Fatalf("import failed: %v", err)
	}
	if len(notes) != 1 {
		t.Fatalf("notes = %d, want 1", len(notes))
	}

	note := notes[0]
	if note.ID == "front" {
		t.Fatalf("header row leaked as note: %+v", note)
	}
	if note.Type != "Reverse" {
		t.Fatalf("note.Type = %q, want Reverse", note.Type)
	}
	if note.DeckID != "B1 German Idioms" {
		t.Fatalf("note.DeckID = %q, want B1 German Idioms", note.DeckID)
	}
	if len(note.Cards) != 2 {
		t.Fatalf("len(note.Cards) = %d, want 2", len(note.Cards))
	}

	backCard := note.Cards[1]
	if backCard.Prompt != "All roads lead to Rome" {
		t.Fatalf("backCard.Prompt = %q, want All roads lead to Rome", backCard.Prompt)
	}
	if backCard.Answer != "Alle Wege führen nach Rom" {
		t.Fatalf("backCard.Answer = %q, want Alle Wege führen nach Rom", backCard.Answer)
	}
	if strings.HasPrefix(backCard.Answer, "Literal:") {
		t.Fatalf("literal explanation became answer: %+v", backCard)
	}
}

func TestImportHeaderedDeckUsesDeckMetadataOverRowDeckLabel(t *testing.T) {
	input := strings.Join([]string{
		"#separator:tab",
		"#deck:B2 Healthcare Systems",
		"front\tback\textra\ttags\tnotetype",
		"das Gesundheitssystem\thealthcare system\tdas Gesundheitssystem, die Gesundheitssysteme\tb2 healthcare systems\tB2 Healthcare\tBasic",
		"",
	}, "\n")

	notes, err := ImportAnkiTSV(strings.NewReader(input), ImportOptions{})
	if err != nil {
		t.Fatalf("import failed: %v", err)
	}
	if notes[0].DeckID != "B2 Healthcare Systems" {
		t.Fatalf("DeckID = %q, want B2 Healthcare Systems", notes[0].DeckID)
	}
}

func TestImportHeaderedMCQWithCommaChoices(t *testing.T) {
	input := strings.Join([]string{
		"#separator:tab",
		"#deck:Grammar",
		"front\tback\textra\ttags\tnotetype",
		"Ich gebe ___ das Buch.\tdem Mann\tDativ\tgrammar\tMCQ:der Mann,den Mann,dem Mann,des Mannes",
		"",
	}, "\n")

	notes, err := ImportAnkiTSV(strings.NewReader(input), ImportOptions{})
	if err != nil {
		t.Fatalf("import failed: %v", err)
	}
	if notes[0].Type != "MCQ" {
		t.Fatalf("Type = %q, want MCQ", notes[0].Type)
	}
	if len(notes[0].Choices) != 4 {
		t.Fatalf("choices = %v, want 4 choices", notes[0].Choices)
	}

	var foundMCQ bool
	for _, card := range notes[0].Cards {
		if card.Kind == core.CardKindMCQ {
			foundMCQ = true
			if len(card.Choices) != 4 {
				t.Fatalf("card choices = %v, want 4 choices", card.Choices)
			}
		}
	}
	if !foundMCQ {
		t.Fatal("no MCQ card generated")
	}
}

func TestExportReverseRoundTrip(t *testing.T) {
	note := core.Note{
		ID:     "rev-1",
		DeckID: "Default",
		Type:   "Reverse",
		Front:  "der Apfel",
		Back:   "apple",
	}
	out, err := ExportAnkiTSVString([]core.Note{note})
	if err != nil {
		t.Fatalf("export failed: %v", err)
	}
	if !strings.Contains(out, "Reverse") {
		t.Fatalf("exported TSV missing Reverse type: %q", out)
	}

	notes, err := ImportAnkiTSV(strings.NewReader(out), ImportOptions{})
	if err != nil {
		t.Fatalf("reimport failed: %v", err)
	}
	if len(notes) != 1 {
		t.Fatalf("len(notes) = %d, want 1", len(notes))
	}
	if notes[0].Type != "Reverse" {
		t.Fatalf("notes[0].Type = %q, want Reverse", notes[0].Type)
	}
	if len(notes[0].Cards) != 2 {
		t.Fatalf("len(notes[0].Cards) = %d, want 2", len(notes[0].Cards))
	}
}

func TestExampleCardGeneration(t *testing.T) {
	note := core.Note{
		ID:       "ex-1",
		Front:    "der Apfel",
		Back:     "apple",
		Examples: []string{"Ich esse einen Apfel."},
	}
	cards := CardsForNote(note)

	// Should have front card and example card
	if len(cards) != 2 {
		t.Fatalf("len(cards) = %d, want 2", len(cards))
	}

	exCard := cards[1]
	if exCard.ID != "ex-1:example" {
		t.Fatalf("id = %q, want ex-1:example", exCard.ID)
	}
	if !strings.Contains(exCard.Prompt, "___") {
		t.Fatalf("prompt missing blank: %q", exCard.Prompt)
	}
	if !strings.Contains(exCard.Prompt, "(apple)") {
		t.Fatalf("prompt missing hint: %q", exCard.Prompt)
	}
	if exCard.Answer != "Apfel" {
		t.Fatalf("answer = %q, want Apfel", exCard.Answer)
	}
}
