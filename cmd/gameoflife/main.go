// Package gameoflife implements Conway's Game of Life.
//
// The board is rendered on the terminal using the Termbox library.
package main

import (
	"flag"
	"math/rand"
	"time"

	"github.com/nsf/termbox-go"

	"github.com/virtualtam/veccell/life"
	"github.com/virtualtam/veccell/render"
)

const (
	DefaultDelay            = 1000 // Milliseconds
	DefaultBorderCellsAlive = false
)

func main() {
	borderCellsAlive := flag.Bool("borderCellsAlive", DefaultBorderCellsAlive, "Whether the border cells are considered alive or dead")
	delay := flag.Int("delay", DefaultDelay, "Delay between two iterations (milliseconds)")
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
	board := life.NewGameOfLife(termHeight, termWidth, *borderCellsAlive)
	board.Randomize()
	//board.RandomizeArea(0, termHeight/2, 0, termWidth/2)
	//board.RandomizeArea(termHeight/4, 3*termHeight/4, termWidth/4, 3*termWidth/4)
	//board.CreateGliderAt(7, 7)
	board.Draw()

	controller := render.NewController(&board, delay, &render.TermboxRenderer{})
	controller.Loop()
}
