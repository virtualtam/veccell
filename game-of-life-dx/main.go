// Package main implements Conway's Game of Lifei.
//
// This implementation comes with additional bells and whistles:
// - each Cell belongs to a Colony;
// - the Colony a Cell belongs to depends on the surrounding live Cells.
//
// The board is rendered on the terminal using the Termbox library.
package main

import (
	"flag"
	"github.com/nsf/termbox-go"
	"math/rand"
	"time"
)

const (
	DefaultDelay        = 1000 // Milliseconds
	DefaultNColonies    = 3
	DefaultShowExplored = false
	DefaultBackground   = termbox.ColorDefault
	OnceAliveBackground = termbox.ColorBlack
)

// A Colony represents a type of population, and can feature distinctive
// attributes.
type Colony struct {
	id     int
	name   string
	symbol rune
	color  termbox.Attribute
}

// The available Colonies.
var Colonies = []Colony{
	{
		id:     0,
		name:   "Zeroth",
		symbol: 'O',
		color:  termbox.ColorDefault,
	},
	{
		id:     1,
		name:   "Crimson Red",
		symbol: '-',
		color:  termbox.ColorRed,
	},
	{
		id:     2,
		name:   "Green Washers",
		symbol: '+',
		color:  termbox.ColorGreen,
	},
	{
		id:     3,
		name:   "Blue Monks",
		symbol: 'X',
		color:  termbox.ColorBlue,
	},
	{
		id:     4,
		name:   "Yellow Cabs",
		symbol: '$',
		color:  termbox.ColorYellow,
	},
	{
		id:     5,
		name:   "Magendas",
		symbol: '*',
		color:  termbox.ColorMagenta,
	},
	{
		id:     6,
		name:   "Cyanides",
		symbol: '%',
		color:  termbox.ColorCyan,
	},
	{
		id:     7,
		name:   "White Walking Bassoons",
		symbol: '8',
		color:  termbox.ColorWhite,
	},
}

// A Cell can be either alive or dead, and belongs to a Colony.
type Cell struct {
	alive     bool
	onceAlive bool
	colony    *Colony
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

// A Board holds the game's parameters and state.
type Board struct {
	nRows     int
	nCols     int
	nColonies int
	colonies  []Colony
	cells     [][]Cell
}

// NewBoard creates and initializes a Board.
func NewBoard(nRows, nCols, nColonies int) Board {
	b := Board{
		nRows:     nRows,
		nCols:     nCols,
		nColonies: nColonies,
	}
	b.colonies = Colonies[:b.nColonies]
	b.cells = make([][]Cell, b.nRows)
	for i := 0; i < b.nRows; i++ {
		b.cells[i] = make([]Cell, b.nCols)
	}

	return b
}

// Randomize sets a board's Cells in a random state (alive|dead) and assigns it
// to a randomly chosen Colony.
func (b *Board) Randomize() {
	for i := 0; i < b.nRows; i++ {
		for j := 0; j < b.nCols; j++ {
			b.cells[i][j].alive = rand.Intn(5) == 1
			b.cells[i][j].onceAlive = b.cells[i][j].alive
			b.cells[i][j].colony = &b.colonies[rand.Intn(b.nColonies)]
		}
	}
}

// LiveNeighboursAt returns the live Cells surrounding the Cell at the given
// position.
func (b *Board) LiveNeighboursAt(row, col int) []*Cell {
	neighbours := []*Cell{}

	for i := -1; i < 2; i++ {
		for j := -1; j < 2; j++ {
			tmpRow := row + i
			tmpCol := col + j

			if tmpRow < 0 || tmpRow >= b.nRows {
				// assume horizontal borders are made of dead cells
				continue
			}
			if tmpCol < 0 || tmpCol >= b.nCols {
				// assume vertical borders are made of dead cells
				continue
			}

			if i == 0 && j == 0 {
				continue
			}

			if b.cells[tmpRow][tmpCol].alive {
				neighbours = append(neighbours, &b.cells[tmpRow][tmpCol])
			}
		}
	}

	return neighbours
}

// Next computes the next iteration of the game, and updates the Board's state.
//
// On to the next iteration!
//
// 1. Any live cell with fewer than two live neighbours dies, as if
// by under population.
//
// 2. Any live cell with two or three live neighbours lives on to
// the the next generation.
//
// 3. Any live cell with more than three live neighbours dies, as if
// by overpopulation.
//
// 4. Any dead cell with exactly three live neighbours becomes a
// live cell, as if by reproduction.
//
// DX bells and whistles:
//
// 5. Any live cell belongs to the colony that has the most live neighbours, as
// if by conquest or conversion (Wololo!).
func (b *Board) Next() {
	nextBoard := NewBoard(b.nRows, b.nCols, b.nColonies)

	for i := 0; i < b.nRows; i++ {
		for j := 0; j < b.nCols; j++ {
			neighbours := b.LiveNeighboursAt(i, j)
			liveNeighbours := len(neighbours)
			nextBoard.cells[i][j].onceAlive = b.cells[i][j].onceAlive

			if b.cells[i][j].alive {
				if liveNeighbours == 2 || liveNeighbours == 3 {
					nextBoard.cells[i][j].alive = true
					nextBoard.cells[i][j].ChooseColony(b.cells[i][j], neighbours)
				}
			} else {
				if liveNeighbours == 3 {
					nextBoard.cells[i][j].alive = true
					nextBoard.cells[i][j].onceAlive = true
					nextBoard.cells[i][j].ChooseColony(b.cells[i][j], neighbours)
				}
			}
		}
	}

	b.cells = nextBoard.cells
}

// Draw renders the Board on the terminal.
func (b *Board) Draw(showExplored bool) {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	for i := 0; i < b.nRows; i++ {
		for j := 0; j < b.nCols; j++ {
			if showExplored {
				if b.cells[i][j].alive {
					termbox.SetCell(
						j,
						i,
						b.cells[i][j].colony.symbol,
						b.cells[i][j].colony.color,
						OnceAliveBackground,
					)
				} else if b.cells[i][j].onceAlive {
					termbox.SetCell(j, i, ' ', termbox.ColorDefault, OnceAliveBackground)
				}
			} else {
				if b.cells[i][j].alive {
					termbox.SetCell(
						j,
						i,
						b.cells[i][j].colony.symbol,
						b.cells[i][j].colony.color,
						DefaultBackground,
					)
				}
			}
		}
	}

	termbox.Flush()
}

func main() {
	// Command-line parameters
	var (
		delay        int
		nColonies    int
		showExplored bool
	)

	flag.IntVar(&delay, "delay", DefaultDelay, "Delay between two iterations (milliseconds)")
	flag.IntVar(&nColonies, "colonies", DefaultNColonies, "Number of colonies")
	flag.BoolVar(&showExplored, "show-explored", DefaultShowExplored, "Show explored regions")
	flag.Parse()

	// Termbox setup
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	termWidth, termHeight := termbox.Size()

	eventQueue := make(chan termbox.Event)
	go func() {
		for {
			eventQueue <- termbox.PollEvent()
		}
	}()

	// Game board setup
	rand.Seed(time.Now().UTC().UnixNano())
	board := NewBoard(termHeight, termWidth, nColonies)
	board.Randomize()
	board.Draw(showExplored)

	drawQueue := make(chan bool)
	go func(delay *int) {
		for {
			time.Sleep(time.Duration(*delay) * time.Millisecond)
			drawQueue <- true
		}
	}(&delay)

mainloop:
	for {
		select {
		case <-drawQueue:
			board.Next()
			board.Draw(showExplored)

		case ev := <-eventQueue:
			switch ev.Type {
			case termbox.EventKey:
				switch ev.Key {
				case termbox.KeyEsc:
					break mainloop
				case termbox.KeyCtrlC:
					break mainloop
				case termbox.KeyArrowUp:
					switch {
					case delay < 10:
						delay++
					case delay < 100:
						delay += 10
					default:
						delay += 100
					}
				case termbox.KeyArrowDown:
					switch {
					case delay > 100:
						delay -= 100
					case delay > 10:
						delay -= 10
					case delay > 2:
						delay--
					}
				}

				switch ev.Ch {
				case 'q':
					break mainloop
				}
			}
		}
	}
}
