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
