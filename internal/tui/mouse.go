package tui

import (
	"strconv"
	"strings"

	tea "charm.land/bubbletea/v2"
)

// handleMouseMsg routes clicks, drags, and wheel scrolls through the hitbox
// registry. Update delegates every tea.MouseMsg here.
func (m *Model) handleMouseMsg(msg tea.MouseMsg) tea.Cmd {
	if m.confirmingDelete {
		return nil
	}
	mouse := msg.Mouse()
	m.mouseX = mouse.X
	m.mouseY = mouse.Y

	switch msg.(type) {
	case tea.MouseClickMsg, *tea.MouseClickMsg:
		if mouse.Button == tea.MouseLeft {
			if hit, ok := m.hitboxAt(mouse.X, mouse.Y); ok {
				m.logger.Debug("Mouse click at (%d, %d) on hitbox %s", mouse.X, mouse.Y, hit.ID)
				if strings.Contains(hit.ID, "-scroll-") {
					m.isDragging = true
					m.dragView = hit.View
					parts := strings.Split(hit.ID, "-")
					if len(parts) > 0 {
						if row, err := strconv.Atoi(parts[len(parts)-1]); err == nil {
							m.dragTrackStartY = mouse.Y - row
						}
					}
					layout := m.activeViewContentLayout()
					switch {
					case hit.View == ViewStatistics:
						m.dragVisible = m.statisticsVisibleLines(layout.Height)
						m.dragTotal = m.statsTotalLines
					case hit.View == ViewBrowser:
						m.dragVisible = m.listVisibleLines(layout.Height)
						m.dragTotal = len(m.browserCards)
					case hit.View == ViewCram:
						m.dragVisible = m.listVisibleLines(layout.Height)
						m.dragTotal = len(m.cramCards)
					case hit.View == ViewSettings:
						m.dragVisible = layout.Height
						m.dragTotal = m.settingsTotalLines
					case strings.HasPrefix(hit.ID, "dict-detail-scroll-"):
						m.dragView = ViewDictionary
						m.dragVisible = m.dictionaryDetailViewportRows(layout)
						m.dragTotal = m.dictionaryDetailTotalLines
						m.isDragging = m.dragTotal > m.dragVisible
					case strings.HasPrefix(hit.ID, "dict-scroll-"):
						m.dragView = ViewDictionary
						m.dragVisible = m.dictionaryListViewportRows(layout)
						m.dragTotal = len(m.dictionaryResults)
						m.isDragging = m.dragTotal > m.dragVisible
					}
				}
				return m.activateHitbox(hit)
			}
		}
	case tea.MouseMotionMsg, *tea.MouseMotionMsg:
		if m.isDragging {
			m.handleMouseDrag(mouse.Y)
		}
	case tea.MouseReleaseMsg, *tea.MouseReleaseMsg:
		if m.isDragging {
			m.logger.Debug("Mouse drag released")
			m.isDragging = false
		}
	case tea.MouseWheelMsg, *tea.MouseWheelMsg:
		m.logger.Debug("Mouse wheel event: button=%v", mouse.Button)
		delta := 0
		if mouse.Button == tea.MouseWheelUp {
			delta = -1
		} else if mouse.Button == tea.MouseWheelDown {
			delta = 1
		}
		if delta != 0 {
			if m.dictionaryOverlayActive || m.activeView == ViewDictionary {
				m.scrollDictionaryWheel(delta)
			} else {
				switch m.activeView {
				case ViewDecks:
					m.moveDeckCursor(delta)
				case ViewStatistics:
					m.scrollStats(delta)
				case ViewBrowser:
					m.moveBrowserCursor(delta)
				case ViewCram:
					m.moveCramCursor(delta)
				case ViewSettings:
					if delta < 0 {
						m.settingsScroll = maxInt(0, m.settingsScroll-1)
					} else {
						maxScroll := maxInt(0, m.settingsTotalLines-1)
						if m.settingsScroll < maxScroll {
							m.settingsScroll++
						}
					}
				case ViewPractice:
					if m.practiceSubView == PracticeSubViewHub {
						if delta < 0 {
							m.practiceScroll = maxInt(0, m.practiceScroll-1)
						} else {
							m.practiceScroll++
						}
					}
				case ViewAnkiWeb:
					m.scrollAnkiWebWheel(delta)
				}
			}
		}
	}
	return nil
}
