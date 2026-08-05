package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"

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
		case lower == ":starred" || lower == ":star" || lower == ":fav" || lower == ":favorite":
			// Handled by TUI filter layer
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

func filterEntries(entries []core.DictionaryEntry, classFilter, genderFilter, langFilter, query string) []core.DictionaryEntry {
	if classFilter == "" && genderFilter == "" && langFilter == "" {
		return entries
	}
	q := strings.ToLower(strings.TrimSpace(query))
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
		if langFilter != "" && q != "" {
			w := strings.ToLower(entry.Word)
			t := strings.ToLower(entry.Translation)
			f := strings.ToLower(entry.Forms)
			switch langFilter {
			case "de":
				if !strings.Contains(w, q) && !strings.Contains(f, q) &&
					!strings.Contains(stripArticle(w), stripArticle(q)) {
					// Keep exact-ish form matches even when the bare query is only a prefix of a form token.
					matchedForm := false
					for _, form := range parseForms(entry.Forms) {
						if strings.Contains(form, q) || strings.HasPrefix(form, q) {
							matchedForm = true
							break
						}
					}
					if !matchedForm {
						continue
					}
				}
			case "en":
				if !strings.Contains(t, q) {
					continue
				}
			}
		}
		filtered = append(filtered, entry)
	}
	return filtered
}

func hasWordChars(s string) bool {
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			return true
		}
	}
	return false
}

func buildFTSMatchQuery(terms []string, langFilter string) string {
	var matchQuery strings.Builder
	written := 0
	for _, term := range terms {
		if !hasWordChars(term) {
			continue
		}
		if written > 0 {
			matchQuery.WriteString(" ")
		}
		safeTerm := strings.ReplaceAll(term, `"`, `""`)
		quoted := fmt.Sprintf(`"%s"*`, safeTerm)
		switch langFilter {
		case "de":
			// Scope to German headwords (forms are UNINDEXED, handled via LIKE separately).
			matchQuery.WriteString("word : " + quoted)
		case "en":
			matchQuery.WriteString("translation : " + quoted)
		default:
			matchQuery.WriteString(quoted)
		}
		written++
	}
	return matchQuery.String()
}

func (s *Store) Search(ctx context.Context, rawQuery string, limit int) ([]core.DictionaryEntry, error) {
	if limit <= 0 {
		limit = 50
	}
	cleanQuery, classFilter, genderFilter, langFilter := parseSearchFilters(rawQuery)
	// Filter-only browse: class/gender/lang pills with no typed text return a
	// sample of matching entries so interactive filter pills are never empty.
	// Language scope alone still returns a sample (every entry is bilingual);
	// the scope matters once the learner types a term.
	if cleanQuery == "" && (classFilter != "" || genderFilter != "" || langFilter != "") {
		// Push class/gender predicates into SQL so browse samples come from the
		// whole dictionary, not only the first 200 alphabetical rows (which
		// silently dropped later matches for rare filters like :pl / :verb).
		var conditions []string
		var args []any
		if classFilter != "" {
			conditions = append(conditions, "lower(word_class) LIKE ?")
			args = append(args, "%"+strings.ToLower(classFilter)+"%")
		}
		if genderFilter != "" {
			conditions = append(conditions, "lower(gender) LIKE ?")
			args = append(args, strings.ToLower(genderFilter)+"%")
		}
		where := ""
		if len(conditions) > 0 {
			where = "WHERE " + strings.Join(conditions, " AND ")
		}
		q := fmt.Sprintf(`
			SELECT id, word, translation, word_class, gender, forms, examples, tags
			FROM dictionary_fts
			%s
			ORDER BY word COLLATE NOCASE, id
			LIMIT ?
		`, where)
		args = append(args, limit)
		entries, err := s.queryDictionaryEntries(ctx, q, args...)
		if err != nil {
			return nil, fmt.Errorf("search dictionary filter-only: %w", err)
		}
		return entries, nil
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
			args = []any{likePattern, likePattern, limit * 4}
		} else if langFilter == "en" {
			qWild = `SELECT id, word, translation, word_class, gender, forms, examples, tags FROM dictionary_fts WHERE translation LIKE ? LIMIT ?`
			args = []any{likePattern, limit * 4}
		} else {
			qWild = `SELECT id, word, translation, word_class, gender, forms, examples, tags FROM dictionary_fts WHERE word LIKE ? OR translation LIKE ? OR forms LIKE ? LIMIT ?`
			args = []any{likePattern, likePattern, likePattern, limit * 4}
		}
		entries, err := s.queryDictionaryEntries(ctx, qWild, args...)
		if err != nil {
			return nil, fmt.Errorf("search dictionary wildcard: %w", err)
		}
		filtered := filterEntries(entries, classFilter, genderFilter, langFilter, cleanQuery)
		filtered = core.ConsolidateDictionaryEntries(filtered)
		sortDictionaryEntries(filtered, cleanQuery)
		if len(filtered) > limit {
			filtered = filtered[:limit]
		}
		return filtered, nil
	}

	matchQuery := buildFTSMatchQuery(terms, langFilter)

	q := `
		SELECT id, word, translation, word_class, gender, forms, examples, tags
		FROM dictionary_fts
		WHERE dictionary_fts MATCH ?
		ORDER BY rank
		LIMIT ?
	`

	likePattern := "%" + cleanQuery + "%"
	var entries []core.DictionaryEntry
	if matchQuery != "" {
		res, err := s.queryDictionaryEntries(ctx, q, matchQuery, limit*4)
		if err == nil {
			entries = res
			// Forms are indexed in a companion FTS5 table because the main dictionary
			// table keeps the denormalized forms column UNINDEXED for cheap imports.
			// This avoids a full table scan on every normal dictionary keystroke.
			if langFilter != "en" {
				qForms := `
					SELECT d.id, d.word, d.translation, d.word_class, d.gender,
					       d.forms, d.examples, d.tags
					FROM dictionary_forms_fts f
					INNER JOIN dictionary_fts d ON d.rowid = f.rowid
					WHERE dictionary_forms_fts MATCH ?
					LIMIT ?
				`
				formMatchQuery := buildFTSMatchQuery(strings.Fields(cleanQuery), "")
				if formMatchQuery != "" {
					formEntries, _ := s.queryDictionaryEntries(ctx, qForms, formMatchQuery, limit*2)
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
				}
			}

			if len(entries) > 0 {
				filtered := filterEntries(entries, classFilter, genderFilter, langFilter, cleanQuery)
				if len(filtered) > 0 {
					filtered = core.ConsolidateDictionaryEntries(filtered)
					sortDictionaryEntries(filtered, cleanQuery)
					if len(filtered) > limit {
						filtered = filtered[:limit]
					}
					return filtered, nil
				}
			}
		}
	}

	// Fallback to LIKE-based substring match
	var qLike string
	var args []any
	if langFilter == "de" {
		qLike = `SELECT id, word, translation, word_class, gender, forms, examples, tags FROM dictionary_fts WHERE word LIKE ? OR forms LIKE ? LIMIT ?`
		args = []any{likePattern, likePattern, limit * 4}
	} else if langFilter == "en" {
		qLike = `SELECT id, word, translation, word_class, gender, forms, examples, tags FROM dictionary_fts WHERE translation LIKE ? LIMIT ?`
		args = []any{likePattern, limit * 4}
	} else {
		qLike = `SELECT id, word, translation, word_class, gender, forms, examples, tags FROM dictionary_fts WHERE word LIKE ? OR translation LIKE ? OR forms LIKE ? LIMIT ?`
		args = []any{likePattern, likePattern, likePattern, limit * 4}
	}
	entries, err := s.queryDictionaryEntries(ctx, qLike, args...)
	if err != nil {
		return nil, fmt.Errorf("search dictionary fallback like: %w", err)
	}

	filtered := filterEntries(entries, classFilter, genderFilter, langFilter, cleanQuery)
	filtered = core.ConsolidateDictionaryEntries(filtered)
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

func (s *Store) FindRelatedEntries(ctx context.Context, word string, limit int) ([]core.DictionaryEntry, error) {
	bareWord := stripArticle(word)
	if len(bareWord) < 3 {
		return nil, nil
	}
	if limit <= 0 {
		limit = 5
	}

	safeWord := strings.ReplaceAll(bareWord, `"`, `""`)
	var entries []core.DictionaryEntry

	if utf8.RuneCountInString(bareWord) >= 4 {
		// Fast FTS5 index prefix match for stems >= 4 runes (e.g. Haustier for Haus)
		matchQuery := fmt.Sprintf(`word:"%s"*`, safeWord)
		q := `
			SELECT id, word, translation, word_class, gender, forms, examples, tags
			FROM dictionary_fts
			WHERE dictionary_fts MATCH ? AND word != ?
			ORDER BY length(word), word COLLATE NOCASE, id
			LIMIT ?
		`
		ftsEntries, _ := s.queryDictionaryEntries(ctx, q, matchQuery, word, limit*3)
		entries = append(entries, ftsEntries...)

		// Suffix matching is a fallback for compounds such as Krankenhaus for
		// Haus. Skip the unindexed scan when the indexed prefix query already
		// filled the requested result set.
		if len(entries) < limit {
			qSuffix := `
				SELECT id, word, translation, word_class, gender, forms, examples, tags
				FROM dictionary_fts
				WHERE word LIKE ? AND word != ?
				LIMIT ?
			`
			suffixEntries, _ := s.queryDictionaryEntries(ctx, qSuffix, "%"+bareWord, word, limit*2)
			entries = append(entries, suffixEntries...)
		}
	} else {
		// Short stems (< 4 runes): only match exact bare word (e.g. "das Haus" -> "Haus") or compounds ending with it
		qExact := `
			SELECT id, word, translation, word_class, gender, forms, examples, tags
			FROM dictionary_fts
			WHERE (lower(word) = lower(?) OR word LIKE ?) AND word != ?
			LIMIT ?
		`
		exactEntries, _ := s.queryDictionaryEntries(ctx, qExact, bareWord, "%"+bareWord, word, limit*2)
		entries = append(entries, exactEntries...)
	}

	var out []core.DictionaryEntry
	seen := make(map[string]bool)
	for _, e := range entries {
		key := strings.ToLower(e.Word)
		if seen[key] || strings.EqualFold(e.Word, word) {
			continue
		}
		seen[key] = true
		out = append(out, e)
		if len(out) >= limit {
			break
		}
	}
	return out, nil
}

func (s *Store) ImportEntries(ctx context.Context, entries []core.DictionaryEntry) error {
	entries = core.ConsolidateDictionaryEntries(entries)
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, "DELETE FROM dictionary_fts"); err != nil {
		return fmt.Errorf("clear dictionary: %w", err)
	}
	if _, err := tx.ExecContext(ctx, "DELETE FROM dictionary_forms_fts"); err != nil {
		return fmt.Errorf("clear dictionary forms: %w", err)
	}

	q := `
		INSERT INTO dictionary_fts (rowid, id, word, translation, word_class, gender, forms, examples, tags)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	stmt, err := tx.PrepareContext(ctx, q)
	if err != nil {
		return fmt.Errorf("prepare insert: %w", err)
	}
	defer stmt.Close()

	formsStmt, err := tx.PrepareContext(ctx, `
		INSERT INTO dictionary_forms_fts (rowid, id, forms)
		VALUES (?, ?, ?)
	`)
	if err != nil {
		return fmt.Errorf("prepare forms insert: %w", err)
	}
	defer formsStmt.Close()

	for i, entry := range entries {
		var examplesStr, tagsStr string
		if len(entry.Examples) > 0 {
			b, _ := json.Marshal(entry.Examples)
			examplesStr = string(b)
		}
		if len(entry.Tags) > 0 {
			b, _ := json.Marshal(entry.Tags)
			tagsStr = string(b)
		}

		rowID := int64(i + 1)
		if _, err := stmt.ExecContext(ctx,
			rowID,
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
		if _, err := formsStmt.ExecContext(ctx, rowID, entry.ID, entry.Forms); err != nil {
			return fmt.Errorf("insert entry forms %s: %w", entry.ID, err)
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

func (s *Store) RandomEntries(ctx context.Context, limit int) ([]core.DictionaryEntry, error) {
	if limit <= 0 {
		limit = 5
	}
	count, err := s.DictionaryCount(ctx)
	if err != nil || count == 0 {
		return nil, err
	}
	// FTS5 content tables use rowid internally. We pick random rowids
	// from the approximate range to avoid a full table scan + ORDER BY RANDOM().
	q := `
		SELECT id, word, translation, word_class, gender, forms, examples, tags
		FROM dictionary_fts
		WHERE rowid IN (
			SELECT abs(random()) % (SELECT max(rowid) FROM dictionary_fts) + 1
			FROM dictionary_fts
			LIMIT ?
		)
		LIMIT ?
	`
	return s.queryDictionaryEntries(ctx, q, limit*3, limit)
}

func (s *Store) Exists(ctx context.Context, word string) (bool, error) {
	w := strings.TrimSpace(word)
	if w == "" {
		return false, nil
	}
	var dummy int
	err := s.db.QueryRowContext(ctx, `
		SELECT 1 FROM dictionary_fts
		WHERE word = ? COLLATE NOCASE
		LIMIT 1
	`, w).Scan(&dummy)
	if err == nil {
		return true, nil
	}
	if err != sql.ErrNoRows {
		return false, fmt.Errorf("exists query: %w", err)
	}

	safeWord := strings.ReplaceAll(w, `"`, `""`)
	err = s.db.QueryRowContext(ctx, `
		SELECT 1 FROM dictionary_fts
		WHERE dictionary_fts MATCH ?
		LIMIT 1
	`, fmt.Sprintf(`"%s"`, safeWord)).Scan(&dummy)
	if err == nil {
		return true, nil
	}
	if err == sql.ErrNoRows {
		return false, nil
	}
	return false, fmt.Errorf("exists fts query: %w", err)
}
