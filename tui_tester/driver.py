import pexpect
import os
from typing import Optional
from .exceptions import TUIProcessExitException

class Driver:
    """Manages the PTY and underlying process."""

    def __init__(self, command: str, columns: int = 80, lines: int = 24, env: Optional[dict] = None):
        self.command = command
        self.columns = columns
        self.lines = lines
        self.env = env or os.environ.copy()
        
        # Force dumb/xterm for headless operation
        self.env["TERM"] = "xterm-256color"
        
        self.child: Optional[pexpect.spawn] = None

    def spawn(self) -> None:
        """Spawn the process in a PTY."""
        self.child = pexpect.spawn(
            self.command,
            env=self.env,
            dimensions=(self.lines, self.columns),
            encoding='utf-8',
            echo=False
        )

    def read_nonblocking(self, timeout: float = 0.1) -> str:
        """Read available data without blocking indefinitely."""
        if not self.child:
            raise TUIProcessExitException("Process not started.")
        
        data = ""
        try:
            while True:
                chunk = self.child.read_nonblocking(size=4096, timeout=timeout)
                if not chunk:
                    break
                data += chunk
        except pexpect.TIMEOUT:
            pass
        except pexpect.EOF:
            pass
            
        return data

    def write(self, data: str) -> None:
        """Write data to the PTY."""
        if not self.child or not self.child.isalive():
            raise TUIProcessExitException("Cannot write, process is dead.")
        self.child.send(data)

    def close(self) -> None:
        """Terminate the process and close the PTY."""
        if self.child:
            self.child.close(force=True)

    def is_alive(self) -> bool:
        """Check if the process is running."""
        return self.child is not None and self.child.isalive()
        
    def wait(self) -> int:
        """Wait for the process to exit and return exit status."""
        if self.child:
            self.child.wait()
            return self.child.exitstatus
        return -1
