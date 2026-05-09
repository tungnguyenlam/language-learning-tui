package content

import (
	"embed"
	"path/filepath"
	"strings"

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
	return StandardDecks(), nil
}

//go:embed testdata/german-decks/*.tsv
var EmbeddedDecks embed.FS

func EmbeddedDeckPaths() []string {
	entries, err := EmbeddedDecks.ReadDir("testdata/german-decks")
	if err != nil {
		return nil
	}
	var paths []string
	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), ".tsv") {
			paths = append(paths, filepath.Join("testdata/german-decks", entry.Name()))
		}
	}
	return paths
}

func AvailableDeckNames() []string {
	paths := EmbeddedDeckPaths()
	names := make([]string, 0, len(paths))
	for _, path := range paths {
		name := filepath.Base(path)
		name = name[:len(name)-4]
		names = append(names, name)
	}
	return names
}
