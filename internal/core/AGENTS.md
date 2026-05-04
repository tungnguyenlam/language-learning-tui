# Core Agent Rules

- **Pure Domain**: Do NOT add any dependencies to external packages like `sqlite`, `bubbletea`, or `resty`.
- **Interface Segregation**: When adding new capabilities (e.g., search, statistics), define the interface here and implement it in the appropriate `internal/` package.
- **Zero Frameworks**: Keep the core models as simple structs and interfaces.
