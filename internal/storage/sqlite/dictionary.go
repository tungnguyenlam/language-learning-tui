package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"deutsch-tui/internal/core"
)

func (s *Store) Search(ctx context.Context, query string, limit int) ([]core.DictionaryEntry, error) {
	terms := strings.Fields(query)
	if len(terms) == 0 {
		return nil, nil
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

	rows, err := s.db.QueryContext(ctx, q, matchQuery.String(), limit)
	if err != nil {
		return nil, fmt.Errorf("search dictionary: %w", err)
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
