package content

// Anki's legacy collection schema (version 11). Every Anki release since 2.1
// can read it, which makes it the safest thing to write: newer versions
// upgrade it on import, older ones read it directly.
//
// The `graves` and `revlog` tables are unused by this app but are part of the
// schema; some importers query them, so they are created empty.
const ankiSchema11 = `
CREATE TABLE col (
	id      integer primary key,
	crt     integer not null,
	mod     integer not null,
	scm     integer not null,
	ver     integer not null,
	dty     integer not null,
	usn     integer not null,
	ls      integer not null,
	conf    text not null,
	models  text not null,
	decks   text not null,
	dconf   text not null,
	tags    text not null
);
CREATE TABLE notes (
	id    integer primary key,
	guid  text not null,
	mid   integer not null,
	mod   integer not null,
	usn   integer not null,
	tags  text not null,
	flds  text not null,
	sfld  integer not null,
	csum  integer not null,
	flags integer not null,
	data  text not null
);
CREATE TABLE cards (
	id     integer primary key,
	nid    integer not null,
	did    integer not null,
	ord    integer not null,
	mod    integer not null,
	usn    integer not null,
	type   integer not null,
	queue  integer not null,
	due    integer not null,
	ivl    integer not null,
	factor integer not null,
	reps   integer not null,
	lapses integer not null,
	left   integer not null,
	odue   integer not null,
	odid   integer not null,
	flags  integer not null,
	data   text not null
);
CREATE TABLE revlog (
	id      integer primary key,
	cid     integer not null,
	usn     integer not null,
	ease    integer not null,
	ivl     integer not null,
	lastIvl integer not null,
	factor  integer not null,
	time    integer not null,
	type    integer not null
);
CREATE TABLE graves (
	usn  integer not null,
	oid  integer not null,
	type integer not null
);
CREATE INDEX ix_notes_usn ON notes (usn);
CREATE INDEX ix_cards_usn ON cards (usn);
CREATE INDEX ix_cards_nid ON cards (nid);
CREATE INDEX ix_cards_sched ON cards (did, queue, due);
CREATE INDEX ix_revlog_cid ON revlog (cid);
CREATE INDEX ix_revlog_usn ON revlog (usn);
CREATE INDEX ix_notes_csum ON notes (csum);
`

// Note type ids. Anki keys note types by id, and treats a same-name type with a
// different field list as a conflict, so these use app-specific names rather
// than shadowing the user's stock "Basic"/"Cloze" types.
const (
	modelIDBasic   int64 = 1700000000001
	modelIDReverse int64 = 1700000000002
	modelIDCloze   int64 = 1700000000003
)

const ankiCardCSS = ".card {\n font-family: arial;\n font-size: 20px;\n text-align: center;\n color: black;\n background-color: white;\n}\n"

// ankiField is one entry of a note type's `flds` array.
type ankiField struct {
	Name   string   `json:"name"`
	Ord    int      `json:"ord"`
	Sticky bool     `json:"sticky"`
	RTL    bool     `json:"rtl"`
	Font   string   `json:"font"`
	Size   int      `json:"size"`
	Media  []string `json:"media"`
}

func ankiFields(names ...string) []ankiField {
	fields := make([]ankiField, len(names))
	for i, name := range names {
		fields[i] = ankiField{Name: name, Ord: i, Font: "Arial", Size: 20, Media: []string{}}
	}
	return fields
}

// ankiTemplate is one entry of a note type's `tmpls` array.
type ankiTemplate struct {
	Name  string  `json:"name"`
	Ord   int     `json:"ord"`
	QFmt  string  `json:"qfmt"`
	AFmt  string  `json:"afmt"`
	BQFmt string  `json:"bqfmt"`
	BAFmt string  `json:"bafmt"`
	DID   *int64  `json:"did"`
	BFont *string `json:"bfont,omitempty"`
}

// ankiModel is one entry of the `col.models` JSON object.
type ankiModel struct {
	ID        int64          `json:"id"`
	Name      string         `json:"name"`
	Type      int            `json:"type"` // 0 = standard, 1 = cloze
	Mod       int64          `json:"mod"`
	USN       int            `json:"usn"`
	SortF     int            `json:"sortf"`
	DID       int64          `json:"did"`
	Fields    []ankiField    `json:"flds"`
	Templates []ankiTemplate `json:"tmpls"`
	CSS       string         `json:"css"`
	LatexPre  string         `json:"latexPre"`
	LatexPost string         `json:"latexPost"`
	Req       [][]any        `json:"req"`
	Tags      []string       `json:"tags"`
	Vers      []any          `json:"vers"`
}

// ankiDeck is one entry of the `col.decks` JSON object.
type ankiDeck struct {
	ID               int64  `json:"id"`
	Name             string `json:"name"`
	Mod              int64  `json:"mod"`
	USN              int    `json:"usn"`
	Desc             string `json:"desc"`
	Conf             int64  `json:"conf"`
	Dyn              int    `json:"dyn"`
	Collapsed        bool   `json:"collapsed"`
	BrowserCollapsed bool   `json:"browserCollapsed"`
	ExtendNew        int    `json:"extendNew"`
	ExtendRev        int    `json:"extendRev"`
	LrnToday         []int  `json:"lrnToday"`
	RevToday         []int  `json:"revToday"`
	NewToday         []int  `json:"newToday"`
	TimeToday        []int  `json:"timeToday"`
}

// ankiModels returns the three note types this app writes, keyed by id string
// the way `col.models` expects.
func ankiModels(mod int64, defaultDeckID int64) map[string]ankiModel {
	basic := ankiModel{
		ID:     modelIDBasic,
		Name:   "Deutsch-TUI Basic",
		Mod:    mod,
		USN:    -1,
		DID:    defaultDeckID,
		Fields: ankiFields("Front", "Back", "Extra"),
		Templates: []ankiTemplate{{
			Name: "Card 1",
			QFmt: "{{Front}}",
			AFmt: "{{FrontSide}}\n\n<hr id=answer>\n\n{{Back}}\n\n<br>{{Extra}}",
		}},
		CSS:  ankiCardCSS,
		Req:  [][]any{{0, "any", []int{0}}},
		Tags: []string{},
		Vers: []any{},
	}

	reverse := basic
	reverse.ID = modelIDReverse
	reverse.Name = "Deutsch-TUI Basic (and reversed card)"
	reverse.Templates = []ankiTemplate{
		{
			Name: "Card 1",
			QFmt: "{{Front}}",
			AFmt: "{{FrontSide}}\n\n<hr id=answer>\n\n{{Back}}\n\n<br>{{Extra}}",
		},
		{
			Name: "Card 2",
			Ord:  1,
			QFmt: "{{Back}}",
			AFmt: "{{FrontSide}}\n\n<hr id=answer>\n\n{{Front}}\n\n<br>{{Extra}}",
		},
	}
	reverse.Req = [][]any{{0, "any", []int{0}}, {1, "any", []int{1}}}

	cloze := ankiModel{
		ID:     modelIDCloze,
		Name:   "Deutsch-TUI Cloze",
		Type:   1,
		Mod:    mod,
		USN:    -1,
		DID:    defaultDeckID,
		Fields: ankiFields("Text", "Extra"),
		Templates: []ankiTemplate{{
			Name: "Cloze",
			QFmt: "{{cloze:Text}}",
			AFmt: "{{cloze:Text}}\n\n<br>{{Extra}}",
		}},
		CSS:  ankiCardCSS,
		Req:  [][]any{{0, "any", []int{0}}},
		Tags: []string{},
		Vers: []any{},
	}

	return map[string]ankiModel{
		itoa(modelIDBasic):   basic,
		itoa(modelIDReverse): reverse,
		itoa(modelIDCloze):   cloze,
	}
}

// ankiDeckConfig is the single deck preset every exported deck points at.
func ankiDeckConfig(mod int64) map[string]any {
	return map[string]any{
		"1": map[string]any{
			"id":       1,
			"name":     "Default",
			"mod":      mod,
			"usn":      -1,
			"maxTaken": 60,
			"autoplay": true,
			"timer":    0,
			"replayq":  true,
			"new": map[string]any{
				"bury":          false,
				"delays":        []float64{1, 10},
				"initialFactor": 2500,
				"ints":          []int{1, 4, 7},
				"order":         1,
				"perDay":        20,
			},
			"rev": map[string]any{
				"bury":       false,
				"ease4":      1.3,
				"ivlFct":     1,
				"maxIvl":     36500,
				"perDay":     200,
				"hardFactor": 1.2,
			},
			"lapse": map[string]any{
				"delays":      []float64{10},
				"leechAction": 1,
				"leechFails":  8,
				"minInt":      1,
				"mult":        0,
			},
			"dyn": false,
		},
	}
}

// ankiCollectionConfig is the `col.conf` blob. Anki tolerates missing keys but
// a bare `{}` leaves it without a current deck or note type.
func ankiCollectionConfig(curDeck int64) map[string]any {
	return map[string]any{
		"nextPos":       1,
		"estTimes":      true,
		"activeDecks":   []int64{curDeck},
		"sortType":      "noteFld",
		"timeLim":       0,
		"sortBackwards": false,
		"addToCur":      true,
		"curDeck":       curDeck,
		"newBury":       true,
		"newSpread":     0,
		"dueCounts":     true,
		"curModel":      itoa(modelIDBasic),
		"collapseTime":  1200,
		"schedVer":      2,
	}
}
