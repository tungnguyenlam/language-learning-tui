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

func TestDictionarySearchInflectedFormsAndFilters(t *testing.T) {
	ctx := context.Background()
	store, err := OpenMemory()
	if err != nil {
		t.Fatalf("open memory store: %v", err)
	}
	defer store.Close()

	entries := []core.DictionaryEntry{
		{ID: "1", Word: "gehen", Translation: "to go", WordClass: "verb", Forms: "ging; gegangen; geht"},
		{ID: "2", Word: "Gingivitis", Translation: "gum infection", WordClass: "noun", Gender: "f"},
		{ID: "3", Word: "Haus", Translation: "house", WordClass: "noun", Gender: "n", Forms: "Häuser; Häusern"},
	}

	if err := store.ImportEntries(ctx, entries); err != nil {
		t.Fatalf("import entries: %v", err)
	}

	// 1. Inflected form exact match: searching "ging" should rank "gehen" first (exact form match) above "Gingivitis" (word prefix)
	res, err := store.Search(ctx, "ging", 10)
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}
	if len(res) < 2 {
		t.Fatalf("expected at least 2 results for 'ging', got %d", len(res))
	}
	if res[0].Word != "gehen" {
		t.Errorf("expected 'gehen' to be ranked first for inflected form 'ging', got '%s'", res[0].Word)
	}

	// 2. Inflected form plural: searching "Häuser" should match "Haus"
	res, err = store.Search(ctx, "Häuser", 10)
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}
	if len(res) < 1 || res[0].Word != "Haus" {
		t.Errorf("expected 'Haus' for inflected form 'Häuser', got %v", res)
	}

	// 3. Filter query: searching "ging :verb" should filter out Gingivitis (noun)
	res, err = store.Search(ctx, "ging :verb", 10)
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}
	if len(res) != 1 || res[0].Word != "gehen" {
		t.Errorf("expected only 'gehen' for 'ging :verb', got %v", res)
	}

	// 4. Gender filter query: searching ":n" should return neuter nouns
	res, err = store.Search(ctx, ":n", 10)
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}
	if len(res) != 1 || res[0].Word != "Haus" {
		t.Errorf("expected 'Haus' for filter ':n', got %v", res)
	}
}

func TestDictionarySearchWildcardAndLanguagePrefix(t *testing.T) {
	ctx := context.Background()
	store, err := OpenMemory()
	if err != nil {
		t.Fatalf("open memory store: %v", err)
	}
	defer store.Close()

	entries := []core.DictionaryEntry{
		{ID: "1", Word: "Hoffnung", Translation: "hope", WordClass: "noun", Gender: "f"},
		{ID: "2", Word: "Zeitung", Translation: "newspaper", WordClass: "noun", Gender: "f"},
		{ID: "3", Word: "Spaziergang", Translation: "walk; stroll", WordClass: "noun", Gender: "m"},
		{ID: "4", Word: "gehen", Translation: "to walk; to go", WordClass: "verb"},
	}

	if err := store.ImportEntries(ctx, entries); err != nil {
		t.Fatalf("import entries: %v", err)
	}

	// 1. Wildcard Suffix search: searching "*ung" should match Hoffnung and Zeitung
	res, err := store.Search(ctx, "*ung", 10)
	if err != nil {
		t.Fatalf("Wildcard search failed: %v", err)
	}
	if len(res) != 2 {
		t.Fatalf("expected 2 results for '*ung', got %d", len(res))
	}

	// 2. Language scope search de: searching "de:walk" should not match German headword walk, but searching "en:walk" should match translation "walk"
	resEn, err := store.Search(ctx, "en:walk", 10)
	if err != nil {
		t.Fatalf("Language scope search failed: %v", err)
	}
	if len(resEn) != 2 {
		t.Errorf("expected 2 results for 'en:walk', got %d", len(resEn))
	}
}

func TestFindRelatedEntries(t *testing.T) {
	ctx := context.Background()
	store, err := OpenMemory()
	if err != nil {
		t.Fatalf("open memory store: %v", err)
	}
	defer store.Close()

	entries := []core.DictionaryEntry{
		{ID: "1", Word: "Haus", Translation: "house"},
		{ID: "2", Word: "das Haus", Translation: "house (with article)"},
		{ID: "3", Word: "Haustier", Translation: "pet"},
		{ID: "4", Word: "Krankenhaus", Translation: "hospital"},
		{ID: "5", Word: "Auto", Translation: "car"},
	}

	if err := store.ImportEntries(ctx, entries); err != nil {
		t.Fatalf("import entries: %v", err)
	}

	res, err := store.FindRelatedEntries(ctx, "das Haus", 10)
	if err != nil {
		t.Fatalf("FindRelatedEntries failed: %v", err)
	}

	// Should not include "das Haus" itself.
	// But it might include "Haus" if "das Haus" is treated as the original word.
	// Actually, the implementation excludes the exact `word` parameter.
	if len(res) != 3 {
		t.Errorf("expected 3 related entries, got %d", len(res))
	}

	expected := map[string]bool{
		"Haus":        true,
		"Haustier":    true,
		"Krankenhaus": true,
	}

	for _, e := range res {
		if !expected[e.Word] {
			t.Errorf("unexpected related entry: %s", e.Word)
		}
	}
}

func TestDictionarySearchLangFilterOnly(t *testing.T) {
	ctx := context.Background()
	store, err := OpenMemory()
	if err != nil {
		t.Fatalf("open memory store: %v", err)
	}
	defer store.Close()

	entries := []core.DictionaryEntry{
		{ID: "1", Word: "Hund", Translation: "dog", WordClass: "noun", Gender: "m"},
		{ID: "2", Word: "gehen", Translation: "to go", WordClass: "verb"},
		{ID: "3", Word: "schnell", Translation: "fast", WordClass: "adj"},
	}
	if err := store.ImportEntries(ctx, entries); err != nil {
		t.Fatalf("import entries: %v", err)
	}

	for _, query := range []string{"de:", "en:", ":de", ":en"} {
		res, err := store.Search(ctx, query, 10)
		if err != nil {
			t.Fatalf("lang-only search %q failed: %v", query, err)
		}
		if len(res) == 0 {
			t.Fatalf("expected browse results for lang-only query %q, got none", query)
		}
	}
}
