# Done Backlog

## 2026-05-06: Bug Hunting & Robustness Improvements (STABILIZED)
- **Bug Fix**: Fixed MCQ spoiler bug where the correct answer was prematurely revealed during animation.
- **Bug Fix**: Resolved key collision between MCQ choices and view switching by implementing correct precedence.
- **Bug Fix**: Fixed Modal Interaction Trap where mouse clicks could trigger background actions during confirmation.
- **Bug Fix**: Improved `truncateLine` to be rune-aware, preventing UTF-8 corruption with German characters.
- **Bug Fix**: Fixed case-insensitive word matching to handle German `ẞ` correctly (byte-length changes).
- **Bug Fix**: Improved FSRS interval prediction accuracy for short-term reviews.
- **UX Improvement**: Added explicit edit mode (`/`) for AI topic input to resolve navigation key collisions.
- **UI Polish**: Unified emojis across views (`🔥` for streak, `✨` for mature) and removed redundant dashboard header lines.
- **Test Coverage**: Verified all 100 E2E tests pass after logic and UI updates.

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
