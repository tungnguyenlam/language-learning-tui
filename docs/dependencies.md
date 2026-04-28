# Dependency Policy

## Current Choices

- `charm.land/bubbletea/v2`: TUI runtime with declarative v2 view capabilities.
- `charm.land/lipgloss/v2`: terminal styling.
- `github.com/open-spaced-repetition/go-fsrs/v3`: FSRS scheduling algorithm.
- `modernc.org/sqlite`: pure Go SQLite driver for local progress storage.

## Upgrade Rules

- Keep dependencies pinned in `go.mod`.
- Upgrade one major dependency family at a time.
- Run `./scripts/verify.sh` after every dependency change.
- For Bubble Tea or Lip Gloss upgrades, test compact, medium, and wide render paths.
- For FSRS upgrades, confirm scheduler tests and review-state persistence still pass.
- For SQLite driver upgrades, run repository tests and check migrations on a fresh temporary database.

## Network Rule

Dependency download is allowed only during explicit dependency work. Tests must not require network access.
