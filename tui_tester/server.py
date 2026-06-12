import sys
import threading
import time
from xmlrpc.server import SimpleXMLRPCServer
from xmlrpc.server import SimpleXMLRPCRequestHandler
from .agent import TUIAgent

class RequestHandler(SimpleXMLRPCRequestHandler):
    rpc_paths = ('/RPC2',)

class TUIServer:
    def __init__(self, command: str, columns: int = 80, lines: int = 24):
        self.agent = TUIAgent(command, columns, lines)
        self.running = True
        
        # Start a background thread to continually drain the PTY buffer.
        # This is critical so the process doesn't block on a full stdout buffer
        # while waiting for the next client command.
        self.drain_thread = threading.Thread(target=self._drain_loop, daemon=True)
        self.drain_thread.start()

    def _drain_loop(self):
        while self.running and not self.agent.done:
            try:
                # _sync will pull nonblocking and update the screen buffer safely
                self.agent.waiter._sync()
            except Exception:
                pass
            time.sleep(0.05)

    def observe(self) -> str:
        return self.agent.observe()

    def act(self, keys: str) -> str:
        try:
            self.agent.act(keys)
            self.agent.wait_until_stable()
            return "OK"
        except Exception as e:
            return f"Error: {str(e)}"

    def act_fast(self, keys: str) -> str:
        try:
            self.agent.act(keys)
            # Very short stability check - just needs one frame of non-change
            self.agent.wait_until_stable(timeout=2.0, min_stable_duration=0.10)
            return "OK"
        except Exception as e:
            return f"Error: {str(e)}"

    def click(self, x: int, y: int, button: int = 0) -> str:
        try:
            self.agent.click(x, y, button)
            self.agent.wait_until_stable()
            return "OK"
        except Exception as e:
            return f"Error: {str(e)}"

    def drag_mouse(self, x1: int, y1: int, x2: int, y2: int, button: int = 0) -> str:
        try:
            self.agent.drag_mouse(x1, y1, x2, y2, button)
            self.agent.wait_until_stable()
            return "OK"
        except Exception as e:
            return f"Error: {str(e)}"

    def wait_for_text(self, text: str, timeout: float) -> str:
        try:
            self.agent.wait_for_text(text, timeout)
            return "OK"
        except Exception as e:
            return f"Error: {str(e)}"
            
    def wait_for_regex(self, pattern: str, timeout: float) -> str:
        try:
            self.agent.wait_for_regex(pattern, timeout)
            return "OK"
        except Exception as e:
            return f"Error: {str(e)}"

    def wait_until_stable(self, timeout: float, min_stable_duration: float) -> str:
        try:
            self.agent.wait_until_stable(timeout, min_stable_duration)
            return "OK"
        except Exception as e:
            return f"Error: {str(e)}"

    def is_done(self) -> bool:
        return self.agent.done

    def stop(self) -> str:
        self.running = False
        self.agent.close()
        return "OK"

def start_server(command: str, port: int = 8765):
    server_instance = TUIServer(command, columns=110, lines=40)
    server = SimpleXMLRPCServer(("localhost", port), requestHandler=RequestHandler, allow_none=True, logRequests=False)
    
    server.register_instance(server_instance)
    server.register_function(lambda: "OK", "ping")
    
    print(f"TUI Tester daemon started on port {port} for command: {command}")
    
    try:
        while server_instance.running and not server_instance.is_done():
            server.handle_request()
    finally:
        server_instance.stop()
        server.server_close()
