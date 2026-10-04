package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// moveCost is a flat gold price for any move offered at a tavern - a
// deliberately simple rule instead of per-ability pricing.
const moveCost = 20

// tavernModel is the rest-stop screen: navigate the shop catalog (and
// two randomly offered moves to learn) with the same arrow-key
// interaction as combat, or leave to continue the journey.
type tavernModel struct {
	hero       *Hero
	moveOffers []Ability
	cursor     int
	log        []string
	left       bool

	// pendingReplace signals a bought move couldn't slot in anywhere
	// (every ability slot is already a real move) - the parent builds
	// a learnModel to ask what to drop, then returns here.
	pendingReplace *Ability
}

func newTavernModel(hero *Hero, dice *Dice) tavernModel {
	return tavernModel{
		hero:       hero,
		moveOffers: randomAbilityChoices(dice, hero.Abilities, 2),
		log:        []string{"You stop by a tavern to rest your weary bones..."},
	}
}

func (m tavernModel) Init() tea.Cmd {
	return nil
}

// moveOffersStart/leaveRow mark where each section of the menu begins:
// potions, then move offers, then "Leave the tavern".
func (m tavernModel) moveOffersStart() int {
	return len(shopCatalog)
}

func (m tavernModel) leaveRow() int {
	return m.moveOffersStart() + len(m.moveOffers)
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
	switch {
	case m.cursor == m.leaveRow():
		m.left = true
	case m.cursor >= m.moveOffersStart():
		m.buyMove(m.cursor - m.moveOffersStart())
	default:
		m.buyItem(m.cursor)
	}
}

func (m *tavernModel) buyItem(i int) {
	item := shopCatalog[i]
	if m.hero.Gold < item.Cost {
		m.appendLog("You can't afford that.")
		return
	}
	m.hero.Gold -= item.Cost
	m.appendLog(item.Use(m.hero))
}

func (m *tavernModel) buyMove(i int) {
	if m.hero.Gold < moveCost {
		m.appendLog("You can't afford that.")
		return
	}
	learned := m.moveOffers[i]
	m.hero.Gold -= moveCost
	m.moveOffers = append(m.moveOffers[:i], m.moveOffers[i+1:]...)
	if m.cursor > m.leaveRow() {
		m.cursor = m.leaveRow()
	}

	for slot, a := range m.hero.Abilities {
		if a.Name == "Punch" {
			m.hero.Abilities[slot] = learned
			m.appendLog(fmt.Sprintf("A trainer teaches you %s!", learned.Name))
			return
		}
	}

	m.pendingReplace = &learned
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
	for i, a := range m.moveOffers {
		line := fmt.Sprintf("Learn: %s (%dg) - %s", a.Name, moveCost, a.Description)
		menu += renderMenuRow(line, m.moveOffersStart()+i == m.cursor)
	}
	menu += renderMenuRow("Leave the tavern", m.cursor == m.leaveRow())
	menu += footerStyle.Render("↑/↓ choose · enter confirm/leave · q quit")

	return lipgloss.JoinVertical(lipgloss.Left, heroPanel, logView, menu)
}
