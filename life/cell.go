package life

import (
	"github.com/virtualtam/veccell/automaton"
)

// A Cell with additional attributes
type DxCell struct {
	automaton.Cell

	onceAlive bool
	colony    *Colony
}

// ChooseColony sets a DXCell's Colony according to its live neighbours.
func (c *DxCell) ChooseColony(ancestor DxCell, neighbours []*DxCell) {
	colonies := make(map[*Colony]int)
	if ancestor.colony != nil {
		colonies[ancestor.colony] = 1
	}
	for _, cell := range neighbours {
		colonies[cell.colony]++
	}
	maxCount := 0
	for colony, count := range colonies {
		if count > maxCount {
			// Let Go's map randomization magic handle the ex-aequo case;
			// a bit of non-determinism sure doesn't hurt?
			maxCount = count
			c.colony = colony
		}
	}
}
