# Active Backlog

Last updated: 2026-06-16

## Current Milestone

Hybrid Dictionary - COMPLETED.

## Exact Next Action

Select next milestone from `docs/backlog/roadmap.md` (e.g., `.apkg` import/export or Backup workflow) after committing the current fixes.

## Top Issues

None active.

## Completed Work

- [x] **Flashcard Loop Enhancement [FIXED]:** Enhanced the AI Drafting flow to capture rich context (word class, gender, forms, examples) from dictionary search results.
- [x] **UI Badge:** Added a visual "★ Dictionary Context Active" badge to the AI View to inform the user when dictionary context is included in the draft request.
- [x] **State Management:** Added `m.draftSource` to the TUI model, correctly clearing it when the user overrides or escapes the drafted query.

## Verification

- `./scripts/verify.sh` passed: all Go unit tests, smoke test, binary build, and 351 E2E tests.
