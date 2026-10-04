package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type screen int

const (
	screenIntro screen = iota
	screenRules
	screenReady
	screenReward
	screenLearnAbility
	screenTavern
	screenCombat
	screenGameOver
	screenVictory
)

// appModel is the root of the game: it tracks which screen is active
// and owns the data (encounter progress, pending flavor text) that
// spans screens, delegating to combatModel/tavernModel while either is
// active.
type appModel struct {
	hero *Hero
	dice *Dice

	screen screen
	combat combatModel
	tavern tavernModel
	learn  learnModel

	encounterIndex int
	encounterCount int
	lastMessage    string
}

func newAppModel(hero *Hero, dice *Dice) appModel {
	return appModel{hero: hero, dice: dice, screen: screenIntro}
}

func (m appModel) Init() tea.Cmd {
	return nil
}

func (m appModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}

	switch m.screen {
	case screenIntro:
		return m.updateIntro(msg)
	case screenRules:
		return m.updateRules(msg)
	case screenReady:
		return m.updateReady(msg)
	case screenReward:
		return m.updateReward(msg)
	case screenLearnAbility:
		return m.updateLearnAbility(msg)
	case screenTavern:
		return m.updateTavern(msg)
	case screenCombat:
		return m.updateCombat(msg)
	case screenGameOver, screenVictory:
		if _, ok := msg.(tea.KeyMsg); ok {
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m appModel) updateIntro(msg tea.Msg) (tea.Model, tea.Cmd) {
	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}

	if keyMsg.String() == "r" {
		m.screen = screenRules
	} else {
		m.screen = screenReady
	}
	return m, nil
}

func (m appModel) updateRules(msg tea.Msg) (tea.Model, tea.Cmd) {
	if _, ok := msg.(tea.KeyMsg); ok {
		m.screen = screenReady
	}
	return m, nil
}

func (m appModel) updateReady(msg tea.Msg) (tea.Model, tea.Cmd) {
	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}

	switch keyMsg.String() {
	case "y":
		m.lastMessage = "Your journey begins!"
	case "n":
		m.lastMessage = "Too bad! Your journey begins anyways!"
	default:
		m.lastMessage = "I have no idea what that input was but ok. Your journey begins anyways!"
	}

	m.encounterCount = 15 + m.dice.D(6)
	m.tavern = newTavernModel(m.hero, m.dice)
	m.screen = screenTavern
	return m, nil
}

func (m appModel) updateReward(msg tea.Msg) (tea.Model, tea.Cmd) {
	if _, ok := msg.(tea.KeyMsg); ok {
		m.lastMessage = ""
		m.advanceEncounter()
	}
	return m, nil
}

func (m appModel) updateLearnAbility(msg tea.Msg) (tea.Model, tea.Cmd) {
	updated, cmd := m.learn.Update(msg)
	m.learn = updated.(learnModel)

	if m.learn.done {
		m.screen = screenTavern
		return m, nil
	}
	return m, cmd
}

func (m appModel) updateTavern(msg tea.Msg) (tea.Model, tea.Cmd) {
	updated, cmd := m.tavern.Update(msg)
	m.tavern = updated.(tavernModel)

	if m.tavern.pendingReplace != nil {
		m.learn = newLearnModel(m.hero, *m.tavern.pendingReplace)
		m.tavern.pendingReplace = nil
		m.screen = screenLearnAbility
		return m, nil
	}
	if m.tavern.left {
		m.startEncounter()
		return m, nil
	}
	return m, cmd
}

func (m appModel) updateCombat(msg tea.Msg) (tea.Model, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.String() == "enter" && m.combat.phase == phaseDone {
		if m.combat.hero.IsAlive() {
			m.applyReward()
			m.screen = screenReward
		} else {
			m.screen = screenGameOver
		}
		return m, nil
	}

	updated, cmd := m.combat.Update(msg)
	m.combat = updated.(combatModel)
	return m, cmd
}

// startEncounter builds the next fight, seeding its log with the
// random flavor line and, right after the ready prompt, the
// journey-start message — the TUI equivalent of the console game's
// two Delay(1000) interstitial prints.
func (m *appModel) startEncounter() {
	enemy := RandomEnemy(m.dice)

	var leadIn []string
	if m.lastMessage != "" {
		leadIn = append(leadIn, m.lastMessage)
		m.lastMessage = ""
	}
	leadIn = append(leadIn, encounterFlavorText[m.dice.D(len(encounterFlavorText))])

	m.combat = newCombatModel(m.hero, enemy, m.dice, leadIn...)
	m.screen = screenCombat
}

// advanceEncounter moves past a won fight: either the journey's over,
// a tavern stop comes next, or straight into the next fight.
func (m *appModel) advanceEncounter() {
	m.encounterIndex++

	switch {
	case m.encounterIndex >= m.encounterCount:
		m.screen = screenVictory
	case m.dice.D(4) == 0:
		m.tavern = newTavernModel(m.hero, m.dice)
		m.screen = screenTavern
	default:
		m.startEncounter()
	}
}

func (m *appModel) applyReward() {
	enemy := m.combat.enemy

	fatigue := m.dice.D(10)
	m.hero.AddFatigue(fatigue)

	gold := enemy.Gold + m.dice.D(5)
	m.hero.Gold += gold

	m.lastMessage = fmt.Sprintf(
		"The enemy %s dies!\ngained %d fatigue\ngained %d gold",
		enemy.Name, fatigue, gold,
	)
}

func (m appModel) View() string {
	switch m.screen {
	case screenIntro:
		return renderMessageScreen(
			"Welcome to the Dankest Dungeon!",
			"Your adventure awaits!",
			"'r' for rules · any other key to continue",
		)
	case screenRules:
		body := "On your adventure, you will encounter enemies.\n" +
			"You must fight enemies to continue your journey.\n\n" +
			"When your HitPoints (HP) reach zero, you die.\n" +
			"You also slowly gain fatigue as you battle.\n\n" +
			"You deal less damage the more fatigue you have.\n\n" +
			"You have no fixed class - your identity is whichever moves you've picked up.\n\n" +
			fmt.Sprintf("Along your journey, you may stop by a tavern to rest, shop, and meet a trainer who'll teach you a new move for gold - you always have exactly %d moves.\n\n", maxAbilities) +
			m.hero.StatusText()
		return renderMessageScreen("Rules", body, "press any key to continue")
	case screenReady:
		return renderMessageScreen(
			"Ready to start your adventure?",
			"(y)es, (n)o, or any other key to continue anyway",
			"",
		)
	case screenReward:
		return renderMessageScreen("Victory!", m.lastMessage+"\n\n"+m.hero.CurrentLine(),
			"press any key to continue your journey...")
	case screenLearnAbility:
		return m.learn.View()
	case screenTavern:
		return m.tavern.View()
	case screenCombat:
		return m.combat.View()
	case screenGameOver:
		return renderMessageScreen(
			"You have died.",
			fmt.Sprintf("Your vision fades to black... you have died on your journey.\nFinal gold collected: %d", m.hero.Gold),
			"press any key to quit",
		)
	case screenVictory:
		return renderMessageScreen(
			"Journey complete!",
			"You have completed your journey!\n\nYou may never be the same again...",
			"press any key to quit",
		)
	}

	return ""
}
