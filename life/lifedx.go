package life

import (
	"math/rand"

	"github.com/nsf/termbox-go"

	"github.com/virtualtam/veccell/render"
)

const (
	OnceAliveBackground = termbox.ColorBlack
)

type GameOfLifeDx struct {
	rows             int
	cols             int
	borderCellsAlive bool
	cells            [][]DxCell
	nColonies        int
	colonies         []Colony
	showExplored     bool
}

// NewGameOfLifeDx creates and initializes a Game Of Life automaton.
func NewGameOfLifeDx(rows, cols, nColonies int, showExplored bool) GameOfLifeDx {
	g := GameOfLifeDx{
		rows:         rows,
		cols:         cols,
		nColonies:    nColonies,
		showExplored: showExplored,
	}
	g.colonies = Colonies[:g.nColonies]
	g.cells = make([][]DxCell, g.rows)
	for i := 0; i < g.rows; i++ {
		g.cells[i] = make([]DxCell, g.cols)
	}

	return g
}

// Randomize sets a board's Cells in a random state .Alive|dead) and assigns
// it to a randomly chosen Colony.
func (g *GameOfLifeDx) Randomize() {
	for i := 0; i < g.rows; i++ {
		for j := 0; j < g.cols; j++ {
			g.cells[i][j].Alive = rand.Intn(5) == 1
			g.cells[i][j].onceAlive = g.cells[i][j].Alive
			g.cells[i][j].colony = &g.colonies[rand.Intn(g.nColonies)]
		}
	}
}

// LiveNeighboursAt returns the live Cells surrounding the Cell at the given
// position.
func (g *GameOfLifeDx) LiveNeighboursAt(row, col int) []*DxCell {
	neighbours := []*DxCell{}

	for i := -1; i < 2; i++ {
		for j := -1; j < 2; j++ {
			tmpRow := row + i
			tmpCol := col + j

			if tmpRow < 0 || tmpRow >= g.rows {
				// skip horizontal borders
				continue
			}
			if tmpCol < 0 || tmpCol >= g.cols {
				// skip vertical borders
				continue
			}

			if i == 0 && j == 0 {
				// skip current Cell
				continue
			}

			if g.cells[tmpRow][tmpCol].Alive {
				neighbours = append(neighbours, &g.cells[tmpRow][tmpCol])
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
func (g *GameOfLifeDx) Next() {
	nextBoard := NewGameOfLifeDx(g.rows, g.cols, g.nColonies, g.showExplored)

	for i := 0; i < g.rows; i++ {
		for j := 0; j < g.cols; j++ {
			neighbours := g.LiveNeighboursAt(i, j)
			liveNeighbours := len(neighbours)
			nextBoard.cells[i][j].onceAlive = g.cells[i][j].onceAlive

			if g.cells[i][j].Alive {
				if liveNeighbours == 2 || liveNeighbours == 3 {
					nextBoard.cells[i][j].Alive = true
					nextBoard.cells[i][j].ChooseColony(g.cells[i][j], neighbours)
				}
			} else {
				if liveNeighbours == 3 {
					nextBoard.cells[i][j].Alive = true
					nextBoard.cells[i][j].onceAlive = true
					nextBoard.cells[i][j].ChooseColony(g.cells[i][j], neighbours)
				}
			}
		}
	}

	g.cells = nextBoard.cells
}

// Draw renders the automaton on the terminal.
func (g *GameOfLifeDx) Draw() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	for i := 0; i < g.rows; i++ {
		for j := 0; j < g.cols; j++ {
			if g.showExplored {
				if g.cells[i][j].Alive {
					termbox.SetCell(
						j,
						i,
						g.cells[i][j].colony.symbol,
						g.cells[i][j].colony.color,
						OnceAliveBackground,
					)
				} else if g.cells[i][j].onceAlive {
					termbox.SetCell(j, i, ' ', termbox.ColorDefault, OnceAliveBackground)
				}
			} else {
				if g.cells[i][j].Alive {
					termbox.SetCell(
						j,
						i,
						g.cells[i][j].colony.symbol,
						g.cells[i][j].colony.color,
						render.DefaultBackground,
					)
				}
			}
		}
	}

	termbox.Flush()
}
