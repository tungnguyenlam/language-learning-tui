package ai

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestOllamaGeneratesDraftsWithoutAPIKey(t *testing.T) {
	mockResponse := `{
		"message": {
			"role": "assistant",
			"content": "{\"cards\":[{\"front\":\"der Rechner\",\"back\":\"the computer\",\"extra\":\"noun, masculine; plural: die Rechner\",\"example\":\"Mein Rechner ist kaputt.\"}]}"
		},
		"done": true
	}`

	var captured *http.Request
	var body []byte
	client := mockHTTPClient(t, http.StatusOK, mockResponse, func(req *http.Request) {
		captured = req
		if req.Body != nil {
			body, _ = io.ReadAll(req.Body)
		}
	})

	p := OllamaProvider{Model: "llama3.1", Client: client}
	drafts, err := p.GenerateDrafts(context.Background(), DraftRequest{SourceText: "Technik", DeckID: "b2-tech"})
	if err != nil {
		t.Fatalf("GenerateDrafts: %v", err)
	}
	if len(drafts) != 1 || drafts[0].Note.Front != "der Rechner" {
		t.Fatalf("unexpected drafts: %+v", drafts)
	}

	// The whole point of this provider is that it works with no credentials.
	if got := captured.Header.Get("Authorization"); got != "" {
		t.Errorf("Ollama request should not send Authorization, got %q", got)
	}
	if !strings.HasSuffix(captured.URL.Path, "/api/chat") {
		t.Errorf("expected /api/chat, got %s", captured.URL.Path)
	}
	if captured.URL.Host != "localhost:11434" {
		t.Errorf("expected default localhost:11434, got %s", captured.URL.Host)
	}

	var sent ollamaRequestBody
	if err := json.Unmarshal(body, &sent); err != nil {
		t.Fatalf("decode sent body: %v", err)
	}
	if sent.Stream {
		t.Error("request must set stream=false; the response decoder reads a single object")
	}
	if sent.Format != "json" {
		t.Errorf("drafting should request JSON mode, got format=%q", sent.Format)
	}
}

func TestOllamaHonorsCustomBaseURLAndModel(t *testing.T) {
	mockResponse := `{"message":{"role":"assistant","content":"{\"cards\":[{\"front\":\"das Haus\",\"back\":\"the house\"}]}"},"done":true}`
	var captured *http.Request
	var body []byte
	client := mockHTTPClient(t, http.StatusOK, mockResponse, func(req *http.Request) {
		captured = req
		if req.Body != nil {
			body, _ = io.ReadAll(req.Body)
		}
	})

	p := OllamaProvider{Model: "mistral", BaseURL: "http://192.168.1.50:11434/", Client: client}
	if _, err := p.GenerateDrafts(context.Background(), DraftRequest{SourceText: "Wohnen", DeckID: "a1"}); err != nil {
		t.Fatalf("GenerateDrafts: %v", err)
	}
	if captured.URL.Host != "192.168.1.50:11434" {
		t.Errorf("custom base URL ignored, got %s", captured.URL.Host)
	}
	// A trailing slash on the configured URL must not produce //api/chat.
	if captured.URL.Path != "/api/chat" {
		t.Errorf("trailing slash not trimmed, path = %s", captured.URL.Path)
	}
	var sent ollamaRequestBody
	if err := json.Unmarshal(body, &sent); err != nil {
		t.Fatalf("decode sent body: %v", err)
	}
	if sent.Model != "mistral" {
		t.Errorf("model = %q, want mistral", sent.Model)
	}
}

func TestOllamaSendChatDoesNotForceJSONMode(t *testing.T) {
	mockResponse := `{"message":{"role":"assistant","content":"Das Wort ist maskulin."},"done":true}`
	var body []byte
	client := mockHTTPClient(t, http.StatusOK, mockResponse, func(req *http.Request) {
		if req.Body != nil {
			body, _ = io.ReadAll(req.Body)
		}
	})

	p := OllamaProvider{Client: client}
	got, err := p.SendChat(context.Background(), "system", "user")
	if err != nil {
		t.Fatalf("SendChat: %v", err)
	}
	if got != "Das Wort ist maskulin." {
		t.Errorf("SendChat = %q", got)
	}

	var sent ollamaRequestBody
	if err := json.Unmarshal(body, &sent); err != nil {
		t.Fatalf("decode sent body: %v", err)
	}
	// Explanations are prose. Forcing JSON here would make them unreadable.
	if sent.Format != "" {
		t.Errorf("SendChat must not force JSON mode, got format=%q", sent.Format)
	}
}

func TestOllamaMissingModelErrorIsActionable(t *testing.T) {
	client := mockHTTPClient(t, http.StatusNotFound, `{"error":"model 'llama3.1' not found"}`, nil)
	p := OllamaProvider{Model: "llama3.1", Client: client}
	_, err := p.GenerateDrafts(context.Background(), DraftRequest{SourceText: "Kaffee", DeckID: "d1"})
	if err == nil {
		t.Fatal("expected an error for a missing model")
	}
	if !strings.Contains(err.Error(), "ollama pull") {
		t.Errorf("error should tell the user how to install the model, got: %v", err)
	}
}

func TestOllamaEmptyCompletionIsRejected(t *testing.T) {
	client := mockHTTPClient(t, http.StatusOK, `{"message":{"role":"assistant","content":""},"done":true}`, nil)
	p := OllamaProvider{Client: client}
	if _, err := p.SendChat(context.Background(), "s", "u"); err == nil {
		t.Fatal("expected an error for an empty completion")
	}
}

func TestOllamaValidatesRequestBeforeCallingOut(t *testing.T) {
	client := mockHTTPClient(t, http.StatusOK, `{}`, func(*http.Request) {
		t.Error("provider should not issue a request for an invalid draft request")
	})
	p := OllamaProvider{Client: client}
	if _, err := p.GenerateDrafts(context.Background(), DraftRequest{DeckID: "d1"}); err == nil {
		t.Fatal("expected an error when source text is empty")
	}
	if _, err := p.GenerateDrafts(context.Background(), DraftRequest{SourceText: "Kaffee"}); err == nil {
		t.Fatal("expected an error when deck id is empty")
	}
}

func TestOllamaSatisfiesChatProvider(t *testing.T) {
	var _ Provider = OllamaProvider{}
	if _, ok := Provider(OllamaProvider{}).(ChatProvider); !ok {
		t.Fatal("OllamaProvider must implement ChatProvider so FixCard and explanations work")
	}
}
