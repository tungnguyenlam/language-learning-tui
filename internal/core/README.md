# internal/core

Domain model for decks, notes, cards, review grades, and scheduling ports.

## Responsibilities

- **Model Definitions**: `Deck`, `Note`, `Card`, `ReviewLog`, `CardKind` (MCQ, Basic, Cloze).
- **Domain Logic**: Scheduling port definitions (`Scheduler` interface), grade constants (`GradeAgain`, `GradeHard`, `GradeGood`, `GradeEasy`).
- **Interfaces**: Defines the contracts for storage (`Repository`, optional `BackupRepository`) and AI generation (`AIProvider`) to maintain loose coupling.

## Architectural Boundaries

This package is the **center** of the modular monolith.
- **NO DEPENDENCIES**: It must NOT import `tui`, `storage`, `ai`, or `content`.
- **Pure Go**: It should primarily contain data structures and interfaces.
- **No Side Effects**: Avoid filesystem, network, or database calls here.
