# Done Backlog

## 2026-05-06: Bug Hunting & Robustness Improvements
- **UI Robustness**: Refactored Review and Dashboard hitboxes to use dynamic line counting and string width calculations instead of hardcoded offsets.
- **Bug Fix**: Fixed Statistics view accuracy color rendering which was previously showing internal object representations.
- **Bug Fix**: Resolved duplicate and unsafe case statements in hitbox activation logic.
- **Data Integrity**: Switched MCQ choice separator from comma to `|||` in TSV import/export to prevent parsing errors with complex choices.
- **Test Coverage**: Added 3 new E2E tests (MCQ robustness, Multiline rendering, AI generation) and verified 100% pass rate (100/100).

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
