package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	defaultOllamaBaseURL = "http://localhost:11434"
	defaultOllamaModel   = "llama3.1"
	// Local models run on the user's own CPU/GPU and a first request often
	// pays for loading weights into memory, so the ceiling is much higher
	// than the hosted providers'.
	defaultOllamaTimeout = 300 * time.Second
)

// OllamaProvider talks to a local Ollama daemon's native API. The hosted
// OpenAIProvider can technically reach Ollama's compatibility endpoint, but
// it demands an API key that a local daemon does not have and does not
// exist, which made "run it locally" impossible in practice. This provider
// requires no credentials at all.
type OllamaProvider struct {
	Model   string
	BaseURL string
	Client  *http.Client
	Timeout time.Duration
}

func (p OllamaProvider) GenerateDrafts(ctx context.Context, request DraftRequest) ([]Draft, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if strings.TrimSpace(request.SourceText) == "" {
		return nil, errors.New("draft source text is required")
	}
	if strings.TrimSpace(request.DeckID) == "" {
		return nil, errors.New("draft deck id is required")
	}

	// Ollama's `format: "json"` constrains decoding to valid JSON, which
	// keeps small local models from wrapping the object in prose.
	content, err := p.chat(ctx, systemPrompt, userPromptFor(request), true)
	if err != nil {
		return nil, err
	}
	rawCards, err := parseCardsJSON(content)
	if err != nil {
		return nil, fmt.Errorf("ollama: %w", err)
	}
	return draftsFromRaw(rawCards, request)
}

// SendChat completes a single turn and returns the raw text. Used by
// FixCard and the explanation flows.
func (p OllamaProvider) SendChat(ctx context.Context, system, user string) (string, error) {
	// JSON mode is deliberately off here. SendChat serves both the card-fix
	// flow (which wants JSON and whose parser already tolerates surrounding
	// prose) and the explain flows (which want plain prose). Forcing JSON
	// would turn every explanation into an unreadable object.
	return p.chat(ctx, system, user, false)
}

func (p OllamaProvider) chat(ctx context.Context, system, user string, jsonMode bool) (string, error) {
	model := strings.TrimSpace(p.Model)
	if model == "" {
		model = defaultOllamaModel
	}
	baseURL := strings.TrimRight(strings.TrimSpace(p.BaseURL), "/")
	if baseURL == "" {
		baseURL = defaultOllamaBaseURL
	}

	body := ollamaRequestBody{
		Model: model,
		Messages: []ollamaMessage{
			{Role: "system", Content: system},
			{Role: "user", Content: user},
		},
		Stream:  false,
		Options: &ollamaOptions{Temperature: 0.4},
	}
	if jsonMode {
		body.Format = "json"
	}
	buf, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("ollama: encode request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURL+"/api/chat", bytes.NewReader(buf))
	if err != nil {
		return "", fmt.Errorf("ollama: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := p.Client
	if client == nil {
		timeout := p.Timeout
		if timeout == 0 {
			timeout = defaultOllamaTimeout
		}
		client = &http.Client{Timeout: timeout}
	}

	resp, err := client.Do(req)
	if err != nil {
		// A refused connection is the single most common failure here, and
		// "connection refused" alone does not tell the user what to do.
		return "", fmt.Errorf("ollama: cannot reach the local daemon at %s (is `ollama serve` running?): %w", baseURL, err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return "", fmt.Errorf("ollama: read response: %w", err)
	}
	if resp.StatusCode == http.StatusNotFound {
		// Ollama returns 404 when the model has never been pulled, which is
		// the second most common first-run failure.
		return "", fmt.Errorf("ollama: model %q is not installed (run `ollama pull %s`)", model, model)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("ollama: %s — %s", resp.Status, truncate(string(respBytes), 400))
	}

	var parsed ollamaResponse
	if err := json.Unmarshal(respBytes, &parsed); err != nil {
		return "", fmt.Errorf("ollama: decode response: %w", err)
	}
	if strings.TrimSpace(parsed.Message.Content) == "" {
		return "", errors.New("ollama: empty completion")
	}
	return parsed.Message.Content, nil
}

type ollamaRequestBody struct {
	Model    string          `json:"model"`
	Messages []ollamaMessage `json:"messages"`
	Stream   bool            `json:"stream"`
	Format   string          `json:"format,omitempty"`
	Options  *ollamaOptions  `json:"options,omitempty"`
}

type ollamaMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ollamaOptions struct {
	Temperature float64 `json:"temperature,omitempty"`
}

type ollamaResponse struct {
	Message ollamaMessage `json:"message"`
	Done    bool          `json:"done"`
}
