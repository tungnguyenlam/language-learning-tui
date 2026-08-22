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

	rawText := strings.TrimSpace(request.SourceText)
	if err := validateDraftRequest(request); err != nil {
		return nil, err
	}
	deckID := strings.TrimSpace(request.DeckID)

	var drafts []Draft
	for _, topic := range splitDraftTopics(rawText) {
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
		drafts = append(drafts, Draft{Note: note})
	}

	if len(drafts) == 0 {
		return nil, errors.New("no valid topics found in input")
	}

	return drafts, nil
}

type TemplateProvider struct {
	Templates map[string]map[string]string
	ActiveSet string
}

func (p TemplateProvider) GenerateDrafts(ctx context.Context, request DraftRequest) ([]Draft, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	if err := validateDraftRequest(request); err != nil {
		return nil, err
	}
	deckID := strings.TrimSpace(request.DeckID)

	activeSet := p.ActiveSet
	if activeSet == "" {
		for k := range p.Templates {
			activeSet = k
			break
		}
	}

	var drafts []Draft
	for _, topic := range splitDraftTopics(strings.TrimSpace(request.SourceText)) {

		front := p.applyTemplate(activeSet, "front", topic, topic)
		back := p.applyTemplate(activeSet, "back", topic, fmt.Sprintf("German prompt for %s", topic))
		example := p.applyTemplate(activeSet, "example", topic, fmt.Sprintf("Practice sentence using %s.", topic))

		idBase := draftIDBase(topic)
		tags := append([]string{"ai-draft", activeSet}, request.Tags...)
		note := core.Note{
			ID:       fmt.Sprintf("ai-%s", idBase),
			DeckID:   deckID,
			Front:    front,
			Back:     back,
			Extra:    fmt.Sprintf("Template draft (%s). Review before keeping.", activeSet),
			Tags:     tags,
			Examples: []string{example},
		}
		note.Cards = content.CardsForNote(note)
		drafts = append(drafts, Draft{Note: note})
	}

	if len(drafts) == 0 {
		return nil, errors.New("no valid topics found in input")
	}

	return drafts, nil
}

func (p TemplateProvider) applyTemplate(set, key, topic, fallback string) string {
	setMap, ok := p.Templates[set]
	if !ok {
		return fallback
	}
	tmpl, ok := setMap[key]
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
		switch r {
		case 'ä':
			b.WriteString("ae")
			lastDash = false
		case 'ö':
			b.WriteString("oe")
			lastDash = false
		case 'ü':
			b.WriteString("ue")
			lastDash = false
		case 'ß':
			b.WriteString("ss")
			lastDash = false
		default:
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

func validateDraftRequest(request DraftRequest) error {
	if strings.TrimSpace(request.SourceText) == "" {
		return errors.New("draft source text is required")
	}
	if strings.TrimSpace(request.DeckID) == "" {
		return errors.New("draft deck id is required")
	}
	return nil
}

func splitDraftTopics(raw string) []string {
	parts := strings.FieldsFunc(raw, func(r rune) bool {
		return r == ',' || r == '\n' || r == '\r'
	})
	topics := make([]string, 0, len(parts))
	for _, topic := range parts {
		if topic = strings.TrimSpace(topic); topic != "" {
			topics = append(topics, topic)
		}
	}
	return topics
}
