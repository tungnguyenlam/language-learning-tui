# Active Backlog

Last updated: 2026-08-23

## Current Milestone

Autonomous workflow improvement: reshape `prompt/improve.md` from a finite
2–3-change session prompt into a continuous, evidence-driven development loop
that can operate safely with little or no human supervision.

## Exact Next Action

Review the rewritten prompt for consistency with `AGENTS.md` and continuity
rules, validate Markdown formatting, update the completed backlog, and commit.

## Completed This Pass

- Rewritten prompt now defines autonomous pickup, evidence-based work selection,
  repeatable implementation/verification/commit cycles, failure recovery,
  engineering guardrails, and explicit pause conditions.

## Top Issues

- View-local state still lives on `Model`; screen files own render + keys only.

## Acceptance Criteria

- The prompt treats verified commits as checkpoints and instructs the agent to
  continue while useful in-scope work remains.
- It minimizes human questions without granting destructive or external scope.
- It leaves deterministic repository-backed handoffs whenever execution pauses.

## Last Verification

- Previous pass: `./scripts/verify.sh` passed on 2026-08-23 (35 E2E tests).
- Current pass: `go test ./internal/tui ./internal/storage/sqlite` passed.
- `./scripts/verify.sh` passed on 2026-08-23: Go tests, vet, offline dict.cc
  import (834,512 entries), smoke test, binary build, core E2E suite (35 passed
  in 37.43s).

## Repository State

- Prompt and continuity documentation are modified and awaiting review/commit.
