package main

import "fmt"

// Item is a one-time tavern purchase: pay gold, get an instant effect.
type Item struct {
	Name        string
	Description string
	Cost        int
	Use         func(hero *Hero) string
}

var shopCatalog = []Item{
	{
		Name:        "Hot Meal",
		Description: "heal 20 HP",
		Cost:        10,
		Use: func(hero *Hero) string {
			hero.Heal(20)
			return fmt.Sprintf("You enjoy a hot meal. %s", hero.CurrentLine())
		},
	},
	{
		Name:        "Mulled Wine",
		Description: "relieve 15 fatigue",
		Cost:        5,
		Use: func(hero *Hero) string {
			hero.ReduceFatigue(15)
			return fmt.Sprintf("You sip some mulled wine. %s", hero.CurrentLine())
		},
	},
	{
		Name:        "Greater Healing Draught",
		Description: "heal 50 HP",
		Cost:        25,
		Use: func(hero *Hero) string {
			hero.Heal(50)
			return fmt.Sprintf("You down the draught in one go. %s", hero.CurrentLine())
		},
	},
}
