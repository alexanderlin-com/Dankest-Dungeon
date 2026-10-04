package main

import "math/rand"

// Dice wraps a single seeded RNG so the whole game shares one source of
// randomness instead of constructing a fresh math/rand.Rand per call.
type Dice struct {
	rng *rand.Rand
}

func NewDice() *Dice {
	return &Dice{rng: rand.New(rand.NewSource(rand.Int63()))}
}

// D returns a value in [0, n).
func (d *Dice) D(n int) int {
	return d.rng.Intn(n)
}

func (d *Dice) Bool() bool {
	return d.rng.Intn(2) == 1
}

// InitiativeWin rolls a d20 for each side; ties favor the player.
func (d *Dice) InitiativeWin() bool {
	playerRoll := d.rng.Intn(20)
	enemyRoll := d.rng.Intn(20)
	return playerRoll >= enemyRoll
}
