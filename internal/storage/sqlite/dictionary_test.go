package sqlite

import (
	"context"
	"fmt"
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
		{ID: "6", Word: "undurchlässig", Translation: "impermeable"},
	}

	if err := store.ImportEntries(ctx, entries); err != nil {
		t.Fatalf("import entries: %v", err)
	}

	res, err := store.FindRelatedEntries(ctx, "das Haus", 10)
	if err != nil {
		t.Fatalf("FindRelatedEntries failed: %v", err)
	}

	expected := map[string]bool{
		"Haus":        true,
		"Haustier":    true,
		"Krankenhaus": true,
	}

	if len(res) != 3 {
		t.Errorf("expected 3 related entries, got %d (%v)", len(res), res)
	}

	for _, e := range res {
		if !expected[e.Word] {
			t.Errorf("unexpected related entry: %s", e.Word)
		}
	}

	// Middle/prefix noise: short stem "und" must not surface undurchlässig.
	noisy, err := store.FindRelatedEntries(ctx, "und", 10)
	if err != nil {
		t.Fatalf("FindRelatedEntries(und) failed: %v", err)
	}
	for _, e := range noisy {
		if e.Word == "undurchlässig" {
			t.Fatalf("middle-substring match should be excluded, got %q", e.Word)
		}
	}
}

func TestDictionaryFilterOnlyBrowseOrdered(t *testing.T) {
	ctx := context.Background()
	store, err := OpenMemory()
	if err != nil {
		t.Fatalf("open memory store: %v", err)
	}
	defer store.Close()

	entries := []core.DictionaryEntry{
		{ID: "3", Word: "Zebra", Translation: "zebra", WordClass: "noun", Gender: "n"},
		{ID: "1", Word: "Apfel", Translation: "apple", WordClass: "noun", Gender: "m"},
		{ID: "2", Word: "Birne", Translation: "pear", WordClass: "noun", Gender: "f"},
		{ID: "4", Word: "gehen", Translation: "to go", WordClass: "verb"},
	}
	if err := store.ImportEntries(ctx, entries); err != nil {
		t.Fatalf("import entries: %v", err)
	}

	res, err := store.Search(ctx, ":noun", 10)
	if err != nil {
		t.Fatalf("filter-only search failed: %v", err)
	}
	if len(res) != 3 {
		t.Fatalf("expected 3 nouns, got %d", len(res))
	}
	if res[0].Word != "Apfel" || res[1].Word != "Birne" || res[2].Word != "Zebra" {
		t.Fatalf("expected alphabetical noun browse, got %v %v %v", res[0].Word, res[1].Word, res[2].Word)
	}
}

func TestDictionaryFilterOnlyBrowseBeyondAlphabeticalSlice(t *testing.T) {
	ctx := context.Background()
	store, err := OpenMemory()
	if err != nil {
		t.Fatalf("open memory store: %v", err)
	}
	defer store.Close()

	// 220 nouns sorting before "ziehen" used to fill the old LIMIT 200 window
	// and hide the only verb when browsing :verb.
	entries := make([]core.DictionaryEntry, 0, 221)
	for i := 0; i < 220; i++ {
		entries = append(entries, core.DictionaryEntry{
			ID:          fmt.Sprintf("n-%03d", i),
			Word:        fmt.Sprintf("Alpha%03d", i),
			Translation: "noun",
			WordClass:   "noun",
			Gender:      "n",
		})
	}
	entries = append(entries, core.DictionaryEntry{
		ID:          "v-1",
		Word:        "ziehen",
		Translation: "to pull",
		WordClass:   "verb",
	})
	if err := store.ImportEntries(ctx, entries); err != nil {
		t.Fatalf("import entries: %v", err)
	}

	res, err := store.Search(ctx, ":verb", 10)
	if err != nil {
		t.Fatalf("filter-only :verb search failed: %v", err)
	}
	if len(res) != 1 || res[0].Word != "ziehen" {
		t.Fatalf("expected ziehen from full-dictionary :verb browse, got %#v", res)
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

func TestDictionarySearchFTSFailureFallback(t *testing.T) {
	ctx := context.Background()
	store, err := OpenMemory()
	if err != nil {
		t.Fatalf("open memory store: %v", err)
	}
	defer store.Close()

	entries := []core.DictionaryEntry{
		{ID: "1", Word: "Handtuch", Translation: "towel", WordClass: "noun", Gender: "n"},
	}
	if err := store.ImportEntries(ctx, entries); err != nil {
		t.Fatalf("import entries: %v", err)
	}

	// Search with special characters that could break raw FTS query string
	res, err := store.Search(ctx, "Handtuch*", 10)
	if err != nil {
		t.Fatalf("search with special wildcards failed: %v", err)
	}
	if len(res) == 0 || res[0].Word != "Handtuch" {
		t.Fatalf("expected Handtuch result, got %v", res)
	}
}

func TestDictionarySearchNormalizesNonPositiveLimit(t *testing.T) {
	ctx := context.Background()
	store, err := OpenMemory()
	if err != nil {
		t.Fatalf("open memory store: %v", err)
	}
	defer store.Close()

	if err := store.ImportEntries(ctx, []core.DictionaryEntry{
		{ID: "1", Word: "Haus", Translation: "house", Forms: "Häuser"},
	}); err != nil {
		t.Fatalf("import entries: %v", err)
	}

	for _, limit := range []int{0, -1} {
		res, err := store.Search(ctx, "Häuser", limit)
		if err != nil {
			t.Fatalf("search with limit %d failed: %v", limit, err)
		}
		if len(res) != 1 || res[0].Word != "Haus" {
			t.Fatalf("search with limit %d = %v, want Haus", limit, res)
		}
	}
}

func TestDictionarySearchPunctuationOnlyNoFTSSyntaxError(t *testing.T) {
	ctx := context.Background()
	store, err := OpenMemory()
	if err != nil {
		t.Fatalf("open memory store: %v", err)
	}
	defer store.Close()

	if err := store.ImportEntries(ctx, []core.DictionaryEntry{
		{ID: "1", Word: "Haus", Translation: "house"},
	}); err != nil {
		t.Fatalf("import entries: %v", err)
	}

	for _, query := range []string{`"`, `:`, `""`, `:::`, `@#$%`} {
		res, err := store.Search(ctx, query, 10)
		if err != nil {
			t.Fatalf("search for punctuation query %q failed: %v", query, err)
		}
		if len(res) != 0 {
			t.Errorf("expected 0 results for punctuation query %q, got %d", query, len(res))
		}
	}
}
