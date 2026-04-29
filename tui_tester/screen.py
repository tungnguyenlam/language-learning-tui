import pyte
from typing import List

class Screen:
    """Manages the virtual screen buffer and ANSI parsing."""

    def __init__(self, columns: int = 80, lines: int = 24):
        self.columns = columns
        self.lines = lines
        self.screen = pyte.Screen(columns, lines)
        self.stream = pyte.Stream(self.screen)

    def feed(self, data: str) -> None:
        """Feed raw ANSI strings to the screen stream."""
        self.stream.feed(data)

    def get_screen_text(self, strip_trailing_empty_lines: bool = True) -> str:
        """
        Returns the plain-text visual content as it appears.
        Strips trailing spaces per line, and optionally strips trailing empty lines
        at the bottom of the screen to save context window space.
        """
        lines = []
        for line in self.screen.display:
            lines.append(line.rstrip())
            
        if strip_trailing_empty_lines:
            while lines and lines[-1] == "":
                lines.pop()
                
        return "\n".join(lines)

    def get_snapshot(self) -> List[str]:
        """Returns a snapshot of the current display buffer."""
        return self.screen.display.copy()

    def clear(self) -> None:
        """Clear the screen buffer."""
        self.screen.reset()
