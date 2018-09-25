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

var (
	delay     int
	randomize bool
	rule      int
)

func init() {
	flag.IntVar(&delay, "delay", DefaultDelay, "Delay between two iterations (milliseconds)")
	flag.BoolVar(&randomize, "randomize", DefaultRandomize, "Randomize initial state")
	flag.IntVar(&rule, "rule", DefaultRule, "Automaton rule")
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

	// Elementary automaton setup
	elementary := automaton.NewElementaryAutomaton(rule, termWidth)
	if randomize {
		elementary.Randomize()
	} else {
		elementary.StartWithCenter()
	}
	history := automaton.NewElementaryAutomatonHistory(termHeight, &elementary)
	history.Draw()

	controller := automaton.NewController(&history, &delay)
	controller.Loop()
}
