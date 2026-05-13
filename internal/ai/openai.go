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
	defaultOpenAIBaseURL = "https://api.openai.com/v1"
	defaultOpenAIModel   = "gpt-4o-mini"
)

// OpenAIProvider talks to any OpenAI-compatible Chat Completions endpoint.
// BaseURL is overridable so users can point at Azure OpenAI, OpenRouter,
// Groq, a local Ollama, or anything else that speaks the same protocol.
type OpenAIProvider struct {
	APIKey  string
	Model   string
	BaseURL string
	Client  *http.Client
	// Timeout caps the request when Client doesn't already set one.
	Timeout time.Duration
}

func (p OpenAIProvider) GenerateDrafts(ctx context.Context, request DraftRequest) ([]Draft, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if strings.TrimSpace(p.APIKey) == "" {
		return nil, errors.New("openai: API key is required (set it in Settings)")
	}
	if strings.TrimSpace(request.SourceText) == "" {
		return nil, errors.New("draft source text is required")
	}
	if strings.TrimSpace(request.DeckID) == "" {
		return nil, errors.New("draft deck id is required")
	}

	model := strings.TrimSpace(p.Model)
	if model == "" {
		model = defaultOpenAIModel
	}
	baseURL := strings.TrimRight(strings.TrimSpace(p.BaseURL), "/")
	if baseURL == "" {
		baseURL = defaultOpenAIBaseURL
	}

	body := openAIRequestBody{
		Model: model,
		Messages: []openAIMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userPromptFor(request)},
		},
		Temperature: 0.4,
		ResponseFormat: &openAIResponseFormat{
			Type: "json_object",
		},
	}
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("openai: encode request: %w", err)
	}

	url := baseURL + "/chat/completions"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(buf))
	if err != nil {
		return nil, fmt.Errorf("openai: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.APIKey)

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
		return nil, fmt.Errorf("openai: request: %w", err)
	}
	defer resp.Body.Close()
	respBytes, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, fmt.Errorf("openai: read response: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("openai: %s — %s", resp.Status, truncate(string(respBytes), 400))
	}

	var parsed openAIResponse
	if err := json.Unmarshal(respBytes, &parsed); err != nil {
		return nil, fmt.Errorf("openai: decode response: %w", err)
	}
	if len(parsed.Choices) == 0 || parsed.Choices[0].Message.Content == "" {
		return nil, errors.New("openai: empty completion")
	}
	rawCards, err := parseCardsJSON(parsed.Choices[0].Message.Content)
	if err != nil {
		return nil, fmt.Errorf("openai: %w", err)
	}
	return draftsFromRaw(rawCards, request)
}

type openAIRequestBody struct {
	Model          string                `json:"model"`
	Messages       []openAIMessage       `json:"messages"`
	Temperature    float64               `json:"temperature,omitempty"`
	ResponseFormat *openAIResponseFormat `json:"response_format,omitempty"`
}

type openAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openAIResponseFormat struct {
	Type string `json:"type"`
}

type openAIResponse struct {
	Choices []struct {
		Message openAIMessage `json:"message"`
	} `json:"choices"`
}
