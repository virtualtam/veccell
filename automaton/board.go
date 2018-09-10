package automaton

import (
	"math/rand"

	"github.com/nsf/termbox-go"
)

// A Board holds the game's parameters and state.
type Board struct {
	rows             int
	cols             int
	borderCellsAlive bool
	cells            [][]Cell
}

// NewBoard creates and initializes a Board's components.
func NewBoard(rows, cols int, borderCellsAlive bool) Board {
	b := Board{
		rows:             rows,
		cols:             cols,
		borderCellsAlive: borderCellsAlive,
	}
	b.cells = make([][]Cell, b.rows)
	for i := 0; i < b.rows; i++ {
		b.cells[i] = make([]Cell, b.cols)
	}
	return b
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
func (b *Board) CreateGliderAt(row, col int) {
	b.cells[row-1][col].alive = true
	b.cells[row][col+1].alive = true
	b.cells[row+1][col-1].alive = true
	b.cells[row+1][col].alive = true
	b.cells[row+1][col+1].alive = true
}

// Randomize sets a board's Cells in a random state (alive|dead).
func (b *Board) Randomize() {
	for i := 0; i < b.rows; i++ {
		for j := 0; j < b.cols; j++ {
			b.cells[i][j].alive = rand.Intn(2) == 1
		}
	}
}

// RandomizeArea is similar to Randomize, expect it acts on a given portion of
// the Board.
func (b *Board) RandomizeArea(startRow, endRow, startCol, endCol int) {
	for i := startRow; i < endRow; i++ {
		for j := startCol; j < endCol; j++ {
			b.cells[i][j].alive = rand.Intn(2) == 1
		}
	}
}

// IsCellAlive returns the liveliness status of a Cell, including virtual Board
// borders.
func (b *Board) IsCellAlive(row, col int) bool {
	if row < 0 || row >= b.rows {
		// horizontal borders
		return b.borderCellsAlive
	}
	if col < 0 || col >= b.cols {
		// vertical borders
		return b.borderCellsAlive
	}
	// actual Cell
	return b.cells[row][col].alive
}

// CountLiveNeighbours returns the number of live Cells surrounding a Cell at a
// given position.
func (b *Board) CountLiveNeighbours(row, col int) int {
	liveNeighbours := 0
	for i := -1; i < 2; i++ {
		for j := -1; j < 2; j++ {
			if i == 0 && j == 0 {
				continue
			}
			if b.IsCellAlive(row+i, col+j) {
				liveNeighbours++
			}
		}
	}
	return liveNeighbours
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
func (b *Board) Next() {
	nextBoard := NewBoard(b.rows, b.cols, b.borderCellsAlive)

	for i := 0; i < b.rows; i++ {
		for j := 0; j < b.cols; j++ {
			liveNeighbours := b.CountLiveNeighbours(i, j)

			if b.cells[i][j].alive {
				if liveNeighbours == 2 || liveNeighbours == 3 {
					nextBoard.cells[i][j].alive = true
				}
			} else {
				if liveNeighbours == 3 {
					nextBoard.cells[i][j].alive = true
				}
			}
		}
	}

	b.cells = nextBoard.cells
}

// Draw renders the Board on the terminal.
func (b *Board) Draw() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	for i := 0; i < b.rows; i++ {
		for j := 0; j < b.cols; j++ {
			if b.cells[i][j].alive {
				termbox.SetCell(j, i, 'O', termbox.ColorDefault, termbox.ColorDefault)
			}
		}
	}

	termbox.Flush()
}
