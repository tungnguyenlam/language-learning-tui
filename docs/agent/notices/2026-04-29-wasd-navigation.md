# WASD Key Navigation Support

This document describes the WASD key navigation support added to the deutsch-tui application.

## Implementation Details

WASD keys have been added to support navigation in the TUI application:

- **W key**: Switches to the previous view (equivalent to left arrow key)
- **S key**: Switches to the next view (equivalent to right arrow key)  
- **A key**: Preserved for existing functions ("Again" grading, "Approve" draft)
- **D key**: Preserved for existing functions ("Discard" draft)

## Key Mapping

The WASD key mappings were added to the main switch statement in `internal/tui/model.go`:

```go
case "left", "shift+tab", "w":
    m.previousView()
case "right", "s", "d":
    m.nextView()
```

## Design Decisions

1. **Consistency with Arrow Keys**: WASD keys behave exactly like arrow keys for view switching to maintain consistency and avoid confusion.

2. **Preservation of Existing Functions**: The 'a' and 'd' keys retain their existing functions rather than being repurposed for navigation to avoid breaking existing workflows.

3. **No Contextual Navigation**: WASD keys always switch views rather than providing contextual navigation within views. This simplifies the implementation and avoids conflicts with existing key mappings.

## Testing

E2E tests have been added in `e2e_tests/test_wasd_navigation.py` to verify:
- WASD keys correctly switch between views
- Existing functions of 'a' and 'd' keys are preserved
- No regression in existing functionality

All tests pass, confirming the implementation works correctly.