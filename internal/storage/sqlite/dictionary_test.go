package sqlite

import (
	"context"
	"testing"

	"deutsch-tui/internal/core"
)

func TestDictionarySearchFallback(t *testing.T) {
	ctx := context.Background()
	store, err := OpenMemory()
	if err != nil {
		t.Fatalf("open memory store: %v", err)
	}
	defer store.Close()

	entries := []core.DictionaryEntry{
		{
			ID:          "1",
			Word:        "Apfelkuchen",
			Translation: "Apple cake",
			WordClass:   "Noun",
			Gender:      "r",
		},
		{
			ID:          "2",
			Word:        "Fahrrad",
			Translation: "bicycle",
			WordClass:   "Noun",
			Gender:      "s",
		},
	}

	if err := store.ImportEntries(ctx, entries); err != nil {
		t.Fatalf("import entries: %v", err)
	}

	// 1. Test standard MATCH query
	// Searching "Apfel" should match "Apfelkuchen" (prefix match)
	res, err := store.Search(ctx, "Apfel", 10)
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}
	if len(res) != 1 || res[0].Word != "Apfelkuchen" {
		t.Errorf("expected Apfelkuchen for prefix MATCH, got: %v", res)
	}

	// 2. Test LIKE-based fallback query
	// Searching "kuch" should NOT match via FTS5 MATCH because it's not a token prefix.
	// It should fall back to LIKE '%kuch%' and match "Apfelkuchen".
	res, err = store.Search(ctx, "kuch", 10)
	if err != nil {
		t.Fatalf("Search fallback failed: %v", err)
	}
	if len(res) != 1 || res[0].Word != "Apfelkuchen" {
		t.Errorf("expected Apfelkuchen for fallback LIKE, got: %v", res)
	}

	// 3. Test searching translation via LIKE fallback
	res, err = store.Search(ctx, "cycl", 10)
	if err != nil {
		t.Fatalf("Search fallback translation failed: %v", err)
	}
	if len(res) != 1 || res[0].Word != "Fahrrad" {
		t.Errorf("expected Fahrrad for translation fallback LIKE, got: %v", res)
	}
}
