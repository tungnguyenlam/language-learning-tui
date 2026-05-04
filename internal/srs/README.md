# internal/srs

Scheduling adapter and FSRS implementation.

## Responsibilities

- **Scheduling Algorithm**: Implements the Spaced Repetition System logic.
- **FSRS Wrapper**: Currently wraps `github.com/open-spaced-repetition/go-fsrs/v3`.
- **State Transition**: Calculates the next state of a card based on the user's grade.

## Key Symbols

- `Scheduler`: Implementation of `core.Scheduler`.
- `NewFSRSScheduler`: Factory for the FSRS-based scheduler.
