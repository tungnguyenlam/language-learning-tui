# Active Backlog

Last updated: 2026-05-02

## Current Milestone

Autonomous Feature Pass 17: Content Richness and Scalability

## Completed Work

### Content Expansion
- [x] Expanded starter deck from 8 to 45 notes with 52 total cards
- [x] Added thematic categories: Greetings, Food, Verbs, Numbers, Colors, Family, Places, Time
- [x] Implemented automatic MCQ generation for nouns (articles) and common verbs

### Backend Scalability & Robustness
- [x] Optimized `currentStreak` calculation to use efficient SQL and Go-based date handling
- [x] Increased card fetch limits (from 50 to 500+) in SQLite and TUI layers
- [x] Fixed potential NULL scan errors in statistics queries

### UI/UX Improvements
- [x] Added vertical panel height stabilization for consistent "rectangular region" appearance
- [x] Implemented pagination for Decks and AI Drafts views
- [x] Improved Statistics scrollbar visualization logic
- [x] Fixed all 68 E2E tests to match expanded content and new UI behaviors

## Planned Features

- [ ] Support for multiple starter decks (e.g. A1, A2, B1)
- [ ] Implement a more advanced "Spaced Repetition" visualization (heatmaps)

## Next Action

Await next user instruction.