import argparse
import sys
import os
import time
import subprocess
import xmlrpc.client

def get_client(port=8765):
    return xmlrpc.client.ServerProxy(f"http://localhost:{port}")

def check_server(port=8765):
    try:
        client = get_client(port)
        client.ping()
        return True
    except ConnectionRefusedError:
        return False

def cmd_start(args):
    if check_server(args.port):
        print(f"Error: Server already running on port {args.port}.")
        sys.exit(1)
        
    print(f"Starting TUI in background: {args.command}")
    
    # Spawn the server as a detached background process
    # We use sys.executable to ensure we use the same Python environment
    log_file = open("tui_tester_daemon.log", "w")
    proc = subprocess.Popen(
        [sys.executable, "-m", "tui_tester.cli", "daemon", args.command, str(args.port)],
        stdout=log_file,
        stderr=subprocess.STDOUT,
        start_new_session=True
    )
    
    # Wait for the server to become responsive
    start_time = time.time()
    while time.time() - start_time < 5.0:
        if check_server(args.port):
            print("Daemon started successfully.")
            return
        time.sleep(0.1)
        
    print("Error: Daemon failed to start or bind to port.")
    sys.exit(1)

def cmd_daemon(args):
    # Hidden command used to actually run the server
    from .server import start_server
    start_server(args.command, int(args.port))

def cmd_observe(args):
    try:
        client = get_client(args.port)
        print(client.observe())
    except ConnectionRefusedError:
        print("Error: Could not connect to daemon. Is it running?")
        sys.exit(1)

def cmd_act(args):
    try:
        client = get_client(args.port)
        res = client.act(args.keys)
        if res != "OK":
            print(res)
    except ConnectionRefusedError:
        print("Error: Could not connect to daemon. Is it running?")
        sys.exit(1)

def cmd_click(args):
    try:
        client = get_client(args.port)
        # Note: TUIAgent.click uses 1-based coordinates.
        # Bubble Tea mouse events are 0-based, but xterm SGR is 1-based.
        # TUIAgent.click handles the SGR protocol.
        res = client.click(args.x, args.y, args.button)
        if res != "OK":
            print(res)
    except ConnectionRefusedError:
        print("Error: Could not connect to daemon. Is it running?")
        sys.exit(1)

def cmd_drag(args):
    try:
        client = get_client(args.port)
        res = client.drag_mouse(args.x1, args.y1, args.x2, args.y2, args.button)
        if res != "OK":
            print(res)
    except ConnectionRefusedError:
        print("Error: Could not connect to daemon. Is it running?")
        sys.exit(1)

def cmd_stop(args):
    try:
        client = get_client(args.port)
        client.stop()
        print("Daemon stopped.")
    except ConnectionRefusedError:
        print("Daemon is not running.")

def main():
    parser = argparse.ArgumentParser(description="TUI Tester CLI")
    subparsers = parser.add_subparsers(dest="action", required=True)

    # Start
    start_p = subparsers.add_parser("start", help="Start the TUI daemon")
    start_p.add_argument("command", help="The terminal command to run")
    start_p.add_argument("--port", type=int, default=8765, help="Port for the daemon")

    # Observe
    obs_p = subparsers.add_parser("observe", help="Print the current screen state")
    obs_p.add_argument("--port", type=int, default=8765, help="Port for the daemon")

    # Act
    act_p = subparsers.add_parser("act", help="Send keys to the TUI (e.g., '<Enter>', '<C-c>', 'hello')")
    act_p.add_argument("keys", help="Keys to send")
    act_p.add_argument("--port", type=int, default=8765, help="Port for the daemon")

    # Click
    click_p = subparsers.add_parser("click", help="Send a mouse click")
    click_p.add_argument("x", type=int, help="X coordinate (1-based)")
    click_p.add_argument("y", type=int, help="Y coordinate (1-based)")
    click_p.add_argument("--button", type=int, default=0, help="Mouse button (0=Left)")
    click_p.add_argument("--port", type=int, default=8765, help="Port for the daemon")

    # Drag
    drag_p = subparsers.add_parser("drag", help="Simulate a mouse drag")
    drag_p.add_argument("x1", type=int)
    drag_p.add_argument("y1", type=int)
    drag_p.add_argument("x2", type=int)
    drag_p.add_argument("y2", type=int)
    drag_p.add_argument("--button", type=int, default=0)
    drag_p.add_argument("--port", type=int, default=8765, help="Port for the daemon")

    # Stop
    stop_p = subparsers.add_parser("stop", help="Stop the TUI daemon")
    stop_p.add_argument("--port", type=int, default=8765, help="Port for the daemon")
    
    # Internal Daemon entrypoint
    daemon_p = subparsers.add_parser("daemon", help=argparse.SUPPRESS)
    daemon_p.add_argument("command")
    daemon_p.add_argument("port")

    args = parser.parse_args()

    if args.action == "start":
        cmd_start(args)
    elif args.action == "daemon":
        cmd_daemon(args)
    elif args.action == "observe":
        cmd_observe(args)
    elif args.action == "act":
        cmd_act(args)
    elif args.action == "click":
        cmd_click(args)
    elif args.action == "drag":
        cmd_drag(args)
    elif args.action == "stop":
        cmd_stop(args)

if __name__ == "__main__":
    main()
