package main

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type combatPhase int

const (
	phaseChoosing combatPhase = iota
	phaseResolving
	phaseDone
)

const maxLogLines = 6

// enemyTurnMsg is delivered after a tea.Tick delay to run the enemy's
// attack — the TUI equivalent of the console game's Delay(1000) pause,
// without blocking the render loop.
type enemyTurnMsg struct{}

type combatModel struct {
	hero  *Hero
	enemy *Enemy
	dice  *Dice

	log    []string
	cursor int
	phase  combatPhase
}

// newCombatModel starts a fight, with its log seeded by any lead-in
// lines (encounter flavor text, a journey-start message, ...) followed
// by the "an enemy appears" line.
func newCombatModel(hero *Hero, enemy *Enemy, dice *Dice, leadIn ...string) combatModel {
	log := append([]string{}, leadIn...)
	log = append(log, fmt.Sprintf("An enemy %s appears!", enemy.Name))

	return combatModel{
		hero:  hero,
		enemy: enemy,
		dice:  dice,
		log:   log,
	}
}

func (m combatModel) Init() tea.Cmd {
	return nil
}

func (m combatModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.phase == phaseChoosing && m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.phase == phaseChoosing && m.cursor < len(m.hero.Abilities)-1 {
				m.cursor++
			}
		case "enter":
			if m.phase == phaseChoosing {
				return m.applyHeroAbility()
			}
		}
	case enemyTurnMsg:
		return m.applyEnemyAttack()
	}

	return m, nil
}

func (m combatModel) applyHeroAbility() (tea.Model, tea.Cmd) {
	ability := m.hero.Abilities[m.cursor]
	m.appendLog(ability.Execute(m.hero, m.enemy))

	if !m.enemy.IsAlive() {
		m.appendLog(fmt.Sprintf("The enemy %s dies! You are victorious!", m.enemy.Name))
		m.phase = phaseDone
		return m, nil
	}

	m.phase = phaseResolving
	return m, tea.Tick(800*time.Millisecond, func(time.Time) tea.Msg {
		return enemyTurnMsg{}
	})
}

func (m combatModel) applyEnemyAttack() (tea.Model, tea.Cmd) {
	dmg, missed := rollEnemyAttack(m.enemy, m.dice)
	if missed {
		m.appendLog(fmt.Sprintf("The enemy %s misses!", m.enemy.Name))
	} else {
		dmg -= m.hero.NextDamageReduction
		m.hero.NextDamageReduction = 0
		if dmg < 0 {
			dmg = 0
		}
		m.hero.TakeDamage(dmg)
		m.appendLog(fmt.Sprintf("The enemy %s hits for %d damage!", m.enemy.Name, dmg))
	}

	if !m.hero.IsAlive() {
		m.appendLog("Your vision fades to black... you have died.")
		m.phase = phaseDone
		return m, nil
	}

	m.phase = phaseChoosing
	return m, nil
}

func (m *combatModel) appendLog(line string) {
	m.log = append(m.log, line)
	if len(m.log) > maxLogLines {
		m.log = m.log[len(m.log)-maxLogLines:]
	}
}

func (m combatModel) View() string {
	heroPanel := statPanelStyle.Render(fmt.Sprintf(
		"%s\n%s\nHP %d/%d\nFatigue: %d\tGold: %d",
		titleStyle.Render(m.hero.Name), renderBar(m.hero.HP, m.hero.MaxHP),
		m.hero.HP, m.hero.MaxHP, m.hero.Fatigue, m.hero.Gold,
	))

	enemyPanel := statPanelStyle.Render(fmt.Sprintf(
		"%s\n%s\nHP %d/%d",
		titleStyle.Render(m.enemy.Name), renderBar(m.enemy.HP, m.enemy.MaxHP),
		m.enemy.HP, m.enemy.MaxHP,
	))

	stats := lipgloss.JoinHorizontal(lipgloss.Top, heroPanel, " ", enemyPanel)

	logBody := ""
	for i, line := range m.log {
		if i > 0 {
			logBody += "\n"
		}
		logBody += line
	}
	logView := logPanelStyle.Render(logBody)

	menu := ""
	switch m.phase {
	case phaseChoosing:
		for i, a := range m.hero.Abilities {
			menu += renderMenuRow(fmt.Sprintf("%s - %s", a.Name, a.Description), i == m.cursor)
		}
		menu += footerStyle.Render("↑/↓ choose · enter confirm · q quit")
	case phaseResolving:
		menu = footerStyle.Render("...")
	case phaseDone:
		menu = footerStyle.Render("press enter to continue · q quit")
	}

	return lipgloss.JoinVertical(lipgloss.Left, stats, logView, menu)
}
