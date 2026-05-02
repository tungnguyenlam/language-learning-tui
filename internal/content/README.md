# Content Module

This module handles all deck content including:
- Starter deck (A1 Survival German) - 52 cards
- German expanded comprehensive deck (A1-B2) - 600+ cards across 600+ notes
- Anki TSV/APKG import/export
- Thematic content (travel, business, food, daily life, emergency)
- 6 ready-to-import TSV decks (500+ total cards)

## Files

- `starter.go` - Basic A1 survival German deck (52 notes)
- `german_expanded.go` - Comprehensive German vocabulary (600+ notes) covering A1-B2
- `anki.go` - Anki TSV import/export
- `apkg.go` - Anki APKG import/export
- `testdata/german-decks/` - 6 importable TSV decks (500+ cards total)
- `testdata/anki/` - Test fixtures and sample Anki decks

## German Content Overview

### A1 Essential German (Starter)
- 52 cards covering greetings, essentials, basic phrases

### A1 Comprehensive (Expanded)
- 200+ cards: greetings, verbs (conjugated), nouns, adjectives, prepositions
- Core verbs: sein, haben, gehen, kommen, machen, sprechen, essen, trinken, lesen, schreiben
- Categories: food, numbers, colors, family, body, clothes, house, places, time, weather

### A2 Level
- Travel & Transportation: 50+ cards - train, plane, car, directions, booking
- Daily Life: 50+ cards - routines, work, hobbies, household activities

### B1 Level  
- Business Professional: 97+ cards - office, meetings, contracts, negotiations
- Abstract concepts and complex expressions

### B2 Level
- Advanced vocabulary: 225+ cards - academic texts, cultural topics, idioms

### Verb Conjugations
- 251 entries: All irregular verbs, modals, and common verbs with full conjugations

## Importable TSV Decks

Located in `testdata/german-decks/`:
- `a1-essential.tsv` - 30+ cards: Core greetings and essentials
- `a1-food-drink.tsv` - 80+ cards: Food, drink, restaurant vocabulary
- `a2-travel.tsv` - 50+ cards: Transportation and travel scenarios
- `a2-daily-life.tsv` - 50+ cards: Routines, work, daily activities
- `b1-business-professional.tsv` - 97+ cards: Business and professional terms
- `b2-advanced.tsv` - 225+ cards: Academic and cultural vocabulary

## Adding New Decks

1. Create TSV files in `testdata/german-decks/` with format:
   ```
   #separator:tab
   #html:false
   #deck:Deck Name
   front\tback\textra\ttags
   ```
2. Import via UI (Import/Export tab) or programmatically
3. Test with `go test ./internal/content/...`

## Testing

Run content-specific tests:
```bash
go test ./internal/content/...

# Test specific import types
go test ./internal/content/... -run TestImportAnkiTSV
go test ./internal/content/... -run TestExportImportAPKG
```

## E2E Tests

Run all E2E tests:
```bash
. tui_tester/venv/bin/activate
python3 -m pytest e2e_tests/ -q
```

Individual test files:
- `test_apkg.py` - APKG import/export cycle
- `test_tui.py` - Core TUI flows and navigation
- `test_tab_load_commands.py` - Tab navigation and deck loading

