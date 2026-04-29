# Active Backlog

Last updated: 2026-04-29

## Current Milestone

Milestone 6: AI Provider Configuration and Prompt Templates.

## Next Action

Implement a configurable AI provider that supports prompt templates.

## Acceptance Criteria

- Settings view displays currently selected AI provider and allows switching (e.g., between "Offline" and "Template").
- New `TemplateProvider` in `internal/ai` supports customizable prompt templates.
- Prompt templates can be used to control how `Front`, `Back`, and `Examples` are generated.
- Settings view allows viewing and basic editing of these templates.
- Unit tests cover template substitution and provider switching logic.

## Blockers

- None.

## Last Verified

- 2026-04-29 13:25 +07: `./scripts/verify.sh`
- 2026-04-29 13:15 +07: `go test ./internal/tui ./internal/storage/sqlite`
