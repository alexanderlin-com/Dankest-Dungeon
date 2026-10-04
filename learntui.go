package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type learnStage int

const (
	stageChoice learnStage = iota
	stageReplace
	stageDone
)

// learnModel is the "learn a move" screen shown after every combat
// win: pick one of two random candidates (or skip), and if the hero's
// already at the ability cap, pick what to drop for it.
type learnModel struct {
	hero *Hero

	stage      learnStage
	candidates []Ability
	chosen     Ability
	cursor     int
	message    string
	done       bool
}

func newLearnModel(hero *Hero, dice *Dice) learnModel {
	return learnModel{
		hero:       hero,
		candidates: randomAbilityChoices(dice, hero.Abilities, 2),
	}
}

// skipRow/keepRow are the trailing "no change" option on each stage's menu.
func (m learnModel) skipRow() int {
	return len(m.candidates)
}

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

	rowCount := m.skipRow() + 1
	if m.stage == stageReplace {
		rowCount = m.keepRow() + 1
	}

	switch keyMsg.String() {
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < rowCount-1 {
			m.cursor++
		}
	case "enter":
		m.selectRow()
	}

	return m, nil
}

func (m *learnModel) selectRow() {
	switch m.stage {
	case stageChoice:
		m.selectChoice()
	case stageReplace:
		m.selectReplace()
	case stageDone:
		m.done = true
	}
}

func (m *learnModel) selectChoice() {
	if m.cursor == m.skipRow() {
		m.message = "You decide not to learn anything new."
		m.stage = stageDone
		return
	}

	candidate := m.candidates[m.cursor]
	if len(m.hero.Abilities) < maxAbilities {
		m.hero.Abilities = append(m.hero.Abilities, candidate)
		m.message = fmt.Sprintf("You learn %s!", candidate.Name)
		m.stage = stageDone
		return
	}

	m.chosen = candidate
	m.cursor = 0
	m.stage = stageReplace
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
	case stageChoice:
		menu := ""
		for i, a := range m.candidates {
			menu += renderMenuRow(fmt.Sprintf("%s - %s", a.Name, a.Description), i == m.cursor)
		}
		menu += renderMenuRow("Skip", m.cursor == m.skipRow())
		body := "You've grown stronger! Choose a new move, or skip:\n\n" + menu
		return messagePanelStyle.Render(body + "\n" + footerStyle.Render("↑/↓ choose · enter confirm"))
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
