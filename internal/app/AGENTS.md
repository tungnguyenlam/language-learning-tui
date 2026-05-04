# App Agent Rules

- **No Business Logic**: Do not put SRS scheduling or TUI rendering logic here.
- **Config Invariants**: Ensure any new config fields have sensible defaults in `DefaultConfig()`.
- **Path Safety**: Always use `filepath.Join` and ensure paths are relative to the data directory where appropriate.
