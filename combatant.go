package main

// Combatant holds the stats shared by heroes and enemies. Having one
// struct as the single source of truth avoids the classic Java-subclass
// trap where a subtype redeclares a field and silently shadows the base.
type Combatant struct {
	Name   string
	HP     int
	MaxHP  int
	Attack int
	Armor  int
	MR     int // magic resistance
}

func (c *Combatant) IsAlive() bool {
	return c.HP > 0
}

func (c *Combatant) TakeDamage(amount int) {
	c.HP -= amount
	if c.HP < 0 {
		c.HP = 0
	}
}

func (c *Combatant) Heal(amount int) {
	c.HP += amount
	if c.HP > c.MaxHP {
		c.HP = c.MaxHP
	}
}
