package ai

import (
	"context"
	"errors"
	"strings"

	"deutsch-tui/internal/core"
)

type DraftRequest struct {
	SourceText string
	DeckID     string
	Tags       []string
}

type Draft struct {
	Note core.Note
}

type Provider interface {
	GenerateDrafts(ctx context.Context, request DraftRequest) ([]Draft, error)
}

type FakeProvider struct {
	Drafts []Draft
	Err    error
}

func (p FakeProvider) GenerateDrafts(context.Context, DraftRequest) ([]Draft, error) {
	if p.Err != nil {
		return nil, p.Err
	}
	return append([]Draft(nil), p.Drafts...), nil
}

func ValidateDraft(draft Draft) error {
	note := draft.Note
	if strings.TrimSpace(note.ID) == "" {
		return errors.New("draft note id is required")
	}
	if strings.TrimSpace(note.DeckID) == "" {
		return errors.New("draft deck id is required")
	}
	if strings.TrimSpace(note.Front) == "" {
		return errors.New("draft front is required")
	}
	if strings.TrimSpace(note.Back) == "" {
		return errors.New("draft back is required")
	}
	for _, card := range note.Cards {
		if err := core.ValidateCard(card); err != nil {
			return err
		}
	}
	return nil
}

func ValidateDrafts(drafts []Draft) error {
	seen := map[string]struct{}{}
	for _, draft := range drafts {
		if err := ValidateDraft(draft); err != nil {
			return err
		}
		if _, ok := seen[draft.Note.ID]; ok {
			return errors.New("duplicate draft note id")
		}
		seen[draft.Note.ID] = struct{}{}
	}
	return nil
}
