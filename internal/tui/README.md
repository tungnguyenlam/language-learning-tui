# internal/tui

Responsive Bubble Tea UI for German language learning.

## Architecture

- **Main Model**: `tui.Model` holds the global state, including the current view, active deck, and shared components like the status line and sidebars.
- **Views**: Each registered screen owns its rendering and key handling. View-local state is moved
  off the main model incrementally when its cross-view dependencies are understood.
- **Input Routing**: Keyboard events are routed to the active view. Mouse events are routed via **Hitboxes**.
- **Responsive Layout**: The UI adapts to "Wide", "Medium", and "Compact" widths using Lip Gloss frame metrics.

## Key Files

- `model.go`: The central `tea.Model` implementation.
- `render_*.go`: View-specific rendering logic.
- `handlers.go`: Command and Message handlers.
- `hitboxes.go`: Mouse interaction detection logic.
- `keys.go`: Global/modal key routing, paste, and Spotlight dictionary overlay intercepts.
- `screen_*.go`: Screen registration boundaries and co-located view-specific key handling.
- `actions_backup.go`: Import/Export progress backup and restore commands (`B` / `U`).

## Mouse Support

The TUI supports clicking and dragging.
- **Tabs**: Clickable to switch views.
- **Buttons**: Clickable in Review and Dashboard.
- **Scrollbars**: Interactive "drag-to-scroll" in the Browser and Statistics views.

## Testing

- **Unit Tests**: Test individual view logic and state transitions in `model_test.go`.
- **E2E Tests**: Use `tui-tester` (Python) to verify visual layout and complex user flows.
