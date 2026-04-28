# Storage Migration Policy

SQLite schema changes are append-only migrations.

## Rules

- Add a new numbered migration ID for every schema change.
- Never edit an old migration after it has shipped, except to fix a syntax error before release.
- Migrations must be idempotent when possible.
- Repository tests must use temporary or in-memory databases.
- Record migration-sensitive changes in `docs/backlog/done.md` when completed.

## Verification

Run:

```sh
GOCACHE=/tmp/deutsch-tui-gocache go test ./internal/storage/sqlite
```
