package main

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const barWidth = 20

var (
	panelStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")).
			Padding(0, 1)

	// statPanelStyle fixes a content height so the hero and enemy panels'
	// borders line up even though the hero panel has an extra stats line.
	statPanelStyle = panelStyle.Height(4)

	logPanelStyle = panelStyle.Width(60)

	messagePanelStyle = panelStyle.Width(60)

	selectedRowStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("229")).
				Bold(true)

	rowStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("250"))

	footerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))

	titleStyle = lipgloss.NewStyle().Bold(true)

	hpBarFullColor  = lipgloss.Color("42")  // green
	hpBarMidColor   = lipgloss.Color("214") // yellow
	hpBarLowColor   = lipgloss.Color("196") // red
	hpBarEmptyColor = lipgloss.Color("238")
)

// renderMessageScreen draws a static full-screen message panel: a
// title, a body, and a dim footer hint. Shared by every screen that's
// just text and a "press a key" prompt (intro, rules, reward, game
// over, victory).
func renderMessageScreen(title, body, hint string) string {
	content := body
	if title != "" {
		content = titleStyle.Render(title) + "\n\n" + body
	}
	if hint != "" {
		content += "\n\n" + footerStyle.Render(hint)
	}
	return messagePanelStyle.Render(content)
}

// renderMenuRow draws one selectable menu line, highlighted when selected.
// Shared by the combat and tavern screens so their menus look identical.
func renderMenuRow(text string, selected bool) string {
	if selected {
		return selectedRowStyle.Render("> "+text) + "\n"
	}
	return rowStyle.Render("  "+text) + "\n"
}

// hpBarColor picks a bar color from the remaining HP fraction.
func hpBarColor(pct float64) lipgloss.Color {
	switch {
	case pct > 0.5:
		return hpBarFullColor
	case pct > 0.2:
		return hpBarMidColor
	default:
		return hpBarLowColor
	}
}

// renderBar draws a fixed-width block bar for current/max.
func renderBar(current, max int) string {
	if max <= 0 {
		max = 1
	}
	pct := float64(current) / float64(max)
	if pct < 0 {
		pct = 0
	}
	filled := int(pct * float64(barWidth))
	if filled > barWidth {
		filled = barWidth
	}

	full := lipgloss.NewStyle().Foreground(hpBarColor(pct))
	empty := lipgloss.NewStyle().Foreground(hpBarEmptyColor)

	return full.Render(strings.Repeat("█", filled)) + empty.Render(strings.Repeat("░", barWidth-filled))
}
