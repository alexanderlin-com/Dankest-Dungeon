package main

import "fmt"

// Ability is a single usable combat move. Representing it as a struct
// with an Execute closure (rather than one Go type per move) keeps new
// abilities a one-liner instead of a new class.
type Ability struct {
	Name        string
	Description string
	Execute     func(user *Hero, target *Enemy) string
}

// rollDamage applies the shared damage formula: base minus the defense
// stat, nudged by a small random modifier, with a 1-in-10 chance to miss.
func rollDamage(d *Dice, base, defense int) (dmg int, missed bool) {
	if d.D(10) == 0 {
		return 0, true
	}
	mod := d.D(5)
	if d.Bool() {
		dmg = base - defense + mod
	} else {
		dmg = base - defense - mod
	}
	if dmg < 0 {
		dmg = 0
	}
	return dmg, false
}

func applyFatigue(user *Hero, dmg int) int {
	dmg -= user.FatiguePenalty()
	if dmg < 0 {
		dmg = 0
	}
	return dmg
}

// NewPhysicalAttack builds an ability that rolls against the target's armor.
func NewPhysicalAttack(name, description, hitVerb string, baseAttack int, d *Dice) Ability {
	return Ability{
		Name:        name,
		Description: description,
		Execute: func(user *Hero, target *Enemy) string {
			dmg, missed := rollDamage(d, baseAttack, target.Armor)
			if missed {
				return "Your attack misses!\n"
			}
			dmg = applyFatigue(user, dmg)
			target.TakeDamage(dmg)
			return fmt.Sprintf("%s the enemy %s, dealing %d damage!\n", hitVerb, target.Name, dmg)
		},
	}
}

// NewMagicAttack builds an ability that rolls against the target's magic resist.
func NewMagicAttack(name, description, hitVerb string, baseAttack int, d *Dice) Ability {
	return Ability{
		Name:        name,
		Description: description,
		Execute: func(user *Hero, target *Enemy) string {
			dmg, missed := rollDamage(d, baseAttack, target.MR)
			if missed {
				return "Your attack misses!\n"
			}
			dmg = applyFatigue(user, dmg)
			target.TakeDamage(dmg)
			return fmt.Sprintf("%s the enemy %s, dealing %d damage!\n", hitVerb, target.Name, dmg)
		},
	}
}

// NewSplitAttack builds an ability that rolls half its damage against
// magic resist and half against armor.
func NewSplitAttack(name, description, hitVerb string, baseAttack int, d *Dice) Ability {
	return Ability{
		Name:        name,
		Description: description,
		Execute: func(user *Hero, target *Enemy) string {
			if d.D(10) == 0 {
				return "Your attack misses!\n"
			}
			half := baseAttack/2 + 3
			mod := d.D(5)
			dmg := half - target.MR + half - target.Armor
			if d.Bool() {
				dmg += mod
			} else {
				dmg -= mod
			}
			if dmg < 0 {
				dmg = 0
			}
			dmg = applyFatigue(user, dmg)
			target.TakeDamage(dmg)
			return fmt.Sprintf("%s the enemy %s, dealing %d damage!\n", hitVerb, target.Name, dmg)
		},
	}
}

// NewHealSelf builds an ability that heals the user for min..min+bonus-1 HP.
func NewHealSelf(name, description string, min, bonus int, d *Dice) Ability {
	return Ability{
		Name:        name,
		Description: description,
		Execute: func(user *Hero, target *Enemy) string {
			healing := min + d.D(bonus)
			user.Heal(healing)
			return fmt.Sprintf("You patch your wounds, healing for %d hp.\n", healing)
		},
	}
}

// NewReduceFatigue builds an ability that relieves min..min+bonus-1 fatigue.
func NewReduceFatigue(name, description string, min, bonus int, d *Dice) Ability {
	return Ability{
		Name:        name,
		Description: description,
		Execute: func(user *Hero, target *Enemy) string {
			recovery := min + d.D(bonus)
			user.ReduceFatigue(recovery)
			return fmt.Sprintf("You meditate for a short period, reducing your fatigue by %d.\n", recovery)
		},
	}
}

// maxAbilities is how many moves a hero can know at once. Past this,
// learning something new means dropping something old.
const maxAbilities = 5

// abilityPool is every move a hero could ever learn. There's no class
// tying any of these together — a hero's identity is just whichever
// of these it's picked up along the way.
var abilityPool = []func(d *Dice) Ability{
	func(d *Dice) Ability {
		return NewPhysicalAttack("Claymore Slash", "standard physical attack", "Your claymore strikes", 25, d)
	},
	func(d *Dice) Ability {
		return NewMagicAttack("Holy Smite", "standard magic attack", "You smite", 25, d)
	},
	func(d *Dice) Ability {
		return NewSplitAttack("Blessed Blade", "physical and magic split attack", "Your blessing-infused Zweihander clobbers", 25, d)
	},
	func(d *Dice) Ability {
		return NewHealSelf("Patch Wounds", "recover some HP", 5, 5, d)
	},
	func(d *Dice) Ability {
		return NewReduceFatigue("Battle Meditation", "relieve some Fatigue", 7, 5, d)
	},
	func(d *Dice) Ability {
		return NewPhysicalAttack("Quick Jab", "fast, lighter physical attack", "You jab", 15, d)
	},
	func(d *Dice) Ability {
		return NewMagicAttack("Fireball", "heavy, reckless magic attack", "You hurl a fireball at", 35, d)
	},
	func(d *Dice) Ability {
		return NewSplitAttack("Guard Break", "lighter physical and magic split attack", "You batter through the guard of", 15, d)
	},
	func(d *Dice) Ability {
		return NewHealSelf("Second Wind", "recover a lot of HP", 15, 10, d)
	},
	func(d *Dice) Ability {
		return NewReduceFatigue("Adrenaline Rush", "relieve a lot of Fatigue", 15, 10, d)
	},
}

// randomAbilityChoices draws up to n distinct abilities from the pool,
// skipping anything already in known. Fewer than n come back if the
// pool doesn't have enough unknown moves left.
func randomAbilityChoices(d *Dice, known []Ability, n int) []Ability {
	knownNames := make(map[string]bool, len(known))
	for _, a := range known {
		knownNames[a.Name] = true
	}

	order := make([]int, len(abilityPool))
	for i := range order {
		order[i] = i
	}
	for i := len(order) - 1; i > 0; i-- {
		j := d.D(i + 1)
		order[i], order[j] = order[j], order[i]
	}

	choices := make([]Ability, 0, n)
	for _, idx := range order {
		if len(choices) == n {
			break
		}
		candidate := abilityPool[idx](d)
		if knownNames[candidate.Name] {
			continue
		}
		choices = append(choices, candidate)
	}
	return choices
}
