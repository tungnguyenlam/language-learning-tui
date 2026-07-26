package content

// KonjunktivExercise is a fill-in-the-blank exercise for the Konjunktiv II
// (Subjunctive II) trainer. Answers are always the correct Konjunktiv II form.
type KonjunktivExercise struct {
	Sentence    string // fill-in-the-blank sentence using {{...}}
	Answer      string // correct Konjunktiv II form
	Meaning     string // English meaning of the full sentence
	Hint        string // short grammar hint (modal / würde / strong verb)
	Explanation string // full explanation shown on reveal
}

// GetKonjunktivExercises returns all Konjunktiv II exercises.
// Exercises cover:
//   - würde + Infinitiv (the productive periphrastic form)
//   - Strong Konjunktiv II of common verbs (wäre, hätte, käme, ginge, …)
//   - Modal verbs in Konjunktiv II (könnte, müsste, dürfte, sollte, wollte)
//   - Conditional/wishful sentences (Wenn …, dann …)
func GetKonjunktivExercises() []KonjunktivExercise {
	return []KonjunktivExercise{
		{
			Sentence:    "Wenn ich mehr Zeit {{...}}, würde ich mehr lesen.",
			Answer:      "hätte",
			Meaning:     "If I had more time, I would read more.",
			Hint:        "Strong Konjunktiv II of haben",
			Explanation: "'haben' → Konjunktiv II: 'hätte'. Used in conditional clauses for unreal/hypothetical situations.",
		},
		{
			Sentence:    "Wenn er reicher {{...}}, würde er ein Haus kaufen.",
			Answer:      "wäre",
			Meaning:     "If he were richer, he would buy a house.",
			Hint:        "Strong Konjunktiv II of sein",
			Explanation: "'sein' → Konjunktiv II: 'wäre'. One of the most common strong forms used in unreal conditions.",
		},
		{
			Sentence:    "Ich {{...}} gerne ins Kino gehen.",
			Answer:      "würde",
			Meaning:     "I would like to go to the cinema.",
			Hint:        "würde + Infinitiv (periphrastic)",
			Explanation: "For most verbs, Konjunktiv II is formed with 'würde' + infinitive. 'würde' is itself the Konjunktiv II of 'werden'.",
		},
		{
			Sentence:    "Das {{...}} ich nicht machen.",
			Answer:      "würde",
			Meaning:     "I wouldn't do that.",
			Hint:        "würde + Infinitiv (periphrastic)",
			Explanation: "'würde' + infinitive expresses hypothetical actions. 'machen' stays as infinitive at the end.",
		},
		{
			Sentence:    "Sie {{...}} kommen, wenn sie nicht so beschäftigt wäre.",
			Answer:      "würde",
			Meaning:     "She would come if she weren't so busy.",
			Hint:        "würde + Infinitiv (periphrastic)",
			Explanation: "The main clause uses 'würde + infinitive' to express an unreal outcome. The wenn-clause uses 'wäre' (sein → Konjunktiv II).",
		},
		{
			Sentence:    "Er {{...}} besser schlafen, wenn er weniger Kaffee tränke.",
			Answer:      "würde",
			Meaning:     "He would sleep better if he drank less coffee.",
			Hint:        "würde + Infinitiv (periphrastic)",
			Explanation: "'würde' + Infinitiv in the main clause. 'tränke' (Konjunktiv II of trinken) appears in the wenn-clause.",
		},
		{
			Sentence:    "Wenn ich du {{...}}, würde ich das Angebot annehmen.",
			Answer:      "wäre",
			Meaning:     "If I were you, I would accept the offer.",
			Hint:        "Strong Konjunktiv II of sein",
			Explanation: "Fixed idiom 'Wenn ich du wäre' — 'wäre' is the Konjunktiv II of 'sein'.",
		},
		{
			Sentence:    "Wir {{...}} früher kommen, aber der Zug hatte Verspätung.",
			Answer:      "hätten",
			Meaning:     "We could have come earlier, but the train was delayed.",
			Hint:        "Konjunktiv II Perfekt (hätten + Partizip II)",
			Explanation: "'hätten' is the Konjunktiv II of 'haben', used with Partizip II ('gekommen') to express an unreal past event.",
		},
		{
			Sentence:    "Das {{...}} er nicht sagen sollen.",
			Answer:      "hätte",
			Meaning:     "He shouldn't have said that.",
			Hint:        "Konjunktiv II Perfekt: hätte + sagen sollen",
			Explanation: "'hätte' + infinitive + 'sollen' = Konjunktiv II with a modal — expresses an unreal obligation in the past.",
		},
		{
			Sentence:    "Ich {{...}} das gerne tun, aber ich habe keine Zeit.",
			Answer:      "würde",
			Meaning:     "I would gladly do that, but I have no time.",
			Hint:        "würde + Infinitiv (polite offer / unfulfilled wish)",
			Explanation: "'würde' + Infinitiv is also used for polite expressions and unfulfilled wishes in present/future.",
		},
		{
			Sentence:    "Wenn es nicht so teuer {{...}}, würde ich es kaufen.",
			Answer:      "wäre",
			Meaning:     "If it weren't so expensive, I would buy it.",
			Hint:        "Strong Konjunktiv II of sein",
			Explanation: "'wäre' (Konjunktiv II of sein) is used here to describe an unreal present property.",
		},
		{
			Sentence:    "Du {{...}} mehr schlafen, dann wärst du nicht so müde.",
			Answer:      "solltest",
			Meaning:     "You should sleep more, then you wouldn't be so tired.",
			Hint:        "Modal Konjunktiv II: sollen → solltest",
			Explanation: "'sollen' → Konjunktiv II 'sollte'. It carries a tone of advice or recommendation.",
		},
		{
			Sentence:    "Man {{...}} öfter Pause machen.",
			Answer:      "sollte",
			Meaning:     "One should take breaks more often.",
			Hint:        "Modal Konjunktiv II: sollen → sollte",
			Explanation: "'sollte' is the Konjunktiv II of 'sollen', used for advice or moral obligation.",
		},
		{
			Sentence:    "Er {{...}} mehr Geld verdienen, wenn er eine bessere Stelle hätte.",
			Answer:      "könnte",
			Meaning:     "He could earn more money if he had a better job.",
			Hint:        "Modal Konjunktiv II: können → könnte",
			Explanation: "'können' → Konjunktiv II 'könnte'. Expresses hypothetical ability.",
		},
		{
			Sentence:    "Ich {{...}} mir vorstellen, in Deutschland zu leben.",
			Answer:      "könnte",
			Meaning:     "I could imagine living in Germany.",
			Hint:        "Modal Konjunktiv II: können → könnte",
			Explanation: "'könnte' is commonly used as a softened, polite expression of possibility or imagination.",
		},
		{
			Sentence:    "Du {{...}} hier nicht rauchen.",
			Answer:      "dürftest",
			Meaning:     "You wouldn't be allowed to smoke here.",
			Hint:        "Modal Konjunktiv II: dürfen → dürftest",
			Explanation: "'dürfen' → Konjunktiv II 'dürfte/dürftest'. Expresses hypothetical permission (often in its negation).",
		},
		{
			Sentence:    "Wenn ich könnte, {{...}} ich jeden Tag reisen.",
			Answer:      "würde",
			Meaning:     "If I could, I would travel every day.",
			Hint:        "würde + Infinitiv (unreal wish)",
			Explanation: "'würde' + Infinitiv in the main clause to express an unreal present wish. The wenn-clause uses 'könnte'.",
		},
		{
			Sentence:    "An deiner Stelle {{...}} ich das Angebot ablehnen.",
			Answer:      "würde",
			Meaning:     "In your position, I would turn down the offer.",
			Hint:        "würde + Infinitiv (advice formulation)",
			Explanation: "'An deiner Stelle würde ich…' is a fixed idiom for giving advice hypothetically using Konjunktiv II.",
		},
		{
			Sentence:    "Ich {{...}} lieber zu Hause bleiben.",
			Answer:      "bliebe",
			Meaning:     "I would rather stay at home.",
			Hint:        "Strong Konjunktiv II of bleiben",
			Explanation: "'bleiben' → Konjunktiv II 'bliebe'. Strong verbs can use either their Konjunktiv II form or 'würde + bleiben'.",
		},
		{
			Sentence:    "Er {{...}} uns helfen, wenn wir ihn darum bitten würden.",
			Answer:      "würde",
			Meaning:     "He would help us if we asked him.",
			Hint:        "würde + Infinitiv (conditional outcome)",
			Explanation: "'würde' + Infinitiv in both main clause and wenn-clause. Using 'würde' in both is acceptable in spoken German.",
		},
		{
			Sentence:    "Das {{...}} wirklich hilfreich sein!",
			Answer:      "wäre",
			Meaning:     "That would really be helpful!",
			Hint:        "Strong Konjunktiv II of sein",
			Explanation: "'sein' → Konjunktiv II 'wäre'. Short exclamatory use of Konjunktiv II for expressing a wish or polite evaluation.",
		},
		{
			Sentence:    "Wenn ich das früher gewusst {{...}}, hätte ich anders gehandelt.",
			Answer:      "hätte",
			Meaning:     "If I had known that earlier, I would have acted differently.",
			Hint:        "Konjunktiv II Perfekt: hätte + Partizip II",
			Explanation: "'hätte' is the Konjunktiv II of 'haben'. 'gewusst' is the Partizip II of 'wissen'. Together they form an unreal past conditional.",
		},
		{
			Sentence:    "Er {{...}} lieber Arzt geworden.",
			Answer:      "wäre",
			Meaning:     "He would rather have become a doctor.",
			Hint:        "Konjunktiv II Perfekt with sein-verb: wäre + Partizip II",
			Explanation: "'wäre' is the Konjunktiv II of 'sein'. Verbs of motion and change-of-state use 'wäre + Partizip II' for unreal past.",
		},
		{
			Sentence:    "Wenn wir früher angefangen {{...}}, wären wir schon fertig.",
			Answer:      "hätten",
			Meaning:     "If we had started earlier, we would already be done.",
			Hint:        "Konjunktiv II Perfekt: hätten + Partizip II",
			Explanation: "'hätten' (Konjunktiv II of haben) + 'angefangen' = unreal past perfect condition. The main clause uses 'wären' (Konjunktiv II of sein).",
		},
		{
			Sentence:    "Das {{...}} ich mir nicht leisten.",
			Answer:      "könnte",
			Meaning:     "I couldn't afford that.",
			Hint:        "Modal Konjunktiv II: können → könnte",
			Explanation: "'könnte' expresses hypothetical inability. 'Das könnte ich mir nicht leisten' is a common colloquial phrase.",
		},
	}
}
