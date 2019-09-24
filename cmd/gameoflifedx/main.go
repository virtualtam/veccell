// Package main implements Conway's Game of Life.
//
// This implementation comes with additional bells and whistles:
// - each Cell belongs to a Colony;
// - the Colony a Cell belongs to depends on the surrounding live Cells.
//
// The board is rendered on the terminal using the Termbox library.
package main

import (
	"flag"
	"math/rand"
	"time"

	"github.com/nsf/termbox-go"

	"github.com/virtualtam/veccell/automaton"
	"github.com/virtualtam/veccell/life"
)

const (
	DefaultDelay        = 1000 // Milliseconds
	DefaultNColonies    = 3
	DefaultShowExplored = false
)

func main() {
	delay := flag.Int("delay", DefaultDelay, "Delay between two iterations (milliseconds)")
	nColonies := flag.Int("colonies", DefaultNColonies, "Number of colonies")
	showExplored := flag.Bool("show-explored", DefaultShowExplored, "Show explored regions")
	flag.Parse()

	rand.Seed(time.Now().UTC().UnixNano())

	// Termbox setup
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	termWidth, termHeight := termbox.Size()

	// Game board setup
	board := life.NewGameOfLifeDx(termHeight, termWidth, *nColonies, *showExplored)
	board.Randomize()
	board.Draw()

	controller := automaton.NewController(&board, delay)
	controller.Loop()
}
