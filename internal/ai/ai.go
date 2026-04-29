package ai

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"deutsch-tui/internal/content"
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

type OfflineProvider struct{}

func (p OfflineProvider) GenerateDrafts(ctx context.Context, request DraftRequest) ([]Draft, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	topic := strings.TrimSpace(request.SourceText)
	if topic == "" {
		return nil, errors.New("draft source text is required")
	}
	deckID := strings.TrimSpace(request.DeckID)
	if deckID == "" {
		return nil, errors.New("draft deck id is required")
	}

	idBase := draftIDBase(topic)
	tags := append([]string{"ai-draft"}, request.Tags...)
	note := core.Note{
		ID:       fmt.Sprintf("ai-%s", idBase),
		DeckID:   deckID,
		Front:    topic,
		Back:     fmt.Sprintf("German prompt for %s", topic),
		Extra:    "Offline draft. Review before keeping.",
		Tags:     tags,
		Examples: []string{fmt.Sprintf("Practice sentence using %s.", topic)},
	}
	note.Cards = content.CardsForNote(note)
	return []Draft{{Note: note}}, nil
}

type TemplateProvider struct {
	Templates map[string]string
}

func (p TemplateProvider) GenerateDrafts(ctx context.Context, request DraftRequest) ([]Draft, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	topic := strings.TrimSpace(request.SourceText)
	if topic == "" {
		return nil, errors.New("draft source text is required")
	}
	deckID := strings.TrimSpace(request.DeckID)
	if deckID == "" {
		return nil, errors.New("draft deck id is required")
	}

	front := p.applyTemplate("front", topic, topic)
	back := p.applyTemplate("back", topic, fmt.Sprintf("German prompt for %s", topic))
	example := p.applyTemplate("example", topic, fmt.Sprintf("Practice sentence using %s.", topic))

	idBase := draftIDBase(topic)
	tags := append([]string{"ai-draft"}, request.Tags...)
	note := core.Note{
		ID:       fmt.Sprintf("ai-%s", idBase),
		DeckID:   deckID,
		Front:    front,
		Back:     back,
		Extra:    "Template draft. Review before keeping.",
		Tags:     tags,
		Examples: []string{example},
	}
	note.Cards = content.CardsForNote(note)
	return []Draft{{Note: note}}, nil
}

func (p TemplateProvider) applyTemplate(key, topic, fallback string) string {
	tmpl, ok := p.Templates[key]
	if !ok || tmpl == "" {
		return fallback
	}
	return strings.ReplaceAll(tmpl, "{{.Topic}}", topic)
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

func draftIDBase(topic string) string {
	topic = strings.ToLower(strings.TrimSpace(topic))
	var b strings.Builder
	lastDash := false
	for _, r := range topic {
		switch {
		case r >= 'a' && r <= 'z':
			b.WriteRune(r)
			lastDash = false
		case r >= '0' && r <= '9':
			b.WriteRune(r)
			lastDash = false
		default:
			if !lastDash && b.Len() > 0 {
				b.WriteByte('-')
				lastDash = true
			}
		}
	}
	out := strings.Trim(b.String(), "-")
	if out == "" {
		return "draft"
	}
	return out
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
