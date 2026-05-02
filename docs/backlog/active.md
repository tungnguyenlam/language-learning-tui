# Active Backlog

Last updated: 2026-05-02

## Current Milestone

Autonomous Feature Pass 18: Final Polish and UI Interactive Elements

## Completed Work

### UI/UX Refinements
- [x] Fixed Statistics scrollbar logic to be more precise and visually representative
- [x] Implemented **interactive scrollbar clicking** for the Statistics view
- [x] Unified all breakpoints (Compact, Medium, Wide) to use a consistent bordered panel for active views
- [x] Added `statsTotalLines` tracking to accurately calculate scroll ratios
- [x] Fixed all Go unit tests for `renderStatistics`

### Bug Fixes
- [x] Added missing `strconv` import in TUI layer
- [x] Resolved regression in `renderCompact` where panels were missing
- [x] Fixed TUI unit tests failing due to signature changes

## Planned Features

- [ ] Interactive scrollbar for Browser and Cram views
- [ ] Drag-to-scroll support for scrollbars

## Next Action

Await next user instruction.