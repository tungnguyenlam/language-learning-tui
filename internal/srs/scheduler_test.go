package srs

import (
	"testing"
	"time"

	"deutsch-tui/internal/core"
)

func TestReviewSchedulesFutureDueDate(t *testing.T) {
	now := time.Date(2026, 4, 29, 10, 0, 0, 0, time.UTC)
	state := core.NewReviewState("card-1", now)

	next, err := NewScheduler(nil).Review(state, core.GradeGood, now)
	if err != nil {
		t.Fatalf("review failed: %v", err)
	}
	if !next.Due.After(now) {
		t.Fatalf("due date should be in the future: %s", next.Due)
	}
	if next.Reviews != 1 {
		t.Fatalf("reviews = %d, want 1", next.Reviews)
	}
}

func TestReviewIntervals(t *testing.T) {
	now := time.Date(2026, 4, 29, 10, 0, 0, 0, time.UTC)
	s := NewScheduler(nil)

	// First review: Good
	state := core.NewReviewState("card-1", now)
	next, _ := s.Review(state, core.GradeGood, now)
	interval1 := next.Interval

	// Second review: Good (1 day later)
	next2, _ := s.Review(next, core.GradeGood, next.Due)
	interval2 := next2.Interval

	if interval2 <= interval1 {
		t.Errorf("interval should increase: %v -> %v", interval1, interval2)
	}

	// Grade Again should reset or significantly decrease interval
	next3, _ := s.Review(next2, core.GradeAgain, next2.Due)
	if next3.Interval >= interval2 {
		t.Errorf("interval should decrease on GradeAgain: %v -> %v", interval2, next3.Interval)
	}
	if next3.Lapses != 1 {
		t.Errorf("lapses should be 1, got %d", next3.Lapses)
	}
}

func TestPredict(t *testing.T) {
	now := time.Date(2026, 4, 29, 10, 0, 0, 0, time.UTC)
	s := NewScheduler(nil)
	state := core.NewReviewState("card-1", now)

	predictions := s.Predict(state, now)
	if len(predictions) != 4 {
		t.Errorf("want 4 predictions, got %d", len(predictions))
	}

	again := predictions[core.GradeAgain]
	easy := predictions[core.GradeEasy]

	if easy <= again {
		t.Errorf("Easy interval (%v) should be longer than Again (%v)", easy, again)
	}
}

func TestIntraDayIntervalNonZero(t *testing.T) {
	now := time.Date(2026, 4, 29, 10, 0, 0, 0, time.UTC)
	s := NewScheduler(nil)
	state := core.NewReviewState("card-intraday", now)

	next, err := s.Review(state, core.GradeAgain, now)
	if err != nil {
		t.Fatalf("review error: %v", err)
	}
	if next.Interval <= 0 {
		t.Errorf("expected non-zero intra-day interval on GradeAgain, got %v", next.Interval)
	}
}

func TestPredictDurationOverflowClamped(t *testing.T) {
	now := time.Date(2026, 4, 29, 10, 0, 0, 0, time.UTC)
	s := NewScheduler(nil)
	state := core.ReviewState{
		CardID:     "card-huge",
		Due:        now.AddDate(100, 0, 0),
		Interval:   time.Duration(100) * 365 * 24 * time.Hour,
		Stability:  100000.0,
		Difficulty: 1.0,
		Reviews:    50,
	}

	predictions := s.Predict(state, now)
	for grade, interval := range predictions {
		if interval < 0 {
			t.Errorf("prediction for %s resulted in negative interval %v", grade, interval)
		}
	}
}
