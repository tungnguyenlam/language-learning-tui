# TTY Exploratory Testing Results - 2026-06-14

## Summary
Exploratory testing was performed via static analysis and unit testing due to a system-wide environmental blocker with `openpty`. Several logic and UX bugs were identified and fixed.

## Environmental Blocker
- **ID:** BLK-001
- **Severity:** Blocker
- **Issue:** `tui-tester` and `pexpect` fail with `OSError: out of pty devices` or `OSError: [Errno 6] Device not configured` when attempting to open a PTY on macOS.
- **Evidence:** `python -c "import os; print(os.openpty())"` fails consistently.
- **Impact:** Prevents automated TTY interaction testing.

## Identified & Fixed Bugs

### 1. Broken Rendering in Dictionary with Highlights
- **ID:** BUG-001
- **Severity:** Major (UI/UX)
- **Issue:** `padString` in `render_dictionary.go` was not ANSI-aware, leading to broken escape sequences and misaligned UI when dictionary results contained highlighted matches near the line limit.
- **Fix:** Replaced `padString` with an ANSI-aware version using `lipgloss.Width` and `truncateLine`.

### 2. Incorrect Meaning in Plural Trainer
- **ID:** BUG-002
- **Severity:** Major (Logic)
- **Issue:** `loadPluralItems` in `loaders.go` always used the `Answer` side as the meaning, even if the `Answer` side was the German side (reverse cards).
- **Fix:** Updated `loadPluralItems` to check which side is German and use the opposite side as the meaning.

### 3. Inconsistent Color Usage in Dashboard
- **ID:** BUG-003
- **Severity:** Minor (Polish)
- **Issue:** `renderDashboard` used hardcoded ANSI color strings instead of the defined style variables in `styles.go`.
- **Fix:** Refactored `progressBar` to take `color.Color` and updated `renderDashboard` to use theme-consistent variables.

## Verification
- Ran `go test ./internal/tui/...` - **PASS**
- Ran `go run ./cmd/deutsch-tui -smoke` - **PASS**
- Unit tests for Dictionary rendering and Highlight mapping - **PASS**

## Next Steps
- Investigate macOS PTY limit/configuration if testing on this machine is required.
- Proceed with `@prompt/improve.md` to further enhance the application.
