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
- Enhanced `tui-tester` with `move_mouse` and `drag_mouse` simulation capabilities.
- Verified all 76 E2E tests passing.
- **Parallelized E2E tests** with `pytest-xdist`, reducing full verification time from ~5.5 minutes to ~45 seconds.
- Verified all Go unit tests passing.
- Updated project notices and index to reflect architectural requirements for reveal synchronization.
