package main

import (
	"flag"
	"math/rand"
	"time"

	"github.com/nsf/termbox-go"

	"github.com/virtualtam/veccell/automaton"
	"github.com/virtualtam/veccell/elementary"
)

const (
	DefaultDelay     = 1000 // Milliseconds
	DefaultRandomize = false
	DefaultRule      = 90
)

func main() {
	delay := flag.Int("delay", DefaultDelay, "Delay between two iterations (milliseconds)")
	randomize := flag.Bool("randomize", DefaultRandomize, "Randomize initial state")
	rule := flag.Int("rule", DefaultRule, "Automaton rule")
	flag.Parse()

	rand.Seed(time.Now().UTC().UnixNano())

	// Termbox setup
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	termWidth, termHeight := termbox.Size()

	// Elementary automaton setup
	a := elementary.NewElementaryAutomaton(*rule, termWidth)
	if *randomize {
		a.Randomize()
	} else {
		a.StartWithCenter()
	}
	elementaryRing := elementary.NewElementaryAutomatonRing(termHeight, &a)
	elementaryRing.Draw()

	controller := automaton.NewController(&elementaryRing, delay)
	controller.Loop()
}
