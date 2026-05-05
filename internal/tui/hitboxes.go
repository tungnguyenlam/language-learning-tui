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

func (m *Model) activateHitbox(id string) tea.Cmd {
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
	case id == "draft-approve":
		return m.approveDraft()
	case strings.HasPrefix(id, "dash-"):
		switch id {
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
		if err == nil && (idx >= 0 && idx <= 3) {
			m.importCursor = idx
			if idx < 2 {
				m.editingImportPath = true
			} else if idx == 2 {
				m.nextExportDeck()
			} else if idx == 3 {
				m.editingExportTag = true
			}
			return nil
		}
		return nil
	case id == "import-tsv":
		return m.importTSV()
	case id == "import-apkg":
		return m.importAPKG()
	case id == "export-tsv":
		return m.exportTSV()
	case id == "export-apkg":
		return m.exportAPKG()
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
	}
	return nil
}
