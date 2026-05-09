package content

import (
	"deutsch-tui/internal/core"
)

func PhilosophyLiteratureDeck() core.Deck {
	notes := []core.Note{
		// --- Philosophy Core Concepts ---
		{ID: "phil-das-sein", DeckID: "philosophy-literature", Front: "das Sein", Back: "being/existence", Extra: "Central concept in existential philosophy. Heidegger's 'Sein und Zeit'.", Tags: []string{"b2", "c1", "philosophy", "noun"}},
		{ID: "phil-das-dasein", DeckID: "philosophy-literature", Front: "das Dasein", Back: "existence/being-there", Extra: "Heidegger's term for human existence, literally 'being-there'.", Tags: []string{"b2", "c1", "philosophy", "noun"}},
		{ID: "phil-die-weltanschauung", DeckID: "philosophy-literature", Front: "die Weltanschauung", Back: "worldview", Extra: "A comprehensive view of the world and human life.", Tags: []string{"b2", "c1", "philosophy", "noun"}},
		{ID: "phil-die-erkenntnis", DeckID: "philosophy-literature", Front: "die Erkenntnis", Back: "knowledge/insight/recognition", Extra: "Key term in epistemology (Erkenntnistheorie).", Tags: []string{"b2", "c1", "philosophy", "noun"}},
		{ID: "phil-die-wahrheit", DeckID: "philosophy-literature", Front: "die Wahrheit", Back: "truth", Extra: "Nietzsche: 'Was ist Wahrheit?'", Tags: []string{"b2", "c1", "philosophy", "noun"}},
		{ID: "phil-die-freiheit", DeckID: "philosophy-literature", Front: "die Freiheit", Back: "freedom/liberty", Extra: "Central to Kant's moral philosophy and existentialism.", Tags: []string{"b2", "c1", "philosophy", "noun"}},
		{ID: "phil-die-gerechtigkeit", DeckID: "philosophy-literature", Front: "die Gerechtigkeit", Back: "justice/fairness", Extra: "Plato's 'Politeia' (The Republic) explores this concept.", Tags: []string{"b2", "c1", "philosophy", "noun"}},
		{ID: "phil-die-moral", DeckID: "philosophy-literature", Front: "die Moral", Back: "morality/morals", Extra: "Kant's 'Kritik der praktischen Vernunft'.", Tags: []string{"b2", "c1", "philosophy", "noun"}},
		{ID: "phil-die-ethik", DeckID: "philosophy-literature", Front: "die Ethik", Back: "ethics", Extra: "Branch of philosophy dealing with right and wrong.", Tags: []string{"b2", "c1", "philosophy", "noun"}},
		{ID: "phil-die-dialektik", DeckID: "philosophy-literature", Front: "die Dialektik", Back: "dialectics", Extra: "Hegel's method: thesis-antithesis-synthesis.", Tags: []string{"b2", "c1", "philosophy", "noun"}},
		{ID: "phil-die-aufklarung", DeckID: "philosophy-literature", Front: "die Aufklärung", Back: "Enlightenment", Extra: "Kant: 'Aufklärung ist der Ausgang des Menschen aus seiner selbstverschuldeten Unmündigkeit.'", Tags: []string{"b2", "c1", "philosophy", "noun"}},
		{ID: "phil-der-idealismus", DeckID: "philosophy-literature", Front: "der Idealismus", Back: "idealism", Extra: "Philosophical tradition from Kant to Hegel.", Tags: []string{"b2", "c1", "philosophy", "noun"}},
		{ID: "phil-der-existentialismus", DeckID: "philosophy-literature", Front: "der Existentialismus", Back: "existentialism", Extra: "Jaspers, Heidegger, Sartre.", Tags: []string{"b2", "c1", "philosophy", "noun"}},
		{ID: "phil-der-nihilismus", DeckID: "philosophy-literature", Front: "der Nihilismus", Back: "nihilism", Extra: "Nietzsche: 'Gott ist tot.'", Tags: []string{"b2", "c1", "philosophy", "noun"}},
		{ID: "phil-die-vernunft", DeckID: "philosophy-literature", Front: "die Vernunft", Back: "reason", Extra: "Kant's 'Kritik der reinen Vernunft'.", Tags: []string{"b2", "c1", "philosophy", "noun"}},
		{ID: "phil-die-vernunft-pflicht", DeckID: "philosophy-literature", Front: "die Pflicht", Back: "duty", Extra: "Kant's categorical imperative acts from duty.", Tags: []string{"b2", "c1", "philosophy", "noun"}},
		{ID: "phil-der-ubermensch", DeckID: "philosophy-literature", Front: "der Übermensch", Back: "overman/superman", Extra: "Nietzsche's concept in 'Also sprach Zarathustra'.", Tags: []string{"b2", "c1", "philosophy", "noun"}},
		{ID: "phil-das-gewissen", DeckID: "philosophy-literature", Front: "das Gewissen", Back: "conscience", Extra: "Central to moral philosophy and ethics.", Tags: []string{"b2", "c1", "philosophy", "noun"}},

		// --- Philosophy Verbs ---
		{ID: "phil-denken", DeckID: "philosophy-literature", Front: "denken", Back: "to think", Tags: []string{"b2", "c1", "philosophy", "verb"}},
		{ID: "phil-reflektieren", DeckID: "philosophy-literature", Front: "reflektieren", Back: "to reflect", Tags: []string{"b2", "c1", "philosophy", "verb"}},
		{ID: "phil-hinterfragen", DeckID: "philosophy-literature", Front: "hinterfragen", Back: "to question/challenge", Tags: []string{"b2", "c1", "philosophy", "verb"}},
		{ID: "phil-begreifen", DeckID: "philosophy-literature", Front: "begreifen", Back: "to comprehend/grasp", Tags: []string{"b2", "c1", "philosophy", "verb"}},
		{ID: "phil-schlussfolgern", DeckID: "philosophy-literature", Front: "schließen", Back: "to conclude/deduce", Tags: []string{"b2", "c1", "philosophy", "verb"}},
		{ID: "phil-ableiten", DeckID: "philosophy-literature", Front: "ableiten", Back: "to derive/deduce", Tags: []string{"b2", "c1", "philosophy", "verb"}},

		// --- Famous German Philosophers ---
		{ID: "phil-kant", DeckID: "philosophy-literature", Front: "Immanuel Kant (1724-1804)", Back: "Critique of Pure Reason; categorical imperative", Extra: " Königsberg philosopher, founder of German idealism.", Tags: []string{"b2", "c1", "philosophy", "philosopher"}},
		{ID: "phil-hegel", DeckID: "philosophy-literature", Front: "Georg W.F. Hegel (1770-1831)", Back: "Phenomenology of Spirit; dialectical method", Extra: "Major figure in German idealism.", Tags: []string{"b2", "c1", "philosophy", "philosopher"}},
		{ID: "phil-nietzsche", DeckID: "philosophy-literature", Front: "Friedrich Nietzsche (1844-1900)", Back: "Thus Spoke Zarathustra; God is dead; will to power", Extra: "Radical critic of morality and religion.", Tags: []string{"b2", "c1", "philosophy", "philosopher"}},
		{ID: "phil-heidegger", DeckID: "philosophy-literature", Front: "Martin Heidegger (1889-1976)", Back: "Being and Time; Dasein; question of Being", Extra: "Most influential 20th-century German philosopher.", Tags: []string{"b2", "c1", "philosophy", "philosopher"}},
		{ID: "phil-schopenhauer", DeckID: "philosophy-literature", Front: "Arthur Schopenhauer (1788-1860)", Back: "The World as Will and Representation", Extra: "Pessimistic philosophy, influenced Nietzsche and Wittgenstein.", Tags: []string{"b2", "c1", "philosophy", "philosopher"}},
		{ID: "phil-leibniz", DeckID: "philosophy-literature", Front: "Gottfried Wilhelm Leibniz (1646-1716)", Back: "Monadology; best of all possible worlds", Extra: "Also mathematician, co-inventor of calculus.", Tags: []string{"b2", "c1", "philosophy", "philosopher"}},

		// --- Literature Core Concepts ---
		{ID: "lit-die-dichtung", DeckID: "philosophy-literature", Front: "die Dichtung", Back: "poetry/literature", Extra: "General term for literary creation.", Tags: []string{"b2", "c1", "literature", "noun"}},
		{ID: "lit-der-roman", DeckID: "philosophy-literature", Front: "der Roman", Back: "novel", Extra: "der Bildungsroman = coming-of-age novel.", Tags: []string{"b2", "c1", "literature", "noun"}},
		{ID: "lit-die-erzahlung", DeckID: "philosophy-literature", Front: "die Erzählung", Back: "narrative/story", Extra: "Short prose fiction, longer than a Kurzgeschichte.", Tags: []string{"b2", "c1", "literature", "noun"}},
		{ID: "lit-die-kurzgeschichte", DeckID: "philosophy-literature", Front: "die Kurzgeschichte", Back: "short story", Extra: "Popular form in post-war German literature.", Tags: []string{"b2", "c1", "literature", "noun"}},
		{ID: "lit-das-drama", DeckID: "philosophy-literature", Front: "das Drama", Back: "drama/play", Extra: "German drama tradition: Schiller, Goethe, Brecht.", Tags: []string{"b2", "c1", "literature", "noun"}},
		{ID: "lit-das-gedicht", DeckID: "philosophy-literature", Front: "das Gedicht", Back: "poem", Extra: "Lyric poetry (Lyrik) is a major German literary tradition.", Tags: []string{"b2", "c1", "literature", "noun"}},
		{ID: "lit-die-lyrik", DeckID: "philosophy-literature", Front: "die Lyrik", Back: "lyric poetry", Extra: "One of the three main literary genres.", Tags: []string{"b2", "c1", "literature", "noun"}},
		{ID: "lit-die-epik", DeckID: "philosophy-literature", Front: "die Epik", Back: "epic/narrative literature", Extra: "Prose fiction genre.", Tags: []string{"b2", "c1", "literature", "noun"}},
		{ID: "lit-die-tragodie", DeckID: "philosophy-literature", Front: "die Tragödie", Back: "tragedy", Extra: "Classical form adapted by Schiller and Goethe.", Tags: []string{"b2", "c1", "literature", "noun"}},
		{ID: "lit-die-komodie", DeckID: "philosophy-literature", Front: "die Komödie", Back: "comedy", Extra: "Brecht's 'Dreigroschenoper' is a famous example.", Tags: []string{"b2", "c1", "literature", "noun"}},
		{ID: "lit-der-protagonist", DeckID: "philosophy-literature", Front: "der Protagonist", Back: "protagonist", Extra: "Main character in a narrative.", Tags: []string{"b2", "c1", "literature", "noun"}},
		{ID: "lit-der-antagonist", DeckID: "philosophy-literature", Front: "der Antagonist", Back: "antagonist", Extra: "Opposing character to the protagonist.", Tags: []string{"b2", "c1", "literature", "noun"}},
		{ID: "lit-die-handlung", DeckID: "philosophy-literature", Front: "die Handlung", Back: "plot/action", Extra: "The sequence of events in a narrative.", Tags: []string{"b2", "c1", "literature", "noun"}},
		{ID: "lit-der-schauplatz", DeckID: "philosophy-literature", Front: "der Schauplatz", Back: "setting/scene", Extra: "Where and when the story takes place.", Tags: []string{"b2", "c1", "literature", "noun"}},
		{ID: "lit-die-metapher", DeckID: "philosophy-literature", Front: "die Metapher", Back: "metaphor", Extra: "A figure of speech making an implicit comparison.", Tags: []string{"b2", "c1", "literature", "noun"}},
		{ID: "lit-die-symbolik", DeckID: "philosophy-literature", Front: "die Symbolik", Back: "symbolism", Extra: "Use of symbols to represent ideas.", Tags: []string{"b2", "c1", "literature", "noun"}},
		{ID: "lit-die-ironie", DeckID: "philosophy-literature", Front: "die Ironie", Back: "irony", Extra: "Saying the opposite of what one means.", Tags: []string{"b2", "c1", "literature", "noun"}},
		{ID: "lit-der-erzahler", DeckID: "philosophy-literature", Front: "der Erzähler", Back: "narrator", Extra: "Person or voice telling the story.", Tags: []string{"b2", "c1", "literature", "noun"}},
		{ID: "lit-die-perspektive", DeckID: "philosophy-literature", Front: "die Perspektive", Back: "perspective/point of view", Extra: "Ich-Erzähler = first-person narrator.", Tags: []string{"b2", "c1", "literature", "noun"}},
		{ID: "lit-das-motiv", DeckID: "philosophy-literature", Front: "das Motiv", Back: "motif", Extra: "Recurring theme or element in literature.", Tags: []string{"b2", "c1", "literature", "noun"}},

		// --- Famous German Authors ---
		{ID: "lit-goethe", DeckID: "philosophy-literature", Front: "Johann Wolfgang von Goethe (1749-1832)", Back: "Faust; Die Leiden des jungen Werthers; Wilhelm Meister", Extra: "Germany's greatest literary figure, Weimar Classicism.", Tags: []string{"b2", "c1", "literature", "author"}},
		{ID: "lit-schiller", DeckID: "philosophy-literature", Front: "Friedrich Schiller (1759-1805)", Back: "Die Räuber; Wilhelm Tell; An die Freude", Extra: "Dramatist and poet, friend of Goethe.", Tags: []string{"b2", "c1", "literature", "author"}},
		{ID: "lit-thomas-mann", DeckID: "philosophy-literature", Front: "Thomas Mann (1875-1955)", Back: "Buddenbrooks; Der Zauberberg; Doktor Faustus", Extra: "Nobel Prize 1929. Exile in USA.", Tags: []string{"b2", "c1", "literature", "author"}},
		{ID: "lit-kafka", DeckID: "philosophy-literature", Front: "Franz Kafka (1883-1924)", Back: "Die Verwandlung; Der Prozess; Das Schloss", Extra: "Prague-born, wrote in German. 'Kafkaesque'.", Tags: []string{"b2", "c1", "literature", "author"}},
		{ID: "lit-hermann-hesse", DeckID: "philosophy-literature", Front: "Hermann Hesse (1877-1962)", Back: "Steppenwolf; Siddhartha; Das Glasperlenspiel", Extra: "Nobel Prize 1946. Themes of self-discovery.", Tags: []string{"b2", "c1", "literature", "author"}},
		{ID: "lit-brecht", DeckID: "philosophy-literature", Front: "Bertolt Brecht (1898-1956)", Back: "Die Dreigroschenoper; Mutter Courage; Episches Theater", Extra: "Revolutionized modern theater.", Tags: []string{"b2", "c1", "literature", "author"}},
		{ID: "lit-remarque", DeckID: "philosophy-literature", Front: "Erich Maria Remarque (1898-1970)", Back: "Im Westen nichts Neues (All Quiet on the Western Front)", Extra: "Anti-war novel, 1929. Exile literature.", Tags: []string{"b2", "c1", "literature", "author"}},
		{ID: "lit-heinrich-boll", DeckID: "philosophy-literature", Front: "Heinrich Böll (1917-1985)", Back: "Die verlorene Ehre der Katharina Blum; Gruppenbild mit Dame", Extra: "Nobel Prize 1972. Trümmerliteratur (rubble literature).", Tags: []string{"b2", "c1", "literature", "author"}},
		{ID: "lit-grass", DeckID: "philosophy-literature", Front: "Günter Grass (1927-2015)", Back: "Die Blechtrommel (The Tin Drum)", Extra: "Nobel Prize 1999. Danzig Trilogy.", Tags: []string{"b2", "c1", "literature", "author"}},

		// --- Literary Periods ---
		{ID: "lit-sturm-und-drang", DeckID: "philosophy-literature", Front: "Sturm und Drang", Back: "Storm and Stress (1767-1785)", Extra: "Youth movement: Goethe, Schiller. Emotion over reason.", Tags: []string{"b2", "c1", "literature", "period"}},
		{ID: "lit-weimarer-klassik", DeckID: "philosophy-literature", Front: "Weimarer Klassik", Back: "Weimar Classicism (1786-1832)", Extra: "Goethe and Schiller in Weimar. Harmony and beauty.", Tags: []string{"b2", "c1", "literature", "period"}},
		{ID: "lit-romantik", DeckID: "philosophy-literature", Front: "die Romantik", Back: "Romanticism (1795-1848)", Extra: "Novalis, E.T.A. Hoffmann, Tieck. Imagination and nature.", Tags: []string{"b2", "c1", "literature", "period"}},
		{ID: "lit-realismus", DeckID: "philosophy-literature", Front: "der Realismus", Back: "Realism (1848-1890)", Extra: "Fontane, Keller. Detailed portrayal of everyday life.", Tags: []string{"b2", "c1", "literature", "period"}},
		{ID: "lit-naturalismus", DeckID: "philosophy-literature", Front: "der Naturalismus", Back: "Naturalism (1880-1900)", Extra: "Hauptmann. Scientific observation of life.", Tags: []string{"b2", "c1", "literature", "period"}},
		{ID: "lit-expressionismus", DeckID: "philosophy-literature", Front: "der Expressionismus", Back: "Expressionism (1910-1925)", Extra: "Distorted reality, inner experience. Trakl, Heym.", Tags: []string{"b2", "c1", "literature", "period"}},
		{ID: "lit-exilliteratur", DeckID: "philosophy-literature", Front: "die Exilliteratur", Back: "exile literature (1933-1945)", Extra: "Thomas Mann, Brecht, Remarque fled Nazi Germany.", Tags: []string{"b2", "c1", "literature", "period"}},
		{ID: "lit-trummerliteratur", DeckID: "philosophy-literature", Front: "die Trümmerliteratur", Back: "rubble literature (1945-1950)", Extra: "Post-war literature: Böll, Borchert. Germany in ruins.", Tags: []string{"b2", "c1", "literature", "period"}},

		// --- Famous Quotes ---
		{ID: "lit-quote-goethe", DeckID: "philosophy-literature", Front: "Edel sei der Mensch, hilfreich und gut!", Back: "Noble be man, helpful and good!", Extra: "Goethe, 'Das Göttliche' (1783).", Tags: []string{"b2", "c1", "literature", "quote"}},
		{ID: "lit-quote-schiller", DeckID: "philosophy-literature", Front: "Alle Menschen werden Brüder.", Back: "All men become brothers.", Extra: "Schiller, 'An die Freude' (1785), set by Beethoven.", Tags: []string{"b2", "c1", "literature", "quote"}},
		{ID: "phil-quote-kant", DeckID: "philosophy-literature", Front: "Handle nur nach derjenigen Maxime, durch die du zugleich wollen kannst, dass sie ein allgemeines Gesetz werde.", Back: "Act only according to that maxim by which you can at the same time will that it should become a universal law.", Extra: "Kant's categorical imperative, Groundwork (1785).", Tags: []string{"b2", "c1", "philosophy", "quote"}},
		{ID: "phil-quote-nietzsche", DeckID: "philosophy-literature", Front: "Was mich nicht umbringt, macht mich stärker.", Back: "What does not kill me makes me stronger.", Extra: "Nietzsche, 'Götzen-Dämmerung' (1889).", Tags: []string{"b2", "c1", "philosophy", "quote"}},
		{ID: "phil-quote-heidegger", DeckID: "philosophy-literature", Front: "Der Mensch ist das Seiende, das sich in seinem Sein zu diesem Sein verhält.", Back: "Man is the entity which relates to its being.", Extra: "Heidegger, 'Sein und Zeit' (1927).", Tags: []string{"b2", "c1", "philosophy", "quote"}},

		// --- Literature Verbs ---
		{ID: "lit-verfassen", DeckID: "philosophy-literature", Front: "verfassen", Back: "to compose/write", Extra: "Formal: ein Werk verfassen.", Tags: []string{"b2", "c1", "literature", "verb"}},
		{ID: "lit-erschaffen", DeckID: "philosophy-literature", Front: "erschaffen", Back: "to create", Extra: "Künstler erschaffen Kunstwerke.", Tags: []string{"b2", "c1", "literature", "verb"}},
		{ID: "lit-deuten", DeckID: "philosophy-literature", Front: "deuten", Back: "to interpret/indicate", Extra: "Ein Text wird gedeutet.", Tags: []string{"b2", "c1", "literature", "verb"}},
		{ID: "lit-interpretieren", DeckID: "philosophy-literature", Front: "interpretieren", Back: "to interpret", Extra: "Ein Gedicht interpretieren.", Tags: []string{"b2", "c1", "literature", "verb"}},
		{ID: "lit-umgestalten", DeckID: "philosophy-literature", Front: "umgestalten", Back: "to transform/reshape", Extra: "Die Gesellschaft umgestalten.", Tags: []string{"b2", "c1", "literature", "verb"}},
		{ID: "lit-kritisieren", DeckID: "philosophy-literature", Front: "kritisieren", Back: "to criticize", Extra: "Gesellschaft kritisieren.", Tags: []string{"b2", "c1", "literature", "verb"}},

		// --- German Intellectual Concepts ---
		{ID: "intel-bildung", DeckID: "philosophy-literature", Front: "die Bildung", Back: "education/formation/self-cultivation", Extra: "Uniquely German concept: holistic personal development.", Tags: []string{"b2", "c1", "culture", "noun"}},
		{ID: "intel-zeitgeist", DeckID: "philosophy-literature", Front: "der Zeitgeist", Back: "spirit of the times", Extra: "Hegel popularized this concept.", Tags: []string{"b2", "c1", "culture", "noun"}},
		{ID: "intel-weltschmerz", DeckID: "philosophy-literature", Front: "der Weltschmerz", Back: "world-weariness/world-pain", Extra: "Romantic-era melancholy about the world's imperfections.", Tags: []string{"b2", "c1", "culture", "noun"}},
		{ID: "intel-gemutlichkeit", DeckID: "philosophy-literature", Front: "die Gemütlichkeit", Back: "coziness/comfort", Extra: "Warm, friendly atmosphere. No direct English equivalent.", Tags: []string{"b2", "c1", "culture", "noun"}},
		{ID: "intel-fingerspitzengefuhl", DeckID: "philosophy-literature", Front: "das Fingerspitzengefühl", Back: "intuitive flair/tact", Extra: "Literally 'fingertip feeling'. Tactful sensitivity.", Tags: []string{"b2", "c1", "culture", "noun"}},
		{ID: "intel-schadenfreude", DeckID: "philosophy-literature", Front: "die Schadenfreude", Back: "malicious joy/gloating", Extra: "Taking pleasure in others' misfortune. Borrowed into English.", Tags: []string{"b2", "c1", "culture", "noun"}},
		{ID: "intel-angst", DeckID: "philosophy-literature", Front: "die Angst", Back: "anxiety/dread", Extra: "Existential concept in Heidegger and Kierkegaard.", Tags: []string{"b2", "c1", "culture", "noun"}},
		{ID: "intel-heimat", DeckID: "philosophy-literature", Front: "die Heimat", Back: "homeland/home", Extra: "Deep emotional attachment to one's region. Untranslatable.", Tags: []string{"b2", "c1", "culture", "noun"}},
		{ID: "intel-wanderlust", DeckID: "philosophy-literature", Front: "die Wanderlust", Back: "desire to travel/roam", Extra: "Strong urge to travel. Borrowed into English.", Tags: []string{"b2", "c1", "culture", "noun"}},
	}

	for i := range notes {
		notes[i].Cards = CardsForNote(notes[i])
	}

	return core.Deck{
		ID:          "philosophy-literature",
		Name:        "Philosophy & Literature (B2-C1)",
		Description: "Advanced German vocabulary for philosophy, literary analysis, intellectual history, and famous German thinkers.",
		Tags:        []string{"german", "b2", "c1", "philosophy", "literature", "culture"},
		Notes:       notes,
	}
}
