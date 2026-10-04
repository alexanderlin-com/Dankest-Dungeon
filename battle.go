package main

// rollEnemyAttack is the enemy AI's damage roll, shared by the combat
// screen's enemy-turn tick.
func rollEnemyAttack(enemy *Enemy, d *Dice) (dmg int, missed bool) {
	if d.D(10) == 0 {
		return 0, true
	}
	mod := d.D(5)
	dmg = enemy.Attack - mod
	if d.Bool() {
		dmg = enemy.Attack + mod
	}
	if dmg < 0 {
		dmg = 0
	}
	return dmg, false
}
