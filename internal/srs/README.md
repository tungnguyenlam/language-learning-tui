# internal/srs

Scheduling adapter. Keep Bubble Tea, SQLite, and content parsing out of this package.

The current implementation wraps `github.com/open-spaced-repetition/go-fsrs/v3` and keeps the public surface behind `core.Scheduler`, so upstream scheduler details do not leak through the app.
