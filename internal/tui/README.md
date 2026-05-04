# internal/tui

Responsive Bubble Tea UI for German language learning.

## Architecture

- **Main Model**: `tui.Model` holds the global state, including the current view, active deck, and shared components like the status line and sidebars.
- **Views**: Each screen (Dashboard, Review, Browser, etc.) is a distinct component with its own `Update` and `View` logic, called by the main model.
- **Input Routing**: Keyboard events are routed to the active view. Mouse events are routed via **Hitboxes**.
- **Responsive Layout**: The UI adapts to "Wide", "Medium", and "Compact" widths using Lip Gloss frame metrics.

## Key Files

- `model.go`: The central `tea.Model` implementation.
- `render_*.go`: View-specific rendering logic.
- `handlers.go`: Command and Message handlers.
- `hitboxes.go`: Mouse interaction detection logic.
- `keys.go`: Keybinding definitions.

## Mouse Support

The TUI supports clicking and dragging.
- **Tabs**: Clickable to switch views.
- **Buttons**: Clickable in Review and Dashboard.
- **Scrollbars**: Interactive "drag-to-scroll" in the Browser and Statistics views.

## Testing

- **Unit Tests**: Test individual view logic and state transitions in `model_test.go`.
- **E2E Tests**: Use `tui-tester` (Python) to verify visual layout and complex user flows.
