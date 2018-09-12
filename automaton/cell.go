package automaton

// A Cell can be alive or dead.
type Cell struct {
	alive     bool
	onceAlive bool    // FIXME Life DX
	colony    *Colony // FIXME Life DX
}

// ChooseColony sets a Cell's Colony according to its live neighbours.
func (c *Cell) ChooseColony(ancestor Cell, neighbours []*Cell) {
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
