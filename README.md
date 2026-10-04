# Dankest Dungeon


## Table of Contents
- [Introduction](#introduction)
- [Game Overview](#game-overview)
- [Installation](#installation)

## Introduction

Welcome to Dankest Dungeon! This is a text-based adventure game. Although its namesake was inspired by Darkest Dungeon, the gameplay is nothing like it at all. Immerse yourself in a somewhat captivating narrative and make strategic choices as you navigate through the lack of storyline.

Originally written in Java during freshman year of my CS degree, then rewritten in Go — Go's composition-over-inheritance approach fits a small text-adventure engine much better than Java's class hierarchies did. The interface is now a proper terminal UI built with [Bubble Tea](https://github.com/charmbracelet/bubbletea) and [Lip Gloss](https://github.com/charmbracelet/lipgloss): bordered panels, color-coded HP bars, and arrow-key menus instead of typing a number and hitting enter.

## Game Overview

Dankest Dungeon is a text-based adventure game that focuses on storytelling and decision-making. Here's what you can experience in the current version:

- **Terminal UI:** Navigate every menu with the arrow keys (or j/k) and Enter; `q` quits from anywhere.

- **Combat:** Fight your way through 10 kinds of enemies using the Knight's 5 abilities — physical, magic, and split attacks, plus self-heal and fatigue relief.

- **Taverns:** Stop between fights to browse a shop catalog and spend gold healing up or shedding fatigue.

## Installation

To experience Dankest Dungeon, follow these simple installation steps:

1. Clone this repository to your local machine.
2. Make sure you have Go installed (1.22+).
3. Run `go run .` from the project root.
4. Have fun (mandatory).
