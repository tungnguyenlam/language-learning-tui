package sqlite

type migration struct {
	ID  int
	SQL string
}

var migrations = []migration{
	{
		ID: 1,
		SQL: `CREATE TABLE IF NOT EXISTS decks (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		description TEXT NOT NULL DEFAULT ''
	)`,
	},
	{
		ID: 2,
		SQL: `CREATE TABLE IF NOT EXISTS notes (
		id TEXT PRIMARY KEY,
		deck_id TEXT NOT NULL REFERENCES decks(id) ON DELETE CASCADE,
		front TEXT NOT NULL,
		back TEXT NOT NULL,
		extra TEXT NOT NULL DEFAULT ''
	)`,
	},
	{
		ID: 3,
		SQL: `CREATE TABLE IF NOT EXISTS cards (
		id TEXT PRIMARY KEY,
		note_id TEXT NOT NULL REFERENCES notes(id) ON DELETE CASCADE,
		deck_id TEXT NOT NULL REFERENCES decks(id) ON DELETE CASCADE,
		kind TEXT NOT NULL,
		prompt TEXT NOT NULL,
		answer TEXT NOT NULL
	)`,
	},
	{
		ID: 4,
		SQL: `CREATE TABLE IF NOT EXISTS review_states (
		card_id TEXT PRIMARY KEY REFERENCES cards(id) ON DELETE CASCADE,
		due_at TIMESTAMP NOT NULL,
		interval_seconds INTEGER NOT NULL DEFAULT 0,
		ease REAL NOT NULL DEFAULT 2.5,
		reviews INTEGER NOT NULL DEFAULT 0,
		lapses INTEGER NOT NULL DEFAULT 0
	)`,
	},
	{
		ID: 5,
		SQL: `CREATE TABLE IF NOT EXISTS reviews (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		card_id TEXT NOT NULL REFERENCES cards(id) ON DELETE CASCADE,
		grade TEXT NOT NULL,
		reviewed_at TIMESTAMP NOT NULL,
		due_at TIMESTAMP NOT NULL,
		interval_seconds INTEGER NOT NULL,
		ease REAL NOT NULL,
		reviews INTEGER NOT NULL,
		lapses INTEGER NOT NULL
	)`,
	},
	{ID: 6, SQL: `CREATE INDEX IF NOT EXISTS idx_cards_deck ON cards(deck_id)`},
	{ID: 7, SQL: `CREATE INDEX IF NOT EXISTS idx_review_states_due ON review_states(due_at)`},
	{ID: 8, SQL: `CREATE INDEX IF NOT EXISTS idx_reviews_card_reviewed ON reviews(card_id, reviewed_at)`},
	{
		ID: 9,
		SQL: `
		ALTER TABLE review_states ADD COLUMN stability REAL NOT NULL DEFAULT 0;
		ALTER TABLE review_states ADD COLUMN difficulty REAL NOT NULL DEFAULT 0;
		ALTER TABLE reviews ADD COLUMN stability REAL NOT NULL DEFAULT 0;
		ALTER TABLE reviews ADD COLUMN difficulty REAL NOT NULL DEFAULT 0;
	`,
	},
	{
		ID: 10,
		SQL: `CREATE TABLE IF NOT EXISTS card_flags (
		card_id TEXT PRIMARY KEY REFERENCES cards(id) ON DELETE CASCADE,
		bookmarked INTEGER NOT NULL DEFAULT 0,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	)`,
	},
	{
		ID: 11,
		SQL: `
		ALTER TABLE card_flags ADD COLUMN leech INTEGER NOT NULL DEFAULT 0;
		ALTER TABLE card_flags ADD COLUMN lapse_streak INTEGER NOT NULL DEFAULT 0;
	`,
	},
	{
		ID:  12,
		SQL: `ALTER TABLE cards ADD COLUMN choices TEXT NOT NULL DEFAULT ''`,
	},
	{
		ID:  13,
		SQL: `ALTER TABLE card_flags ADD COLUMN suspended INTEGER NOT NULL DEFAULT 0`,
	},
	{
		ID: 14,
		SQL: `CREATE TABLE IF NOT EXISTS app_settings (
			key TEXT PRIMARY KEY,
			value TEXT NOT NULL,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`,
	},
	{
		ID: 15,
		SQL: `
			ALTER TABLE cards ADD COLUMN audio TEXT NOT NULL DEFAULT '';
			ALTER TABLE notes ADD COLUMN audio TEXT NOT NULL DEFAULT '';
		`,
	},
	{
		ID:  16,
		SQL: `ALTER TABLE decks ADD COLUMN tags TEXT NOT NULL DEFAULT ''`,
	},
	{
		ID: 17,
		SQL: `
			ALTER TABLE notes ADD COLUMN tags TEXT NOT NULL DEFAULT '';
			ALTER TABLE cards ADD COLUMN tags TEXT NOT NULL DEFAULT '';
		`,
	},
	{
		ID:  18,
		SQL: `CREATE INDEX IF NOT EXISTS idx_reviews_reviewed_at ON reviews(reviewed_at)`,
	},
	{
		ID: 19,
		SQL: `
			ALTER TABLE decks ADD COLUMN new_cards_per_day INTEGER NOT NULL DEFAULT 20;
			ALTER TABLE decks ADD COLUMN review_limit_per_day INTEGER NOT NULL DEFAULT 200;
	`,
	},
	{
		ID: 20,
		SQL: `
			DELETE FROM notes
			WHERE id = 'front'
			  AND front = 'back'
			  AND back = 'extra';

			UPDATE cards
			SET answer = 'Alle Wege führen nach Rom'
			WHERE deck_id = 'b1_idioms'
			  AND prompt = 'All roads lead to Rome'
			  AND answer = 'Literal: all ways lead to Rome';

			UPDATE notes
			SET front = 'Alle Wege führen nach Rom',
			    back = 'All roads lead to Rome',
			    extra = 'Literal: all ways lead to Rome',
			    tags = 'idiom b1'
			WHERE id = 'Alle Wege führen nach Rom'
			  AND front = 'All roads lead to Rome'
			  AND back = 'Literal: all ways lead to Rome';
		`,
	},
}
