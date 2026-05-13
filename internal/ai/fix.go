package ai

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"deutsch-tui/internal/core"
)

// FixRequest asks the AI provider to correct a flashcard whose user
// has reported as wrong. The provider returns a *proposed* fix; the TUI
// shows the diff and only persists it if the user accepts.
type FixRequest struct {
	Note       core.Note
	UserReport string
}

// FixedNote is the AI's proposed correction. The TUI maps these fields
// directly onto the existing note (preserving ID/DeckID/Cards), so SRS
// scheduling and history are not affected.
type FixedNote struct {
	Front   string
	Back    string
	Extra   string
	Example string
	Reason  string
}

// systemPromptFix is reused across all providers so prompt caching can
// share the prefix. It pins the model to strict JSON output with a
// fixed shape, so we don't need to fight prose stripping.
const systemPromptFix = `You are reviewing a German vocabulary flashcard that a learner reported as wrong.

Fix any factual, grammatical, or translation errors. Keep the same word/concept being taught — do NOT swap to a different vocabulary item. If the front already has an article (der/die/das), keep one and make sure it is correct.

Output ONLY a single JSON object with this exact shape — no prose, no markdown:
{
  "front": "der Hund",
  "back": "the dog",
  "extra": "noun, masculine; plural: die Hunde",
  "example": "Ich habe einen Hund.",
  "reason": "Brief explanation of what was wrong and what was changed."
}

Rules:
- front: the German word or phrase, with correct article for nouns.
- back: the English translation.
- extra: 1 short line about grammar (gender, plural, tense, case, register).
- example: ONE natural German sentence using the word.
- reason: ONE short sentence explaining the fix; "Looks correct, no change needed." is acceptable.
- Never invent unrelated content. If you cannot improve anything, return the same values.`

// FixCard asks the provider to correct a card and returns the proposed fix.
// The provider's existing GenerateDrafts path is not reused, because we
// want a single object back, not a list. Each transport (OpenAI, Anthropic)
// exposes a lower-level chat call via the SendChat method below.
func FixCard(ctx context.Context, provider Provider, req FixRequest) (FixedNote, error) {
	if provider == nil {
		return FixedNote{}, errors.New("AI provider is disabled; enable one in Settings to fix cards")
	}
	chat, ok := provider.(ChatProvider)
	if !ok {
		// Fall back to template/offline: synthesize a "no-op" fix so the
		// flow still completes deterministically in tests and offline mode.
		return FixedNote{
			Front:   req.Note.Front,
			Back:    req.Note.Back,
			Extra:   req.Note.Extra,
			Example: firstExample(req.Note.Examples),
			Reason:  "Offline/template provider cannot review cards; left unchanged.",
		}, nil
	}
	user := userPromptForFix(req)
	raw, err := chat.SendChat(ctx, systemPromptFix, user)
	if err != nil {
		return FixedNote{}, err
	}
	return parseFixedNoteJSON(raw)
}

// ChatProvider is implemented by online providers that can complete a
// single chat turn with system+user messages and return the raw text.
type ChatProvider interface {
	SendChat(ctx context.Context, system, user string) (string, error)
}

func userPromptForFix(req FixRequest) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Card front: %s\n", req.Note.Front)
	fmt.Fprintf(&b, "Card back: %s\n", req.Note.Back)
	if strings.TrimSpace(req.Note.Extra) != "" {
		fmt.Fprintf(&b, "Card extra: %s\n", req.Note.Extra)
	}
	if ex := firstExample(req.Note.Examples); ex != "" {
		fmt.Fprintf(&b, "Card example: %s\n", ex)
	}
	if len(req.Note.Tags) > 0 {
		fmt.Fprintf(&b, "Tags: %s\n", strings.Join(req.Note.Tags, ", "))
	}
	if strings.TrimSpace(req.UserReport) != "" {
		fmt.Fprintf(&b, "User says: %s\n", strings.TrimSpace(req.UserReport))
	} else {
		b.WriteString("User says: This card looks wrong. Please review and fix.\n")
	}
	return b.String()
}

func firstExample(examples []string) string {
	for _, e := range examples {
		if s := strings.TrimSpace(e); s != "" {
			return s
		}
	}
	return ""
}

type rawFixedNote struct {
	Front   string `json:"front"`
	Back    string `json:"back"`
	Extra   string `json:"extra"`
	Example string `json:"example"`
	Reason  string `json:"reason"`
}

func parseFixedNoteJSON(raw string) (FixedNote, error) {
	body := extractJSON(raw)
	if body == "" {
		return FixedNote{}, fmt.Errorf("no JSON object found in fix response: %q", truncate(raw, 200))
	}
	parsed, err := decodeFixedNote(body)
	if err != nil {
		return FixedNote{}, fmt.Errorf("decode fix: %w", err)
	}
	if strings.TrimSpace(parsed.Front) == "" || strings.TrimSpace(parsed.Back) == "" {
		return FixedNote{}, errors.New("fix response missing front or back")
	}
	return FixedNote{
		Front:   strings.TrimSpace(parsed.Front),
		Back:    strings.TrimSpace(parsed.Back),
		Extra:   strings.TrimSpace(parsed.Extra),
		Example: strings.TrimSpace(parsed.Example),
		Reason:  strings.TrimSpace(parsed.Reason),
	}, nil
}

func decodeFixedNote(body string) (rawFixedNote, error) {
	var r rawFixedNote
	err := json.Unmarshal([]byte(body), &r)
	return r, err
}
