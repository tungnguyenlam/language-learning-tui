# TUI Agent Rules

- **Responsive Design**: Always check `tea.WindowSizeMsg`. Test with compact (80 cols) and wide (160+ cols) layouts.
- **Mouse Safety**: If you add a clickable element, register it in `hitboxes.go`. Do NOT hardcode coordinate checks in `Update`.
- **Status Line**: Keep the main status message on the stable `status:` surface for `tui-tester` compatibility.
- **Performance**: Avoid expensive computations in `View()`. Cache results in the model if necessary.
- **Decoupling**: Do NOT call `storage` or `ai` directly from the TUI. Use the `tea.Cmd` pattern to trigger actions that interact with external services.
- **Navigation**: Support both mouse clicks and keyboard shortcuts (e.g., WASD, Arrow keys, Number keys).
