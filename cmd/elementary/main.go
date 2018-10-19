package main

import (
	"flag"
	"math/rand"
	"time"

	"github.com/nsf/termbox-go"

	"github.com/virtualtam/veccell/automaton"
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
	elementary := automaton.NewElementaryAutomaton(*rule, termWidth)
	if *randomize {
		elementary.Randomize()
	} else {
		elementary.StartWithCenter()
	}
	elementaryRing := automaton.NewElementaryAutomatonRing(termHeight, &elementary)
	elementaryRing.Draw()

	controller := automaton.NewController(&elementaryRing, delay)
	controller.Loop()
}
