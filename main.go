package main

import (
	"flag"
	"github.com/nsf/termbox-go"
	"math/rand"
	"time"
)

const (
	DefaultDelay            = 1000 // Milliseconds
	DefaultBorderCellsAlive = false
)

type Cell struct {
	alive bool
}

type Board struct {
	height           int
	width            int
	borderCellsAlive bool
	cells            [][]Cell
}

func (b *Board) Init() {
	b.cells = make([][]Cell, b.height)
	for i := 0; i < b.height; i++ {
		b.cells[i] = make([]Cell, b.width)
	}
}

func (b *Board) CreateGliderAt(row, col int) {
	b.cells[row-1][col].alive = true
	b.cells[row][col+1].alive = true
	b.cells[row+1][col-1].alive = true
	b.cells[row+1][col].alive = true
	b.cells[row+1][col+1].alive = true
}

func (b *Board) Randomize() {
	for i := 0; i < b.height; i++ {
		for j := 0; j < b.width; j++ {
			b.cells[i][j].alive = rand.Intn(2) == 1
		}
	}
}

func (b *Board) RandomizeArea(startRow, endRow, startCol, endCol int) {
	for i := startRow; i < endRow; i++ {
		for j := startCol; j < endCol; j++ {
			b.cells[i][j].alive = rand.Intn(2) == 1
		}
	}
}

func (b *Board) IsCellAlive(row, col int) bool {
	if row < 0 || row >= b.height {
		return b.borderCellsAlive
	}
	if col < 0 || col >= b.width {
		return b.borderCellsAlive
	}
	return b.cells[row][col].alive
}

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
	nextBoard := Board{
		height:           b.height,
		width:            b.width,
		borderCellsAlive: b.borderCellsAlive,
	}
	nextBoard.Init()

	for i := 0; i < b.height; i++ {
		for j := 0; j < b.width; j++ {
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

func (b *Board) Draw() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	for i := 0; i < b.height; i++ {
		for j := 0; j < b.width; j++ {
			if b.cells[i][j].alive {
				termbox.SetCell(j, i, 'O', termbox.ColorDefault, termbox.ColorDefault)
			}
		}
	}

	termbox.Flush()
}

func main() {
	// Command-line parameters
	var borderCellsAlive bool
	var delay int

	flag.BoolVar(&borderCellsAlive, "borderCellsAlive", DefaultBorderCellsAlive, "Whether the border cells are considered alive or dead")
	flag.IntVar(&delay, "delay", DefaultDelay, "Delay between two iterations (milliseconds)")
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
	board := Board{
		height:           termHeight,
		width:            termWidth,
		borderCellsAlive: borderCellsAlive,
	}
	board.Init()
	board.Randomize()
	//board.RandomizeArea(0, termHeight/2, 0, termWidth/2)
	//board.RandomizeArea(termHeight/4, 3*termHeight/4, termWidth/4, 3*termWidth/4)
	//board.CreateGliderAt(7, 7)
	board.Draw()

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
			}

		case <-drawQueue:
			board.Next()
			board.Draw()
		}
	}
}
