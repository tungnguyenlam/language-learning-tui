package content

type ConjunctionExercise struct {
	Sentence    string
	Answer      string
	Meaning     string
	Explanation string
}

func GetConjunctionExercises() []ConjunctionExercise {
	return []ConjunctionExercise{
		{
			Sentence:    "Ich bleibe heute zu Hause, {{...}} es regnet.",
			Answer:      "weil",
			Meaning:     "because",
			Explanation: "weil (because) is a subordinating conjunction: it sends the verb (regnet) to the end of the clause.",
		},
		{
			Sentence:    "Er lernt Deutsch, {{...}} er in Deutschland leben möchte.",
			Answer:      "weil",
			Meaning:     "because",
			Explanation: "weil (because) is a subordinating conjunction: it sends the verb (möchte) to the end of the clause.",
		},
		{
			Sentence:    "Es regnet, {{...}} ich gehe trotzdem spazieren.",
			Answer:      "aber",
			Meaning:     "but",
			Explanation: "aber (but) is a coordinating conjunction (ADUSO): it occupies position 0 and does not change the word order (verb 'gehe' is in position 2).",
		},
		{
			Sentence:    "Ich gehe ins Bett, {{...}} ich bin sehr müde.",
			Answer:      "denn",
			Meaning:     "because",
			Explanation: "denn (because/for) is a coordinating conjunction (ADUSO): it occupies position 0 and does not change the word order (verb 'bin' is in position 2).",
		},
		{
			Sentence:    "Es ist kalt, {{...}} ziehe ich einen Mantel an.",
			Answer:      "deshalb",
			Meaning:     "therefore",
			Explanation: "deshalb (therefore) is a conjunctive adverb: it triggers inversion (verb 'ziehe' comes immediately after it).",
		},
		{
			Sentence:    "Sie weiß nicht, {{...}} er morgen kommt.",
			Answer:      "ob",
			Meaning:     "if/whether",
			Explanation: "ob (whether) is a subordinating conjunction: it sends the verb (kommt) to the end of the clause.",
		},
		{
			Sentence:    "Ich rufe dich an, {{...}} ich am Bahnhof ankomme.",
			Answer:      "wenn",
			Meaning:     "when/if",
			Explanation: "wenn (when/if) is a subordinating conjunction: it sends the verb (ankomme) to the end of the clause.",
		},
		{
			Sentence:    "Er hat viel gelernt, {{...}} hat er die Prüfung nicht bestanden.",
			Answer:      "trotzdem",
			Meaning:     "nevertheless",
			Explanation: "trotzdem (nevertheless) is a conjunctive adverb: it triggers inversion (verb 'hat' comes immediately after it).",
		},
		{
			Sentence:    "Wir essen Pizza {{...}} wir gehen ins Restaurant.",
			Answer:      "oder",
			Meaning:     "or",
			Explanation: "oder (or) is a coordinating conjunction (ADUSO): it occupies position 0 and does not change the word order.",
		},
		{
			Sentence:    "Ich hoffe, {{...}} du dich bald besser fühlst.",
			Answer:      "dass",
			Meaning:     "that",
			Explanation: "dass (that) is a subordinating conjunction: it sends the verb (fühlst) to the end of the clause.",
		},
		{
			Sentence:    "Sie geht zum Deutschkurs, {{...}} sie möchte schneller lernen.",
			Answer:      "denn",
			Meaning:     "because",
			Explanation: "denn (because/for) is a coordinating conjunction (ADUSO): it does not change the word order (verb 'möchte' remains in position 2).",
		},
		{
			Sentence:    "Er hat kein Auto, {{...}} fährt er mit dem Fahrrad.",
			Answer:      "deshalb",
			Meaning:     "therefore",
			Explanation: "deshalb (therefore) is a conjunctive adverb: it triggers inversion (verb 'fährt' immediately follows it).",
		},
		{
			Sentence:    "Wir bleiben drinnen, {{...}} der Sturm vorbei ist.",
			Answer:      "bis",
			Meaning:     "until",
			Explanation: "bis (until) is a subordinating conjunction: it sends the verb (ist) to the end of the clause.",
		},
		{
			Sentence:    "Ich trinke keinen Kaffee, {{...}} lieber Tee.",
			Answer:      "sondern",
			Meaning:     "but rather",
			Explanation: "sondern (but rather) is a coordinating conjunction (ADUSO): it introduces a correction and does not change the word order.",
		},
		{
			Sentence:    "Obwohl es schneit, {{...}} die Kinder draußen.",
			Answer:      "spielen",
			Meaning:     "are playing",
			Explanation: "Since 'obwohl' is a subordinating conjunction starting the sentence, the entire subordinating clause occupies position 1, forcing the main clause verb (spielen) to position 2 (immediately after the comma).",
		},
	}
}
