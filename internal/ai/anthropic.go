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
	defaultAnthropicModel   = "claude-haiku-4-5-20251001"
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
	if strings.TrimSpace(p.APIKey) == "" {
		return nil, errors.New("anthropic: API key is required (set it in Settings)")
	}
	if strings.TrimSpace(request.SourceText) == "" {
		return nil, errors.New("draft source text is required")
	}
	if strings.TrimSpace(request.DeckID) == "" {
		return nil, errors.New("draft deck id is required")
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
		Model:     model,
		MaxTokens: defaultAnthropicMaxTok,
		System:    systemPrompt,
		Messages: []anthropicMessage{
			{Role: "user", Content: userPromptFor(request)},
		},
		Temperature: 0.4,
	}
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("anthropic: encode request: %w", err)
	}

	url := baseURL + "/v1/messages"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(buf))
	if err != nil {
		return nil, fmt.Errorf("anthropic: build request: %w", err)
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
		return nil, fmt.Errorf("anthropic: request: %w", err)
	}
	defer resp.Body.Close()
	respBytes, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, fmt.Errorf("anthropic: read response: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("anthropic: %s — %s", resp.Status, truncate(string(respBytes), 400))
	}

	var parsed anthropicResponse
	if err := json.Unmarshal(respBytes, &parsed); err != nil {
		return nil, fmt.Errorf("anthropic: decode response: %w", err)
	}
	text := joinAnthropicText(parsed.Content)
	if text == "" {
		return nil, errors.New("anthropic: empty message")
	}
	rawCards, err := parseCardsJSON(text)
	if err != nil {
		return nil, fmt.Errorf("anthropic: %w", err)
	}
	return draftsFromRaw(rawCards, request)
}

// SendChat performs a single Anthropic Messages turn with system+user
// content and returns the raw text content. Used by FixCard.
func (p AnthropicProvider) SendChat(ctx context.Context, system, user string) (string, error) {
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
		Temperature: 0.2,
	}
	buf, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("anthropic: encode chat request: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURL+"/v1/messages", bytes.NewReader(buf))
	if err != nil {
		return "", fmt.Errorf("anthropic: build chat request: %w", err)
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
		return "", fmt.Errorf("anthropic: chat request: %w", err)
	}
	defer resp.Body.Close()
	respBytes, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return "", fmt.Errorf("anthropic: read chat response: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("anthropic: %s — %s", resp.Status, truncate(string(respBytes), 400))
	}
	var parsed anthropicResponse
	if err := json.Unmarshal(respBytes, &parsed); err != nil {
		return "", fmt.Errorf("anthropic: decode chat response: %w", err)
	}
	text := joinAnthropicText(parsed.Content)
	if text == "" {
		return "", errors.New("anthropic: empty chat response")
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
