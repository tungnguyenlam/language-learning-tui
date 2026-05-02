from .driver import Driver
import re

class Actions:
    """Handles input simulation for the TUI."""

    # Map of special keys to their escape sequences
    SPECIAL_KEYS = {
        '<Up>': '\x1b[A',
        '<Down>': '\x1b[B',
        '<Right>': '\x1b[C',
        '<Left>': '\x1b[D',
        '<Enter>': '\r',
        '<Tab>': '\t',
        '<Esc>': '\x1b',
        '<Space>': ' ',
        '<Backspace>': '\x7f',
    }

    def __init__(self, driver: Driver):
        self.driver = driver

    def send_keys(self, keys: str) -> None:
        """
        Send a string of keys. Supports special keys formatted like `<Up>`.
        Also supports control characters formatted like `<C-c>` or `<Ctrl-c>`.
        """
        parsed_keys = ""
        i = 0
        while i < len(keys):
            if keys[i] == '<':
                end = keys.find('>', i)
                if end != -1:
                    special = keys[i:end+1]
                    
                    # Check for Ctrl sequence like <C-c> or <Ctrl-c>
                    ctrl_match = re.match(r'^<(?:C|Ctrl)-([a-zA-Z\[\\\]^_])>$', special, re.IGNORECASE)
                    if ctrl_match:
                        char = ctrl_match.group(1).upper()
                        ctrl_code = chr(ord(char) - ord('A') + 1)
                        parsed_keys += ctrl_code
                        i = end + 1
                        continue
                    
                    if special in self.SPECIAL_KEYS:
                        parsed_keys += self.SPECIAL_KEYS[special]
                        i = end + 1
                        continue
            parsed_keys += keys[i]
            i += 1
            
        self.driver.write(parsed_keys)

    def send_ctrl(self, char: str) -> None:
        """
        Send a control character. E.g., 'c' for Ctrl-C.
        Preferred to use send_keys('<C-c>') instead.
        """
        if len(char) != 1:
            raise ValueError("send_ctrl expects a single character.")
        
        char = char.upper()
        if not ('A' <= char <= 'Z' or char in ('[', '\\', ']', '^', '_')):
            raise ValueError(f"Invalid control character: {char}")

        ctrl_code = chr(ord(char) - ord('A') + 1)
        self.driver.write(ctrl_code)

    def click(self, x: int, y: int, button: int = 0) -> None:
        """
        Send an xterm SGR mouse click at 1-based terminal coordinates.
        Button 0 is the primary mouse button.
        """
        if x < 1 or y < 1:
            raise ValueError("mouse coordinates are 1-based and must be positive")
        self.driver.write(f"\x1b[<{button};{x};{y}M")
        self.driver.write(f"\x1b[<{button};{x};{y}m")

    def move_mouse(self, x: int, y: int, button: int = 0) -> None:
        """
        Send an xterm SGR mouse motion at 1-based terminal coordinates.
        Button 0 is the primary mouse button.
        """
        if x < 1 or y < 1:
            raise ValueError("mouse coordinates are 1-based and must be positive")
        # Motion bit is 32
        self.driver.write(f"\x1b[<{button + 32};{x};{y}M")

    def drag_mouse(self, start_x: int, start_y: int, end_x: int, end_y: int, button: int = 0, steps: int = 5) -> None:
        """
        Simulate a mouse drag from start to end coordinates.
        """
        import time
        self.driver.write(f"\x1b[<{button};{start_x};{start_y}M")
        time.sleep(0.05)
        for i in range(1, steps + 1):
            x = start_x + (end_x - start_x) * i // steps
            y = start_y + (end_y - start_y) * i // steps
            self.move_mouse(x, y, button)
            time.sleep(0.05)
        self.driver.write(f"\x1b[<{button};{end_x};{end_y}m")
        time.sleep(0.05)
