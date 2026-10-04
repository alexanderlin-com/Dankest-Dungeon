package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	d := NewDice()
	hero := NewHero(d)

	if _, err := tea.NewProgram(newAppModel(hero, d)).Run(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
