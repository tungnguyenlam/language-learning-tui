package content

import (
	"deutsch-tui/internal/core"
)

func B2MusicDeck() core.Deck {
	notes := []core.Note{
		{ID: "b2-msc-musik", DeckID: "b2-music", Front: "die Musik", Back: "music", Tags: []string{"b2", "music", "entertainment"}},
		{ID: "b2-msc-musiker", DeckID: "b2-music", Front: "der Musiker / die Musikerin", Back: "musician", Tags: []string{"b2", "music", "entertainment"}},
		{ID: "b2-msc-band", DeckID: "b2-music", Front: "die Band", Back: "band", Tags: []string{"b2", "music", "entertainment"}},
		{ID: "b2-msc-konzert", DeckID: "b2-music", Front: "das Konzert", Back: "concert", Tags: []string{"b2", "music", "entertainment"}},
		{ID: "b2-msc-künstler", DeckID: "b2-music", Front: "der Künstler / die Künstlerin", Back: "artist", Tags: []string{"b2", "music", "entertainment"}},
		{ID: "b2-msc-album", DeckID: "b2-music", Front: "das Album", Back: "album", Tags: []string{"b2", "music", "entertainment"}},
		{ID: "b2-msc-lied", DeckID: "b2-music", Front: "das Lied", Back: "song", Tags: []string{"b2", "music", "entertainment"}},
		{ID: "b2-msc-hit", DeckID: "b2-music", Front: "der Hit", Back: "hit (song)", Tags: []string{"b2", "music", "entertainment"}},
		{ID: "b2-msc-titel", DeckID: "b2-music", Front: "der Titel", Back: "title / track", Tags: []string{"b2", "music", "entertainment"}},
		{ID: "b2-msc-single", DeckID: "b2-music", Front: "die Single", Back: "single", Tags: []string{"b2", "music", "entertainment"}},
		{ID: "b2-msc-komponist", DeckID: "b2-music", Front: "der Komponist / die Komponistin", Back: "composer", Tags: []string{"b2", "music", "entertainment"}},
		{ID: "b2-msc-dirigent", DeckID: "b2-music", Front: "der Dirigent / die Dirigentin", Back: "conductor", Tags: []string{"b2", "music", "entertainment"}},
		{ID: "b2-msc-orchester", DeckID: "b2-music", Front: "das Orchester", Back: "orchestra", Tags: []string{"b2", "music", "entertainment"}},
		{ID: "b2-msc-symphonie", DeckID: "b2-music", Front: "die Symphonie", Back: "symphony", Tags: []string{"b2", "music", "entertainment"}},
		{ID: "b2-msc-sonate", DeckID: "b2-music", Front: "die Sonate", Back: "sonata", Tags: []string{"b2", "music", "entertainment"}},
		{ID: "b2-msc-konzertsaal", DeckID: "b2-music", Front: "der Konzertsaal", Back: "concert hall", Tags: []string{"b2", "music", "entertainment"}},
		{ID: "b2-msc-bühne", DeckID: "b2-music", Front: "die Bühne", Back: "stage", Tags: []string{"b2", "music", "entertainment"}},
		{ID: "b2-msc-vorhang", DeckID: "b2-music", Front: "der Vorhang", Back: "curtain", Tags: []string{"b2", "music", "entertainment"}},
		{ID: "b2-msc-sänger", DeckID: "b2-music", Front: "der Sänger / die Sängerin", Back: "singer / vocalist", Tags: []string{"b2", "music", "entertainment"}},
		{ID: "b2-msc-gitarrist", DeckID: "b2-music", Front: "der Gitarrist / die Gitarristin", Back: "guitarist", Tags: []string{"b2", "music", "entertainment"}},
		{ID: "b2-msc-schlagzeug", DeckID: "b2-music", Front: "das Schlagzeug", Back: "drums / percussion", Tags: []string{"b2", "music", "entertainment"}},
		{ID: "b2-msc-tasteninstr", DeckID: "b2-music", Front: "das Tasteninstrument", Back: "keyboard instrument", Tags: []string{"b2", "music", "entertainment"}},
		{ID: "b2-msc-studio", DeckID: "b2-music", Front: "das Studio", Back: "studio", Tags: []string{"b2", "music", "entertainment"}},
		{ID: "b2-msc-produzent", DeckID: "b2-music", Front: "der Produzent / die Produzentin", Back: "producer", Tags: []string{"b2", "music", "entertainment"}},
		{ID: "b2-msc-tontechnik", DeckID: "b2-music", Front: "die Tontechnik", Back: "sound engineering", Tags: []string{"b2", "music", "entertainment"}},
		{ID: "b2-msc-mischpult", DeckID: "b2-music", Front: "das Mischpult", Back: "mixing console", Tags: []string{"b2", "music", "entertainment"}},
		{ID: "b2-msc-plattenfirma", DeckID: "b2-music", Front: "die Plattenfirma", Back: "record label", Tags: []string{"b2", "music", "entertainment"}},
		{ID: "b2-msc-charts", DeckID: "b2-music", Front: "die Charts", Back: "charts (music charts)", Tags: []string{"b2", "music", "entertainment"}},
		{ID: "b2-msc-rhythmus", DeckID: "b2-music", Front: "der Rhythmus", Back: "rhythm", Tags: []string{"b2", "music", "entertainment"}},
		{ID: "b2-msc-melodie", DeckID: "b2-music", Front: "die Melodie", Back: "melody", Tags: []string{"b2", "music", "entertainment"}},
		{ID: "b2-msc-harmonie", DeckID: "b2-music", Front: "die Harmonie", Back: "harmony", Tags: []string{"b2", "music", "entertainment"}},
		{ID: "b2-msc-takt", DeckID: "b2-music", Front: "der Takt", Back: "beat / measure", Tags: []string{"b2", "music", "entertainment"}},
		{ID: "b2-msc-note", DeckID: "b2-music", Front: "die Note", Back: "note (musical)", Tags: []string{"b2", "music", "entertainment"}},
		{ID: "b2-msc-partitur", DeckID: "b2-music", Front: "die Partitur", Back: "score (musical)", Tags: []string{"b2", "music", "entertainment"}},
		{ID: "b2-msc-instrument", DeckID: "b2-music", Front: "das Instrument", Back: "instrument", Tags: []string{"b2", "music", "entertainment"}},
		{ID: "b2-msc-kino", DeckID: "b2-music", Front: "das Kino", Back: "cinema / movie theater", Tags: []string{"b2", "music", "entertainment"}},
		{ID: "b2-msc-film", DeckID: "b2-music", Front: "der Film", Back: "film / movie", Tags: []string{"b2", "music", "entertainment"}},
		{ID: "b2-msc-regisseur", DeckID: "b2-music", Front: "der Regisseur / die Regisseurin", Back: "director (film)", Tags: []string{"b2", "music", "entertainment"}},
		{ID: "b2-msc-schauspieler", DeckID: "b2-music", Front: "der Schauspieler / die Schauspielerin", Back: "actor / actress", Tags: []string{"b2", "music", "entertainment"}},
		{ID: "b2-msc-drehbuch", DeckID: "b2-music", Front: "das Drehbuch", Back: "screenplay", Tags: []string{"b2", "music", "entertainment"}},
		{ID: "b2-msc-szene", DeckID: "b2-music", Front: "die Szene", Back: "scene", Tags: []string{"b2", "music", "entertainment"}},
		{ID: "b2-msc-premiere", DeckID: "b2-music", Front: "die Premiere", Back: "premiere", Tags: []string{"b2", "music", "entertainment"}},
		{ID: "b2-msc-festival", DeckID: "b2-music", Front: "das Festival", Back: "festival", Tags: []string{"b2", "music", "entertainment"}},
		{ID: "b2-msc-veranstaltung", DeckID: "b2-music", Front: "die Veranstaltung", Back: "event", Tags: []string{"b2", "music", "entertainment"}},
		{ID: "b2-msc-besucher", DeckID: "b2-music", Front: "der Besucher / die Besucherin", Back: "visitor / audience", Tags: []string{"b2", "music", "entertainment"}},
		{ID: "b2-msc-publikum", DeckID: "b2-music", Front: "das Publikum", Back: "audience", Tags: []string{"b2", "music", "entertainment"}},
		{ID: "b2-msc-applaus", DeckID: "b2-music", Front: "der Applaus", Back: "applause", Tags: []string{"b2", "music", "entertainment"}},
		{ID: "b2-msc-ovalbum", DeckID: "b2-music", Front: "das Vinylalbum", Back: "vinyl record", Tags: []string{"b2", "music", "entertainment"}},
		{ID: "b2-msc-streaming", DeckID: "b2-music", Front: "das Streaming", Back: "streaming", Tags: []string{"b2", "music", "entertainment"}},
		{ID: "b2-msc-podcast", DeckID: "b2-music", Front: "der Podcast", Back: "podcast", Tags: []string{"b2", "music", "entertainment"}},
		{ID: "b2-msc-live-auftritt", DeckID: "b2-music", Front: "der Live-Auftritt", Back: "live performance", Tags: []string{"b2", "music", "entertainment"}},
	}

	for i := range notes {
		notes[i].Cards = CardsForNote(notes[i])
	}

	return core.Deck{
		ID:          "b2-music",
		Name:        "German B2 Music & Entertainment",
		Description: "Advanced vocabulary for music, concerts, and entertainment.",
		Tags:        []string{"german", "b2", "music", "entertainment"},
		Notes:       notes,
	}
}
