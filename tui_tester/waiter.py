import time
import re
from .driver import Driver
from .screen import Screen
from .exceptions import TUITimeoutException, TUIAssertionError

class Waiter:
    """Provides smart synchronization and wait mechanisms."""

    def __init__(self, driver: Driver, screen: Screen):
        self.driver = driver
        self.screen = screen

    def _sync(self) -> None:
        """Pull new data from driver and feed to screen."""
        if self.driver.is_alive():
            data = self.driver.read_nonblocking(timeout=0.05)
            if data:
                self.screen.feed(data)

    def wait_for_text(self, text: str, timeout: float = 5.0) -> bool:
        """Wait until specific text appears on the screen."""
        start_time = time.time()
        while time.time() - start_time < timeout:
            self._sync()
            if text in self.screen.get_screen_text():
                return True
            time.sleep(0.1)
        raise TUITimeoutException(f"Timed out waiting for text: '{text}'", self.screen.get_screen_text())

    def wait_for_regex(self, pattern: str, timeout: float = 5.0) -> bool:
        """Wait until a regex pattern matches on the screen."""
        start_time = time.time()
        prog = re.compile(pattern)
        while time.time() - start_time < timeout:
            self._sync()
            if prog.search(self.screen.get_screen_text()):
                return True
            time.sleep(0.1)
        raise TUITimeoutException(f"Timed out waiting for regex: '{pattern}'", self.screen.get_screen_text())

    def wait_for_stable(self, timeout: float = 5.0, min_stable_duration: float = 0.5) -> bool:
        """
        Wait until screen content stops changing for `min_stable_duration`.
        Throws TUITimeoutException if it doesn't stabilize within `timeout`.
        """
        start_time = time.time()
        stable_start = time.time()
        last_screen_state = self.screen.get_screen_text()

        while time.time() - start_time < timeout:
            self._sync()
            current_screen_state = self.screen.get_screen_text()

            if current_screen_state != last_screen_state:
                stable_start = time.time()
                last_screen_state = current_screen_state
            elif time.time() - stable_start >= min_stable_duration:
                return True

            time.sleep(0.05)

        raise TUITimeoutException(f"Screen did not stabilize within {timeout} seconds.", self.screen.get_screen_text())

    def assert_text(self, text: str) -> None:
        """Assert that the text is currently visible on the screen."""
        self._sync()
        if text not in self.screen.get_screen_text():
            raise TUIAssertionError(f"Expected text '{text}' not found on screen.", self.screen.get_screen_text())

    def assert_not_text(self, text: str) -> None:
        """Assert that the text is NOT currently visible on the screen."""
        self._sync()
        if text in self.screen.get_screen_text():
            raise TUIAssertionError(f"Unexpected text '{text}' found on screen.", self.screen.get_screen_text())

    def assert_regex(self, pattern: str) -> None:
        """Assert that the regex pattern matches on the screen."""
        self._sync()
        if not re.search(pattern, self.screen.get_screen_text()):
            raise TUIAssertionError(f"Expected regex '{pattern}' not found on screen.", self.screen.get_screen_text())

    def assert_not_regex(self, pattern: str) -> None:
        """Assert that the regex pattern does NOT match on the screen."""
        self._sync()
        if re.search(pattern, self.screen.get_screen_text()):
            raise TUIAssertionError(f"Unexpected regex '{pattern}' found on screen.", self.screen.get_screen_text())
