package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// tavernModel is the rest-stop screen: navigate the shop catalog with
// the same arrow-key interaction as combat, or leave to continue the
// journey.
type tavernModel struct {
	hero   *Hero
	cursor int
	log    []string
	left   bool
}

func newTavernModel(hero *Hero) tavernModel {
	return tavernModel{
		hero: hero,
		log:  []string{"You stop by a tavern to rest your weary bones..."},
	}
}

func (m tavernModel) Init() tea.Cmd {
	return nil
}

// leaveRow is the index of the trailing "Leave the tavern" option,
// just past the last catalog item.
func (m tavernModel) leaveRow() int {
	return len(shopCatalog)
}

func (m tavernModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}

	switch keyMsg.String() {
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < m.leaveRow() {
			m.cursor++
		}
	case "enter":
		m.selectRow()
	}

	return m, nil
}

func (m *tavernModel) selectRow() {
	if m.cursor == m.leaveRow() {
		m.left = true
		return
	}

	item := shopCatalog[m.cursor]
	if m.hero.Gold < item.Cost {
		m.appendLog("You can't afford that.")
		return
	}
	m.hero.Gold -= item.Cost
	m.appendLog(item.Use(m.hero))
}

func (m *tavernModel) appendLog(line string) {
	m.log = append(m.log, line)
	if len(m.log) > maxLogLines {
		m.log = m.log[len(m.log)-maxLogLines:]
	}
}

func (m tavernModel) View() string {
	heroPanel := statPanelStyle.Render(fmt.Sprintf(
		"%s\n%s\nHP %d/%d\nFatigue: %d\tGold: %d",
		titleStyle.Render(m.hero.Name), renderBar(m.hero.HP, m.hero.MaxHP),
		m.hero.HP, m.hero.MaxHP, m.hero.Fatigue, m.hero.Gold,
	))

	logBody := ""
	for i, line := range m.log {
		if i > 0 {
			logBody += "\n"
		}
		logBody += line
	}
	logView := logPanelStyle.Render(logBody)

	menu := ""
	for i, item := range shopCatalog {
		line := fmt.Sprintf("%s (%dg): %s", item.Name, item.Cost, item.Description)
		menu += renderMenuRow(line, i == m.cursor)
	}
	menu += renderMenuRow("Leave the tavern", m.cursor == m.leaveRow())
	menu += footerStyle.Render("↑/↓ choose · enter confirm/leave · q quit")

	return lipgloss.JoinVertical(lipgloss.Left, heroPanel, logView, menu)
}
