package srs

import (
	"errors"
	"time"

	"deutsch-tui/internal/core"

	fsrs "github.com/open-spaced-repetition/go-fsrs/v3"
)

type Scheduler struct {
	fsrs *fsrs.FSRS
}

func NewScheduler() Scheduler {
	return Scheduler{fsrs: fsrs.NewFSRS(fsrs.DefaultParam())}
}

func (s Scheduler) Review(state core.ReviewState, grade core.ReviewGrade, now time.Time) (core.ReviewState, error) {
	if state.CardID == "" {
		return core.ReviewState{}, errors.New("card id is required")
	}

	rating, err := ratingForGrade(grade)
	if err != nil {
		return core.ReviewState{}, errors.New("unsupported review grade")
	}

	info := s.fsrs.Next(fsrsCardFromState(state, now), now, rating)
	next := stateFromFSRS(state.CardID, info.Card, now)
	return next, nil
}

func (s Scheduler) Predict(state core.ReviewState, now time.Time) map[core.ReviewGrade]time.Duration {
	card := fsrsCardFromState(state, now)
	schedulingCards := s.fsrs.Repeat(card, now)

	predictions := make(map[core.ReviewGrade]time.Duration)
	for rating, sc := range schedulingCards {
		grade := gradeForRating(rating)
		if grade == "" {
			continue
		}
		interval := time.Duration(sc.Card.ScheduledDays) * 24 * time.Hour
		if sc.Card.Due.After(now) && (interval == 0 || sc.Card.Due.Sub(now) < interval) {
			interval = sc.Card.Due.Sub(now)
		}
		predictions[grade] = interval
	}

	return predictions
}

func gradeForRating(rating fsrs.Rating) core.ReviewGrade {
	switch rating {
	case fsrs.Again:
		return core.GradeAgain
	case fsrs.Hard:
		return core.GradeHard
	case fsrs.Good:
		return core.GradeGood
	case fsrs.Easy:
		return core.GradeEasy
	default:
		return ""
	}
}

func ratingForGrade(grade core.ReviewGrade) (fsrs.Rating, error) {
	switch grade {
	case core.GradeAgain:
		return fsrs.Again, nil
	case core.GradeHard:
		return fsrs.Hard, nil
	case core.GradeGood:
		return fsrs.Good, nil
	case core.GradeEasy:
		return fsrs.Easy, nil
	default:
		return 0, errors.New("unsupported review grade")
	}
}

func fsrsCardFromState(state core.ReviewState, now time.Time) fsrs.Card {
	card := fsrs.NewCard()
	card.Due = state.Due
	card.Stability = state.Stability
	card.Difficulty = state.Difficulty
	card.Reps = uint64(maxInt(0, state.Reviews))
	card.Lapses = uint64(maxInt(0, state.Lapses))
	card.ScheduledDays = uint64(maxInt(0, int(state.Interval.Hours()/24)))
	card.LastReview = state.LastReview
	if card.LastReview.IsZero() && state.Interval > 0 {
		card.LastReview = now.Add(-state.Interval)
	}
	if state.Reviews > 0 {
		card.State = fsrs.Review
	}
	return card
}

func stateFromFSRS(cardID string, card fsrs.Card, now time.Time) core.ReviewState {
	days := card.ScheduledDays
	if days > 36500 { // clamp to 100 years to prevent time.Duration overflow
		days = 36500
	}
	interval := time.Duration(days) * 24 * time.Hour
	if card.Due.After(now) && interval == 0 {
		interval = card.Due.Sub(now)
	}
	return core.ReviewState{
		CardID:     cardID,
		Due:        card.Due,
		LastReview: card.LastReview,
		Interval:   interval,
		Stability:  card.Stability,
		Difficulty: card.Difficulty,
		Ease:       card.Difficulty,
		Reviews:    int(card.Reps),
		Lapses:     int(card.Lapses),
	}
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
