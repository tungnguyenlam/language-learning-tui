package content

import (
	"strings"
	"testing"
)

func TestGetPassiveExercises(t *testing.T) {
	exercises := GetPassiveExercises()
	if len(exercises) == 0 {
		t.Fatal("Expected non-empty passive exercises slice")
	}

	for i, ex := range exercises {
		if ex.Sentence == "" {
			t.Errorf("Exercise %d: empty sentence", i)
		}
		if !strings.Contains(ex.Sentence, "{{...}}") {
			t.Errorf("Exercise %d: sentence missing blank marker {{...}}: %q", i, ex.Sentence)
		}
		if ex.Answer == "" {
			t.Errorf("Exercise %d: empty answer", i)
		}
		if ex.Meaning == "" {
			t.Errorf("Exercise %d: empty meaning", i)
		}
		if ex.Hint == "" {
			t.Errorf("Exercise %d: empty hint", i)
		}
		if ex.Explanation == "" {
			t.Errorf("Exercise %d: empty explanation", i)
		}
	}
}
