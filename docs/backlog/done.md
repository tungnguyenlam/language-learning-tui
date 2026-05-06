# Done Backlog

## 2026-05-06: AI Enhancements & Deck Statistics (STABILIZED)
- **Feature**: Multi-topic AI drafting. Users can now enter multiple topics separated by commas or newlines to generate a batch of drafts.
- **Feature**: Deck-level statistics. The Statistics view now displays metrics for the currently selected deck, providing granular progress tracking.
- **UX Improvement**: Added 'v' shortcut in Decks view to instantly view statistics for the highlighted deck.
- **Content Expansion**: Added 5 new German grammar tips covering imperative, reflexive verbs, separable verbs, and comparisons.
- **Refactoring**: Unified streak calculation logic and removed duplication in SQLite storage layer.
- **Test Coverage**: Added 2 new E2E tests for multi-topic drafting and deck stats shortcut, bringing total verified tests to 102.

## 2026-05-05: TUI Stabilisation & UX Polish
- **E2E Performance**: Built binary once in `verify.sh`, reducing test suite time and resource contention.
- **Safety**: Added confirmation dialogs for all destructive actions (card/deck delete).
- **UI Integrity**: Aligned all TUI hitboxes with Lip Gloss metrics for pixel-perfect mouse interaction.
- **Data Integrity**: Implemented transactional bulk updates in the SQLite layer.
- **Test Coverage**: Added 4 new E2E tests and verified 100% pass rate (97/97).

## 2026-05-03: Template Editing & Configuration
- [x] Added revert-on-cancel for template editing.
- [x] Stored original values in `originalTemplateValue`.
- [x] Fixed template set switching logic.
