# internal/ai

AI-powered content generation for German flashcards.

## Responsibilities

- **Drafting**: Generating potential flashcards based on a topic or vocabulary word.
- **Validation**: Ensuring AI-generated content follows the required format (MCQ, Basic, etc.).
- **Provider Abstraction**: Supporting multiple LLM providers (e.g., OpenAI, Anthropic, Ollama) via a unified interface.

## Core Principles

- **Local-First**: The app remains functional without AI. AI is an optional enhancement.
- **Human-in-the-loop**: AI content is presented as a "Draft" that the user must review and approve before it's saved to a deck.

## Key Symbols

- `Provider`: Interface for all AI backends.
- `Draft`: Intermediary structure for unapproved content.
- `OfflineProvider`: Generates stub flashcards locally — no network.
- `TemplateProvider`: Substitutes a user-defined template per topic.
- `OpenAIProvider`: Talks to any OpenAI-compatible Chat Completions endpoint (`/v1/chat/completions`). BaseURL is overridable, so the same code targets Azure OpenAI, OpenRouter, Groq, Ollama, etc.
- `AnthropicProvider`: Talks to the Anthropic Messages API (`/v1/messages`). Sets `x-api-key` and `anthropic-version` headers.

## Credentials

API keys live in `secrets.json` next to `config.json`, written at file mode `0600` and never logged. The settings UI in the TUI displays the key masked (last 4 chars only). See `internal/app/secrets.go`.

## Prompt / Response Shape

Both providers share a system prompt that asks for strict JSON:

```json
{"cards": [{"front": "...", "back": "...", "extra": "...", "example": "..."}]}
```

`parseCardsJSON` tolerates leading prose, markdown fences, and braces inside example strings; it extracts the first balanced JSON object.

## Testing

- **NO NETWORK**: Unit tests must use a `FakeProvider`, `httptest.NewServer`, or a stubbed `http.Client` — never the real APIs.
- **Validation Tests**: Ensure the parser can handle various AI output quirks (extra whitespace, different JSON shapes, markdown fences).
