package ai

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"deutsch-tui/internal/content"
	"deutsch-tui/internal/core"
)

// systemPrompt instructs the model to produce a strict JSON object so we
// can parse it without prose stripping. Kept short and stable so prompt
// caching (where supported) gets reused across requests.
const systemPrompt = `You generate German vocabulary flashcards.

Output ONLY a single JSON object with this exact shape — no prose, no markdown:
{
  "cards": [
    {
      "front": "der Hund",
      "back": "the dog",
      "extra": "noun, masculine; plural: die Hunde",
      "example": "Ich habe einen Hund."
    }
  ]
}

Rules:
- front: the German word or phrase. For nouns include the article (der/die/das).
- back: the English translation.
- extra: 1 short line about grammar/usage (gender, plural for nouns; infinitive and principal parts for verbs).
- example: ONE natural German sentence using the word.
- Generate between 5 and 10 cards for the user's topic.
- If a level (A1/A2/B1/B2/C1/C2) is mentioned, calibrate vocabulary to that level.
- Ensure high-quality, practical vocabulary.`

// userPromptFor produces the per-request user message. We embed the topic
// and any tags so the model can specialise the cards.
func userPromptFor(req DraftRequest) string {
	topic := strings.TrimSpace(req.SourceText)
	var b strings.Builder
	fmt.Fprintf(&b, "Topic: %s", topic)
	if len(req.Tags) > 0 {
		fmt.Fprintf(&b, "\nTags: %s", strings.Join(req.Tags, ", "))
	}
	return b.String()
}

// rawCard mirrors the JSON shape we ask the model for.
type rawCard struct {
	Front   string `json:"front"`
	Back    string `json:"back"`
	Extra   string `json:"extra"`
	Example string `json:"example"`
}

type rawCards struct {
	Cards []rawCard `json:"cards"`
}

// parseCardsJSON pulls the {cards:[...]} object out of the model's response.
// Models sometimes wrap JSON in ```json fences or prepend a sentence; we
// extract the first balanced JSON object and decode that.
func parseCardsJSON(raw string) ([]rawCard, error) {
	body := extractJSON(raw)
	if body == "" {
		return nil, fmt.Errorf("no JSON object found in response: %q", truncate(raw, 200))
	}
	var parsed rawCards
	if err := json.Unmarshal([]byte(body), &parsed); err != nil {
		return nil, fmt.Errorf("decode cards: %w (body=%q)", err, truncate(body, 200))
	}
	if len(parsed.Cards) == 0 {
		return nil, errors.New("model returned no cards")
	}
	return parsed.Cards, nil
}

// extractJSON returns the first balanced {...} block in s, or "".
// It ignores braces inside strings (with backslash escapes) so we don't
// get fooled by example sentences that contain "{".
func extractJSON(s string) string {
	start := strings.Index(s, "{")
	if start < 0 {
		return ""
	}
	depth := 0
	inString := false
	escaped := false
	for i := start; i < len(s); i++ {
		c := s[i]
		if escaped {
			escaped = false
			continue
		}
		if inString {
			switch c {
			case '\\':
				escaped = true
			case '"':
				inString = false
			}
			continue
		}
		switch c {
		case '"':
			inString = true
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				return s[start : i+1]
			}
		}
	}
	return ""
}

func truncate(s string, n int) string {
	runes := []rune(s)
	if len(runes) <= n {
		return s
	}
	return string(runes[:n]) + "…"
}

func generateDraftsViaChat(ctx context.Context, request DraftRequest, prefix string, chat func(context.Context, string, string) (string, error)) ([]Draft, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if err := validateDraftRequest(request); err != nil {
		return nil, err
	}
	text, err := chat(ctx, systemPrompt, userPromptFor(request))
	if err != nil {
		return nil, err
	}
	rawCards, err := parseCardsJSON(text)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", prefix, err)
	}
	return draftsFromRaw(rawCards, request)
}

// draftsFromRaw maps the parsed JSON into validated ai.Draft values.
// Cards with empty front/back are skipped rather than failing the whole
// batch — the model occasionally emits a stub.
func draftsFromRaw(raws []rawCard, req DraftRequest) ([]Draft, error) {
	deckID := strings.TrimSpace(req.DeckID)
	if deckID == "" {
		return nil, errors.New("draft deck id is required")
	}
	var drafts []Draft
	seen := map[string]int{}
	for _, r := range raws {
		front := strings.TrimSpace(r.Front)
		back := strings.TrimSpace(r.Back)
		if front == "" || back == "" {
			continue
		}
		idBase := draftIDBase(front)
		if idBase == "" {
			idBase = "draft"
		}
		// Disambiguate duplicates so ValidateDrafts is happy.
		if n := seen[idBase]; n > 0 {
			idBase = fmt.Sprintf("%s-%d", idBase, n+1)
		}
		seen[r.Front]++
		tags := append([]string{"ai-draft"}, req.Tags...)
		note := core.Note{
			ID:       "ai-" + idBase,
			DeckID:   deckID,
			Front:    front,
			Back:     back,
			Extra:    strings.TrimSpace(r.Extra),
			Tags:     tags,
			Examples: nil,
		}
		if ex := strings.TrimSpace(r.Example); ex != "" {
			note.Examples = []string{ex}
		}
		note.Cards = content.CardsForNote(note)
		drafts = append(drafts, Draft{Note: note})
	}
	if len(drafts) == 0 {
		return nil, errors.New("model returned no usable cards")
	}
	return drafts, nil
}
