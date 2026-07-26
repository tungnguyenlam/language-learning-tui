package content

import (
	"testing"
)

func TestGetRelativeExercises(t *testing.T) {
	exs := GetRelativeExercises()
	if len(exs) < 25 {
		t.Fatalf("Expected at least 25 relative exercises, got %d", len(exs))
	}
	for i, ex := range exs {
		if ex.Sentence == "" {
			t.Errorf("Exercise %d: empty sentence", i)
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
