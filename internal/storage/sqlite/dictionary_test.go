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

func TestDictionarySearchExactMatchOrdering(t *testing.T) {
	ctx := context.Background()
	store, err := OpenMemory()
	if err != nil {
		t.Fatalf("open memory store: %v", err)
	}
	defer store.Close()

	entries := []core.DictionaryEntry{
		{ID: "1", Word: "Krankenhaus", Translation: "hospital"},
		{ID: "2", Word: "Haus", Translation: "house"},
		{ID: "3", Word: "Haustür", Translation: "front door"},
	}

	if err := store.ImportEntries(ctx, entries); err != nil {
		t.Fatalf("import entries: %v", err)
	}

	res, err := store.Search(ctx, "Haus", 10)
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}
	if len(res) != 2 {
		t.Fatalf("expected 2 results, got %d", len(res))
	}
	if res[0].Word != "Haus" {
		t.Errorf("expected exact match 'Haus' to be first, got '%s'", res[0].Word)
	}
	if res[1].Word != "Haustür" {
		t.Errorf("expected prefix match 'Haustür' to be second, got '%s'", res[1].Word)
	}
}

func TestDictionarySearchArticleAndMultiTranslation(t *testing.T) {
	ctx := context.Background()
	store, err := OpenMemory()
	if err != nil {
		t.Fatalf("open memory store: %v", err)
	}
	defer store.Close()

	entries := []core.DictionaryEntry{
		{ID: "1", Word: "das Haus", Translation: "house; home; building"},
		{ID: "2", Word: "die Katze", Translation: "cat"},
		{ID: "3", Word: "Häuser", Translation: "houses"},
	}

	if err := store.ImportEntries(ctx, entries); err != nil {
		t.Fatalf("import entries: %v", err)
	}

	// 1. Searching "Haus" (without article) should rank "das Haus" first (exact match ignoring article)
	res, err := store.Search(ctx, "Haus", 10)
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}
	if len(res) < 1 {
		t.Fatalf("expected results for 'Haus'")
	}
	if res[0].Word != "das Haus" {
		t.Errorf("expected 'das Haus' to match 'Haus' at rank 0, got '%s'", res[0].Word)
	}

	// 2. Searching "home" should rank "das Haus" first because "home" is an exact match in the semicolon list
	res, err = store.Search(ctx, "home", 10)
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}
	if len(res) < 1 {
		t.Fatalf("expected results for 'home'")
	}
	if res[0].Word != "das Haus" {
		t.Errorf("expected 'das Haus' to match 'home' at rank 0, got '%s'", res[0].Word)
	}
}
