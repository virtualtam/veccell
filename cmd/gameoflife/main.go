// Package gameoflife implements Conway's Game of Life.
//
// The board is rendered on the terminal using the Termbox library.
package main

import (
	"flag"
	"math/rand"
	"time"

	"github.com/nsf/termbox-go"

	"github.com/virtualtam/veccell/automaton"
)

const (
	DefaultDelay            = 1000 // Milliseconds
	DefaultBorderCellsAlive = false
)

var (
	borderCellsAlive bool
	delay            int
)

func init() {
	flag.BoolVar(&borderCellsAlive, "borderCellsAlive", DefaultBorderCellsAlive, "Whether the border cells are considered alive or dead")
	flag.IntVar(&delay, "delay", DefaultDelay, "Delay between two iterations (milliseconds)")
	flag.Parse()

}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	// Termbox setup
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	termWidth, termHeight := termbox.Size()

	// Game board setup
	board := automaton.NewGameOfLife(termHeight, termWidth, borderCellsAlive)
	board.Randomize()
	//board.RandomizeArea(0, termHeight/2, 0, termWidth/2)
	//board.RandomizeArea(termHeight/4, 3*termHeight/4, termWidth/4, 3*termWidth/4)
	//board.CreateGliderAt(7, 7)
	board.Draw()

	controller := automaton.NewController(&board, &delay)
	controller.Loop()
}
