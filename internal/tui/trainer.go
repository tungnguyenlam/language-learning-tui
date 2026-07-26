package tui

import (
	"fmt"
	"math/rand/v2"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// trainerItem is the unified, display-ready model for a single practice
// question. Per-trainer specialization lives in the data produced by each
// loader (see loaders.go), not in bespoke render/key code — so adding a new
// text-input trainer is a config entry plus a loader, nothing more.
//
// The MCQ Gender trainer is deliberately NOT modeled here: it uses click
// hitboxes with hand-tuned coordinates and a fixed der/die/das choice set, so
// it keeps its own render/key path (render_gender_trainer.go).
type trainerItem struct {
	Title       string // bold accent line, always shown (word / verb / sentence-with-blank / question)
	Subtitle    string // muted line under Title, always shown when non-empty (meaning / english / help / verb info)
	PromptLine  string // bold instruction shown in both states (e.g. "Conjugate for: ich")
	Instruction string // plain instruction shown only before answering (e.g. "Fill in the blank:")
	Answer      string // correct answer, compared against user input

	RevealTitle string // bold accent shown on reveal; when empty an answer box with Answer is shown instead
	HintText    string // toggleable hint content (shown via 'h' when the trainer enables hints)
	Context     string // "Grammar Context:" block shown on reveal
	Example     string // "Example:" block shown on reveal
	Explanation string // colored explanation shown on reveal
}

// trainerConfig describes a single text-input trainer. Entries live in
// trainerConfigs, keyed by PracticeSubView.
type trainerConfig struct {
	Title      string // header, e.g. "CASE ENDING TRAINER"
	ItemNoun   string // used in status line, e.g. "case exercises"
	NextLabel  string // trailing word in "Press any key for <NextLabel>"
	EmptyMsg   string // shown when there are no items
	InputWidth int    // input/answer box width (defaults to 30)
	HintKey    bool   // whether 'h' toggles the hint line

	// match overrides the default case-insensitive equality check (used by the
	// plural trainer for umlaut/article-flexible matching).
	match func(input, answer string) bool
	Load  func(m *Model) tea.Cmd
}

// trainerState holds the live state for one trainer. A single map of these on
// the Model replaces the ~50 parallel fields the trainers used to carry.
type trainerState struct {
	config   trainerConfig
	items    []trainerItem
	index    int
	round    int // completed passes over items; drives the reshuffle
	correct  int
	total    int
	revealed bool
	lastOK   bool
	input    string
	showHint bool
}

func (st *trainerState) matches(answer string) bool {
	if st.config.match != nil {
		return st.config.match(st.input, answer)
	}
	return strings.TrimSpace(strings.ToLower(st.input)) == strings.TrimSpace(strings.ToLower(answer))
}

// advance clears the answer state and moves to the next item.
//
// Exercise sets are small (15 items for several trainers) and used to cycle in
// a fixed order forever, so a returning learner recalled the position rather
// than the grammar. Each completed pass is therefore reshuffled. The first pass
// keeps its authored order, which is pedagogically sequenced and which the E2E
// trainer tests rely on.
func (st *trainerState) advance() {
	st.revealed = false
	st.showHint = false
	st.input = ""
	if len(st.items) == 0 {
		st.index = 0
		return
	}
	st.index++
	if st.index >= len(st.items) {
		st.index = 0
		st.round++
		rand.Shuffle(len(st.items), func(i, j int) {
			st.items[i], st.items[j] = st.items[j], st.items[i]
		})
	}
}

// trainerItemsMsg delivers freshly loaded items to the named trainer.
type trainerItemsMsg struct {
	kind  PracticeSubView
	items []trainerItem
}

// trainerState returns (lazily creating) the live state for a trainer kind.
func (m *Model) trainerStateFor(kind PracticeSubView) *trainerState {
	if m.trainers == nil {
		m.trainers = make(map[PracticeSubView]*trainerState)
	}
	st, ok := m.trainers[kind]
	if !ok {
		st = &trainerState{config: trainerConfigs[kind]}
		m.trainers[kind] = st
	}
	return st
}

// isGenericTrainer reports whether a sub-view is handled by the generic
// text-input trainer (everything except the Hub and the MCQ Gender trainer).
func isGenericTrainer(kind PracticeSubView) bool {
	_, ok := trainerConfigs[kind]
	return ok
}

// updateTrainerKey is the single key handler shared by every text-input
// trainer. It preserves the exact per-trainer behavior the bespoke handlers
// had: empty input is a no-op on Enter, 'h' toggles hints only where enabled
// (and is otherwise typed), Esc clears input then exits, and any key advances
// while revealed.
func (m *Model) updateTrainerKey(kind PracticeSubView, msg tea.KeyPressMsg) (tea.Cmd, bool) {
	st := m.trainerStateFor(kind)
	if len(st.items) == 0 {
		return nil, false
	}

	key := msg.String()

	if key == "esc" {
		if st.input != "" {
			st.input = ""
		} else {
			m.practiceSubView = PracticeSubViewHub
		}
		return nil, true
	}

	if st.revealed {
		st.advance()
		return nil, true
	}

	switch key {
	case "h":
		if st.config.HintKey {
			st.showHint = !st.showHint
			return nil, true
		}
	case "enter", "\r", "\n":
		if st.input == "" {
			return nil, true
		}
		st.total++
		st.revealed = true
		if st.matches(st.items[st.index].Answer) {
			st.correct++
			st.lastOK = true
		} else {
			st.lastOK = false
		}
		return nil, true
	case "backspace":
		if len(st.input) > 0 {
			st.input = trimLastRune(st.input)
		}
		return nil, true
	case "ctrl+u":
		st.input = ""
		return nil, true
	}

	if ch, ok := singlePrintableInput(key); ok {
		st.input += ch
		return nil, true
	}

	return nil, false
}

// renderTrainer is the single renderer shared by every text-input trainer.
func (m *Model) renderTrainer(kind PracticeSubView, layout viewportLayout) string {
	st := m.trainerStateFor(kind)
	cfg := st.config

	center := func(s string) string {
		return lipgloss.PlaceHorizontal(layout.Width, lipgloss.Center, s)
	}

	if len(st.items) == 0 {
		return lipgloss.NewStyle().
			Width(layout.Width).
			Height(layout.Height).
			Align(lipgloss.Center, lipgloss.Center).
			Render(cfg.EmptyMsg)
	}

	item := st.items[st.index]
	width := cfg.InputWidth
	if width == 0 {
		width = 30
	}

	var b strings.Builder

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("212")).Underline(true)
	b.WriteString(center(titleStyle.Render(" "+cfg.Title+" ")) + "\n\n")

	accuracy := 0.0
	if st.total > 0 {
		accuracy = float64(st.correct) / float64(st.total) * 100
	}
	// "Score: x/y" stays contiguous — E2E trainer tests match on that substring.
	progress := fmt.Sprintf("Item %d/%d  •  Score: %d/%d (%.0f%%)",
		st.index+1, len(st.items), st.correct, st.total, accuracy)
	if st.round > 0 {
		progress += fmt.Sprintf("  •  Round %d", st.round+1)
	}
	b.WriteString(center(mutedStyle.Render(progress)) + "\n\n")

	accentStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("81"))
	b.WriteString(center(accentStyle.Render(item.Title)) + "\n")
	if item.Subtitle != "" {
		b.WriteString(center(mutedStyle.Render(item.Subtitle)) + "\n")
	}
	b.WriteString("\n")
	if item.PromptLine != "" {
		b.WriteString(center(lipgloss.NewStyle().Bold(true).Render(item.PromptLine)) + "\n\n")
	}

	if !st.revealed {
		if item.Instruction != "" {
			b.WriteString(center(item.Instruction) + "\n\n")
		}

		inputBox := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Padding(0, 1).
			Width(width).
			Align(lipgloss.Center).
			Render(st.input + "▌")
		b.WriteString(center(inputBox) + "\n")

		if cfg.HintKey {
			b.WriteString("\n")
			if st.showHint && item.HintText != "" {
				hintStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("86")).Italic(true)
				b.WriteString(center(hintStyle.Render("Hint: "+item.HintText)) + "\n")
			} else {
				b.WriteString(center(mutedStyle.Render("Press 'h' for a hint")) + "\n")
			}
		}
		return b.String()
	}

	// Revealed state.
	resultStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("42"))
	resultText := "CORRECT!"
	if !st.lastOK {
		resultStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("196"))
		resultText = "INCORRECT"
	}
	b.WriteString(center(resultStyle.Render(resultText)) + "\n\n")

	if !st.lastOK {
		b.WriteString(center("You typed: "+st.input) + "\n")
	}

	if item.RevealTitle != "" {
		b.WriteString(center(accentStyle.Render(item.RevealTitle)) + "\n")
	} else {
		answerBox := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Padding(0, 1).
			Width(width).
			Align(lipgloss.Center)
		b.WriteString(center(answerBox.Render(item.Answer)) + "\n")
	}

	if item.Context != "" {
		b.WriteString("\n" + center("Grammar Context:") + "\n")
		b.WriteString(center(mutedStyle.Render(item.Context)) + "\n")
	}
	if item.Example != "" {
		b.WriteString("\n" + center("Example:") + "\n")
		b.WriteString(center(mutedStyle.Render(item.Example)) + "\n")
	}
	if item.Explanation != "" {
		explanationStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("220")).Italic(true)
		b.WriteString("\n" + center(explanationStyle.Render(item.Explanation)) + "\n")
	}

	b.WriteString("\n" + center("Press any key for "+cfg.NextLabel) + "\n")

	m.hitboxes = append(m.hitboxes, Hitbox{
		ID:     "trainer-next",
		View:   ViewPractice,
		X:      layout.X,
		Y:      layout.Y,
		Width:  layout.Width,
		Height: layout.Height,
		Action: func() tea.Cmd {
			st.advance()
			return nil
		},
	})

	return b.String()
}
