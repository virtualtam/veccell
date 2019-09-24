package automaton

import (
	"math/rand"

	"github.com/nsf/termbox-go"
)

// A GameOfLife holds the game's parameters and state.
type GameOfLife struct {
	rows             int
	cols             int
	borderCellsAlive bool
	cells            [][]Cell
}

// NewGameOfLife creates and initializes a Game Of Life automaton.
func NewGameOfLife(rows, cols int, borderCellsAlive bool) GameOfLife {
	g := GameOfLife{
		rows:             rows,
		cols:             cols,
		borderCellsAlive: borderCellsAlive,
	}
	g.cells = make([][]Cell, g.rows)
	for i := 0; i < g.rows; i++ {
		g.cells[i] = make([]Cell, g.cols)
	}
	return g
}

// CreateGliderAt creates a glider object centered on the provided coordinates.
//
// Glider pattern:
//
//     O
//      O
//    OOO
//
// See https://en.wikipedia.org/wiki/Glider_(Conway%27s_Life)
func (g *GameOfLife) CreateGliderAt(row, col int) {
	g.cells[row-1][col].Alive = true
	g.cells[row][col+1].Alive = true
	g.cells[row+1][col-1].Alive = true
	g.cells[row+1][col].Alive = true
	g.cells[row+1][col+1].Alive = true
}

// Randomize sets a board's Cells in a random state .Alive|dead).
func (g *GameOfLife) Randomize() {
	for i := 0; i < g.rows; i++ {
		for j := 0; j < g.cols; j++ {
			g.cells[i][j].Alive = rand.Intn(2) == 1
		}
	}
}

// RandomizeArea is similar to Randomize, expect it acts on a given portion of
// the Board.
func (g *GameOfLife) RandomizeArea(startRow, endRow, startCol, endCol int) {
	for i := startRow; i < endRow; i++ {
		for j := startCol; j < endCol; j++ {
			g.cells[i][j].Alive = rand.Intn(2) == 1
		}
	}
}

// LiveNeighboursAt returns the live Cells surrounding the Cell at the given
// position.
func (g *GameOfLife) LiveNeighboursAt(row, col int) []*Cell {
	neighbours := []*Cell{}

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
// 1. any live cell with fewer than two live neighbours dies, as if
// by under population
//
// 2. any live cell with two or three live neighbours lives on to
// the the next generation
//
// 3. any live cell with more than three live neighbours dies, as if
// by overpopulation
//
// 4. any dead cell with exactly three live neighbours becomes a
// live cell, as if by reproduction
func (g *GameOfLife) Next() {
	nextBoard := NewGameOfLife(g.rows, g.cols, g.borderCellsAlive)

	for i := 0; i < g.rows; i++ {
		for j := 0; j < g.cols; j++ {
			neighbours := g.LiveNeighboursAt(i, j)
			liveNeighbours := len(neighbours)

			if g.borderCellsAlive {
				if i == 0 || i == g.rows-1 {
					if j == 0 || j == g.cols-1 {
						liveNeighbours += 5
					} else {
						liveNeighbours += 3
					}
				} else {
					if j == 0 || j == g.cols-1 {
						liveNeighbours += 3
					}
				}
			}

			if g.cells[i][j].Alive {
				if liveNeighbours == 2 || liveNeighbours == 3 {
					nextBoard.cells[i][j].Alive = true
				}
			} else {
				if liveNeighbours == 3 {
					nextBoard.cells[i][j].Alive = true
				}
			}
		}
	}

	g.cells = nextBoard.cells
}

// Draw renders the automaton on the terminal.
func (g *GameOfLife) Draw() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	for i := 0; i < g.rows; i++ {
		for j := 0; j < g.cols; j++ {
			if g.cells[i][j].Alive {
				termbox.SetCell(j, i, 'O', termbox.ColorDefault, termbox.ColorDefault)
			}
		}
	}

	termbox.Flush()
}
