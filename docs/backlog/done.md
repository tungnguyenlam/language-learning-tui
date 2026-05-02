# Done Backlog

## 2026-05-02 - Review Reveal Sync and Interactive Testing

- Fixed race condition in Review view where grading hints were displayed before hitboxes and answers were ready.
- Refactored `renderReview` to strictly synchronize visual state with interaction state.
- Implemented **drag-to-scroll support** for scrollbars in Statistics, Browser, and Cram views.
- **Enhanced Decks view** with separate counts for New, Due, and Total cards, providing better progress visibility.
- **Added Session Statistics** to the status line, showing accuracy and cards reviewed in the current session.
- **Improved AI view UX** with interactive [Approve] and [Discard] buttons and a detailed draft preview section.
- **Refactored Card Browser** with a modal **Search Mode** (triggered by `/`) and direct card management via **b** (bookmark) and **x** (suspend).
- **Enhanced Settings view** with interactive hitboxes and clickable [+] and [-] buttons for adjusting the daily goal.
- **Enhanced Import view** with a modern button-like layout and interactive hitboxes for path editing and triggering actions.
- **Added interactive Dashboard hitboxes** linking to Review, Browser, and Statistics views.
- **Added interactive Cram filters** allowing mouse-based switching of card sets.
- **Fixed multiple edge-case bugs** in session statistics (undo support), state management (resetting dragging on view/size change), and MCQ interaction state leakage.
- **Implemented robust ID-based deck synchronization** preventing visual and functional desyncs across all collection-refresh actions.
- **Resolved navigation "traps"** by ensuring global keyboard shortcuts (1-9, ?) take precedence over unhandled local view keys.
- **Fixed deck selection bug** and **secured Cram mode cursor**, improving the reliability of filtered collection management.
- **Fixed audio auto-play and error reporting**, ensuring reliable playback feedback for learners.
- **Documented TUI Tester** as a standalone utility with detailed `README.md` and `AGENTS.md` guides, including instructions for publishing as a Gemini Skill or MCP server.
- **Cleaned up TUI Tester directory** by removing redundant internal `e2e_tests` to ensure a clean standalone repository.
- Enhanced `tui-tester` with `move_mouse` and `drag_mouse` simulation capabilities.
- Verified all 76 E2E tests passing.
- **Parallelized E2E tests** with `pytest-xdist`, reducing full verification time from ~5.5 minutes to ~45 seconds.
- Verified all Go unit tests passing.
- Updated project notices and index to reflect architectural requirements for reveal synchronization.

## Animation and Appearance Polish

- **Smother reveal animation**: Increased from 5 to 10 steps (20% → 10% per tick) for smoother text reveal
- **Extended reveal timing**: Increased tick interval from 50ms to 60ms for more comfortable pacing (600ms total)
- **Modern block characters**: Replaced `█` with `▌` in reveal animations for a more elegant look
- **Color palette refresh**: Updated all view titles to use accent color (159) for better consistency
- **Session accuracy colors**: Added dynamic color coding (green/yellow/red) based on performance thresholds
- **Dashboard border colors**: Harmonized with accent color scheme for better visual cohesion
- **Browser search active state**: Enhanced with accent color when search is active
- **Streak/maturity indicators**: Changed from 🔥/⭐ to ✨ for a more modern appearance
- **Refined typography**: Updated muted text from 244 to 248, headers from 229 to 231 for better contrast