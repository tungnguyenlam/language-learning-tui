package tui

import (
	"fmt"
	"strings"
	"unicode"

	tea "charm.land/bubbletea/v2"
)

// dictionaryScreen wraps the Dictionary view to satisfy the screen interface.
type dictionaryScreen struct{}

func (dictionaryScreen) Render(m *Model, layout viewportLayout) string {
	return m.renderDictionary(layout)
}

func (dictionaryScreen) HandleKey(m *Model, msg tea.KeyPressMsg) (tea.Cmd, bool) {
	oldCursor := m.dictionaryCursor
	oldDetailView := m.dictionaryDetailView

	cmd, handled := m.doUpdateDictionaryKey(msg)
	if handled {
		if m.dictionaryDetailVisible() && (m.dictionaryCursor != oldCursor || m.dictionaryDetailView != oldDetailView) {
			if m.dictionaryCursor >= 0 && m.dictionaryCursor < len(m.dictionaryResults) {
				entry := m.dictionaryResults[m.dictionaryCursor]
				m.recordDictionaryView(entry)
				relatedCmd := m.findRelatedEntries(entry.Word)
				if cmd == nil {
					cmd = relatedCmd
				} else {
					cmd = tea.Batch(cmd, relatedCmd)
				}
			}
		}
	}
	return cmd, handled
}

func (m *Model) doUpdateDictionaryKey(msg tea.KeyPressMsg) (tea.Cmd, bool) {
	key := msg.String()
	switch key {
	case "up", "k":
		if key == "k" && !m.dictionaryFocusResults && !m.dictionaryDetailView {
			break
		}
		if m.dictionaryFocusResults || m.dictionaryDetailView {
			if m.dictionaryCursor == 0 && !m.dictionaryDetailView {
				m.dictionaryFocusResults = false
				return nil, true
			}
			if m.dictionaryDetailView {
				if m.dictionaryDetailScroll > 0 {
					m.dictionaryDetailScroll--
				}
				return nil, true
			}
			if m.dictionaryCursor > 0 {
				m.dictionaryCursor--
				m.dictionaryDetailScroll = 0
			}
			return nil, true
		}
		if key == "up" && len(m.dictionarySearchHistory) > 0 {
			return m.cycleDictionaryHistory(-1), true
		}
		return nil, false
	case "down", "j":
		if !m.dictionaryFocusResults && !m.dictionaryDetailView {
			if key == "j" {
				break
			}
			if len(m.dictionaryResults) > 0 {
				m.dictionaryFocusResults = true
				return nil, true
			}
			if key == "down" {
				if len(m.dictionarySearchHistory) > 0 {
					return m.cycleDictionaryHistory(1), true
				}
				return nil, false
			}
			break
		}
		if m.dictionaryDetailView {
			maxScroll := maxInt(0, m.dictionaryDetailTotalLines-m.dictionaryDetailViewportRows(m.activeViewContentLayout()))
			if m.dictionaryDetailScroll < maxScroll {
				m.dictionaryDetailScroll++
			}
			return nil, true
		}
		if m.dictionaryCursor < len(m.dictionaryResults)-1 {
			m.dictionaryCursor++
			m.dictionaryDetailScroll = 0
		}
		return nil, true
	case "d":
		// Toggle detail only when results/detail are focused — otherwise "d"
		// must type into the search bar (der, das, denken, …).
		if m.dictionaryDetailView {
			m.dictionaryDetailView = false
			return nil, true
		}
		if m.dictionaryFocusResults && len(m.dictionaryResults) > 0 {
			m.dictionaryDetailView = true
			m.dictionaryDetailScroll = 0
			return nil, true
		}
	case "shift+up":
		if m.dictionaryDetailScroll > 0 {
			m.dictionaryDetailScroll--
		}
		return nil, true
	case "shift+down":
		maxScroll := maxInt(0, m.dictionaryDetailTotalLines-m.dictionaryDetailViewportRows(m.activeViewContentLayout()))
		if m.dictionaryDetailScroll < maxScroll {
			m.dictionaryDetailScroll++
		}
		return nil, true
	case "pgdown":
		if m.dictionaryDetailView {
			maxScroll := maxInt(0, m.dictionaryDetailTotalLines-m.dictionaryDetailViewportRows(m.activeViewContentLayout()))
			m.dictionaryDetailScroll += 10
			if m.dictionaryDetailScroll > maxScroll {
				m.dictionaryDetailScroll = maxScroll
			}
			return nil, true
		}
		m.dictionaryCursor += 10
		if m.dictionaryCursor >= len(m.dictionaryResults) {
			m.dictionaryCursor = len(m.dictionaryResults) - 1
		}
		if m.dictionaryCursor < 0 {
			m.dictionaryCursor = 0
		}
		m.dictionaryDetailScroll = 0
		return nil, true
	case "pgup":
		if m.dictionaryDetailView {
			m.dictionaryDetailScroll -= 10
			if m.dictionaryDetailScroll < 0 {
				m.dictionaryDetailScroll = 0
			}
			return nil, true
		}
		m.dictionaryCursor -= 10
		if m.dictionaryCursor < 0 {
			m.dictionaryCursor = 0
		}
		m.dictionaryDetailScroll = 0
		return nil, true
	case "ctrl+d":
		if len(m.dictionaryResults) > 0 {
			m.dictionaryDetailView = !m.dictionaryDetailView
			m.dictionaryDetailScroll = 0
		}
		return nil, true
	case "ctrl+p":
		if m.dictionaryCursor >= 0 && m.dictionaryCursor < len(m.dictionaryResults) {
			m.recordDictionarySearch(m.dictionarySearch)
			entry := m.dictionaryResults[m.dictionaryCursor]
			return m.playDictionaryAudio(entry.Word), true
		}
		return nil, true
	case "ctrl+a":
		if m.dictionaryCursor >= 0 && m.dictionaryCursor < len(m.dictionaryResults) {
			m.recordDictionarySearch(m.dictionarySearch)
			entry := m.dictionaryResults[m.dictionaryCursor]
			return m.addDictionaryEntryCmd(entry), true
		}
		return nil, true
	case "ctrl+s":
		if len(m.dictionaryResults) > 0 {
			m.recordDictionarySearch(m.dictionarySearch)
			return m.addDictionaryEntriesBatchCmd(m.dictionaryResults), true
		}
		return nil, true
	case "ctrl+b":
		if m.dictionaryCursor >= 0 && m.dictionaryCursor < len(m.dictionaryResults) {
			m.recordDictionarySearch(m.dictionarySearch)
			entry := m.dictionaryResults[m.dictionaryCursor]
			return m.toggleStarDictionaryEntry(entry), true
		}
		return nil, true
	case "b":
		if (m.dictionaryFocusResults || m.dictionaryDetailView) && m.dictionaryCursor >= 0 && m.dictionaryCursor < len(m.dictionaryResults) {
			m.recordDictionarySearch(m.dictionarySearch)
			entry := m.dictionaryResults[m.dictionaryCursor]
			return m.toggleStarDictionaryEntry(entry), true
		}
	case "a":
		// Audio only when results/detail are focused — otherwise "a" must type
		// into the search bar (Auto, haben, Fahrrad). Overlay alone is not enough:
		// Spotlight keeps prior results while the search input is active.
		if (m.dictionaryFocusResults || m.dictionaryDetailView) && m.dictionaryCursor >= 0 && m.dictionaryCursor < len(m.dictionaryResults) {
			m.recordDictionarySearch(m.dictionarySearch)
			entry := m.dictionaryResults[m.dictionaryCursor]
			return m.playDictionaryAudio(entry.Word), true
		}
	case "c":
		// Note: ctrl+c is reserved for global quit; cloze is only on plain 'c'
		// while results/detail are focused so typing in the search bar still works.
		if (m.dictionaryFocusResults || m.dictionaryDetailView) && m.dictionaryCursor >= 0 && m.dictionaryCursor < len(m.dictionaryResults) {
			m.recordDictionarySearch(m.dictionarySearch)
			entry := m.dictionaryResults[m.dictionaryCursor]
			m.recordDictionaryView(entry)
			return m.addDictionaryClozeEntryCmd(entry), true
		}
	case "r":
		if (m.dictionaryFocusResults || m.dictionaryDetailView) && m.dictionaryCursor >= 0 && m.dictionaryCursor < len(m.dictionaryResults) {
			m.recordDictionarySearch(m.dictionarySearch)
			entry := m.dictionaryResults[m.dictionaryCursor]
			m.recordDictionaryView(entry)
			return m.addDictionaryReverseEntryCmd(entry), true
		}
	case "1", "2", "3", "4", "5":
		if (m.dictionaryFocusResults || m.dictionaryDetailView) && m.dictionaryCursor >= 0 && m.dictionaryCursor < len(m.dictionaryResults) {
			entry := m.dictionaryResults[m.dictionaryCursor]
			compoundParts := m.getCompoundBreakdown(entry.Word)
			targetWord := ""
			idx := int(key[0] - '1')
			if len(compoundParts) >= 2 {
				if idx >= 0 && idx < len(compoundParts) {
					targetWord = compoundParts[idx].Word
				} else if relIdx := idx - len(compoundParts); relIdx >= 0 && relIdx < len(m.dictionaryRelatedEntries) {
					targetWord = m.dictionaryRelatedEntries[relIdx].Word
				}
			} else if idx < len(m.dictionaryRelatedEntries) {
				targetWord = m.dictionaryRelatedEntries[idx].Word
			}

			if targetWord != "" {
				m.recordDictionarySearch(m.dictionarySearch)
				m.dictionarySearch = targetWord
				m.dictionaryFocusResults = false
				m.dictionaryDetailView = false
				return m.searchDictionary(), true
			}
		}
	case "ctrl+r":
		if m.dictionaryCursor >= 0 && m.dictionaryCursor < len(m.dictionaryResults) {
			m.recordDictionarySearch(m.dictionarySearch)
			entry := m.dictionaryResults[m.dictionaryCursor]
			m.recordDictionaryView(entry)
			return m.addDictionaryReverseEntryCmd(entry), true
		}
		return nil, true
	case "ctrl+i":
		if m.dictionaryCursor >= 0 && m.dictionaryCursor < len(m.dictionaryResults) {
			m.recordDictionarySearch(m.dictionarySearch)
			entry := m.dictionaryResults[m.dictionaryCursor]
			m.recordDictionaryView(entry)
			return m.addDictionaryInflectionEntryCmd(entry), true
		}
		return nil, true
	case "i":
		if (m.dictionaryFocusResults || m.dictionaryDetailView) && m.dictionaryCursor >= 0 && m.dictionaryCursor < len(m.dictionaryResults) {
			m.recordDictionarySearch(m.dictionarySearch)
			entry := m.dictionaryResults[m.dictionaryCursor]
			m.recordDictionaryView(entry)
			return m.addDictionaryInflectionEntryCmd(entry), true
		}
	case "ctrl+g":
		return m.cycleDictionaryTargetDeck(), true
	case "ctrl+f":
		if m.dictionaryCursor >= 0 && m.dictionaryCursor < len(m.dictionaryResults) {
			m.recordDictionarySearch(m.dictionarySearch)
			entry := m.dictionaryResults[m.dictionaryCursor]
			if m.dictionaryOverlayActive {
				m.closeDictionaryOverlay()
			}
			// updateView(ViewBrowser) resets transient browser filters. Apply
			// the dictionary lookup after that reset, then issue one load so
			// the requested word is not discarded and no duplicate query races.
			m.updateView(ViewBrowser)
			m.browserSearch = entry.Word
			m.browserTag = ""
			m.searchingBrowser = false
			m.searchingTags = false
			return m.loadBrowserCards(), true
		}
		return nil, true
	case "ctrl+u":
		m.resetDictionarySearchState()
		return nil, true
	case "ctrl+x":
		if len(m.dictionarySearchHistory) > 0 {
			m.dictionarySearchHistory = nil
			m.saveDictionaryHistory()
			m.status = "Cleared recent searches"
			return nil, true
		}
		if len(m.dictionaryRecentlyViewed) > 0 {
			m.clearDictionaryRecentlyViewed()
			return nil, true
		}
		return nil, true
	case "ctrl+o":
		return m.exportDictionaryResultsTSVCmd(), true
	case "ctrl+e":
		if m.dictionaryCursor >= 0 && m.dictionaryCursor < len(m.dictionaryResults) {
			m.recordDictionarySearch(m.dictionarySearch)
			entry := m.dictionaryResults[m.dictionaryCursor]
			m.recordDictionaryView(entry)
			m.aiInput = entry.Word
			if entry.Translation != "" {
				m.aiInput = entry.Word + " — " + entry.Translation
			}
			m.draftSource = ""
			m.drafts = nil
			m.draftCursor = 0
			m.explanation = ""
			m.explainError = ""
			if m.dictionaryOverlayActive {
				m.closeDictionaryOverlay()
			}
			m.updateView(ViewAI)
			return m.explainDictionaryEntry(entry), true
		}
		return nil, true
	case "enter":
		if m.dictionaryCursor >= 0 && m.dictionaryCursor < len(m.dictionaryResults) {
			m.recordDictionarySearch(m.dictionarySearch)
			entry := m.dictionaryResults[m.dictionaryCursor]

			m.aiInput = entry.Word + " - " + entry.Translation

			var source strings.Builder
			source.WriteString(fmt.Sprintf("Word: %s\nTranslation: %s\n", entry.Word, entry.Translation))
			if entry.WordClass != "" {
				source.WriteString(fmt.Sprintf("Class: %s\n", entry.WordClass))
			}
			if entry.Gender != "" {
				source.WriteString(fmt.Sprintf("Gender: %s\n", entry.Gender))
			}
			if entry.Forms != "" {
				source.WriteString(fmt.Sprintf("Forms: %s\n", entry.Forms))
			}
			if len(entry.Examples) > 0 {
				source.WriteString("Examples:\n")
				for _, ex := range entry.Examples {
					source.WriteString(fmt.Sprintf("- %s\n", ex))
				}
			}

			m.draftSource = source.String()

			if m.dictionaryOverlayActive {
				m.closeDictionaryOverlay()
			}
			m.updateView(ViewAI)
			m.status = "Drafting flashcard from dictionary entry"
			return m.startDrafting(), true
		}
		return nil, true
	case "backspace", "\x7f", "\x08":
		if m.dictionaryDetailView {
			return nil, true // Swallowed while viewing details
		}
		if len(m.dictionarySearch) > 0 {
			m.dictionarySearch = trimLastRune(m.dictionarySearch)
			return m.debounceSearch(ViewDictionary), true
		}
		return nil, true
	case "esc":
		if m.dictionaryDetailView {
			m.dictionaryDetailView = false
			return nil, true
		}
		m.resetDictionarySearchState()
		destView := ViewDashboard
		if m.dictionaryPreviousView != "" && m.dictionaryPreviousView != ViewDictionary {
			destView = m.dictionaryPreviousView
		}
		return m.updateView(destView), true
	case "tab":
		if len(m.dictionaryResults) > 0 {
			m.dictionaryFocusResults = !m.dictionaryFocusResults
			return nil, true
		}
		return nil, false
	case "shift+tab", "left", "right":
		return nil, false
	}

	if ch, ok := singlePrintableInput(key); ok {
		if m.dictionaryDetailView {
			if ch == "i" && m.dictionaryCursor >= 0 && m.dictionaryCursor < len(m.dictionaryResults) {
				entry := m.dictionaryResults[m.dictionaryCursor]
				return m.addDictionaryInflectionEntryCmd(entry), true
			}
			if ch == " " {
				maxScroll := maxInt(0, m.dictionaryDetailTotalLines-m.dictionaryDetailViewportRows(m.activeViewContentLayout()))
				m.dictionaryDetailScroll += 5
				if m.dictionaryDetailScroll > maxScroll {
					m.dictionaryDetailScroll = maxScroll
				}
				return nil, true
			}
			return nil, true // Swallowed while viewing details
		}
		// If search is empty and it's a number key, don't handle it here to allow global navigation.
		// Standard German/English searches rarely start with a single digit, and nav parity is more important.
		if m.dictionarySearch == "" && unicode.IsDigit([]rune(ch)[0]) {
			return nil, false
		}

		m.dictionarySearch += ch
		return m.debounceSearch(ViewDictionary), true
	}

	return nil, false
}
