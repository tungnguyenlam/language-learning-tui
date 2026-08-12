package sqlite

import (
	"database/sql"
	"strings"
	"time"

	"deutsch-tui/internal/core"
)

// rowScanner is implemented by both *sql.Row and *sql.Rows. Keeping card
// decoding behind this small interface gives every card query one canonical
// mapping from SQL columns to core.Card fields.
type rowScanner interface {
	Scan(dest ...any) error
}

func scanCard(row rowScanner) (core.Card, error) {
	var card core.Card
	var kind string
	var choices, tags string
	var bookmarked, leech, suspended, intervalSeconds, reviews, lapses int
	var ease float64
	var due, lastReviewed sql.NullTime

	err := row.Scan(
		&card.ID,
		&card.NoteID,
		&card.DeckID,
		&kind,
		&card.Prompt,
		&card.Answer,
		&card.Extra,
		&card.Hint,
		&choices,
		&card.Audio,
		&tags,
		&bookmarked,
		&leech,
		&suspended,
		&intervalSeconds,
		&reviews,
		&lapses,
		&ease,
		&due,
		&lastReviewed,
	)
	if err != nil {
		return core.Card{}, err
	}

	card.Kind = core.CardKind(kind)
	card.Bookmarked = bookmarked != 0
	card.Leech = leech != 0
	card.Suspended = suspended != 0
	card.Interval = time.Duration(intervalSeconds) * time.Second
	card.Reviews = reviews
	card.Lapses = lapses
	card.Ease = ease
	card.Mature = intervalSeconds >= 21*24*60*60
	card.Choices = parseChoices(choices)
	if tags != "" {
		card.Tags = strings.Fields(tags)
	}
	if due.Valid {
		card.Due = due.Time
	}
	if lastReviewed.Valid {
		card.LastReviewed = lastReviewed.Time
	}

	return card, nil
}

func scanCards(rows *sql.Rows) ([]core.Card, error) {
	var cards []core.Card
	for rows.Next() {
		card, err := scanCard(rows)
		if err != nil {
			return nil, err
		}
		cards = append(cards, card)
	}
	return cards, rows.Err()
}
