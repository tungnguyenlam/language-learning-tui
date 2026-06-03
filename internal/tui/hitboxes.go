package tui

import (
	"strconv"
	"strings"

	"deutsch-tui/internal/core"

	tea "charm.land/bubbletea/v2"
)

type Hitbox struct {
	ID     string
	View   View
	X      int
	Y      int
	Width  int
	Height int
	Action func() tea.Cmd
}

func (h Hitbox) Contains(x, y int) bool {
	return x >= h.X && x < h.X+h.Width && y >= h.Y && y < h.Y+h.Height
}

func (m *Model) hitboxAt(x, y int) (Hitbox, bool) {
	for _, hitbox := range m.hitboxes {
		if hitbox.Contains(x, y) {
			return hitbox, true
		}
	}
	return Hitbox{}, false
}

func (m *Model) activateHitbox(h Hitbox) tea.Cmd {
	if h.Action != nil {
		return h.Action()
	}
	return m.activateHitboxByID(h.ID)
}

func (m *Model) activateHitboxByID(id string) tea.Cmd {
	switch {
	case strings.HasPrefix(id, "nav-"):
		return m.updateView(View(strings.TrimPrefix(id, "nav-")))
	case strings.HasPrefix(id, "tab-"):
		return m.updateView(View(strings.TrimPrefix(id, "tab-")))
	case id == "grade-again":
		return m.gradeCard(core.GradeAgain)
	case id == "grade-hard":
		return m.gradeCard(core.GradeHard)
	case id == "grade-good":
		return m.gradeCard(core.GradeGood)
	case id == "grade-easy":
		return m.gradeCard(core.GradeEasy)
	case strings.HasPrefix(id, "draft-approve-"):
		idx, err := strconv.Atoi(strings.TrimPrefix(id, "draft-approve-"))
		if err == nil && idx >= 0 && idx < len(m.drafts) {
			m.draftCursor = idx
			return m.approveDraft()
		}
		return nil
	case strings.HasPrefix(id, "draft-discard-"):
		idx, err := strconv.Atoi(strings.TrimPrefix(id, "draft-discard-"))
		if err == nil && idx >= 0 && idx < len(m.drafts) {
			m.draftCursor = idx
			return m.discardDraft()
		}
		return nil
	case strings.HasPrefix(id, "ai-topic-"):
		topic := strings.TrimPrefix(id, "ai-topic-")
		m.aiInput = topic
		m.drafts = nil
		m.draftCursor = 0
		return m.startDrafting()
	case strings.HasPrefix(id, "practice-"):
		mode := strings.TrimPrefix(id, "practice-")
		switch mode {
		case "gender":
			return m.enterPracticeMode(PracticeSubViewGender)
		case "conjugation":
			return m.enterPracticeMode(PracticeSubViewConjugation)
		case "case":
			return m.enterPracticeMode(PracticeSubViewCase)
		case "adjective":
			return m.enterPracticeMode(PracticeSubViewAdjective)
		case "preposition":
			return m.enterPracticeMode(PracticeSubViewPreposition)
		case "plural":
			return m.enterPracticeMode(PracticeSubViewPlural)
		case "separable":
			return m.enterPracticeMode(PracticeSubViewSeparable)
		}
		return nil
	case strings.HasPrefix(id, "dash-"):
		if strings.HasPrefix(id, "dash-recent-") {
			idx, err := strconv.Atoi(strings.TrimPrefix(id, "dash-recent-"))
			if err == nil && idx < len(m.recentDecks) {
				m.selectDeckByID(m.recentDecks[idx])
				return m.updateView(ViewReview)
			}
		}
		switch id {
		case "dash-tip":
			return m.searchGrammarTipInBrowser()
		case "dash-verb":
			return m.practiceVerbOfTheDay()
		case "dash-word":
			return m.addWordOfTheDayToCollection()
		case "dash-review":
			return m.updateView(ViewReview)
		case "dash-collection":
			return m.updateView(ViewBrowser)
		case "dash-progress":
			return m.updateView(ViewStatistics)
		case "dash-digest":
			return m.updateView(ViewDecks)
		}
		return nil
	case strings.HasPrefix(id, "deck-select-"):
		idx, err := strconv.Atoi(strings.TrimPrefix(id, "deck-select-"))
		if err == nil {
			filtered := m.filteredDecks()
			if idx >= 0 && idx < len(filtered) {
				m.selectDeckByID(filtered[idx].ID)
				return m.updateView(ViewDashboard)
			}
		}
		return nil
	case strings.HasPrefix(id, "deck-study-"):
		idx, err := strconv.Atoi(strings.TrimPrefix(id, "deck-study-"))
		if err == nil {
			filtered := m.filteredDecks()
			if idx >= 0 && idx < len(filtered) {
				m.selectDeckByID(filtered[idx].ID)
				return m.updateView(ViewReview)
			}
		}
		return nil
	case strings.HasPrefix(id, "deck-cram-"):
		idx, err := strconv.Atoi(strings.TrimPrefix(id, "deck-cram-"))
		if err == nil {
			filtered := m.filteredDecks()
			if idx >= 0 && idx < len(filtered) {
				m.selectDeckByID(filtered[idx].ID)
				return m.updateView(ViewCram)
			}
		}
		return nil
	case strings.HasPrefix(id, "cram-filter-"):
		idx, err := strconv.Atoi(strings.TrimPrefix(id, "cram-filter-"))
		if err == nil {
			return m.setCramFilter(idx)
		}
		return nil
	case strings.HasPrefix(id, "settings-"):
		if id == "settings-goal-minus" {
			return m.setDailyGoal(m.stats.DailyGoal - 1)
		}
		if id == "settings-goal-plus" {
			return m.setDailyGoal(m.stats.DailyGoal + 1)
		}
		idx, err := strconv.Atoi(strings.TrimPrefix(id, "settings-"))
		if err == nil && idx >= 0 && idx <= 5 {
			m.settingsCursor = idx
			return m.handleSettingsEnter()
		}
	case strings.HasPrefix(id, "import-path-"):
		idx, err := strconv.Atoi(strings.TrimPrefix(id, "import-path-"))
		if err == nil && (idx >= 0 && idx <= 4) {
			m.importCursor = idx
			if idx < 2 {
				m.editingImportPath = true
			} else if idx == 2 {
				m.nextExportDeck()
			} else if idx == 3 {
				m.editingExportTag = true
			} else if idx == 4 {
				m.cycleExportFilter(true)
			}
			return nil
		}
		return nil
	case id == "import-tsv":
		return m.importTSV()
	case id == "import-apkg":
		return m.importAPKG()
	case id == "seed-std":
		return m.seedStandardContent()
	case id == "export-tsv":
		return m.exportTSV()
	case id == "export-apkg":
		return m.exportAPKG()
	case id == "reset-db":
		return m.handleResetDatabase()
	case strings.HasPrefix(id, "stats-scroll-"):
		line, err := strconv.Atoi(strings.TrimPrefix(id, "stats-scroll-"))
		if err == nil {
			visible := m.statisticsVisibleLines(m.activeViewContentLayout().Height)
			m.setStatsScroll(scrollOffsetForTrackRow(m.statsTotalLines, visible, line))
		}
		return nil
	case strings.HasPrefix(id, "browser-scroll-"):
		line, err := strconv.Atoi(strings.TrimPrefix(id, "browser-scroll-"))
		if err == nil && len(m.browserCards) > 0 {
			visible := m.listVisibleLines(m.activeViewContentLayout().Height)
			next := selectedIndexForTrackRow(len(m.browserCards), visible, line)
			if next != m.browserCursor {
				m.browserCursor = next
				m.clearReviewHistory()
			}
		}
		return nil
	case strings.HasPrefix(id, "cram-scroll-"):
		line, err := strconv.Atoi(strings.TrimPrefix(id, "cram-scroll-"))
		if err == nil && len(m.cramCards) > 0 {
			visible := m.listVisibleLines(m.activeViewContentLayout().Height)
			m.cramCursor = selectedIndexForTrackRow(len(m.cramCards), visible, line)
		}
		return nil
	case strings.HasPrefix(id, "deck-scroll-"):
		line, err := strconv.Atoi(strings.TrimPrefix(id, "deck-scroll-"))
		if err == nil {
			filtered := m.filteredDecks()
			if len(filtered) > 0 {
				layout := m.activeViewContentLayout()
				maxVisible := 10
				if layout.Height > 25 {
					maxVisible = (layout.Height - 10) / 2
				}
				if maxVisible < 5 {
					maxVisible = 5
				}
				m.deckCursor = selectedIndexForTrackRow(len(filtered), maxVisible, line)
			}
		}
		return nil
	}
	return nil
}
