// Elementary Cellular Automaton implementation.
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
//
// The board is rendered on the terminal using the Termbox library.
package elementary

import (
	"container/ring"
	"fmt"
	"math/rand"
	"strconv"

	"github.com/nsf/termbox-go"

	"github.com/virtualtam/veccell/automaton"
)

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
	cells []automaton.Cell
}

func NewElementaryAutomaton(ruleNumber, size int) ElementaryAutomaton {
	a := ElementaryAutomaton{
		rule: NewRule(ruleNumber),
		size: size,
	}
	a.cells = make([]automaton.Cell, a.size)
	return a
}

func (a *ElementaryAutomaton) Randomize() {
	for i := 0; i < len(a.cells); i++ {
		a.cells[i].Alive = rand.Intn(2) == 1
	}
}

func (a *ElementaryAutomaton) StartWithCenter() {
	a.cells[len(a.cells)/2].Alive = true
}

func (a *ElementaryAutomaton) Next() {
	nextState := make([]automaton.Cell, a.size)

	pattern := 0

	for i := 0; i < len(a.cells); i++ {
		pattern = 0

		// left
		if i > 0 {
			if a.cells[i-1].Alive {
				pattern += 4
			}
		}

		// center
		if a.cells[i].Alive {
			pattern += 2
		}

		// right
		if i < len(a.cells)-1 {
			if a.cells[i+1].Alive {
				pattern += 1
			}
		}

		nextState[i].Alive = a.rule.transitions[pattern]
	}

	a.cells = nextState
}

func (a *ElementaryAutomaton) Draw() {
	for i := 0; i < len(a.cells); i++ {
		if a.cells[i].Alive {
			fmt.Printf("+")
		} else {
			fmt.Printf(" ")
		}
	}
	fmt.Println()
}

type ElementaryAutomatonRing struct {
	size      int
	automaton *ElementaryAutomaton
	history   *ring.Ring
}

func NewElementaryAutomatonRing(size int, a *ElementaryAutomaton) ElementaryAutomatonRing {
	h := ElementaryAutomatonRing{size: size, automaton: a}
	h.history = ring.New(h.size)
	for i := 0; i < h.history.Len(); i++ {
		h.history.Value = make([]automaton.Cell, len(h.automaton.cells))
		copy(h.history.Value.([]automaton.Cell), h.automaton.cells)
		h.history = h.history.Next()
		h.automaton.Next()
	}
	return h
}

func (h *ElementaryAutomatonRing) Next() {
	h.automaton.Next()
	copy(h.history.Value.([]automaton.Cell), h.automaton.cells)
	h.history = h.history.Move(1)
}

func (h *ElementaryAutomatonRing) Randomize() {
	h.automaton.Randomize()
}

func (h *ElementaryAutomatonRing) Draw() {
	row := 0
	h.history.Do(func(p interface{}) {
		cells := p.([]automaton.Cell)
		for col := 0; col < len(cells); col++ {
			if cells[col].Alive {
				termbox.SetCell(col, row, '+', termbox.ColorDefault, termbox.ColorDefault)
			}
		}
		row++
	})
}
