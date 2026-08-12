# internal/storage

SQLite-based persistence layer for flashcards, progress, and review history.

## Architecture

- **SQLite**: Local-first storage using a single `.db` file in the data directory.
- **Migrations**: Incremental schema updates using a numbered migration system.
- **Repository Pattern**: All database access is encapsulated in `sqlite.Store`, which implements
  `core.Repository`.

## Migrations

Migrations are stored as Go strings in `migrations.go` and applied automatically on startup.
- Follow the naming convention: `vN_description.sql` (if they were separate files, but here they are indexed).
- **NEVER** modify an existing migration. Always add a new one.

## Key Symbols

- `Store`: Main implementation of the storage interface.
- `Open` / `OpenMemory`: Initialize a file-backed or in-memory store and run migrations.
- `migrations`: Ordered schema history in `migrations.go`.
- `scanCard` / `scanCards`: Canonical persisted-card row decoding shared by card queries.
- `Store.withTx`: Commit-on-success transaction boundary for repository mutations.

## Testing

- Use `sqlite_test.go` for integration tests.
- Always use a temporary file or `:memory:` database for tests to avoid polluting the user's data.
