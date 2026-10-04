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

// finalizeDamage is the single choke point every damage-dealing
// ability funnels its raw roll through: it applies the fatigue
// penalty and adds (then clears) a pending NextAttackBonus from Brace.
func finalizeDamage(user *Hero, dmg int) int {
	dmg -= user.FatiguePenalty()
	dmg += user.NextAttackBonus
	user.NextAttackBonus = 0
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
			dmg = finalizeDamage(user, dmg)
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
			dmg = finalizeDamage(user, dmg)
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
			dmg = finalizeDamage(user, dmg)
			target.TakeDamage(dmg)
			return fmt.Sprintf("%s the enemy %s, dealing %d damage!\n", hitVerb, target.Name, dmg)
		},
	}
}

// NewPureAttack builds an ability that ignores the target's defenses
// entirely - no armor, no magic resist.
func NewPureAttack(name, description, hitVerb string, baseAttack int, d *Dice) Ability {
	return Ability{
		Name:        name,
		Description: description,
		Execute: func(user *Hero, target *Enemy) string {
			dmg, missed := rollDamage(d, baseAttack, 0)
			if missed {
				return "Your attack misses!\n"
			}
			dmg = finalizeDamage(user, dmg)
			target.TakeDamage(dmg)
			return fmt.Sprintf("%s the enemy %s, dealing %d damage!\n", hitVerb, target.Name, dmg)
		},
	}
}

// NewRecklessAttack is a NewPureAttack that also costs the caster
// self-inflicted fatigue - power with a price.
func NewRecklessAttack(name, description, hitVerb string, baseAttack, selfFatigue int, d *Dice) Ability {
	return Ability{
		Name:        name,
		Description: description,
		Execute: func(user *Hero, target *Enemy) string {
			dmg, missed := rollDamage(d, baseAttack, 0)
			user.AddFatigue(selfFatigue)
			if missed {
				return fmt.Sprintf("%s misses, but the backlash still costs you %d fatigue!\n", hitVerb, selfFatigue)
			}
			dmg = finalizeDamage(user, dmg)
			target.TakeDamage(dmg)
			return fmt.Sprintf("%s the enemy %s, dealing %d damage (and %d fatigue to you)!\n", hitVerb, target.Name, dmg, selfFatigue)
		},
	}
}

// executeThreshold and executeMultiplier govern NewExecuteAttack: at or
// below this fraction of max HP, damage is multiplied.
const (
	executeThreshold  = 0.25
	executeMultiplier = 2
)

// NewExecuteAttack builds a physical attack that deals bonus damage to
// a badly wounded target - reads the target's own HP, already visible
// on the HUD, rather than needing any new state.
func NewExecuteAttack(name, description, hitVerb string, baseAttack int, d *Dice) Ability {
	return Ability{
		Name:        name,
		Description: description,
		Execute: func(user *Hero, target *Enemy) string {
			dmg, missed := rollDamage(d, baseAttack, target.Armor)
			if missed {
				return "Your attack misses!\n"
			}
			dmg = finalizeDamage(user, dmg)
			finishing := target.MaxHP > 0 && target.HP <= int(float64(target.MaxHP)*executeThreshold)
			if finishing {
				dmg *= executeMultiplier
			}
			target.TakeDamage(dmg)
			if finishing {
				return fmt.Sprintf("%s the weakened enemy %s, dealing a brutal %d damage!\n", hitVerb, target.Name, dmg)
			}
			return fmt.Sprintf("%s the enemy %s, dealing %d damage!\n", hitVerb, target.Name, dmg)
		},
	}
}

// NewDoubleAttack builds an ability that rolls two weaker physical
// hits in one turn instead of one strong one.
func NewDoubleAttack(name, description, hitVerb string, baseAttack int, d *Dice) Ability {
	return Ability{
		Name:        name,
		Description: description,
		Execute: func(user *Hero, target *Enemy) string {
			total := 0
			hits := 0
			for i := 0; i < 2; i++ {
				dmg, missed := rollDamage(d, baseAttack, target.Armor)
				if missed {
					continue
				}
				total += dmg
				hits++
			}
			if hits == 0 {
				return "Both of your strikes miss!\n"
			}
			total = finalizeDamage(user, total)
			target.TakeDamage(total)
			return fmt.Sprintf("%s the enemy %s %d time(s), dealing %d damage!\n", hitVerb, target.Name, hits, total)
		},
	}
}

// NewMarkingMove builds an ability that deals no real damage but marks
// the target for NewMarkConsumingAttack's bonus.
func NewMarkingMove(name, description, verb string) Ability {
	return Ability{
		Name:        name,
		Description: description,
		Execute: func(user *Hero, target *Enemy) string {
			target.Marked = true
			return fmt.Sprintf("%s %s - it's marked for a devastating follow-up!\n", verb, target.Name)
		},
	}
}

// NewMarkConsumingAttack builds an ability that's weak on its own but
// hits hard and clears the mark if the target is marked.
func NewMarkConsumingAttack(name, description, hitVerb string, baseAttack, markedBonus int, d *Dice) Ability {
	return Ability{
		Name:        name,
		Description: description,
		Execute: func(user *Hero, target *Enemy) string {
			dmg, missed := rollDamage(d, baseAttack, target.Armor)
			if missed {
				return "Your attack misses!\n"
			}
			if target.Marked {
				dmg += markedBonus
				target.Marked = false
			}
			dmg = finalizeDamage(user, dmg)
			target.TakeDamage(dmg)
			return fmt.Sprintf("%s the enemy %s, dealing %d damage!\n", hitVerb, target.Name, dmg)
		},
	}
}

// NewLifestealAttack builds a magic attack that heals the caster for a
// percentage of the damage it deals.
func NewLifestealAttack(name, description, hitVerb string, baseAttack, healPct int, d *Dice) Ability {
	return Ability{
		Name:        name,
		Description: description,
		Execute: func(user *Hero, target *Enemy) string {
			dmg, missed := rollDamage(d, baseAttack, target.MR)
			if missed {
				return "Your attack misses!\n"
			}
			dmg = finalizeDamage(user, dmg)
			target.TakeDamage(dmg)
			healing := dmg * healPct / 100
			user.Heal(healing)
			return fmt.Sprintf("%s the enemy %s, dealing %d damage and healing you for %d!\n", hitVerb, target.Name, dmg, healing)
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

// NewNextAttackBuff builds an ability that skips the turn and buffs
// the user's next outgoing hit.
func NewNextAttackBuff(name, description string, bonus int) Ability {
	return Ability{
		Name:        name,
		Description: description,
		Execute: func(user *Hero, target *Enemy) string {
			user.NextAttackBonus = bonus
			return "You ready yourself - your next strike will land harder.\n"
		},
	}
}

// NewGuardUp builds an ability that skips the turn and reduces the
// next hit the user takes.
func NewGuardUp(name, description string, reduction int) Ability {
	return Ability{
		Name:        name,
		Description: description,
		Execute: func(user *Hero, target *Enemy) string {
			user.NextDamageReduction = reduction
			return "You brace yourself for the next blow.\n"
		},
	}
}

// maxAbilities is how many ability slots a hero has. There's no
// growing past this - a hero always has exactly this many moves, each
// one either Punch or something learned to replace it.
const maxAbilities = 4

// newPunch is the weak, universal starting move every slot begins as.
// Deliberately worse than anything in abilityPool (the weakest real
// move, Shield Bash, is base 18) so learning something is always an
// upgrade. Not part of abilityPool - you'd never pay to learn Punch.
func newPunch(d *Dice) Ability {
	return NewPhysicalAttack("Punch", "a weak, untrained strike", "You weakly punch", 8, d)
}

// abilityPool is every move a hero could ever learn: 8 archetypes, 2
// moves each. There's no class tying the archetypes to the hero - a
// hero's identity is just whichever of these it's picked up along the
// way.
var abilityPool = []func(d *Dice) Ability{
	// Fighter
	func(d *Dice) Ability {
		return NewPhysicalAttack("Cleave", "standard physical attack", "Your blade cleaves into", 25, d)
	},
	func(d *Dice) Ability {
		return NewNextAttackBuff("Brace", "skip this turn, your next hit lands harder", 15)
	},
	// Tank
	func(d *Dice) Ability {
		return NewPhysicalAttack("Shield Bash", "a reliable, if unspectacular, physical attack", "Your shield bash slams into", 18, d)
	},
	func(d *Dice) Ability {
		return NewGuardUp("Guard Up", "skip this turn, reduce the next hit you take", 15)
	},
	// Rogue
	func(d *Dice) Ability {
		return NewExecuteAttack("Finishing Blow", "physical attack, brutal against a weakened enemy", "Your blade finds a weak point in", 20, d)
	},
	func(d *Dice) Ability {
		return NewDoubleAttack("Twin Daggers", "two quick, weaker physical strikes", "Your daggers flash twice against", 14, d)
	},
	// Wizard
	func(d *Dice) Ability {
		return NewMagicAttack("Arcane Bolt", "standard magic attack", "Your arcane bolt strikes", 25, d)
	},
	func(d *Dice) Ability {
		return NewPureAttack("Mind Spike", "magic damage that ignores resistances", "A spike of pure thought pierces", 18, d)
	},
	// Sorcerer
	func(d *Dice) Ability {
		return NewRecklessAttack("Chaos Bolt", "reckless magic damage, ignores resistances, costs fatigue", "Your chaos bolt detonates against", 30, 8, d)
	},
	func(d *Dice) Ability {
		return NewRecklessAttack("Overload", "an all-out magic blast, ignores resistances, costs heavy fatigue", "Your overloaded spell obliterates", 45, 15, d)
	},
	// Witch
	func(d *Dice) Ability {
		return NewMarkingMove("Hex", "marks the enemy for a devastating follow-up", "You lay a hex upon")
	},
	func(d *Dice) Ability {
		return NewMarkConsumingAttack("Curse Strike", "weak on its own, brutal against a marked enemy", "Your cursed blade strikes", 15, 25, d)
	},
	// Cleric
	func(d *Dice) Ability {
		return NewHealSelf("Mend", "recover some HP", 10, 8, d)
	},
	func(d *Dice) Ability {
		return NewLifestealAttack("Smite the Wicked", "magic damage that heals you for part of the damage dealt", "You smite", 22, 40, d)
	},
	// Ranger
	func(d *Dice) Ability {
		return NewPhysicalAttack("Steady Shot", "standard physical attack", "Your arrow strikes", 25, d)
	},
	func(d *Dice) Ability {
		return NewReduceFatigue("Second Breath", "relieve some Fatigue", 10, 8, d)
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
