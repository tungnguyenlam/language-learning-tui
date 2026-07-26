package content

// PassiveExercise is a fill-in-the-blank exercise for the Passive Voice
// (Passiv) trainer. Answers are the correct passive voice form.
type PassiveExercise struct {
	Sentence    string // fill-in-the-blank sentence using {{...}}
	Answer      string // correct passive form
	Meaning     string // English meaning of the full sentence
	Hint        string // short grammar hint (Vorgangspassiv / Zustandspassiv / modal)
	Explanation string // full explanation shown on reveal
}

// GetPassiveExercises returns all Passive Voice exercises.
// Exercises cover:
//   - Vorgangspassiv (werden + Partizip II): process passive
//   - Zustandspassiv (sein + Partizip II): state passive
//   - Modal passive (modal verb + Partizip II + werden)
//   - Passive with agent (von + Dative / durch + Accusative)
//   - Passive in different tenses (Präsens, Präteritum, Perfekt, Futur)
//   - Passive alternatives (man, sich lassen, -bar adjectives)
func GetPassiveExercises() []PassiveExercise {
	return []PassiveExercise{
		// --- Vorgangspassiv (Präsens) ---
		{
			Sentence:    "Das Buch {{...}} von vielen Leuten gelesen.",
			Answer:      "wird",
			Meaning:     "The book is read by many people.",
			Hint:        "Vorgangspassiv Präsens: werden + Partizip II",
			Explanation: "Vorgangspassiv (process passive) uses 'werden' + Partizip II. Subject 'das Buch' is 3rd person singular → 'wird'.",
		},
		{
			Sentence:    "Die Fenster {{...}} jeden Morgen geöffnet.",
			Answer:      "werden",
			Meaning:     "The windows are opened every morning.",
			Hint:        "Vorgangspassiv Präsens: werden + Partizip II (plural)",
			Explanation: "'Die Fenster' is plural → 'werden' (not 'wird'). The action of opening happens regularly.",
		},
		{
			Sentence:    "Der Brief {{...}} heute geschrieben.",
			Answer:      "wird",
			Meaning:     "The letter is being written today.",
			Hint:        "Vorgangspassiv Präsens: werden + Partizip II",
			Explanation: "'Der Brief' is singular masculine → 'wird'. Vorgangspassiv emphasizes the process/action.",
		},
		{
			Sentence:    "Die Kinder {{...}} von der Lehrerin unterrichtet.",
			Answer:      "werden",
			Meaning:     "The children are taught by the teacher.",
			Hint:        "Vorgangspassiv Präsens with agent (von + Dative)",
			Explanation: "Plural subject 'die Kinder' → 'werden'. The agent is introduced with 'von' + Dative.",
		},
		{
			Sentence:    "Das Essen {{...}} von meiner Mutter gekocht.",
			Answer:      "wird",
			Meaning:     "The food is cooked by my mother.",
			Hint:        "Vorgangspassiv Präsens with agent (von + Dative)",
			Explanation: "'von' introduces the personal agent (who does the action) in the Dative case.",
		},
		// --- Vorgangspassiv (Präteritum) ---
		{
			Sentence:    "Das Auto {{...}} letzte Woche repariert.",
			Answer:      "wurde",
			Meaning:     "The car was repaired last week.",
			Hint:        "Vorgangspassiv Präteritum: wurde/wurden + Partizip II",
			Explanation: "Präteritum passive uses 'wurde' (singular) / 'wurden' (plural) + Partizip II. 'Das Auto' is singular → 'wurde'.",
		},
		{
			Sentence:    "Die Häuser {{...}} im 19. Jahrhundert gebaut.",
			Answer:      "wurden",
			Meaning:     "The houses were built in the 19th century.",
			Hint:        "Vorgangspassiv Präteritum: wurden (plural)",
			Explanation: "'Die Häuser' is plural → 'wurden'. Präteritum passive is the most common tense for narrative passive.",
		},
		{
			Sentence:    "Der Dieb {{...}} von der Polizei gefasst.",
			Answer:      "wurde",
			Meaning:     "The thief was caught by the police.",
			Hint:        "Vorgangspassiv Präteritum with agent",
			Explanation: "Singular subject → 'wurde'. 'von der Polizei' names the agent in the Dative.",
		},
		// --- Vorgangspassiv (Perfekt) ---
		{
			Sentence:    "Die E-Mail {{...}} schon geschickt worden.",
			Answer:      "ist",
			Meaning:     "The email has already been sent.",
			Hint:        "Vorgangspassiv Perfekt: sein + Partizip II + worden",
			Explanation: "Perfekt passive uses 'sein' + Partizip II + 'worden' (not 'geworden'!). Singular → 'ist'.",
		},
		{
			Sentence:    "Die Pakete {{...}} gestern geliefert worden.",
			Answer:      "sind",
			Meaning:     "The packages were delivered yesterday.",
			Hint:        "Vorgangspassiv Perfekt: sein + Partizip II + worden (plural)",
			Explanation: "Plural subject 'die Pakete' → 'sind'. Note: passive Perfekt uses 'worden', not 'geworden'.",
		},
		// --- Zustandspassiv ---
		{
			Sentence:    "Die Tür {{...}} geschlossen.",
			Answer:      "ist",
			Meaning:     "The door is closed (state).",
			Hint:        "Zustandspassiv: sein + Partizip II",
			Explanation: "Zustandspassiv (state passive) uses 'sein' + Partizip II. It describes a resulting state, not a process.",
		},
		{
			Sentence:    "Der Tisch {{...}} schon gedeckt.",
			Answer:      "ist",
			Meaning:     "The table is already set.",
			Hint:        "Zustandspassiv: sein + Partizip II",
			Explanation: "Zustandspassiv emphasizes the result: the table IS set (completed state). Compare: 'Der Tisch wird gedeckt' (is being set — process).",
		},
		{
			Sentence:    "Die Fenster {{...}} geöffnet.",
			Answer:      "sind",
			Meaning:     "The windows are open (state).",
			Hint:        "Zustandspassiv: sein + Partizip II (plural)",
			Explanation: "Plural subject → 'sind'. This describes the current state (windows are open), not the act of opening them.",
		},
		{
			Sentence:    "Das Museum {{...}} montags geschlossen.",
			Answer:      "ist",
			Meaning:     "The museum is closed on Mondays.",
			Hint:        "Zustandspassiv: sein + Partizip II",
			Explanation: "Zustandspassiv describes a state. The museum is in the state of being closed — not the act of closing it.",
		},
		// --- Modal Passive ---
		{
			Sentence:    "Das Formular {{...}} ausgefüllt werden.",
			Answer:      "muss",
			Meaning:     "The form must be filled out.",
			Hint:        "Modal passive: modal verb + Partizip II + werden",
			Explanation: "Modal passive puts the modal verb in conjugated position: 'muss ... ausgefüllt werden'. 'werden' stays as infinitive at the end.",
		},
		{
			Sentence:    "Die Rechnung {{...}} sofort bezahlt werden.",
			Answer:      "muss",
			Meaning:     "The bill must be paid immediately.",
			Hint:        "Modal passive: muss + Partizip II + werden",
			Explanation: "Modal verb 'müssen' is conjugated ('muss'), Partizip II + 'werden' go to the end of the clause.",
		},
		{
			Sentence:    "Das Projekt {{...}} bis Freitag abgeschlossen werden.",
			Answer:      "soll",
			Meaning:     "The project should be completed by Friday.",
			Hint:        "Modal passive: sollen + Partizip II + werden",
			Explanation: "'sollen' in passive: 'soll ... abgeschlossen werden'. The modal verb carries the conjugation.",
		},
		{
			Sentence:    "Die Dateien {{...}} nicht gelöscht werden.",
			Answer:      "dürfen",
			Meaning:     "The files must not be deleted.",
			Hint:        "Modal passive: dürfen + Partizip II + werden (plural)",
			Explanation: "'dürfen nicht' = 'must not'. Plural 'die Dateien' → 'dürfen'. Modal passive: modal + Partizip II + werden.",
		},
		{
			Sentence:    "Der Text {{...}} ins Deutsche übersetzt werden.",
			Answer:      "kann",
			Meaning:     "The text can be translated into German.",
			Hint:        "Modal passive: können + Partizip II + werden",
			Explanation: "'können' expresses ability/possibility in passive. Singular → 'kann'. Structure: modal + ... + Partizip II + werden.",
		},
		// --- Passive with durch ---
		{
			Sentence:    "Die Stadt wurde {{...}} ein Erdbeben zerstört.",
			Answer:      "durch",
			Meaning:     "The city was destroyed by an earthquake.",
			Hint:        "Agent with 'durch' (non-personal cause/force)",
			Explanation: "'durch' + Accusative is used for non-personal agents, forces, or means (earthquake, storm, fire). Compare: 'von' for personal agents.",
		},
		{
			Sentence:    "Die Nachricht wurde {{...}} das Radio verbreitet.",
			Answer:      "durch",
			Meaning:     "The news was spread through the radio.",
			Hint:        "Agent with 'durch' (means/instrument)",
			Explanation: "'durch' is used for means/instruments/methods. 'von' would be used for a person: 'von dem Reporter verbreitet'.",
		},
		// --- Passive alternatives ---
		{
			Sentence:    "Dieses Problem {{...}} sich leicht lösen.",
			Answer:      "lässt",
			Meaning:     "This problem can be easily solved.",
			Hint:        "Passive alternative: sich lassen + Infinitiv",
			Explanation: "'sich lassen + Infinitiv' is a common passive alternative meaning 'can be done'. 'Das lässt sich machen' = 'That can be done'.",
		},
		{
			Sentence:    "Hier {{...}} nicht geraucht werden.",
			Answer:      "darf",
			Meaning:     "Smoking is not allowed here.",
			Hint:        "Impersonal passive with modal verb",
			Explanation: "Impersonal passive (no subject agent): 'es darf hier nicht geraucht werden' or with 'hier' in first position, 'es' is dropped.",
		},
		// --- Futur I Passive ---
		{
			Sentence:    "Das Gebäude {{...}} nächstes Jahr renoviert werden.",
			Answer:      "wird",
			Meaning:     "The building will be renovated next year.",
			Hint:        "Futur I Passive: werden + Partizip II + werden",
			Explanation: "Futur I passive: 'wird' (future auxiliary) + Partizip II + 'werden' (passive). Two forms of 'werden' in one sentence!",
		},
		{
			Sentence:    "Die Ergebnisse {{...}} morgen veröffentlicht werden.",
			Answer:      "werden",
			Meaning:     "The results will be published tomorrow.",
			Hint:        "Futur I Passive: werden + Partizip II + werden (plural)",
			Explanation: "Plural subject → 'werden' as future auxiliary. Structure: werden + ... + Partizip II + werden (passive infinitive).",
		},
	}
}
