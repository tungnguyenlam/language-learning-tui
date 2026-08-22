package tui

import (
	"charm.land/lipgloss/v2"
)

var (
	// Core colors
	colorPink      = lipgloss.Color("212") // Primary active
	colorBlue      = lipgloss.Color("81")  // Action/Secondary
	colorGreen     = lipgloss.Color("46")  // Success/Progress
	colorYellow    = lipgloss.Color("226") // Warning/Highlight
	colorGold      = lipgloss.Color("220") // Milestone/Goal
	colorMuted     = lipgloss.Color("244") // Subtext
	colorHeader    = lipgloss.Color("229") // Header text
	colorPanel     = lipgloss.Color("240") // Borders
	colorAccent    = lipgloss.Color("205") // Titles
	colorStatusBg  = lipgloss.Color("236") // Status bar background
	colorProgress  = lipgloss.Color("238") // Empty progress bar
	colorSecondary = lipgloss.Color("105") // Secondary highlight
	colorOrange    = lipgloss.Color("208") // Warning/Alt highlight
	colorPurple    = lipgloss.Color("99")  // Accent/Section
	colorCyan      = lipgloss.Color("159") // Active deck highlight
	colorAITitle   = lipgloss.Color("39")  // AI view titles
	colorSuccess   = lipgloss.Color("76")  // New success indicator

	// Global styles
	headerStyle = lipgloss.NewStyle().Bold(true).Foreground(colorHeader)
	mutedStyle  = lipgloss.NewStyle().Foreground(colorMuted)
	panelStyle  = lipgloss.NewStyle().Padding(1, 2).Border(lipgloss.RoundedBorder()).BorderForeground(colorPanel)

	navStyle       = lipgloss.NewStyle().PaddingRight(2)
	navActiveStyle = lipgloss.NewStyle().Bold(true).Foreground(colorPink).PaddingRight(2)
	tabStyle       = lipgloss.NewStyle().Padding(0, 1).Foreground(lipgloss.Color("250"))
	tabActiveStyle = lipgloss.NewStyle().Padding(0, 1).Bold(true).Foreground(colorHeader).Background(lipgloss.Color("57"))
	statusStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("250")).Background(colorStatusBg).Padding(0, 1)

	// Dashboard specific
	dashTitleStyle      = lipgloss.NewStyle().Bold(true).Foreground(colorAccent)
	dashStatsStyle      = lipgloss.NewStyle().Foreground(colorYellow).Bold(true)
	dashReviewStyle     = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(colorBlue).Padding(0, 1)
	dashCollectionStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(colorOrange).Padding(0, 1)
	dashMixStyle        = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(colorSecondary).Padding(0, 1)
	dashProgressStyle   = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(0, 1)
	dashDigestStyle     = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("213")).Padding(0, 1)
	dashActivityStyle   = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(colorAITitle).Padding(0, 1)
	dashRecentStyle     = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(colorPurple).Padding(0, 1)
	dashTipStyle        = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(colorAITitle).Padding(0, 1)
	dashVerbStyle       = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(colorPink).Padding(0, 1)
	dashWordStyle       = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(colorGold).Padding(0, 1)
	dashActionsStyle    = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(colorBlue).Padding(0, 2)

	keyStyle = lipgloss.NewStyle().Foreground(colorPink).Bold(true)
	btnStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("255")).Background(lipgloss.Color("62")).Padding(0, 1).MarginRight(1)

	// Generic UI
	titleStyle   = lipgloss.NewStyle().Bold(true).Foreground(colorAccent)
	infoStyle    = lipgloss.NewStyle().Foreground(colorBlue)
	warnStyle    = lipgloss.NewStyle().Foreground(colorOrange)
	editStyle    = lipgloss.NewStyle().Bold(true).Background(lipgloss.Color("62"))
	successStyle = lipgloss.NewStyle().Foreground(colorSuccess).Bold(true)
	boldStyle    = lipgloss.NewStyle().Bold(true)
)
