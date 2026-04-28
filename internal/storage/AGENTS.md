# Storage Agent Rules

- Schema changes must be migrations, not ad hoc table edits.
- Add a new numbered migration; do not edit shipped migrations.
- Repository tests must use temporary SQLite databases.
- Keep TUI and AI provider details out of storage.
