package tui

import (
	"context"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"

	"deutsch-tui/internal/content"
)

// trainerConfigs is the registry of every text-input trainer. Adding a trainer
// means adding an entry here plus a loader that returns its trainerItems — no
// new state fields, key handler, or render function required.
var trainerConfigs = map[PracticeSubView]trainerConfig{
	PracticeSubViewConjugation: {
		Title:     "VERB CONJUGATION TRAINER",
		ItemNoun:  "verbs for conjugation practice",
		NextLabel: "next verb",
		EmptyMsg:  "No verbs found for practice.\nTry adding some verbs to your decks!",
		Load:      (*Model).loadConjugationItems,
	},
	PracticeSubViewCase: {
		Title:     "CASE ENDING TRAINER",
		ItemNoun:  "case exercises",
		NextLabel: "next exercise",
		EmptyMsg:  "No case exercises found.\nTry adding more grammar content!",
		HintKey:   true,
		Load:      (*Model).loadCaseItems,
	},
	PracticeSubViewAdjective: {
		Title:     "ADJECTIVE ENDING TRAINER",
		ItemNoun:  "adjective exercises",
		NextLabel: "next exercise",
		EmptyMsg:  "No adjective exercises found.\nTry adding more grammar content!",
		HintKey:   true,
		Load:      (*Model).loadAdjectiveItems,
	},
	PracticeSubViewPreposition: {
		Title:     "PREPOSITION TRAINER",
		ItemNoun:  "preposition exercises",
		NextLabel: "next exercise",
		EmptyMsg:  "No preposition exercises found.\nTry adding more grammar content!",
		HintKey:   true,
		Load:      (*Model).loadPrepositionItems,
	},
	PracticeSubViewPlural: {
		Title:     "NOUN PLURAL TRAINER",
		ItemNoun:  "plural exercises",
		NextLabel: "next noun",
		EmptyMsg:  "No nouns found for plural practice.\nTry adding more grammar content!",
		match:     pluralMatch,
		Load:      (*Model).loadPluralItems,
	},
	PracticeSubViewSeparable: {
		Title:     "SEPARABLE VERB TRAINER",
		ItemNoun:  "separable verb exercises",
		NextLabel: "next exercise",
		EmptyMsg:  "No separable verb exercises found.",
		Load:      (*Model).loadSeparableItems,
	},
	PracticeSubViewNumbers: {
		Title:      "NUMBER & TIME TRAINER",
		ItemNoun:   "number exercises",
		NextLabel:  "next exercise",
		EmptyMsg:   "No number exercises found.",
		InputWidth: 40,
		Load:       (*Model).loadNumberItems,
	},
	PracticeSubViewConjunctions: {
		Title:     "CONJUNCTIONS & WORD ORDER",
		ItemNoun:  "conjunction exercises",
		NextLabel: "next exercise",
		EmptyMsg:  "No conjunction exercises found.",
		HintKey:   true,
		Load:      (*Model).loadConjItems,
	},
}

// blanked renders a fill-in-the-blank sentence, replacing the {{...}} marker.
func blanked(sentence, with string) string {
	return strings.Replace(sentence, "{{...}}", with, 1)
}

// pluralMatch accepts the plural with or without a leading "die " article and
// tolerates umlaut spelling (ä/ae etc.).
func pluralMatch(input, target string) bool {
	u := strings.TrimSpace(strings.ToLower(input))
	t := strings.TrimSpace(strings.ToLower(target))
	tn := strings.TrimPrefix(t, "die ")
	return u == t || u == tn ||
		normalizeUmlauts(u) == normalizeUmlauts(t) ||
		normalizeUmlauts(u) == normalizeUmlauts(tn)
}

var conjugationPersons = []string{"ich", "du", "er/sie/es", "wir", "ihr", "sie/Sie"}

func conjugationForm(v content.DailyVerb, person int) string {
	switch person {
	case 1:
		return v.Du
	case 2:
		return v.ErSieEs
	case 3:
		return v.Wir
	case 4:
		return v.Ihr
	case 5:
		return v.SieSie
	default:
		return v.Ich
	}
}

func (m *Model) loadConjugationItems() tea.Cmd {
	return func() tea.Msg {
		verbs := content.AllDailyVerbs()
		items := make([]trainerItem, len(verbs))
		// Cycle through grammatical persons so a full pass covers all forms,
		// starting at "ich" to match the trainer's original first prompt.
		for i, v := range verbs {
			person := i % len(conjugationPersons)
			items[i] = trainerItem{
				Title:      v.German,
				Subtitle:   "(" + v.English + ")",
				PromptLine: "Conjugate for: " + conjugationPersons[person],
				Answer:     conjugationForm(v, person),
				Example:    v.Example,
			}
		}
		return trainerItemsMsg{kind: PracticeSubViewConjugation, items: items}
	}
}

func (m *Model) loadCaseItems() tea.Cmd {
	return func() tea.Msg {
		raw := []struct{ sentence, answer, context string }{
			{"Ich gehe mit {{...}} Hund.", "dem", "m, Dative (after mit)"},
			{"Ich sehe {{...}} Mann.", "den", "m, Accusative (direct object)"},
			{"Das ist das Buch {{...}} Frau.", "der", "f, Genitive (possession)"},
			{"Wir wohnen in {{...}} Stadt.", "der", "f, Dative (location)"},
			{"Er wartet auf {{...}} Bus.", "den", "m, Accusative (movement/direction)"},
			{"Sie gibt {{...}} Kind einen Apfel.", "dem", "n, Dative (indirect object)"},
			{"Ohne {{...}} Hilfe schaffe ich es nicht.", "deine", "f, Accusative (after ohne)"},
			{"Das ist das Haus {{...}} Mannes.", "des", "m, Genitive (possession)"},
			{"Ich komme aus {{...}} Schweiz.", "der", "f, Dative (after aus)"},
			{"Für {{...}} Mutter kaufe ich Blumen.", "meine", "f, Accusative (after für)"},
			{"Neben {{...}} Tisch steht ein Stuhl.", "dem", "m, Dative (location)"},
			{"Stell die Lampe auf {{...}} Tisch.", "den", "m, Accusative (movement)"},
			{"Wegen {{...}} Wetters bleiben wir zu Hause.", "des", "m, Genitive (after wegen)"},
			{"Ich danke {{...}} Lehrer.", "dem", "m, Dative (verb danken)"},
			{"Hilfst du {{...}} Bruder?", "deinem", "m, Dative (verb helfen)"},
		}
		return trainerItemsMsg{kind: PracticeSubViewCase, items: blankItems(raw)}
	}
}

func (m *Model) loadAdjectiveItems() tea.Cmd {
	return func() tea.Msg {
		raw := []struct{ sentence, answer, context string }{
			{"Ich trinke ein {{...}} (kalt) Bier.", "kaltes", "n, Accusative, mixed declension (ein-word)"},
			{"Der {{...}} (groß) Hund bellt.", "große", "m, Nominative, weak declension (der-word)"},
			{"Ich wohne in einem {{...}} (alt) Haus.", "alten", "n, Dative, weak declension"},
			{"Das ist eine {{...}} (schön) Blume.", "schöne", "f, Nominative, mixed declension"},
			{"Wir mögen {{...}} (deutsch) Wein.", "deutschen", "m, Accusative, strong declension (no article)"},
			{"Mit {{...}} (freundlich) Grüßen", "freundlichen", "plural, Dative, strong declension"},
			{"Sie kauft die {{...}} (rot) Schuhe.", "roten", "plural, Accusative, weak declension"},
			{"Ein {{...}} (gut) Freund hilft immer.", "guter", "m, Nominative, mixed declension"},
			{"Er trinkt gerne {{...}} (schwarz) Tee.", "schwarzen", "m, Accusative, strong declension"},
			{"Das ist das Auto des {{...}} (reich) Mannes.", "reichen", "m, Genitive, weak declension"},
			{"Eine {{...}} (jung) Frau sucht einen Job.", "junge", "f, Nominative, mixed declension"},
			{"Sie trägt ein {{...}} (blau) Kleid.", "blaues", "n, Accusative, mixed declension"},
			{"Wir gehen durch den {{...}} (dunkel) Wald.", "dunklen", "m, Accusative, weak declension"},
			{"Er arbeitet mit {{...}} (neu) Kollegen.", "neuen", "plural, Dative, strong declension"},
			{"Das ist ein {{...}} (schwierig) Rätsel.", "schwieriges", "n, Nominative, mixed declension"},
		}
		items := blankItems(raw)
		for i := range items {
			items[i].Instruction = "Fill in the blank (enter the correct adjective ending):"
		}
		return trainerItemsMsg{kind: PracticeSubViewAdjective, items: items}
	}
}

func (m *Model) loadPrepositionItems() tea.Cmd {
	return func() tea.Msg {
		exercises := content.GetPrepositionExercises()
		items := make([]trainerItem, len(exercises))
		for i, ex := range exercises {
			items[i] = trainerItem{
				Title:       blanked(ex.Sentence, "_____"),
				Answer:      ex.Answer,
				RevealTitle: blanked(ex.Sentence, ex.Answer),
				Instruction: "Fill in the blank (enter the preposition or article):",
				HintText:    ex.Context,
				Context:     ex.Context,
			}
		}
		return trainerItemsMsg{kind: PracticeSubViewPreposition, items: items}
	}
}

// blankItems builds fill-in-the-blank trainer items (case/adjective shape):
// blanked sentence as the prompt, filled sentence on reveal, grammar context
// used both as the toggleable hint and the post-reveal explanation block.
func blankItems(raw []struct{ sentence, answer, context string }) []trainerItem {
	items := make([]trainerItem, len(raw))
	for i, r := range raw {
		items[i] = trainerItem{
			Title:       blanked(r.sentence, "_____"),
			Answer:      r.answer,
			RevealTitle: blanked(r.sentence, r.answer),
			Instruction: "Fill in the blank:",
			HintText:    r.context,
			Context:     r.context,
		}
	}
	return items
}

func (m *Model) loadPluralItems() tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		cards, err := m.repo.Cards(ctx, "", "", "")
		if err != nil {
			return err
		}

		type pluralRow struct{ singular, plural, meaning string }
		var rows []pluralRow
		for _, card := range cards {
			pluralVal := extractPlural(card.Extra)
			if pluralVal == "" {
				continue
			}
			singular := card.Prompt
			meaning := card.Answer
			info := content.AnalyzeCard(card.Prompt, card.Answer)
			if info.Kind == content.KindNoun {
				singular = info.Display
				if strings.Contains(card.Answer, info.Base) {
					meaning = card.Prompt
				}
			}
			rows = append(rows, pluralRow{singular, pluralVal, meaning})
		}

		if len(rows) < 10 {
			rows = append(rows, []pluralRow{
				{"das Buch", "die Bücher", "book"},
				{"der Hund", "die Hunde", "dog"},
				{"das Kind", "die Kinder", "child"},
				{"das Auto", "die Autos", "car"},
				{"das Haus", "die Häuser", "house"},
				{"der Tisch", "die Tische", "table"},
				{"die Hand", "die Hände", "hand"},
				{"der Apfel", "die Äpfel", "apple"},
				{"das Ei", "die Eier", "egg"},
				{"die Stadt", "die Städte", "city"},
				{"das Bild", "die Bilder", "picture"},
				{"die Blume", "die Blumen", "flower"},
				{"der Freund", "die Freunde", "friend"},
				{"der Mann", "die Männer", "man"},
				{"die Frau", "die Frauen", "woman"},
			}...)
		}

		items := make([]trainerItem, len(rows))
		for i, r := range rows {
			subtitle := ""
			if r.meaning != "" {
				subtitle = "(" + r.meaning + ")"
			}
			items[i] = trainerItem{
				Title:       r.singular,
				Subtitle:    subtitle,
				Answer:      r.plural,
				Instruction: "Enter the plural form (with or without article):",
			}
		}
		return trainerItemsMsg{kind: PracticeSubViewPlural, items: items}
	}
}

func (m *Model) loadSeparableItems() tea.Cmd {
	return func() tea.Msg {
		raw := []struct{ sentence, verb, answer, meaning string }{
			{"Ich stehe um 7 Uhr {{...}}.", "aufstehen", "auf", "to get up"},
			{"Wann fängt der Film {{...}}?", "anfangen", "an", "to begin"},
			{"Komm bitte {{...}}!", "mitkommen", "mit", "to come along"},
			{"Ich rufe dich morgen {{...}}.", "anrufen", "an", "to call"},
			{"Wir kaufen im Supermarkt {{...}}.", "einkaufen", "ein", "to shop"},
			{"Der Zug kommt um 10 Uhr {{...}}.", "ankommen", "an", "to arrive"},
			{"Mach bitte das Licht {{...}}.", "ausmachen", "aus", "to turn off"},
			{"Er sieht jeden Abend {{...}}.", "fernsehen", "fern", "to watch TV"},
			{"Wir ziehen nächste Woche {{...}}.", "umziehen", "um", "to move"},
			{"Ich bereite das Essen {{...}}.", "vorbereiten", "vor", "to prepare"},
			{"Stell bitte die Milch {{...}}.", "wegstellen", "weg", "to put away"},
			{"Hör mir bitte {{...}}!", "zuhören", "zu", "to listen"},
			{"Wir laden euch zur Party {{...}}.", "einladen", "ein", "to invite"},
			{"Er gibt das Geld {{...}}.", "ausgeben", "aus", "to spend"},
			{"Ich schlage ein Treffen {{...}}.", "vorschlagen", "vor", "to suggest"},
		}
		items := make([]trainerItem, len(raw))
		for i, r := range raw {
			items[i] = trainerItem{
				Title:       r.sentence,
				Subtitle:    "Verb: " + r.verb + " (" + r.meaning + ")",
				Answer:      r.answer,
				Instruction: "Enter the missing prefix:",
			}
		}
		return trainerItemsMsg{kind: PracticeSubViewSeparable, items: items}
	}
}

func (m *Model) loadNumberItems() tea.Cmd {
	return func() tea.Msg {
		exercises := content.GetNumberExercises()
		items := make([]trainerItem, len(exercises))
		for i, ex := range exercises {
			items[i] = trainerItem{
				Title:       ex.Question,
				Subtitle:    ex.Help,
				Answer:      ex.Answer,
				Instruction: "Enter the German translation:",
			}
		}
		return trainerItemsMsg{kind: PracticeSubViewNumbers, items: items}
	}
}

func (m *Model) loadConjItems() tea.Cmd {
	return func() tea.Msg {
		exercises := content.GetConjunctionExercises()
		items := make([]trainerItem, len(exercises))
		for i, ex := range exercises {
			items[i] = trainerItem{
				Title:       ex.Sentence,
				Subtitle:    "Meaning: " + ex.Meaning,
				Answer:      ex.Answer,
				Instruction: "Enter the missing word:",
				HintText:    ex.Hint,
				Explanation: ex.Explanation,
			}
		}
		return trainerItemsMsg{kind: PracticeSubViewConjunctions, items: items}
	}
}
