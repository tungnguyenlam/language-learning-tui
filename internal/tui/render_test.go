package tui

import (
	"fmt"
	"strings"
	"testing"

	"charm.land/lipgloss/v2"
	"deutsch-tui/internal/ankiweb"
	"deutsch-tui/internal/core"
)

func TestRenderDecksBug(t *testing.T) {
	decks := []core.Deck{
		{
			ID:          "a1-travel",
			Name:        "A1 Travel & Transport",
			Description: "Essential A1 vocabulary for traveling and transportation.",
			Tags:        []string{"german", "a1", "travel"},
		},
		{
			ID:          "a2-transport-directions",
			Name:        "A2 Transport & Directions",
			Description: "Imported from Anki TSV.",
		},
	}

	model := NewModel(&mockRepo{}, &mockScheduler{})
	model.decks = decks
	model.activeView = ViewDecks
	model.deckFilter = ""
	model.width = 120
	model.height = 60
	model.breakpoint = BreakpointWide

	layout := model.activeViewContentLayout()
	view := model.renderDecks(layout)
	t.Logf("layout width: %d, height: %d", layout.Width, layout.Height)
	lines := strings.Split(view, "\n")
	for i, l := range lines {
		t.Logf("Line %d: %q", i, l)
	}
}

func TestRenderSettingsScrollbarAlignment(t *testing.T) {
	m := NewModel(&mockRepo{}, &mockScheduler{})
	m.activeView = ViewSettings
	m.width = 60 // Small width to force truncation/wrapping
	m.height = 20
	m.breakpoint = BreakpointCompact

	// Add a very long AI template
	m.aiTemplates = map[string]map[string]string{
		"custom": {
			"front":   "This is a very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very long template",
			"back":    "short",
			"example": "short",
		},
	}
	m.aiTemplateSets = []string{"custom"}
	m.aiTemplateIndex = 0

	// Force scrollbar by having many items or small height
	m.height = 10

	view := m.renderSettings(0, 0)
	lines := strings.Split(view, "\n")

	// Check if scrollbar character '█' or '│' is at the same column for all lines that have it
	scrollbarCol := -1
	for i, l := range lines {
		plain := stripANSI(l)
		if strings.Contains(l, "█") || strings.Contains(l, "│") {
			idx := strings.LastIndex(plain, "█")
			if idx == -1 {
				idx = strings.LastIndex(plain, "│")
			}
			visualCol := lipgloss.Width(plain[:idx])
			if scrollbarCol == -1 {
				scrollbarCol = visualCol
			} else if visualCol != scrollbarCol {
				t.Errorf("Line %d: scrollbar at visual col %d, expected %d. Plain line: %q", i, visualCol, scrollbarCol, plain)
			}
		}
	}
}

func TestRenderStatisticsScrollbarAlignment(t *testing.T) {
	m := NewModel(&mockRepo{}, &mockScheduler{})
	m.activeView = ViewStatistics
	m.width = 60
	m.height = 20
	m.breakpoint = BreakpointCompact

	// Fill with some data
	m.stats.TotalCards = 1000
	m.stats.TotalDecks = 10
	m.decks = make([]core.Deck, 10)
	for i := 0; i < 10; i++ {
		m.decks[i] = core.Deck{
			ID:          fmt.Sprintf("deck-%d", i),
			Name:        fmt.Sprintf("Deck %d with a very very very very very very very very very long name", i),
			SuccessRate: 0.85,
		}
	}

	layout := m.activeViewContentLayout()
	view := m.renderStatistics(layout)
	lines := strings.Split(view, "\n")

	scrollbarCol := -1
	for i, l := range lines {
		plain := stripANSI(l)
		if strings.Contains(l, "█") || strings.Contains(l, "│") {
			idx := strings.LastIndex(plain, "█")
			if idx == -1 {
				idx = strings.LastIndex(plain, "│")
			}
			visualCol := lipgloss.Width(plain[:idx])
			if scrollbarCol == -1 {
				scrollbarCol = visualCol
			} else if visualCol != scrollbarCol {
				t.Errorf("Line %d: scrollbar at visual col %d, expected %d. Plain line: %q", i, visualCol, scrollbarCol, plain)
			}
		}
	}
}

func TestScrollbarThumbSmoothnessAndRounding(t *testing.T) {
	// Verify that scrollbarThumb distributes steps evenly without lagging at scroll boundary
	totalLines := 50
	visibleLines := 10
	maxScroll := totalLines - visibleLines                    // 40
	thumbHeight := (visibleLines * visibleLines) / totalLines // 2
	maxThumbStart := visibleLines - thumbHeight               // 8

	// At start
	start, _ := scrollbarThumb(totalLines, visibleLines, 0)
	if start != 0 {
		t.Fatalf("expected thumb start 0 at offset 0, got %d", start)
	}

	// At end
	startEnd, _ := scrollbarThumb(totalLines, visibleLines, maxScroll)
	if startEnd != maxThumbStart {
		t.Fatalf("expected thumb start %d at maxScroll, got %d", maxThumbStart, startEnd)
	}

	// Near end (one line before maxScroll), rounded thumb math should reach maxThumbStart smoothly
	startNearEnd, _ := scrollbarThumb(totalLines, visibleLines, maxScroll-1)
	if startNearEnd != maxThumbStart {
		t.Fatalf("expected thumb start near end to reach %d smoothly, got %d", maxThumbStart, startNearEnd)
	}
}

func TestRenderScrollbarColumnUnification(t *testing.T) {
	lines := []string{"Line 1", "Line 2", "Line 3", "Line 4", "Line 5"}

	// When totalLines <= visibleHeight, appends space
	unscrolled := renderScrollbarColumn(lines, 10, 5, 0)
	if len(unscrolled) != 5 || unscrolled[0] != "Line 1 " {
		t.Fatalf("unexpected unscrolled result: %v", unscrolled)
	}

	// When totalLines > visibleHeight, appends track/thumb characters
	scrolled := renderScrollbarColumn(lines, 5, 20, 0)
	if len(scrolled) != 5 {
		t.Fatalf("expected 5 scrolled lines, got %d", len(scrolled))
	}
	plain0 := stripANSI(scrolled[0])
	if !strings.HasSuffix(plain0, "█") {
		t.Fatalf("expected top line thumb █, got %q", plain0)
	}
	plainLast := stripANSI(scrolled[4])
	if !strings.HasSuffix(plainLast, "│") {
		t.Fatalf("expected bottom line track │, got %q", plainLast)
	}
}

func TestPracticeHubAutoScrollbarAndFullHeight(t *testing.T) {
	m := NewModel(nil, nil)
	m.width = 120
	m.height = 35
	m.activeView = ViewPractice
	m.practiceSubView = PracticeSubViewHub

	_, h := m.activePanelSize()
	if h < 25 {
		t.Fatalf("expected panel height to maximize terminal space (want >= 25, got %d)", h)
	}

	view := m.renderPractice(m.activeViewContentLayout())
	plainView := stripANSI(view)

	// Verify scrollbar characters rendered on right edge
	if !strings.Contains(plainView, "│") && !strings.Contains(plainView, "█") {
		t.Fatal("expected practice hub to automatically render scrollbar characters")
	}

	// Verify scrollbar hitboxes registered
	var foundPracticeScroll bool
	for _, hb := range m.hitboxes {
		if strings.HasPrefix(hb.ID, "practice-scroll-") {
			foundPracticeScroll = true
			break
		}
	}
	if !foundPracticeScroll {
		t.Fatal("expected practice-scroll- hitboxes to be registered")
	}
}

func TestFormatWrappedFooterNoTruncation(t *testing.T) {
	parts := []string{
		"tab/arrows views",
		"0-9 views",
		"? help",
		"q quit",
		"Nav: j/k",
		"Select: m",
		"History: enter",
		"Search: /",
	}

	style := lipgloss.NewStyle()
	// Test wrapping on 60-character terminal width
	wrapped := formatWrappedFooter(parts, 60, style)
	lines := strings.Split(wrapped, "\n")

	if len(lines) < 2 {
		t.Fatalf("expected footer to wrap across at least 2 lines, got %d", len(lines))
	}

	for _, p := range parts {
		if !strings.Contains(wrapped, p) {
			t.Errorf("expected wrapped footer to contain part %q without truncation", p)
		}
	}
}

func TestFilterPillDisplayNames(t *testing.T) {
	m := NewModel(nil, nil)
	m.width = 120
	m.height = 30
	pillsRow := renderFilterPillsRow(m, 0, 0, ViewDictionary, 120)
	plain := stripANSI(pillsRow)

	expectedNames := []string{"★ Starred", "DE", "EN", "Verb", "Noun", "Adj", "Adv", "Der", "Die", "Das", "Pl"}
	for _, name := range expectedNames {
		if !strings.Contains(plain, name) {
			t.Errorf("expected filter pills row to display name %q, got: %s", name, plain)
		}
	}

	rawTagsThatShouldNotBePrimaryLabels := []string{":starred", "de:", "en:", ":verb", ":noun", ":adj", ":adv", ":m", ":f", ":n", ":pl"}
	for _, tag := range rawTagsThatShouldNotBePrimaryLabels {
		if strings.Contains(plain, tag) {
			t.Errorf("expected filter pills row NOT to display raw tag %q as label, got: %s", tag, plain)
		}
	}
}

func TestWideDashboardVerbConjugationLabels(t *testing.T) {
	m := NewModel(nil, nil)
	m.width = 140
	m.height = 55
	m.activeView = ViewDashboard

	dash := m.renderDashboard(m.activeViewContentLayout())
	plain := stripANSI(dash)

	if !strings.Contains(plain, "er/sie/es") {
		t.Errorf("expected wide dashboard verb box to contain 'er/sie/es', got: %s", plain)
	}
	if !strings.Contains(plain, "sie/Sie") {
		t.Errorf("expected wide dashboard verb box to contain 'sie/Sie', got: %s", plain)
	}
}

func TestReviewGradeHitboxWithExtraContext(t *testing.T) {
	m := NewModel(nil, nil)
	m.width = 100
	m.height = 30
	m.activeView = ViewReview
	m.dueCards = []core.Card{
		{
			ID:     "c1",
			Prompt: "Hund",
			Answer: "dog",
			Extra:  "Line 1 context\nLine 2 context\nLine 3 context",
		},
	}
	m.revealState = RevealRevealed

	m.hitboxes = nil
	_ = m.renderReview(0, 0)

	var againHitbox *Hitbox
	for i, hb := range m.hitboxes {
		if hb.ID == "grade-again" {
			againHitbox = &m.hitboxes[i]
			break
		}
	}

	if againHitbox == nil {
		t.Fatal("expected grade-again hitbox to be registered")
	}

	// Calculate baseline without extra
	m.dueCards[0].Extra = ""
	m.hitboxes = nil
	_ = m.renderReview(0, 0)

	var baselineHitbox *Hitbox
	for i, hb := range m.hitboxes {
		if hb.ID == "grade-again" {
			baselineHitbox = &m.hitboxes[i]
			break
		}
	}

	if baselineHitbox == nil {
		t.Fatal("expected baseline grade-again hitbox")
	}

	// Extra content had "\n\n💡 CONTEXT: Line 1\nLine 2\nLine 3" (4 newlines total), so againHitbox.Y should be baselineHitbox.Y + 4
	expectedDiff := 4
	actualDiff := againHitbox.Y - baselineHitbox.Y
	if actualDiff != expectedDiff {
		t.Fatalf("expected Y difference of %d due to extra context lines, got %d (with extra Y=%d, baseline Y=%d)", expectedDiff, actualDiff, againHitbox.Y, baselineHitbox.Y)
	}
}

func TestCramGradeHitboxes(t *testing.T) {
	m := NewModel(nil, nil)
	m.width = 100
	m.height = 30
	m.activeView = ViewCram
	m.cramActive = true
	m.cramRevealed = true
	m.cramCards = []core.Card{
		{ID: "card1", Prompt: "Kaffee", Answer: "coffee"},
	}

	m.hitboxes = nil
	_ = m.renderCramAt(m.activeViewContentLayout())

	gradeIDs := []string{"cram-grade-again", "cram-grade-hard", "cram-grade-good", "cram-grade-easy"}
	for _, id := range gradeIDs {
		var found bool
		for _, hb := range m.hitboxes {
			if hb.ID == id {
				found = true
				if hb.Action == nil {
					t.Errorf("expected hitbox %s to have non-nil Action callback", id)
				}
				break
			}
		}
		if !found {
			t.Errorf("expected Cram mode hitbox %s to be registered when revealed", id)
		}
	}
}

func TestGenderTrainerOptionsHitboxYOffset(t *testing.T) {
	m := NewModel(nil, nil)
	m.width = 100
	m.height = 30
	m.activeView = ViewPractice
	m.practiceSubView = PracticeSubViewGender
	m.practiceItems = []practiceItem{
		{Word: "Tisch", Article: "der"},
	}
	m.practiceRevealed = false

	layout := m.activeViewContentLayout()
	m.hitboxes = nil
	_ = m.renderGenderTrainer(layout)

	for _, hb := range m.hitboxes {
		if strings.HasPrefix(hb.ID, "gender-opt-") {
			expectedY := layout.Y + 8
			if hb.Y != expectedY {
				t.Errorf("expected hitbox %s Y to be %d, got %d", hb.ID, expectedY, hb.Y)
			}
		}
	}
}

func TestAnkiWebSearchResultHitboxes(t *testing.T) {
	m := NewModel(nil, nil)
	m.width = 100
	m.height = 30
	m.activeView = ViewAnkiWeb
	m.ankiWebScreen.results = []ankiweb.Deck{
		{ID: 101, Title: "A1 German"},
		{ID: 102, Title: "A2 German"},
	}

	layout := m.activeViewContentLayout()
	m.hitboxes = nil
	_ = m.ankiWebScreen.Render(m, layout)

	var found0, found1 bool
	for _, hb := range m.hitboxes {
		if hb.ID == "ankiweb-result-0" {
			found0 = true
		}
		if hb.ID == "ankiweb-result-1" {
			found1 = true
		}
	}
	if !found0 || !found1 {
		t.Errorf("expected AnkiWeb search result hitboxes to be registered (found0=%v, found1=%v)", found0, found1)
	}
}

func TestResponsiveHelpLayout(t *testing.T) {
	m := NewModel(nil, nil)
	m.width = 80
	m.height = 30

	helpNarrow := m.renderHelp(viewportLayout{Width: 80, Height: 30})
	plainNarrow := stripANSI(helpNarrow)

	if !strings.Contains(plainNarrow, "Keyboard Shortcuts") {
		t.Error("expected help output to contain header")
	}
	// On 80-width screen, help should render in multi-line stacked/2x2 grid layout
	if !strings.Contains(plainNarrow, "Dictionary:") || !strings.Contains(plainNarrow, "Review:") {
		t.Error("expected 80-width help to contain Dictionary and Review shortcut sections")
	}
}
