package content

import (
	"embed"
	"path/filepath"

	"deutsch-tui/internal/core"
)

var DefaultRegistry *Registry

func init() {
	DefaultRegistry = NewRegistry()
	DefaultRegistry.Register(NewEmbeddedSource(EmbeddedDecks, "testdata/german-decks", 10))
	DefaultRegistry.Register(&GoSource{priority: 20})
}

func AllDecks() ([]core.Deck, error) {
	if DefaultRegistry == nil {
		return nil, nil
	}
	return DefaultRegistry.AllDecks()
}

func DeckByID(id string) (*core.Deck, error) {
	if DefaultRegistry == nil {
		return nil, nil
	}
	return DefaultRegistry.DeckByID(id)
}

func DeckIDs() ([]string, error) {
	if DefaultRegistry == nil {
		return nil, nil
	}
	return DefaultRegistry.DeckIDs()
}

type GoSource struct {
	priority int
}

func (s *GoSource) Name() string {
	return "go"
}

func (s *GoSource) Priority() int {
	return s.priority
}

func (s *GoSource) LoadDecks() ([]core.Deck, error) {
	decks := []core.Deck{}

	starter := StarterDeck()
	if starter.ID != "" {
		decks = append(decks, starter)
	}

	return decks, nil
}

//go:embed testdata/german-decks/*.tsv
var EmbeddedDecks embed.FS

func EmbeddedDeckPaths() []string {
	return []string{
		filepath.Join("testdata/german-decks", "a1-essential.tsv"),
		filepath.Join("testdata/german-decks", "a1-food-drink.tsv"),
		filepath.Join("testdata/german-decks", "a1-health-body.tsv"),
		filepath.Join("testdata/german-decks", "a1-travel.tsv"),
		filepath.Join("testdata/german-decks", "a2-daily-life.tsv"),
		filepath.Join("testdata/german-decks", "a2-grammar-essentials.tsv"),
		filepath.Join("testdata/german-decks", "a2-shopping-services.tsv"),
		filepath.Join("testdata/german-decks", "a2-transport-directions.tsv"),
		filepath.Join("testdata/german-decks", "b1-business-professional.tsv"),
		filepath.Join("testdata/german-decks", "b1-apartment-housing.tsv"),
		filepath.Join("testdata/german-decks", "b1-emotions-feelings.tsv"),
		filepath.Join("testdata/german-decks", "b1-idioms.tsv"),
		filepath.Join("testdata/german-decks", "b1-false-friends.tsv"),
		filepath.Join("testdata/german-decks", "b1-phrasal-verbs.tsv"),
		filepath.Join("testdata/german-decks", "b1-workplace-office.tsv"),
		filepath.Join("testdata/german-decks", "b1-travel-tourism.tsv"),
		filepath.Join("testdata/german-decks", "b1-technology-internet.tsv"),
		filepath.Join("testdata/german-decks", "b2-advanced.tsv"),
	}
}

func AvailableDeckNames() []string {
	names := []string{}
	for _, path := range EmbeddedDeckPaths() {
		name := filepath.Base(path)
		name = name[:len(name)-4]
		names = append(names, name)
	}
	return names
}
