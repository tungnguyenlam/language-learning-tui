# SRS Agent Rules

- **Encapsulation**: Do NOT let `go-fsrs` types leak into other packages. Map them to `core` types.
- **Determinism**: Scheduling logic should be deterministic and easy to unit test.
- **Stability**: Avoid changing the core scheduling algorithm without a clear requirement and migration plan for existing card states.
