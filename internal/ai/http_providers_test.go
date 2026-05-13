package ai

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// roundTripFunc lets us mock the HTTP client without spinning up a server.
type roundTripFunc func(req *http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func mockHTTPClient(t *testing.T, status int, body string, capture func(req *http.Request)) *http.Client {
	t.Helper()
	return &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			if capture != nil {
				capture(req)
			}
			return &http.Response{
				StatusCode: status,
				Body:       io.NopCloser(strings.NewReader(body)),
				Header:     http.Header{"Content-Type": []string{"application/json"}},
			}, nil
		}),
	}
}

func TestOpenAIRequiresAPIKey(t *testing.T) {
	p := OpenAIProvider{}
	_, err := p.GenerateDrafts(context.Background(), DraftRequest{SourceText: "Kaffee", DeckID: "d1"})
	if err == nil {
		t.Fatal("expected error when API key is empty")
	}
	if !strings.Contains(err.Error(), "API key") {
		t.Errorf("error should mention API key, got: %v", err)
	}
}

func TestAnthropicRequiresAPIKey(t *testing.T) {
	p := AnthropicProvider{}
	_, err := p.GenerateDrafts(context.Background(), DraftRequest{SourceText: "Kaffee", DeckID: "d1"})
	if err == nil {
		t.Fatal("expected error when API key is empty")
	}
	if !strings.Contains(err.Error(), "API key") {
		t.Errorf("error should mention API key, got: %v", err)
	}
}

func TestOpenAIGeneratesDraftsFromMockServer(t *testing.T) {
	mockResponse := `{
		"choices": [
			{"message": {"role": "assistant", "content": "{\"cards\":[{\"front\":\"der Hund\",\"back\":\"the dog\",\"extra\":\"masculine noun\",\"example\":\"Ich habe einen Hund.\"},{\"front\":\"die Katze\",\"back\":\"the cat\",\"extra\":\"feminine noun\",\"example\":\"Die Katze schläft.\"}]}"}}
		]
	}`
	var capturedReq *http.Request
	var capturedBody []byte
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedReq = r
		capturedBody, _ = io.ReadAll(r.Body)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(mockResponse))
	}))
	defer srv.Close()

	p := OpenAIProvider{
		APIKey:  "sk-test-key",
		Model:   "gpt-4o-mini",
		BaseURL: srv.URL,
	}
	drafts, err := p.GenerateDrafts(context.Background(), DraftRequest{
		SourceText: "A1 animals",
		DeckID:     "a1-animals",
		Tags:       []string{"a1"},
	})
	if err != nil {
		t.Fatalf("GenerateDrafts: %v", err)
	}
	if len(drafts) != 2 {
		t.Fatalf("len(drafts) = %d, want 2", len(drafts))
	}
	if drafts[0].Note.Front != "der Hund" {
		t.Errorf("front = %q, want der Hund", drafts[0].Note.Front)
	}
	if drafts[1].Note.Back != "the cat" {
		t.Errorf("back = %q, want the cat", drafts[1].Note.Back)
	}

	if capturedReq.Header.Get("Authorization") != "Bearer sk-test-key" {
		t.Errorf("Authorization header = %q, want Bearer sk-test-key", capturedReq.Header.Get("Authorization"))
	}
	if capturedReq.URL.Path != "/chat/completions" {
		t.Errorf("path = %q, want /chat/completions", capturedReq.URL.Path)
	}
	var body openAIRequestBody
	if err := json.Unmarshal(capturedBody, &body); err != nil {
		t.Fatalf("request body decode: %v", err)
	}
	if body.Model != "gpt-4o-mini" {
		t.Errorf("model = %q, want gpt-4o-mini", body.Model)
	}
	if len(body.Messages) != 2 {
		t.Errorf("messages len = %d, want 2 (system+user)", len(body.Messages))
	}
}

func TestAnthropicGeneratesDraftsFromMockServer(t *testing.T) {
	mockResponse := `{
		"content": [
			{"type": "text", "text": "{\"cards\":[{\"front\":\"lernen\",\"back\":\"to learn\",\"extra\":\"verb, weak\",\"example\":\"Ich lerne Deutsch.\"}]}"}
		]
	}`
	var capturedReq *http.Request
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedReq = r
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(mockResponse))
	}))
	defer srv.Close()

	p := AnthropicProvider{
		APIKey:  "sk-ant-test",
		Model:   "claude-haiku-4-5-20251001",
		BaseURL: srv.URL,
	}
	drafts, err := p.GenerateDrafts(context.Background(), DraftRequest{
		SourceText: "verbs",
		DeckID:     "a1-verbs",
	})
	if err != nil {
		t.Fatalf("GenerateDrafts: %v", err)
	}
	if len(drafts) != 1 {
		t.Fatalf("len(drafts) = %d, want 1", len(drafts))
	}
	if drafts[0].Note.Front != "lernen" {
		t.Errorf("front = %q, want lernen", drafts[0].Note.Front)
	}

	if capturedReq.Header.Get("x-api-key") != "sk-ant-test" {
		t.Errorf("x-api-key header = %q, want sk-ant-test", capturedReq.Header.Get("x-api-key"))
	}
	if capturedReq.Header.Get("anthropic-version") == "" {
		t.Error("anthropic-version header should be set")
	}
	if capturedReq.URL.Path != "/v1/messages" {
		t.Errorf("path = %q, want /v1/messages", capturedReq.URL.Path)
	}
}

func TestOpenAIReportsServerError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error":{"message":"invalid api key"}}`))
	}))
	defer srv.Close()

	p := OpenAIProvider{APIKey: "bad", BaseURL: srv.URL}
	_, err := p.GenerateDrafts(context.Background(), DraftRequest{SourceText: "x", DeckID: "d"})
	if err == nil {
		t.Fatal("expected error on 401")
	}
	if !strings.Contains(err.Error(), "401") {
		t.Errorf("error should mention 401: %v", err)
	}
}

func TestAnthropicReportsServerError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":{"type":"invalid_request_error","message":"bad model"}}`))
	}))
	defer srv.Close()

	p := AnthropicProvider{APIKey: "k", BaseURL: srv.URL}
	_, err := p.GenerateDrafts(context.Background(), DraftRequest{SourceText: "x", DeckID: "d"})
	if err == nil {
		t.Fatal("expected error on 400")
	}
}

func TestExtractJSONHandlesMarkdownFences(t *testing.T) {
	// Models sometimes wrap JSON in ```json blocks.
	raw := "Here are your cards:\n```json\n{\"cards\":[{\"front\":\"a\",\"back\":\"b\"}]}\n```"
	got := extractJSON(raw)
	if got == "" {
		t.Fatal("extractJSON returned empty for markdown-fenced JSON")
	}
	var parsed rawCards
	if err := json.Unmarshal([]byte(got), &parsed); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(parsed.Cards) != 1 {
		t.Errorf("len(cards) = %d, want 1", len(parsed.Cards))
	}
}

func TestParseCardsJSONHandlesEmbeddedBraces(t *testing.T) {
	raw := `{"cards":[{"front":"x","back":"y","extra":"see {{c1::...}}"}]}`
	cards, err := parseCardsJSON(raw)
	if err != nil {
		t.Fatalf("parseCardsJSON: %v", err)
	}
	if len(cards) != 1 {
		t.Fatalf("len(cards) = %d, want 1", len(cards))
	}
	if cards[0].Extra != "see {{c1::...}}" {
		t.Errorf("extra not preserved: %q", cards[0].Extra)
	}
}
