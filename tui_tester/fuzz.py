import time
import random
import traceback
from tui_tester.agent import TUIAgent

def main():
    agent = TUIAgent(command="go run ../cmd/deutsch-tui -test-mode")
    print("Started TUI")
    
    keys = ["<Tab>", "<Enter>", "1", "2", "3", "4", "j", "k", "h", "l", "[", "]", "e", "C"]
    try:
        for i in range(500):
            if agent.done:
                print("App exited prematurely!")
                break
            key = random.choice(keys)
            agent.act(key)
            try:
                agent.wait_until_stable(timeout=0.1, min_stable_duration=0.01)
            except Exception:
                pass
    except Exception as e:
        print(f"Crashed! {e}")
        traceback.print_exc()
    finally:
        out = agent.observe()
        agent.close()
        print(out)
        
if __name__ == "__main__":
    main()
