package content

import (
	"testing"
)

func TestGetPrepositionExercises(t *testing.T) {
	exercises := GetPrepositionExercises()

	// Original had 11, we added 12, so total should be 23.
	if len(exercises) < 20 {
		t.Errorf("expected at least 20 preposition exercises, got %d", len(exercises))
	}

	// Verify that we have some two-way prepositions in our set
	twoWayPreps := map[string]int{
		"an":       0,
		"hinter":   0,
		"neben":    0,
		"über":     0,
		"vor":      0,
		"zwischen": 0,
	}

	for _, ex := range exercises {
		if _, ok := twoWayPreps[ex.Preposition]; ok {
			twoWayPreps[ex.Preposition]++
		}
	}

	for prep, count := range twoWayPreps {
		if count < 2 {
			t.Errorf("expected at least 2 exercises for two-way preposition %q (dative & accusative), got %d", prep, count)
		}
	}
}
