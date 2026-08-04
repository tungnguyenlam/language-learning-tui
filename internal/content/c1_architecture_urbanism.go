package content

import (
	"deutsch-tui/internal/core"
)

func C1ArchitectureUrbanismDeck() core.Deck {
	notes := []core.Note{
		{ID: "c1-arch-architektur", DeckID: "c1-architecture-urbanism", Front: "die Architektur", Back: "architecture", Tags: []string{"c1", "architecture"}},
		{ID: "c1-arch-architekt", DeckID: "c1-architecture-urbanism", Front: "der Architekt / die Architektin", Back: "architect", Tags: []string{"c1", "architecture"}},
		{ID: "c1-arch-staedtebau", DeckID: "c1-architecture-urbanism", Front: "der Städtebau", Back: "urban design", Tags: []string{"c1", "urbanism"}},
		{ID: "c1-arch-stadtplanung", DeckID: "c1-architecture-urbanism", Front: "die Stadtplanung", Back: "urban planning", Tags: []string{"c1", "urbanism"}},
		{ID: "c1-arch-bebauungsplan", DeckID: "c1-architecture-urbanism", Front: "der Bebauungsplan", Back: "development / zoning plan", Tags: []string{"c1", "urbanism"}},
		{ID: "c1-arch-baugenehmigung", DeckID: "c1-architecture-urbanism", Front: "die Baugenehmigung", Back: "building permit", Tags: []string{"c1", "architecture"}},
		{ID: "c1-arch-bauvorschrift", DeckID: "c1-architecture-urbanism", Front: "die Bauvorschrift", Back: "building regulation", Tags: []string{"c1", "architecture"}},
		{ID: "c1-arch-entwurf", DeckID: "c1-architecture-urbanism", Front: "der Entwurf", Back: "design / draft", Tags: []string{"c1", "architecture"}},
		{ID: "c1-arch-grundriss", DeckID: "c1-architecture-urbanism", Front: "der Grundriss", Back: "floor plan", Tags: []string{"c1", "architecture"}},
		{ID: "c1-arch-aufriss", DeckID: "c1-architecture-urbanism", Front: "der Aufriss", Back: "elevation (drawing)", Tags: []string{"c1", "architecture"}},
		{ID: "c1-arch-fassade", DeckID: "c1-architecture-urbanism", Front: "die Fassade", Back: "facade", Tags: []string{"c1", "architecture"}},
		{ID: "c1-arch-tragwerk", DeckID: "c1-architecture-urbanism", Front: "das Tragwerk", Back: "load-bearing structure", Tags: []string{"c1", "architecture"}},
		{ID: "c1-arch-fundament", DeckID: "c1-architecture-urbanism", Front: "das Fundament", Back: "foundation", Tags: []string{"c1", "architecture"}},
		{ID: "c1-arch-traeger", DeckID: "c1-architecture-urbanism", Front: "der Träger", Back: "girder / beam", Tags: []string{"c1", "architecture"}},
		{ID: "c1-arch-saeule", DeckID: "c1-architecture-urbanism", Front: "die Säule", Back: "column / pillar", Tags: []string{"c1", "architecture"}},
		{ID: "c1-arch-bogen", DeckID: "c1-architecture-urbanism", Front: "der Bogen", Back: "arch", Tags: []string{"c1", "architecture"}},
		{ID: "c1-arch-kuppel", DeckID: "c1-architecture-urbanism", Front: "die Kuppel", Back: "dome", Tags: []string{"c1", "architecture"}},
		{ID: "c1-arch-gewoelbe", DeckID: "c1-architecture-urbanism", Front: "das Gewölbe", Back: "vault", Tags: []string{"c1", "architecture"}},
		{ID: "c1-arch-giebel", DeckID: "c1-architecture-urbanism", Front: "der Giebel", Back: "gable", Tags: []string{"c1", "architecture"}},
		{ID: "c1-arch-erker", DeckID: "c1-architecture-urbanism", Front: "der Erker", Back: "bay window / oriel", Tags: []string{"c1", "architecture"}},
		{ID: "c1-arch-dachgeschoss", DeckID: "c1-architecture-urbanism", Front: "das Dachgeschoss", Back: "top floor / attic storey", Tags: []string{"c1", "architecture"}},
		{ID: "c1-arch-erdgeschoss", DeckID: "c1-architecture-urbanism", Front: "das Erdgeschoss", Back: "ground floor", Tags: []string{"c1", "architecture"}},
		{ID: "c1-arch-grundstueck", DeckID: "c1-architecture-urbanism", Front: "das Grundstück", Back: "plot of land", Tags: []string{"c1", "urbanism"}},
		{ID: "c1-arch-bausubstanz", DeckID: "c1-architecture-urbanism", Front: "die Bausubstanz", Back: "building fabric / structural condition", Tags: []string{"c1", "architecture"}},
		{ID: "c1-arch-denkmalschutz", DeckID: "c1-architecture-urbanism", Front: "der Denkmalschutz", Back: "heritage protection", Tags: []string{"c1", "architecture"}},
		{ID: "c1-arch-sanierung", DeckID: "c1-architecture-urbanism", Front: "die Sanierung", Back: "refurbishment / redevelopment", Tags: []string{"c1", "architecture"}},
		{ID: "c1-arch-umbau", DeckID: "c1-architecture-urbanism", Front: "der Umbau", Back: "conversion / remodelling", Tags: []string{"c1", "architecture"}},
		{ID: "c1-arch-abriss", DeckID: "c1-architecture-urbanism", Front: "der Abriss", Back: "demolition", Tags: []string{"c1", "architecture"}},
		{ID: "c1-arch-nachverdichtung", DeckID: "c1-architecture-urbanism", Front: "die Nachverdichtung", Back: "urban infill densification", Tags: []string{"c1", "urbanism"}},
		{ID: "c1-arch-flaechennutzung", DeckID: "c1-architecture-urbanism", Front: "die Flächennutzung", Back: "land use", Tags: []string{"c1", "urbanism"}},
		{ID: "c1-arch-gruenflaeche", DeckID: "c1-architecture-urbanism", Front: "die Grünfläche", Back: "green space", Tags: []string{"c1", "urbanism"}},
		{ID: "c1-arch-fussgaengerzone", DeckID: "c1-architecture-urbanism", Front: "die Fußgängerzone", Back: "pedestrian zone", Tags: []string{"c1", "urbanism"}},
		{ID: "c1-arch-wohnquartier", DeckID: "c1-architecture-urbanism", Front: "das Wohnquartier", Back: "residential quarter", Tags: []string{"c1", "urbanism"}},
		{ID: "c1-arch-stadtteil", DeckID: "c1-architecture-urbanism", Front: "der Stadtteil", Back: "city district", Tags: []string{"c1", "urbanism"}},
		{ID: "c1-arch-zersiedelung", DeckID: "c1-architecture-urbanism", Front: "die Zersiedelung", Back: "urban sprawl", Tags: []string{"c1", "urbanism"}},
		{ID: "c1-arch-gentrifizierung", DeckID: "c1-architecture-urbanism", Front: "die Gentrifizierung", Back: "gentrification", Tags: []string{"c1", "urbanism"}},
		{ID: "c1-arch-leerstand", DeckID: "c1-architecture-urbanism", Front: "der Leerstand", Back: "vacancy (unoccupied buildings)", Tags: []string{"c1", "urbanism"}},
		{ID: "c1-arch-verkehrsplanung", DeckID: "c1-architecture-urbanism", Front: "die Verkehrsplanung", Back: "transport planning", Tags: []string{"c1", "urbanism"}},
		{ID: "c1-arch-baustoff", DeckID: "c1-architecture-urbanism", Front: "der Baustoff", Back: "building material", Tags: []string{"c1", "architecture"}},
		{ID: "c1-arch-beton", DeckID: "c1-architecture-urbanism", Front: "der Beton", Back: "concrete", Tags: []string{"c1", "architecture"}},
		{ID: "c1-arch-stahlbeton", DeckID: "c1-architecture-urbanism", Front: "der Stahlbeton", Back: "reinforced concrete", Tags: []string{"c1", "architecture"}},
		{ID: "c1-arch-daemmung", DeckID: "c1-architecture-urbanism", Front: "die Dämmung", Back: "insulation", Tags: []string{"c1", "architecture"}},
		{ID: "c1-arch-statik", DeckID: "c1-architecture-urbanism", Front: "die Statik", Back: "structural analysis / statics", Tags: []string{"c1", "architecture"}},
		{ID: "c1-arch-barrierefreiheit", DeckID: "c1-architecture-urbanism", Front: "die Barrierefreiheit", Back: "accessibility (barrier-free design)", Tags: []string{"c1", "architecture"}},
		{ID: "c1-arch-hochhaus", DeckID: "c1-architecture-urbanism", Front: "das Hochhaus", Back: "high-rise building", Tags: []string{"c1", "architecture"}},
		{ID: "c1-arch-wolkenkratzer", DeckID: "c1-architecture-urbanism", Front: "der Wolkenkratzer", Back: "skyscraper", Tags: []string{"c1", "architecture"}},
		{ID: "c1-arch-entwerfen", DeckID: "c1-architecture-urbanism", Front: "entwerfen", Back: "to design", Tags: []string{"c1", "architecture", "verb"}},
		{ID: "c1-arch-errichten", DeckID: "c1-architecture-urbanism", Front: "errichten", Back: "to erect / to construct", Tags: []string{"c1", "architecture", "verb"}},
		{ID: "c1-arch-sanieren", DeckID: "c1-architecture-urbanism", Front: "sanieren", Back: "to refurbish / to renovate", Tags: []string{"c1", "architecture", "verb"}},
		{ID: "c1-arch-abreissen", DeckID: "c1-architecture-urbanism", Front: "abreißen", Back: "to demolish / to tear down", Tags: []string{"c1", "architecture", "verb"}},
		{ID: "c1-arch-umbauen", DeckID: "c1-architecture-urbanism", Front: "umbauen", Back: "to convert / to remodel", Tags: []string{"c1", "architecture", "verb"}},
		{ID: "c1-arch-daemmen", DeckID: "c1-architecture-urbanism", Front: "dämmen", Back: "to insulate", Tags: []string{"c1", "architecture", "verb"}},
		{ID: "c1-arch-genehmigen", DeckID: "c1-architecture-urbanism", Front: "genehmigen", Back: "to approve / to authorise", Tags: []string{"c1", "architecture", "verb"}},
		{ID: "c1-arch-gestalten", DeckID: "c1-architecture-urbanism", Front: "gestalten", Back: "to shape / to design", Tags: []string{"c1", "architecture", "verb"}},
		{ID: "c1-arch-denkmalgeschuetzt", DeckID: "c1-architecture-urbanism", Front: "denkmalgeschützt", Back: "listed / heritage-protected", Tags: []string{"c1", "architecture", "adjective"}},
		{ID: "c1-arch-tragend", DeckID: "c1-architecture-urbanism", Front: "tragend", Back: "load-bearing", Tags: []string{"c1", "architecture", "adjective"}},
		{ID: "c1-arch-freistehend", DeckID: "c1-architecture-urbanism", Front: "freistehend", Back: "detached / free-standing", Tags: []string{"c1", "architecture", "adjective"}},
		{ID: "c1-arch-barrierefrei", DeckID: "c1-architecture-urbanism", Front: "barrierefrei", Back: "barrier-free / accessible", Tags: []string{"c1", "architecture", "adjective"}},
		{ID: "c1-arch-massstabsgetreu", DeckID: "c1-architecture-urbanism", Front: "maßstabsgetreu", Back: "true to scale", Tags: []string{"c1", "architecture", "adjective"}},
	}

	for i := range notes {
		notes[i].Cards = CardsForNote(notes[i])
	}

	return core.Deck{
		ID:          "c1-architecture-urbanism",
		Name:        "C1 Architektur und Städtebau",
		Description: "Advanced architecture and urban planning vocabulary.",
		Tags:        []string{"german", "c1", "architecture", "urbanism"},
		Notes:       notes,
	}
}
