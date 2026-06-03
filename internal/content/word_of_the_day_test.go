package content

import "testing"

func TestGetWordOfTheDay(t *testing.T) {
	word := GetWordOfTheDay()
	if word.German == "" {
		t.Fatal("GetWordOfTheDay returned empty German word")
	}
	if word.English == "" {
		t.Fatal("GetWordOfTheDay returned empty English translation")
	}
	if word.Article == "" {
		t.Fatal("GetWordOfTheDay returned empty article")
	}
	if word.Example == "" {
		t.Fatal("GetWordOfTheDay returned empty example")
	}
}

func TestGetWordOfTheDayConsistency(t *testing.T) {
	w1 := GetWordOfTheDay()
	w2 := GetWordOfTheDay()
	if w1.German != w2.German {
		t.Fatalf("GetWordOfTheDay should return the same word within the same day: %s vs %s", w1.German, w2.German)
	}
}

func TestAllWordsOfDayHaveValidData(t *testing.T) {
	for i, w := range wordsOfDay {
		if w.German == "" {
			t.Errorf("word %d has empty German", i)
		}
		if w.English == "" {
			t.Errorf("word %d has empty English", i)
		}
		if w.Article == "" {
			t.Errorf("word %d has empty Article", i)
		}
		if w.Example == "" {
			t.Errorf("word %d has empty Example", i)
		}
	}
}
