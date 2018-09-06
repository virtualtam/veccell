// Package main implements an Elementary Cellular Automaton
//
// The automaton is made of an array of Cells, each Cell having two possible
// states.
//
// Every cell changes states based on its current state, and the state of its
// two neighbors.
//
// 111 110 101 100 011 010 001 000
// -------------------------------
//  0   0   0   1   1   1   1   0   Rule  30
//  0   1   0   1   1   0   1   0   Rule  90
//  0   1   1   0   1   1   1   0   Rule 110
package main

import (
	"flag"
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

const (
	DefaultDelay     = 1000 // Milliseconds
	DefaultRandomize = false
	DefaultRule      = 90
)

type Cell struct {
	alive bool
}

type Rule struct {
	number      int
	transitions [8]bool
}

func NewRule(number int) Rule {
	rule := Rule{number: number}

	binaryRule := strconv.FormatInt(int64(rule.number), 2)
	for index, bit := range binaryRule {
		rule.transitions[len(binaryRule)-1-index], _ = strconv.ParseBool(string(bit))
	}

	return rule
}

type ElementaryAutomaton struct {
	rule  Rule
	size  int
	cells []Cell
}

func NewElementaryAutomaton(ruleNumber, size int) ElementaryAutomaton {
	automaton := ElementaryAutomaton{
		rule: NewRule(ruleNumber),
		size: size,
	}
	automaton.cells = make([]Cell, automaton.size)
	return automaton
}

func (a *ElementaryAutomaton) Randomize() {
	for i := 0; i < len(a.cells); i++ {
		a.cells[i].alive = rand.Intn(2) == 1
	}
}

func (a *ElementaryAutomaton) StartWithCenter() {
	a.cells[len(a.cells)/2].alive = true
}

func (a *ElementaryAutomaton) Next() {
	nextState := make([]Cell, a.size)

	pattern := 0

	for i := 0; i < len(a.cells); i++ {
		pattern = 0

		// left
		if i > 0 {
			if a.cells[i-1].alive {
				pattern += 4
			}
		}

		// center
		if a.cells[i].alive {
			pattern += 2
		}

		// right
		if i < len(a.cells)-1 {
			if a.cells[i+1].alive {
				pattern += 1
			}
		}

		nextState[i].alive = a.rule.transitions[pattern]
	}

	a.cells = nextState
}

func (a *ElementaryAutomaton) Draw() {
	for i := 0; i < len(a.cells); i++ {
		if a.cells[i].alive {
			fmt.Printf("+")
		} else {
			fmt.Printf(" ")
		}
	}
	fmt.Println()
}

func main() {
	var (
		delay     int
		randomize bool
		rule      int
	)

	flag.IntVar(&delay, "delay", DefaultDelay, "Delay between two iterations (milliseconds)")
	flag.BoolVar(&randomize, "randomize", DefaultRandomize, "Randomize initial state")
	flag.IntVar(&rule, "rule", DefaultRule, "Automaton rule")
	flag.Parse()

	rand.Seed(time.Now().UTC().UnixNano())

	automaton := NewElementaryAutomaton(rule, 90)
	if randomize {
		automaton.Randomize()
	} else {
		automaton.StartWithCenter()
	}
	automaton.Draw()

	for {
		time.Sleep(time.Duration(delay) * time.Millisecond)
		automaton.Next()
		automaton.Draw()
	}
}
