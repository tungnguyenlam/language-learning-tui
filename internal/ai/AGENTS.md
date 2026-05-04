# AI Agent Rules

- **Offline First**: Tests must NEVER require network access or API keys. Use `FakeProvider`.
- **Drafting Safety**: AI content must be validated before being returned to the TUI. Use `core.Note` types for validation.
- **Secrets**: NEVER log or hardcode API keys. Use `internal/app.Config` for credential management.
- **Prompt Stability**: If you update prompts, ensure they are versioned or at least well-commented to explain the intent.
