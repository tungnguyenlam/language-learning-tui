package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"deutsch-tui/internal/core"
)

func (s *Store) queryDictionaryEntries(ctx context.Context, q string, args ...any) ([]core.DictionaryEntry, error) {
	rows, err := s.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []core.DictionaryEntry
	for rows.Next() {
		var entry core.DictionaryEntry
		var examplesStr, tagsStr string

		err := rows.Scan(
			&entry.ID,
			&entry.Word,
			&entry.Translation,
			&entry.WordClass,
			&entry.Gender,
			&entry.Forms,
			&examplesStr,
			&tagsStr,
		)
		if err != nil {
			return nil, fmt.Errorf("scan dictionary entry: %w", err)
		}

		if examplesStr != "" {
			if err := json.Unmarshal([]byte(examplesStr), &entry.Examples); err != nil {
				entry.Examples = []string{examplesStr}
			}
		}
		if tagsStr != "" {
			if err := json.Unmarshal([]byte(tagsStr), &entry.Tags); err != nil {
				entry.Tags = []string{tagsStr}
			}
		}

		entries = append(entries, entry)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}
	return entries, nil
}

func parseSearchFilters(query string) (cleanQuery string, classFilter string, genderFilter string, langFilter string) {
	terms := strings.Fields(query)
	var textTerms []string

	for _, term := range terms {
		lower := strings.ToLower(term)
		switch {
		case strings.HasPrefix(lower, "de:"):
			langFilter = "de"
			if len(term) > 3 {
				textTerms = append(textTerms, term[3:])
			}
		case strings.HasPrefix(lower, "en:"):
			langFilter = "en"
			if len(term) > 3 {
				textTerms = append(textTerms, term[3:])
			}
		case lower == "lang:de" || lower == ":de":
			langFilter = "de"
		case lower == "lang:en" || lower == ":en":
			langFilter = "en"
		case lower == ":verb" || lower == ":v" || lower == "class:verb":
			classFilter = "verb"
		case lower == ":noun" || lower == "class:noun":
			classFilter = "noun"
		case lower == ":adj" || lower == ":adjective" || lower == "class:adj":
			classFilter = "adj"
		case lower == ":adv" || lower == ":adverb" || lower == "class:adv":
			classFilter = "adv"
		case lower == ":m" || lower == ":masc" || lower == "gender:m":
			genderFilter = "m"
		case lower == ":f" || lower == ":fem" || lower == "gender:f":
			genderFilter = "f"
		case lower == ":n" || lower == ":neut" || lower == "gender:n":
			genderFilter = "n"
		case lower == ":pl" || lower == ":plural" || lower == "gender:pl":
			genderFilter = "pl"
		default:
			textTerms = append(textTerms, term)
		}
	}
	return strings.Join(textTerms, " "), classFilter, genderFilter, langFilter
}

func filterEntries(entries []core.DictionaryEntry, classFilter, genderFilter, langFilter string) []core.DictionaryEntry {
	if classFilter == "" && genderFilter == "" && langFilter == "" {
		return entries
	}
	filtered := make([]core.DictionaryEntry, 0, len(entries))
	for _, entry := range entries {
		if classFilter != "" {
			wc := strings.ToLower(entry.WordClass)
			if !strings.Contains(wc, classFilter) {
				continue
			}
		}
		if genderFilter != "" {
			g := strings.ToLower(entry.Gender)
			if !strings.HasPrefix(g, genderFilter) {
				continue
			}
		}
		filtered = append(filtered, entry)
	}
	return filtered
}

func (s *Store) Search(ctx context.Context, rawQuery string, limit int) ([]core.DictionaryEntry, error) {
	cleanQuery, classFilter, genderFilter, langFilter := parseSearchFilters(rawQuery)
	if cleanQuery == "" && (classFilter != "" || genderFilter != "" || langFilter != "") {
		q := `
			SELECT id, word, translation, word_class, gender, forms, examples, tags
			FROM dictionary_fts
			LIMIT 200
		`
		entries, err := s.queryDictionaryEntries(ctx, q)
		if err != nil {
			return nil, fmt.Errorf("search dictionary filter-only: %w", err)
		}
		filtered := filterEntries(entries, classFilter, genderFilter, langFilter)
		if len(filtered) > limit {
			filtered = filtered[:limit]
		}
		return filtered, nil
	}

	terms := strings.Fields(cleanQuery)
	if len(terms) == 0 {
		return nil, nil
	}

	// Direct LIKE search for wildcard queries (* or ?)
	if strings.ContainsAny(cleanQuery, "*?") {
		likePattern := strings.ReplaceAll(cleanQuery, "*", "%")
		likePattern = strings.ReplaceAll(likePattern, "?", "_")
		var qWild string
		var args []any
		if langFilter == "de" {
			qWild = `SELECT id, word, translation, word_class, gender, forms, examples, tags FROM dictionary_fts WHERE word LIKE ? OR forms LIKE ? LIMIT ?`
			args = []any{likePattern, likePattern, limit * 2}
		} else if langFilter == "en" {
			qWild = `SELECT id, word, translation, word_class, gender, forms, examples, tags FROM dictionary_fts WHERE translation LIKE ? LIMIT ?`
			args = []any{likePattern, limit * 2}
		} else {
			qWild = `SELECT id, word, translation, word_class, gender, forms, examples, tags FROM dictionary_fts WHERE word LIKE ? OR translation LIKE ? OR forms LIKE ? LIMIT ?`
			args = []any{likePattern, likePattern, likePattern, limit * 2}
		}
		entries, err := s.queryDictionaryEntries(ctx, qWild, args...)
		if err != nil {
			return nil, fmt.Errorf("search dictionary wildcard: %w", err)
		}
		filtered := filterEntries(entries, classFilter, genderFilter, langFilter)
		sortDictionaryEntries(filtered, cleanQuery)
		if len(filtered) > limit {
			filtered = filtered[:limit]
		}
		return filtered, nil
	}

	var matchQuery strings.Builder
	for i, term := range terms {
		if i > 0 {
			matchQuery.WriteString(" ")
		}
		safeTerm := strings.ReplaceAll(term, `"`, `""`)
		matchQuery.WriteString(fmt.Sprintf(`"%s"*`, safeTerm))
	}

	q := `
		SELECT id, word, translation, word_class, gender, forms, examples, tags
		FROM dictionary_fts
		WHERE dictionary_fts MATCH ?
		ORDER BY rank
		LIMIT ?
	`

	entries, err := s.queryDictionaryEntries(ctx, q, matchQuery.String(), limit*2)
	if err != nil {
		return nil, fmt.Errorf("search dictionary fts: %w", err)
	}

	// Also check forms via LIKE to catch unindexed inflected form matches
	likePattern := "%" + cleanQuery + "%"
	qForms := `
		SELECT id, word, translation, word_class, gender, forms, examples, tags
		FROM dictionary_fts
		WHERE forms LIKE ?
		LIMIT ?
	`
	formEntries, _ := s.queryDictionaryEntries(ctx, qForms, likePattern, limit)
	if len(formEntries) > 0 {
		seen := make(map[string]bool, len(entries))
		for _, e := range entries {
			seen[e.ID] = true
		}
		for _, e := range formEntries {
			if !seen[e.ID] {
				entries = append(entries, e)
				seen[e.ID] = true
			}
		}
	}

	if len(entries) > 0 {
		filtered := filterEntries(entries, classFilter, genderFilter, langFilter)
		if len(filtered) > 0 {
			sortDictionaryEntries(filtered, cleanQuery)
			if len(filtered) > limit {
				filtered = filtered[:limit]
			}
			return filtered, nil
		}
	}

	// Fallback to LIKE-based substring match
	var qLike string
	var args []any
	if langFilter == "de" {
		qLike = `SELECT id, word, translation, word_class, gender, forms, examples, tags FROM dictionary_fts WHERE word LIKE ? OR forms LIKE ? LIMIT ?`
		args = []any{likePattern, likePattern, limit * 2}
	} else if langFilter == "en" {
		qLike = `SELECT id, word, translation, word_class, gender, forms, examples, tags FROM dictionary_fts WHERE translation LIKE ? LIMIT ?`
		args = []any{likePattern, limit * 2}
	} else {
		qLike = `SELECT id, word, translation, word_class, gender, forms, examples, tags FROM dictionary_fts WHERE word LIKE ? OR translation LIKE ? OR forms LIKE ? LIMIT ?`
		args = []any{likePattern, likePattern, likePattern, limit * 2}
	}
	entries, err = s.queryDictionaryEntries(ctx, qLike, args...)
	if err != nil {
		return nil, fmt.Errorf("search dictionary fallback like: %w", err)
	}

	filtered := filterEntries(entries, classFilter, genderFilter, langFilter)
	sortDictionaryEntries(filtered, cleanQuery)
	if len(filtered) > limit {
		filtered = filtered[:limit]
	}
	return filtered, nil
}

func sortDictionaryEntries(entries []core.DictionaryEntry, query string) {
	q := strings.ToLower(strings.TrimSpace(query))
	if q == "" {
		return
	}
	sort.SliceStable(entries, func(i, j int) bool {
		scoreI := entryMatchScore(entries[i], q)
		scoreJ := entryMatchScore(entries[j], q)
		return scoreI < scoreJ
	})
}

func stripArticle(word string) string {
	w := strings.TrimSpace(strings.ToLower(word))
	for _, art := range []string{"der ", "die ", "das ", "den ", "dem ", "des ", "the ", "a ", "an "} {
		if strings.HasPrefix(w, art) {
			return strings.TrimSpace(w[len(art):])
		}
	}
	return w
}

func parseForms(formsStr string) []string {
	if strings.TrimSpace(formsStr) == "" {
		return nil
	}
	parts := strings.FieldsFunc(formsStr, func(r rune) bool {
		return r == ';' || r == ','
	})
	var res []string
	for _, p := range parts {
		trimmed := stripArticle(strings.TrimSpace(p))
		if trimmed != "" {
			res = append(res, trimmed)
		}
	}
	return res
}

func entryMatchScore(entry core.DictionaryEntry, q string) int {
	wRaw := strings.ToLower(strings.TrimSpace(entry.Word))
	tRaw := strings.ToLower(strings.TrimSpace(entry.Translation))
	qClean := strings.ToLower(strings.TrimSpace(q))
	wClean := stripArticle(wRaw)
	qBare := stripArticle(qClean)

	// 0: Exact word match (with or without German/English article)
	if wRaw == qClean || wClean == qBare || wClean == qClean {
		return 0
	}

	// 1: Exact translation match (or exact match on semicolon-separated translation)
	if tRaw == qClean || tRaw == qBare {
		return 1
	}
	for _, part := range strings.Split(tRaw, ";") {
		p := strings.TrimSpace(part)
		if p == qClean || p == qBare {
			return 1
		}
	}

	// 2: Exact word form match (e.g. "ging" or "gegangen" matching entry "gehen")
	forms := parseForms(entry.Forms)
	for _, f := range forms {
		if f == qClean || f == qBare {
			return 2
		}
	}

	// 3: Word prefix match
	if strings.HasPrefix(wRaw, qClean) || strings.HasPrefix(wClean, qBare) {
		return 3
	}

	// 4: Translation prefix match
	if strings.HasPrefix(tRaw, qClean) || strings.HasPrefix(tRaw, qBare) {
		return 4
	}

	// 5: Form prefix match
	for _, f := range forms {
		if strings.HasPrefix(f, qClean) || strings.HasPrefix(f, qBare) {
			return 5
		}
	}

	// 6: Word contains query substring
	if strings.Contains(wRaw, qClean) || strings.Contains(wClean, qBare) {
		return 6
	}

	// 7: Translation contains query substring
	if strings.Contains(tRaw, qClean) || strings.Contains(tRaw, qBare) {
		return 7
	}

	// 8: Form contains query substring
	for _, f := range forms {
		if strings.Contains(f, qClean) || strings.Contains(f, qBare) {
			return 8
		}
	}

	return 9
}

func (s *Store) GetEntry(ctx context.Context, id string) (core.DictionaryEntry, error) {
	q := `
		SELECT id, word, translation, word_class, gender, forms, examples, tags
		FROM dictionary_fts
		WHERE id = ?
		LIMIT 1
	`

	var entry core.DictionaryEntry
	var examplesStr, tagsStr string

	err := s.db.QueryRowContext(ctx, q, id).Scan(
		&entry.ID,
		&entry.Word,
		&entry.Translation,
		&entry.WordClass,
		&entry.Gender,
		&entry.Forms,
		&examplesStr,
		&tagsStr,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return core.DictionaryEntry{}, fmt.Errorf("entry not found: %w", err)
		}
		return core.DictionaryEntry{}, fmt.Errorf("get dictionary entry: %w", err)
	}

	if examplesStr != "" {
		if err := json.Unmarshal([]byte(examplesStr), &entry.Examples); err != nil {
			entry.Examples = []string{examplesStr}
		}
	}
	if tagsStr != "" {
		if err := json.Unmarshal([]byte(tagsStr), &entry.Tags); err != nil {
			entry.Tags = []string{tagsStr}
		}
	}

	return entry, nil
}

func (s *Store) ImportEntries(ctx context.Context, entries []core.DictionaryEntry) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, "DELETE FROM dictionary_fts"); err != nil {
		return fmt.Errorf("clear dictionary: %w", err)
	}

	q := `
		INSERT INTO dictionary_fts (id, word, translation, word_class, gender, forms, examples, tags)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
	stmt, err := tx.PrepareContext(ctx, q)
	if err != nil {
		return fmt.Errorf("prepare insert: %w", err)
	}
	defer stmt.Close()

	for _, entry := range entries {
		var examplesStr, tagsStr string
		if len(entry.Examples) > 0 {
			b, _ := json.Marshal(entry.Examples)
			examplesStr = string(b)
		}
		if len(entry.Tags) > 0 {
			b, _ := json.Marshal(entry.Tags)
			tagsStr = string(b)
		}

		if _, err := stmt.ExecContext(ctx,
			entry.ID,
			entry.Word,
			entry.Translation,
			entry.WordClass,
			entry.Gender,
			entry.Forms,
			examplesStr,
			tagsStr,
		); err != nil {
			return fmt.Errorf("insert entry %s: %w", entry.ID, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}
	return nil
}

func (s *Store) DictionaryCount(ctx context.Context) (int, error) {
	var count int
	err := s.db.QueryRowContext(ctx, "SELECT count(*) FROM dictionary_fts").Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("dictionary count: %w", err)
	}
	return count, nil
}
