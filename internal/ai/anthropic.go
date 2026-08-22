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
	defaultAnthropicBaseURL = "https://api.anthropic.com"
	defaultAnthropicVersion = "2023-06-01"
	defaultAnthropicModel   = "claude-3-5-haiku-latest"
	defaultAnthropicMaxTok  = 1500
)

// AnthropicProvider talks to the Anthropic Messages API. BaseURL is
// overridable so the same code can target a local proxy or staging
// endpoint while still using the Messages protocol shape.
type AnthropicProvider struct {
	APIKey  string
	Model   string
	BaseURL string
	Client  *http.Client
	Timeout time.Duration
}

func (p AnthropicProvider) GenerateDrafts(ctx context.Context, request DraftRequest) ([]Draft, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if err := validateDraftRequest(request); err != nil {
		return nil, err
	}
	text, err := p.chat(ctx, systemPrompt, userPromptFor(request), 0.4)
	if err != nil {
		return nil, err
	}
	rawCards, err := parseCardsJSON(text)
	if err != nil {
		return nil, fmt.Errorf("anthropic: %w", err)
	}
	return draftsFromRaw(rawCards, request)
}

func (p AnthropicProvider) SendChat(ctx context.Context, system, user string) (string, error) {
	return p.chat(ctx, system, user, 0.2)
}

func (p AnthropicProvider) chat(ctx context.Context, system, user string, temperature float64) (string, error) {
	if strings.TrimSpace(p.APIKey) == "" {
		return "", errors.New("anthropic: API key is required (set it in Settings)")
	}
	model := strings.TrimSpace(p.Model)
	if model == "" {
		model = defaultAnthropicModel
	}
	baseURL := strings.TrimRight(strings.TrimSpace(p.BaseURL), "/")
	if baseURL == "" {
		baseURL = defaultAnthropicBaseURL
	}

	body := anthropicRequestBody{
		Model:       model,
		MaxTokens:   defaultAnthropicMaxTok,
		System:      system,
		Messages:    []anthropicMessage{{Role: "user", Content: user}},
		Temperature: temperature,
	}
	buf, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("anthropic: encode request: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURL+"/v1/messages", bytes.NewReader(buf))
	if err != nil {
		return "", fmt.Errorf("anthropic: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", p.APIKey)
	req.Header.Set("anthropic-version", defaultAnthropicVersion)

	client := p.Client
	if client == nil {
		timeout := p.Timeout
		if timeout == 0 {
			timeout = 60 * time.Second
		}
		client = &http.Client{Timeout: timeout}
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("anthropic: request: %w", err)
	}
	defer resp.Body.Close()
	respBytes, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return "", fmt.Errorf("anthropic: read response: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("anthropic: %s — %s", resp.Status, truncate(string(respBytes), 400))
	}

	var parsed anthropicResponse
	if err := json.Unmarshal(respBytes, &parsed); err != nil {
		return "", fmt.Errorf("anthropic: decode response: %w", err)
	}
	text := joinAnthropicText(parsed.Content)
	if text == "" {
		return "", errors.New("anthropic: empty message")
	}
	return text, nil
}

type anthropicRequestBody struct {
	Model       string             `json:"model"`
	MaxTokens   int                `json:"max_tokens"`
	System      string             `json:"system,omitempty"`
	Messages    []anthropicMessage `json:"messages"`
	Temperature float64            `json:"temperature,omitempty"`
}

type anthropicMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type anthropicResponse struct {
	Content []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"content"`
}

func joinAnthropicText(blocks []struct {
	Type string `json:"type"`
	Text string `json:"text"`
}) string {
	var parts []string
	for _, b := range blocks {
		if b.Type == "text" && b.Text != "" {
			parts = append(parts, b.Text)
		}
	}
	return strings.Join(parts, "\n")
}
