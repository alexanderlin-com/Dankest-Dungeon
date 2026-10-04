package main

import "testing"

func TestRandomAbilityChoicesExcludesKnown(t *testing.T) {
	d := NewDice()
	known := randomAbilityChoices(d, nil, 3)
	if len(known) != 3 {
		t.Fatalf("expected 3 abilities with an empty known list, got %d", len(known))
	}

	choices := randomAbilityChoices(d, known, len(abilityPool)-len(known))
	knownNames := map[string]bool{}
	for _, a := range known {
		knownNames[a.Name] = true
	}
	for _, c := range choices {
		if knownNames[c.Name] {
			t.Fatalf("randomAbilityChoices returned %q, which was already known", c.Name)
		}
	}
}

func TestRandomAbilityChoicesAreDistinct(t *testing.T) {
	d := NewDice()
	choices := randomAbilityChoices(d, nil, len(abilityPool))

	seen := map[string]bool{}
	for _, c := range choices {
		if seen[c.Name] {
			t.Fatalf("duplicate ability %q in a single draw", c.Name)
		}
		seen[c.Name] = true
	}
	if len(choices) != len(abilityPool) {
		t.Fatalf("expected all %d pool abilities back, got %d", len(abilityPool), len(choices))
	}
}

func TestRandomAbilityChoicesStopsWhenPoolExhausted(t *testing.T) {
	d := NewDice()
	all := randomAbilityChoices(d, nil, len(abilityPool))

	// Every pool ability is now "known" - nothing left to offer.
	none := randomAbilityChoices(d, all, 2)
	if len(none) != 0 {
		t.Fatalf("expected 0 choices once the whole pool is known, got %d", len(none))
	}
}

func TestNewHeroStartsWithPunchesAndGold(t *testing.T) {
	d := NewDice()
	hero := NewHero(d)

	if len(hero.Abilities) != maxAbilities {
		t.Fatalf("expected %d starting abilities, got %d", maxAbilities, len(hero.Abilities))
	}
	for _, a := range hero.Abilities {
		if a.Name != "Punch" {
			t.Fatalf("expected every starting slot to be Punch, got %q", a.Name)
		}
	}
	if hero.Gold != startingGold {
		t.Fatalf("expected starting gold %d, got %d", startingGold, hero.Gold)
	}
}
