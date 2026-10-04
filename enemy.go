package main

// Enemy is anything the hero fights.
type Enemy struct {
	Combatant
	Gold int
}

// EnemyTemplate is the stat block for one kind of enemy.
type EnemyTemplate struct {
	Name   string
	HP     int
	Attack int
	Armor  int
	MR     int
	Gold   int
}

// enemyTemplates replaces the original switch statement (whose missing
// break on case 0 meant Bandit could never actually spawn).
var enemyTemplates = []EnemyTemplate{
	{Name: "Bandit", HP: 15, Attack: 10, Armor: 5, MR: 5, Gold: 5},
	{Name: "Cultist", HP: 15, Attack: 10, Armor: 0, MR: 10, Gold: 10},
	{Name: "Undead", HP: 10, Attack: 10, Armor: 0, MR: 0, Gold: 2},
	{Name: "Knight", HP: 50, Attack: 5, Armor: 10, MR: 0, Gold: 20},
	{Name: "Demon", HP: 5, Attack: 25, Armor: 10, MR: 10, Gold: 25},
	{Name: "Slime", HP: 50, Attack: 5, Armor: 15, MR: -10, Gold: 5},
	{Name: "Big Slime", HP: 100, Attack: 10, Armor: 30, MR: -20, Gold: 10},
	{Name: "King Slime", HP: 200, Attack: 20, Armor: 60, MR: -40, Gold: 20},
	{Name: "Goblin", HP: 20, Attack: 10, Armor: 10, MR: 10, Gold: 50},
	{Name: "Skeleton", HP: 20, Attack: 10, Armor: -10, MR: -10, Gold: 0},
}

// RandomEnemy picks a random template and builds an Enemy from it.
func RandomEnemy(d *Dice) *Enemy {
	t := enemyTemplates[d.D(len(enemyTemplates))]
	return &Enemy{
		Combatant: Combatant{
			Name:   t.Name,
			HP:     t.HP,
			MaxHP:  t.HP,
			Attack: t.Attack,
			Armor:  t.Armor,
			MR:     t.MR,
		},
		Gold: t.Gold,
	}
}

var encounterFlavorText = []string{
	"You hear a rustle among the bushes...",
	"You hear footsteps in the distance...",
	"A shadowy figure approaches...",
	"You feel a hostile presence...",
	"You realize there's someone in the trees...",
	"You realize you're not alone...",
	"There's trouble afoot!",
	"An enemy monster was summoned in attack position!",
}
