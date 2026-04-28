# Release Checklist

Use this before tagging or sharing a build.

1. Run `./scripts/verify.sh`.
2. Confirm migrations run on a fresh database.
3. Confirm `docs/backlog/active.md` has the next action.
4. Add completed work to `docs/backlog/done.md`.
5. Review dependency changes in `go.mod` and `go.sum`.
6. Back up personal app data before testing migration-heavy changes.
