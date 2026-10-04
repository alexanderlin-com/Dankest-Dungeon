package main

import "fmt"

// Hero is the player character. There's no class system — a hero
// starts with one random move from the ability pool and picks up more
// as it goes, so its identity is entirely whatever moves it's learned.
type Hero struct {
	Combatant
	Fatigue   int
	Gold      int
	Abilities []Ability
}

// NewHero builds a bare hero with a single random starting ability.
func NewHero(d *Dice) *Hero {
	return &Hero{
		Combatant: Combatant{Name: "Hero", HP: 100, MaxHP: 100},
		Abilities: randomAbilityChoices(d, nil, 1),
	}
}

// FatiguePenalty is how much damage falloff a hero suffers from fatigue.
func (h *Hero) FatiguePenalty() int {
	return h.Fatigue / 10
}

func (h *Hero) AddFatigue(amount int) {
	h.Fatigue += amount
}

func (h *Hero) ReduceFatigue(amount int) {
	h.Fatigue -= amount
	if h.Fatigue < 0 {
		h.Fatigue = 0
	}
}

// StatusText is the full ability listing shown on the rules screen.
func (h *Hero) StatusText() string {
	text := fmt.Sprintf("You are playing as %s:\n", h.Name)
	text += fmt.Sprintf("HP: %d\n", h.HP)
	text += fmt.Sprintf("Fatigue: %d\n", h.Fatigue)
	for i, a := range h.Abilities {
		text += fmt.Sprintf("%d ability: %s - %s\n", i+1, a.Name, a.Description)
	}
	return text
}

// CurrentLine is the HP/Fatigue/Gold strip shown every turn.
func (h *Hero) CurrentLine() string {
	return fmt.Sprintf("HP: %d\tFatigue: %d\tGold: %d", h.HP, h.Fatigue, h.Gold)
}
