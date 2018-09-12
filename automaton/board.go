package automaton

import (
	"math/rand"

	"github.com/nsf/termbox-go"
)

const (
	DefaultBackground   = termbox.ColorDefault
	OnceAliveBackground = termbox.ColorBlack
)

// A Board holds the game's parameters and state.
type Board struct {
	rows             int
	cols             int
	borderCellsAlive bool
	cells            [][]Cell

	nColonies    int      // FIXME: Life DX
	colonies     []Colony // FIXME: Life DX
	showExplored bool     // FIXME: Life DX
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

// NewDxBoard creates and initializes a Board.
// FIXME: Life DX, interface{}, API
func NewDxBoard(rows, cols, nColonies int, showExplored bool) Board {
	b := Board{
		rows:         rows,
		cols:         cols,
		nColonies:    nColonies,
		showExplored: showExplored,
	}
	b.colonies = Colonies[:b.nColonies]
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
	if b.colonies != nil {
		// FIXME: Life DX, API
		b.DxRandomize()
		return
	}
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

// DxRandomize sets a board's Cells in a random state (alive|dead) and assigns
// it to a randomly chosen Colony.
// FIXME: Life DX, interface{}, API
func (b *Board) DxRandomize() {
	for i := 0; i < b.rows; i++ {
		for j := 0; j < b.cols; j++ {
			b.cells[i][j].alive = rand.Intn(5) == 1
			b.cells[i][j].onceAlive = b.cells[i][j].alive
			b.cells[i][j].colony = &b.colonies[rand.Intn(b.nColonies)]
		}
	}
}

// IsCellAlive returns the liveliness status of a Cell, including virtual Board
// borders.
// FIXME: merge with LiveNeighboursAt from Life DX
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
// FIXME: merge with LiveNeighboursAt from Life DX
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

// LiveNeighboursAt returns the live Cells surrounding the Cell at the given
// position.
func (b *Board) LiveNeighboursAt(row, col int) []*Cell {
	neighbours := []*Cell{}

	for i := -1; i < 2; i++ {
		for j := -1; j < 2; j++ {
			tmpRow := row + i
			tmpCol := col + j

			if tmpRow < 0 || tmpRow >= b.rows {
				// assume horizontal borders are made of dead cells
				continue
			}
			if tmpCol < 0 || tmpCol >= b.cols {
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
	if b.colonies != nil {
		// FIXME: Life DX, API
		b.DxNext()
		return
	}
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

// DxNext computes the next iteration of the game, and updates the Board's state.
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
func (b *Board) DxNext() {
	nextBoard := NewDxBoard(b.rows, b.cols, b.nColonies, b.showExplored)

	for i := 0; i < b.rows; i++ {
		for j := 0; j < b.cols; j++ {
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
func (b *Board) Draw() {
	if b.colonies != nil {
		// FIXME: Life DX, API
		b.DxDraw()
		return
	}

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

// DxDraw renders the Board on the terminal.
func (b *Board) DxDraw() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	for i := 0; i < b.rows; i++ {
		for j := 0; j < b.cols; j++ {
			if b.showExplored {
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
