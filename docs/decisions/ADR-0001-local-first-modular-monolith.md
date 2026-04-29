# ADR-0001: Local-First Modular Monolith

Date: 2026-04-29

## Status

Accepted.

## Decision

Build `deutsch-tui` as a local-first Go modular monolith. Keep domain logic, storage, content interop, SRS, AI provider adapters, and TUI concerns in separate packages under `internal/`.

Use small indexed continuity docs instead of a large root `AGENTS.md`.

## Consequences

- Future agents can start from `AGENTS.md` and inspect only the relevant package.
- Feature work stays easier to test without launching the full TUI.
- Cross-package contracts must be documented when they change.
