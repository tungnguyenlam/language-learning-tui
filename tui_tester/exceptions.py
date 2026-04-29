from typing import Optional

class TUITesterException(Exception):
    """Base exception for TUI Tester."""
    def __init__(self, message: str, screen_state: Optional[str] = None):
        self.screen_state = screen_state
        if screen_state:
            message = f"{message}\n\nCurrent screen state:\n{'-'*40}\n{screen_state}\n{'-'*40}"
        super().__init__(message)

class TUITimeoutException(TUITesterException):
    """Raised when an operation times out."""
    pass

class TUIProcessExitException(TUITesterException):
    """Raised when the target process exits unexpectedly."""
    pass

class TUIAssertionError(TUITesterException):
    """Raised when an assertion fails."""
    pass
