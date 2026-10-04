package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type learnStage int

const (
	stageReplace learnStage = iota
	stageDone
)

// learnModel is the "what do you drop" screen: it only ever appears
// when a tavern purchase can't slot into an open ability slot because
// every slot is already a real (non-Punch) move.
type learnModel struct {
	hero *Hero

	stage   learnStage
	chosen  Ability
	cursor  int
	message string
	done    bool
}

func newLearnModel(hero *Hero, learned Ability) learnModel {
	return learnModel{
		hero:   hero,
		chosen: learned,
		stage:  stageReplace,
	}
}

// keepRow is the trailing "no change" option, just past the hero's
// current abilities.
func (m learnModel) keepRow() int {
	return len(m.hero.Abilities)
}

func (m learnModel) Init() tea.Cmd {
	return nil
}

func (m learnModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
		if m.cursor < m.keepRow() {
			m.cursor++
		}
	case "enter":
		m.selectRow()
	}

	return m, nil
}

func (m *learnModel) selectRow() {
	switch m.stage {
	case stageReplace:
		m.selectReplace()
	case stageDone:
		m.done = true
	}
}

func (m *learnModel) selectReplace() {
	if m.cursor == m.keepRow() {
		m.message = "You keep your current moves."
		m.stage = stageDone
		return
	}

	old := m.hero.Abilities[m.cursor]
	m.hero.Abilities[m.cursor] = m.chosen
	m.message = fmt.Sprintf("You swap out %s for %s!", old.Name, m.chosen.Name)
	m.stage = stageDone
}

func (m learnModel) View() string {
	switch m.stage {
	case stageReplace:
		menu := ""
		for i, a := range m.hero.Abilities {
			menu += renderMenuRow(fmt.Sprintf("%s - %s", a.Name, a.Description), i == m.cursor)
		}
		menu += renderMenuRow("Keep current kit", m.cursor == m.keepRow())
		body := fmt.Sprintf("Your kit is full (%d moves). Replace which move with %s?\n\n", maxAbilities, m.chosen.Name) + menu
		return messagePanelStyle.Render(body + "\n" + footerStyle.Render("↑/↓ choose · enter confirm"))
	default:
		return renderMessageScreen("", m.message, "press any key to continue")
	}
}
