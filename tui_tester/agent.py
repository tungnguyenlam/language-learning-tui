from typing import Optional
from .driver import Driver
from .screen import Screen
from .actions import Actions
from .waiter import Waiter

class TUIAgent:
    """
    High-level agent loop wrapper for interacting with a TUI application.
    """

    def __init__(self, command: str, columns: int = 80, lines: int = 24, env: Optional[dict] = None):
        self.command = command
        self.driver = Driver(command, columns, lines, env)
        self.screen = Screen(columns, lines)
        self.actions = Actions(self.driver)
        self.waiter = Waiter(self.driver, self.screen)
        
        self._done = False
        self.driver.spawn()
        # Initial stabilization
        try:
            self.wait_until_stable(timeout=2.0)
        except Exception:
            pass # First run might be quick or fail to stabilize if app exits immediately

    @property
    def done(self) -> bool:
        """Returns True if the underlying process has exited or the agent was closed."""
        if self._done:
            return True
        if not self.driver.is_alive():
            self._done = True
        return self._done

    def observe(self, strip_trailing_empty_lines: bool = True) -> str:
        """Captures and returns the current plain-text screen state."""
        self.waiter._sync()
        return self.screen.get_screen_text(strip_trailing_empty_lines)

    def act(self, keys: str) -> None:
        """Send keys to the TUI."""
        self.actions.send_keys(keys)

    def act_ctrl(self, char: str) -> None:
        """Send a control character to the TUI. Prefer act('<C-c>') instead."""
        self.actions.send_ctrl(char)

    def click(self, x: int, y: int, button: int = 0) -> None:
        """Send a mouse click at 1-based terminal coordinates."""
        self.actions.click(x, y, button)

    def move_mouse(self, x: int, y: int, button: int = 0) -> None:
        """Send a mouse motion at 1-based terminal coordinates."""
        self.actions.move_mouse(x, y, button)

    def drag_mouse(self, start_x: int, start_y: int, end_x: int, end_y: int, button: int = 0, steps: int = 5) -> None:
        """Simulate a mouse drag from start to end coordinates."""
        self.actions.drag_mouse(start_x, start_y, end_x, end_y, button, steps)

    def wait_until_stable(self, timeout: float = 5.0, min_stable_duration: float = 0.5) -> None:
        """Wait for the screen to stabilize after an action."""
        self.waiter.wait_for_stable(timeout, min_stable_duration)

    def wait_for_text(self, text: str, timeout: float = 5.0) -> None:
        """Wait for specific text to appear."""
        self.waiter.wait_for_text(text, timeout)

    def wait_for_regex(self, pattern: str, timeout: float = 5.0) -> None:
        """Wait for a regex pattern to appear."""
        self.waiter.wait_for_regex(pattern, timeout)

    def assert_text(self, text: str) -> None:
        """Assert text is on screen."""
        self.waiter.assert_text(text)

    def assert_not_text(self, text: str) -> None:
        """Assert text is not on screen."""
        self.waiter.assert_not_text(text)

    def assert_regex(self, pattern: str) -> None:
        """Assert a regex pattern matches on screen."""
        self.waiter.assert_regex(pattern)

    def assert_not_regex(self, pattern: str) -> None:
        """Assert a regex pattern does not match on screen."""
        self.waiter.assert_not_regex(pattern)

    def close(self) -> None:
        """Close the TUI application."""
        self.driver.close()
        self._done = True
