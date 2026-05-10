# Active Backlog

Last updated: 2026-05-10 (Post-Improvements Phase 2)

## Current Milestone

Stability & Navigation Polish - **COMPLETED**

## Next Actions

- [x] Add "Recently Studied" section to Dashboard
- [x] Add "German Business Vocabulary" deck
- [x] Add "Card Preview" in Browser view
- [x] Add "Debug Log" view (accessible via Ctrl+D)
- [x] Implement "Focus Mode" in Review
- [x] Add Grammar Tip of the Day persistence
- [x] Add 8 new E2E tests covering new features
- [x] Improve AI prompt for better German explanations
- [ ] Add Audio pronunciation support (requires library integration)
- [ ] Implement deck merging/splitting UI
- [ ] Add local-LLM provider support (Ollama/llama.cpp)
- [ ] Implement custom card templates UI

## Recent Completed (2026-05-10)

### UI/UX
- [x] Renamed Decks view title to "DECK LIST" for better differentiation
- [x] Auto-select first matching deck on search Enter in Decks view
- [x] Added Backspace shortcut to Browser view help hints
- [x] Added session duration and speed to Statistics view

### Learning Content
- [x] Added "Common German Verbs (A1-A2)" deck
- [x] Expanded Grammar Tips with advanced C1-C2 level tips

### Core & Stability
- [x] Fixed critical panic in Decks view search/filter
- [x] Non-blocking audio playback with `playAudio`
- [x] Added database reset action with confirmation
- [x] Fixed E2E test `test_new_content_visibility.py` wrapping and matching issues
