package content

import (
	"strings"

	"deutsch-tui/internal/core"
)

func stripArticles(s string) string {
	s = strings.TrimSpace(s)
	lower := strings.ToLower(s)
	if strings.HasPrefix(lower, "der ") {
		return s[4:]
	}
	if strings.HasPrefix(lower, "die ") {
		return s[4:]
	}
	if strings.HasPrefix(lower, "das ") {
		return s[4:]
	}
	return s
}

func findWordInSentence(sentence, word string) (int, int) {
	w := strings.ToLower(word)
	s := strings.ToLower(sentence)

	idx := strings.Index(s, w)
	if idx == -1 {
		return -1, -1
	}

	// Because ToLower can change byte length (e.g. ẞ -> ß),
	// we need to find the corresponding byte range in the original sentence.
	// We can do this by iterating through the original string and comparing ToLowered versions.
	for i := 0; i <= len(sentence); i++ {
		// We use ToLower on the prefix and check its length.
		// This is still a bit complex. A simpler way is to check each position.
		if strings.HasPrefix(strings.ToLower(sentence[i:]), w) {
			// Find how many bytes in the original sentence correspond to 'w' in lowered sentence.
			// We can do this by increasing length until ToLower(sentence[i:i+l]) starts with w and is at least as long as needed.
			for l := 1; i+l <= len(sentence); l++ {
				if strings.ToLower(sentence[i:i+l]) == w {
					return i, i + l
				}
				if len(strings.ToLower(sentence[i:i+l])) > len(w) {
					break
				}
			}
		}
	}

	return -1, -1
}

func DecksFromNotes(notes []core.Note) []core.Deck {
	byID := make(map[string]int)
	decks := make([]core.Deck, 0, 1)
	for _, note := range notes {
		deckID := strings.TrimSpace(note.DeckID)
		if deckID == "" {
			deckID = "Imported"
			note.DeckID = deckID
		}
		index, ok := byID[deckID]
		if !ok {
			index = len(decks)
			byID[deckID] = index
			decks = append(decks, core.Deck{
				ID:          deckID,
				Name:        deckID,
				Description: "Imported from Anki TSV.",
			})
		}
		decks[index].Notes = append(decks[index].Notes, note)
	}
	return decks
}
