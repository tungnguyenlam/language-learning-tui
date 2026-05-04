# internal/storage

SQLite-based persistence layer for flashcards, progress, and review history.

## Architecture

- **SQLite**: Local-first storage using a single `.db` file in the data directory.
- **Migrations**: Incremental schema updates using a numbered migration system.
- **Repository Pattern**: All database access is encapsulated in the `SqliteRepository`, which implements `core.Repository`.

## Migrations

Migrations are stored as Go strings in `migrations.go` and applied automatically on startup.
- Follow the naming convention: `vN_description.sql` (if they were separate files, but here they are indexed).
- **NEVER** modify an existing migration. Always add a new one.

## Key Symbols

- `SqliteRepository`: Main implementation of the storage interface.
- `NewSqliteRepository`: Initializes the database and runs migrations.
- `Schema`: The initial base schema.

## Testing

- Use `sqlite_test.go` for integration tests.
- Always use a temporary file or `:memory:` database for tests to avoid polluting the user's data.
