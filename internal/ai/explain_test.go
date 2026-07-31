package ai

import (
	"context"
	"testing"

	"deutsch-tui/internal/core"
)

type mockChatProvider struct {
	response string
	err      error
}

func (p mockChatProvider) GenerateDrafts(ctx context.Context, request DraftRequest) ([]Draft, error) {
	return nil, nil
}

func (p mockChatProvider) SendChat(ctx context.Context, system, user string) (string, error) {
	return p.response, p.err
}

func TestExplainCard(t *testing.T) {
	card := core.Card{
		Prompt: "Hund",
		Answer: "dog",
		Extra:  "masculine",
	}

	t.Run("Chat provider returns explanation", func(t *testing.T) {
		provider := mockChatProvider{response: "This is an explanation."}
		explanation, err := ExplainCard(context.Background(), provider, card)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if explanation != "This is an explanation." {
			t.Errorf("expected 'This is an explanation.', got %q", explanation)
		}
	})

	t.Run("Offline provider returns placeholder", func(t *testing.T) {
		provider := OfflineProvider{}
		explanation, err := ExplainCard(context.Background(), provider, card)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if explanation != "Offline/template provider cannot provide detailed explanations." {
			t.Errorf("unexpected explanation: %q", explanation)
		}
	})
}

func TestExplainDictionaryEntry(t *testing.T) {
	entry := core.DictionaryEntry{
		Word:        "geheim",
		Translation: "secret",
		WordClass:   "adj",
	}

	t.Run("Chat provider returns explanation", func(t *testing.T) {
		provider := mockChatProvider{response: "geheim is an adjective."}
		explanation, err := ExplainDictionaryEntry(context.Background(), provider, entry)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if explanation != "geheim is an adjective." {
			t.Errorf("expected explanation, got %q", explanation)
		}
	})

	t.Run("Offline provider returns placeholder", func(t *testing.T) {
		explanation, err := ExplainDictionaryEntry(context.Background(), OfflineProvider{}, entry)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if explanation != "Offline/template provider cannot provide detailed explanations." {
			t.Errorf("unexpected explanation: %q", explanation)
		}
	})
}
