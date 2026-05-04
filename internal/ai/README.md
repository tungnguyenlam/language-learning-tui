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

- `AIProvider`: Interface for all AI backends.
- `DraftNote`: Intermediary structure for unapproved content.
- `OpenAIProvider`: (Example) implementation for OpenAI's API.

## Testing

- **NO NETWORK**: Unit tests must use a `FakeProvider` to simulate AI responses.
- **Validation Tests**: Ensure the parser can handle various AI output quirks (extra whitespace, different JSON shapes).
