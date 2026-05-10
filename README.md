<h1 align="center">🇩🇪 deutsch-tui</h1>

<p align="center"><em>Your terminal. Your data. German flashcards and MCQs with FSRS scheduling—fast, local-first, and absurdly pleasant to use.</em></p>

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.26-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go 1.26"/>
  <img src="https://img.shields.io/badge/Bubble%20Tea-TUI-FF69B4?style=for-the-badge&logo=terminal&logoColor=white" alt="Bubble Tea"/>
  <img src="https://img.shields.io/badge/FSRS-Spaced%20Repetition-7C3AED?style=for-the-badge&logo=clockify&logoColor=white" alt="FSRS"/>
  <img src="https://img.shields.io/badge/SQLite-Local%20First-003B57?style=for-the-badge&logo=sqlite&logoColor=white" alt="SQLite"/>
</p>

<p align="center">
  <a href="https://github.com/tungnguyenlam/language-learning-tui/blob/main/LICENSE"><img src="https://img.shields.io/github/license/tungnguyenlam/language-learning-tui?style=for-the-badge&label=License" alt="License"/></a>
  <a href="https://github.com/tungnguyenlam/language-learning-tui/stargazers"><img src="https://img.shields.io/github/stars/tungnguyenlam/language-learning-tui?style=for-the-badge&logo=github&label=Stars" alt="GitHub stars"/></a>
  <a href="https://github.com/tungnguyenlam/language-learning-tui/commits/main/"><img src="https://img.shields.io/github/last-commit/tungnguyenlam/language-learning-tui?style=for-the-badge&logo=git&label=Last%20commit" alt="Last commit"/></a>
  <a href="https://github.com/tungnguyenlam/language-learning-tui/pulls"><img src="https://img.shields.io/badge/PRs-welcome-brightgreen?style=for-the-badge&logo=github-sponsors" alt="PRs welcome"/></a>
</p>

<p align="center">
  <img src="https://via.placeholder.com/700x380/0d1117/58a6ff?text=Drop+a+asciinema+GIF+or+screenshot+here+%28docs%2Fassets%2Fdemo.gif%29" alt="deutsch-tui demo placeholder — replace with your recording" width="700"/>
</p>

<p align="center">
  <a href="#features">Features</a>
  ·
  <a href="#getting-started">Getting Started</a>
  ·
  <a href="#usage">Usage</a>
  ·
  <a href="#configuration">Configuration</a>
  ·
  <a href="#roadmap">Roadmap</a>
  ·
  <a href="#contributing">Contributing</a>
  ·
  <a href="#license">License</a>
</p>

---

<a id="features"></a>

## ✨ Features

- **🧠 FSRS scheduling** — Modern spaced repetition via [`go-fsrs`](https://github.com/open-spaced-repetition/go-fsrs), persisted in SQLite (`modernc.org/sqlite`).
- **🏠 Local-first** — Progress, decks, and settings stay on disk under one data directory; optional smoke init for automation.
- **🖱️ Serious TUI polish** — [Bubble Tea](https://github.com/charmbracelet/bubbletea) + [Lip Gloss](https://github.com/charmbracelet/lipgloss): mouse support, hitboxes, drag-to-scroll, wide/medium/compact layouts.
- **🎴 Rich card models** — Flashcards, multiple-choice, and cloze (`{{c1::…}}`) with reveal animations, typing practice, and grading (`again` / `hard` / `good` / `easy`).
- **📚 Collection tools** — Dashboard, deck list with search & limits, card browser (bookmark / suspend / tags / bulk actions), **Cram** mode for targeted drills.
- **📊 Stats & streaks** — Session metrics, accuracy, streaks, CSV export from Statistics.
- **📥 Interop** — TSV import/export and **Anki `.apkg`** import; embedded starter decks from editable TSV sources (A1–advanced tracks).
- **✍️ AI drafting (optional)** — Human-in-the-loop draft cards; works offline with template/offline providers—extend with real LLM backends when you wire them in.
- **🔊 Audio hooks** — Play pronunciation when cards carry audio paths; autoplay toggle in settings.

> **Tip:** Press `?` anytime for the in-app help overlay—keyboard-first, but the mouse works where it counts.

<a id="getting-started"></a>

## 🚀 Getting Started

### Prerequisites

| Requirement | Notes |
|-------------|--------|
| **Go** | **1.26+** (see `go.mod`). |
| **Terminal** | True-color friendly; resize to see layout breakpoints. |
| **Python + venv** *(dev only)* | For parallel E2E tests via `./scripts/verify.sh` (`tui_tester/venv`). |

### Installation

Clone and run from source:

```bash
git clone https://github.com/tungnguyenlam/language-learning-tui.git
cd language-learning-tui
go run ./cmd/deutsch-tui
```

Custom data directory (everything stays under this path):

```bash
go run ./cmd/deutsch-tui -data-dir ./my-deutsch-data
```

Initialize data files and exit (great for CI smoke checks):

```bash
go run ./cmd/deutsch-tui -data-dir ./my-deutsch-data -smoke
```

Build a static binary:

```bash
go build -o deutsch-tui ./cmd/deutsch-tui
./deutsch-tui
```

<a id="usage"></a>

## 🎮 Usage

### CLI flags

| Flag | Description |
|------|-------------|
| `-data-dir <path>` | Store `learning.db`, `config.json`, logs, and default `import.tsv` / `export.tsv` here. If omitted, uses the OS config directory + `deutsch-tui`. |
| `-smoke` | Create data dir & config, open DB, seed starter deck if empty, then exit successfully. |

### Global shortcuts

| Key | Action |
|-----|--------|
| `Tab` / `Shift+Tab` · `→` / `←` · `s` / `w` | Cycle views (when not editing text). |
| `0` `1` | Dashboard |
| `2` | Deck list |
| `3` | Review queue |
| `4` | Statistics |
| `5` | Import / export |
| `6` | AI drafts |
| `7` | Settings |
| `8` | Card browser |
| `9` | Cram |
| `[` / `]` | Previous / next deck (context-aware). |
| `?` | Toggle help overlay. |
| `Ctrl+D` | Toggle debug view. |
| `q` | Quit (or exit cram session when active). |
| `Ctrl+C` | Quit immediately. |

### Review mode (highlights)

| Key | Action |
|-----|--------|
| `Enter` / `Space` | Reveal / advance reveal / start grading. |
| `a` `h` `g` `e` | Grade **Again** / **Hard** / **Good** / **Easy** after reveal. |
| `1`–`4` | Answer MCQ choices when shown. |
| `t` | Cloze typing mode. |
| `p` | Play audio for current card. |
| `b` / `B` | Bookmark · filter bookmarks. |
| `d` | Dictionary lookup (token from prompt). |
| `x` | Suspend card. |
| `u` | Undo last review. |
| `r` | Toggle review history. |
| `f` | Focus mode. |

> **Tip:** Deck list supports `/` to search, `Enter` to jump to first match, `L` for daily limits, `c` for cram, `v` for stats—explore the footer hints per view.

<a id="configuration"></a>

## ⚙️ Configuration

All persistent settings live beside the database as **`config.json`** inside your data directory.

| Field | Purpose |
|-------|---------|
| `theme` | UI palette (`system` default—cycle in Settings). |
| `keymap` | Layout preset (`default`). |
| `ai_provider` | e.g. `disabled`, `offline`, `template`—controls draft generation behavior. |
| `log_level` | Logging verbosity (`info`, etc.). |
| `autoplay_audio` | Automatically play audio when available. |
| `strict_normalization` | Stricter answer matching for typed/cloze checks. |
| `ai_templates` | Per-topic template sets (`front` / `back` / `example`) for template-based drafting. |

Other notable paths (under the same data directory):

| Path | Role |
|------|------|
| `learning.db` | SQLite progress & cards. |
| `deutsch-tui.log` | Application log. |
| `import.tsv` / `export.tsv` | Default TSV paths surfaced in the Import view (override by editing in UI). |

<a id="roadmap"></a>

## 🗺️ Roadmap

Planned and in-flight ideas (from project backlog):

- Recently studied section on the dashboard & grammar tip of the day persistence.
- New curated decks (e.g. business vocabulary) and richer embedded content.
- Card preview in the browser, dedicated debug log view, stronger AI prompts, local LLM providers (Ollama / llama.cpp).
- Audio pronunciation stack, deck merge/split UX, custom card template UI, and expanded E2E coverage.

<a id="contributing"></a>

## 🤝 Contributing

1. Fork the repo and create a focused branch.
2. Run the full gate before pushing:

   ```bash
   ./scripts/verify.sh
   ```

   This runs `gofmt` checks, `go test ./...`, `go vet`, a TUI smoke script, and (with `tui_tester/venv`) parallel pytest over `e2e_tests/`.
3. AI-facing workflow notes live in [`AGENTS.md`](AGENTS.md); package-level contracts are in nested `AGENTS.md` files—read the area you touch before large edits.

> **Callout:** Keep contributions tight and tested; the project optimizes for maintainable agents *and* humans.

<a id="license"></a>

## 📄 License

MIT License © 2026 tungnguyen — see [`LICENSE`](LICENSE).
