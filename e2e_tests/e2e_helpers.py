"""Shared helpers for E2E tests.

Tests must not hard-code the StarterDeck() size. Read the live counts from a
freshly seeded app instead, so starter content can grow without breaking
assertions. On a fresh database every starter card is new and due, so the
dashboard due count also equals the starter collection size.

See docs/agent/notices/2026-08-22-seeded-content-e2e-due-count.md.
"""

import re


def read_due_count(agent) -> int:
    """Parse the dashboard 'Due cards: N' count from the current screen.

    Call while the Dashboard is visible (e.g. right after start_agent on a
    fresh data dir).
    """
    return _read_count(agent, r"Due cards:\s+(\d+)")


def read_cards_found(agent) -> int:
    """Parse the Browser 'N cards found' count from the current screen."""
    return _read_count(agent, r"(\d+) cards found")


def _read_count(agent, pattern: str) -> int:
    screen = agent.observe()
    match = re.search(pattern, screen)
    if not match:
        raise AssertionError(f"no match for {pattern!r} on screen:\n{screen}")
    return int(match.group(1))
