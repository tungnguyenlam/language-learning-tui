package content

import (
	"embed"
	"path/filepath"
	"sort"
	"strings"

	"deutsch-tui/internal/core"
)

type ContentSource interface {
	Name() string
	LoadDecks() ([]core.Deck, error)
	Priority() int
}

type Registry struct {
	sources []ContentSource
}

func NewRegistry() *Registry {
	return &Registry{}
}

func (r *Registry) Register(source ContentSource) {
	r.sources = append(r.sources, source)
}

func (r *Registry) AllDecks() ([]core.Deck, error) {
	var allDecks []core.Deck
	seenDecks := make(map[string]bool)

	sort.Slice(r.sources, func(i, j int) bool {
		return r.sources[i].Priority() < r.sources[j].Priority()
	})

	for _, source := range r.sources {
		decks, err := source.LoadDecks()
		if err != nil {
			continue
		}
		for _, deck := range decks {
			if !seenDecks[deck.ID] {
				seenDecks[deck.ID] = true
				allDecks = append(allDecks, deck)
			}
		}
	}

	return allDecks, nil
}

func (r *Registry) DeckByID(id string) (*core.Deck, error) {
	decks, err := r.AllDecks()
	if err != nil {
		return nil, err
	}
	for _, deck := range decks {
		if deck.ID == id {
			return &deck, nil
		}
	}
	return nil, nil
}

func (r *Registry) DeckIDs() ([]string, error) {
	decks, err := r.AllDecks()
	if err != nil {
		return nil, err
	}
	ids := make([]string, len(decks))
	for i, deck := range decks {
		ids[i] = deck.ID
	}
	return ids, nil
}

type EmbeddedSource struct {
	fs       embed.FS
	dir      string
	priority int
}

func NewEmbeddedSource(fs embed.FS, dir string, priority int) *EmbeddedSource {
	return &EmbeddedSource{fs: fs, dir: dir, priority: priority}
}

func (s *EmbeddedSource) Name() string {
	return "embedded"
}

func (s *EmbeddedSource) Priority() int {
	return s.priority
}

func (s *EmbeddedSource) LoadDecks() ([]core.Deck, error) {
	entries, err := s.fs.ReadDir(s.dir)
	if err != nil {
		return nil, err
	}

	var decks []core.Deck
	deckMap := make(map[string]*core.Deck)

	for _, entry := range entries {
		if !strings.HasSuffix(entry.Name(), ".tsv") {
			continue
		}

		path := filepath.Join(s.dir, entry.Name())
		data, err := s.fs.ReadFile(path)
		if err != nil {
			continue
		}

		notes, err := ImportAnkiTSV(strings.NewReader(string(data)), ImportOptions{})
		if err != nil {
			continue
		}

		for i := range notes {
			notes[i].Cards = CardsForNote(notes[i])
		}

		deckName := strings.TrimSuffix(entry.Name(), ".tsv")
		deckID := ToDeckID(deckName)

		if deck, ok := deckMap[deckID]; ok {
			deck.Notes = append(deck.Notes, notes...)
		} else {
			deckMap[deckID] = &core.Deck{
				ID:    deckID,
				Name:  deckName,
				Notes: notes,
			}
		}
	}

	for _, deck := range deckMap {
		decks = append(decks, *deck)
	}

	return decks, nil
}

func ToDeckID(name string) string {
	id := strings.ReplaceAll(name, "-", "_")
	id = strings.ReplaceAll(id, " ", "_")
	return strings.ToLower(id)
}
