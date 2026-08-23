package tui

import (
	"fmt"
	"strings"
	"testing"
	"time"

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

var benchmarkRenderedDecks string

func BenchmarkRenderStatisticsLargeDeckCollection(b *testing.B) {
	model := NewModel(&mockRepo{}, &mockScheduler{})
	model.activeView = ViewStatistics
	model.stats.TotalCards = 500_000
	model.stats.TotalDecks = 5_000
	model.stats.ActiveDecks = 5_000
	model.stats.NewCards = 100_000
	model.stats.YoungCards = 200_000
	model.stats.MatureCards = 200_000
	model.stats.TotalReviews = 1_000_000
	model.stats.DailyGoal = 20
	model.decks = make([]core.Deck, 5_000)
	for i := range model.decks {
		model.decks[i] = core.Deck{
			ID:          fmt.Sprintf("deck-%d", i),
			Name:        fmt.Sprintf("Deutsch Übungsdeck %d", i),
			SuccessRate: 0.85,
		}
	}
	model.statsScroll = len(model.decks) / 2
	layout := viewportLayout{Width: 120, Height: 40}

	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		model.hitboxes = model.hitboxes[:0]
		benchmarkRenderedDecks = model.renderStatistics(layout)
	}
}

func BenchmarkRenderDecksLargeCollection(b *testing.B) {
	model := NewModel(&mockRepo{}, &mockScheduler{})
	model.activeView = ViewDecks
	model.width = 120
	model.height = 40
	model.breakpoint = BreakpointWide
	model.decks = make([]core.Deck, 5_000)
	for i := range model.decks {
		model.decks[i] = core.Deck{
			ID:           fmt.Sprintf("deck-%d", i),
			Name:         fmt.Sprintf("Deutsch Übungsdeck %d", i),
			Description:  "A representative deck description with Unicode: Straße und Grüße.",
			Tags:         []string{"german", "practice", fmt.Sprintf("level-%d", i%6)},
			TotalCards:   200,
			NewCards:     20,
			DueCards:     35,
			ReviewsToday: 12,
			SuccessRate:  0.85,
		}
	}
	model.deckCursor = len(model.decks) / 2
	layout := model.activeViewContentLayout()

	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		model.hitboxes = model.hitboxes[:0]
		benchmarkRenderedDecks = model.renderDecks(layout)
	}
}

func BenchmarkRenderContextWrite(b *testing.B) {
	model := NewModel(&mockRepo{}, &mockScheduler{})
	layout := viewportLayout{X: 2, Y: 1, Width: 100, Height: 40}
	sampleSingle := "A single line of text with some formatting and words"
	sampleMulti := "Line 1 of multiline text\nLine 2 of multiline text\nLine 3 of multiline text\nLine 4 of multiline text"

	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		ctx := NewRenderContext(model, layout, ViewBrowser)
		for j := 0; j < 50; j++ {
			ctx.Write(sampleSingle)
			ctx.NewLine()
			ctx.Write(sampleMulti)
			ctx.NewLine()
		}
		_ = ctx.String()
	}
}

func BenchmarkRenderBrowserLargeCollection(b *testing.B) {
	model := NewModel(&mockRepo{}, &mockScheduler{})
	model.activeView = ViewBrowser
	model.width = 120
	model.height = 40
	model.breakpoint = BreakpointWide
	model.browserCards = make([]core.Card, 5_000)
	for i := range model.browserCards {
		model.browserCards[i] = core.Card{
			ID:           fmt.Sprintf("card-%d", i),
			DeckID:       "deck-1",
			Prompt:       fmt.Sprintf("German prompt phrase %d: Straße und Grüße", i),
			Answer:       fmt.Sprintf("English translation %d: Street and greetings", i),
			Tags:         []string{"noun", "a1", "travel"},
			Reviews:      5,
			Interval:     14,
			Mature:       i%2 == 0,
			LastReviewed: time.Now().UTC(),
		}
	}
	model.browserCursor = len(model.browserCards) / 2
	layout := model.activeViewContentLayout()

	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		model.hitboxes = model.hitboxes[:0]
		benchmarkRenderedDecks = model.renderBrowserAt(layout)
	}
}

func TestRenderDecksLargeCollectionKeepsVisibleUnicodeRowsAndHitboxes(t *testing.T) {
	for _, tc := range []struct {
		name   string
		width  int
		height int
	}{
		{name: "compact", width: 48, height: 10},
		{name: "wide", width: 120, height: 14},
	} {
		t.Run(tc.name, func(t *testing.T) {
			model := NewModel(&mockRepo{}, &mockScheduler{})
			model.activeView = ViewDecks
			model.decks = make([]core.Deck, 1_000)
			for i := range model.decks {
				model.decks[i] = core.Deck{
					ID:          fmt.Sprintf("deck-%04d", i),
					Name:        fmt.Sprintf("Sammlung %04d", i),
					Description: fmt.Sprintf("Straße und Grüße: Beschreibung %04d", i),
					Tags:        []string{"Deutsch", fmt.Sprintf("Niveau-%d", i%6)},
					TotalCards:  100,
					DueCards:    10,
				}
			}
			model.deckCursor = 500
			layout := viewportLayout{X: 3, Y: 2, Width: tc.width, Height: tc.height}

			view := stripANSI(model.renderDecks(layout))
			for _, want := range []string{"Sammlung 0500", "Beschreibung 0500", "Straße", "Grüße", "Deutsch"} {
				if !strings.Contains(view, want) {
					t.Fatalf("rendered view missing visible Unicode content %q:\n%s", want, view)
				}
			}
			if strings.Contains(view, "Sammlung 0000") || strings.Contains(view, "Beschreibung 0000") {
				t.Fatalf("rendered view contains off-screen first deck:\n%s", view)
			}
			if got, want := model.deckTotalLines, len(model.decks)*4; got != want {
				t.Fatalf("deckTotalLines = %d, want %d", got, want)
			}

			var selectedHitbox *Hitbox
			scrollHitboxes := 0
			for i := range model.hitboxes {
				hit := &model.hitboxes[i]
				if hit.ID == "deck-select-500" {
					selectedHitbox = hit
				}
				if strings.HasPrefix(hit.ID, "deck-scroll-") {
					scrollHitboxes++
					if hit.X != layout.X+layout.Width-1 {
						t.Fatalf("scrollbar hitbox X = %d, want %d", hit.X, layout.X+layout.Width-1)
					}
				}
			}
			if selectedHitbox == nil || selectedHitbox.Action == nil {
				t.Fatal("selected visible deck has no actionable hitbox")
			}
			availableHeight := maxInt(5, layout.Height-2-3)
			if scrollHitboxes != availableHeight {
				t.Fatalf("scrollbar hitboxes = %d, want %d", scrollHitboxes, availableHeight)
			}

			_ = selectedHitbox.Action()
			if model.deck.ID != "deck-0500" || model.activeView != ViewDashboard {
				t.Fatalf("visible hitbox selected deck %q and view %q", model.deck.ID, model.activeView)
			}
		})
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

func TestRenderListUsesExplicitContentOffset(t *testing.T) {
	for _, tc := range []struct {
		name          string
		content       string
		contentOffset int
		scroll        int
		want          string
	}{
		{name: "full buffer with trailing line mismatch", content: "zero\none\ntwo\n", contentOffset: 0, scroll: 1, want: "one"},
		{name: "viewport slice", content: "two\nthree\n", contentOffset: 2, scroll: 2, want: "two"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			model := NewModel(&mockRepo{}, &mockScheduler{})
			totalLines := 4
			got := model.RenderList(viewportLayout{Width: 20, Height: 1}, tc.content, ListOptions{
				HitboxPrefix:  "test",
				View:          ViewStatistics,
				ContentOffset: tc.contentOffset,
				ScrollOffset:  &tc.scroll,
				TotalLines:    &totalLines,
			})
			if line := strings.TrimSpace(stripANSI(got)); !strings.HasPrefix(line, tc.want) {
				t.Fatalf("rendered line = %q, want content %q", line, tc.want)
			}
		})
	}
}

func TestRenderListCallsLazyProviderOnlyForVisibleLines(t *testing.T) {
	model := NewModel(&mockRepo{}, &mockScheduler{})
	totalLines := 100
	scroll := 40
	var called []int
	got := model.RenderList(viewportLayout{Width: 20, Height: 5}, "", ListOptions{
		HitboxPrefix: "lazy",
		View:         ViewStatistics,
		ScrollOffset: &scroll,
		TotalLines:   &totalLines,
		LineAt: func(lineIndex int) (string, bool) {
			called = append(called, lineIndex)
			return fmt.Sprintf("lazy-%d", lineIndex), true
		},
	})

	if got, want := fmt.Sprint(called), "[40 41 42 43 44]"; got != want {
		t.Fatalf("lazy provider calls = %s, want %s", got, want)
	}
	for i := 40; i <= 44; i++ {
		if !strings.Contains(stripANSI(got), fmt.Sprintf("lazy-%d", i)) {
			t.Fatalf("render missing lazy line %d:\n%s", i, stripANSI(got))
		}
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

func TestGroupedClozeAnswersRenderInOrder(t *testing.T) {
	m := NewModel(&mockRepo{}, &mockScheduler{})
	m.activeView = ViewReview
	m.width = 100
	m.height = 30
	m.breakpoint = BreakpointWide
	m.dueCards = []core.Card{{
		ID:      "cloze-1",
		Kind:    core.CardKindCloze,
		Prompt:  "Ich gebe [...] [...].",
		Answer:  "Ich gebe dem Mann.",
		Choices: []string{"dem", "Mann"},
	}}
	m.revealState = RevealRevealed

	plain := stripANSI(m.renderReview(0, 0))
	if !strings.Contains(plain, "Ich gebe dem Mann.") {
		t.Fatalf("grouped cloze answers were not rendered in order:\n%s", plain)
	}
}
