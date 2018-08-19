package main

import (
	"github.com/nsf/termbox-go"
	"math/rand"
	"time"
)

type Cell struct {
	alive bool
}

type Board struct {
	height int
	width  int
	cells  [][]Cell
}

func (b *Board) Init() {
	b.cells = make([][]Cell, b.height)
	for i := 0; i < b.height; i++ {
		b.cells[i] = make([]Cell, b.width)
	}
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
		return false
	}
	if col < 0 || col >= b.width {
		return false
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

func (b *Board) Next() {
	nextBoard := Board{height: b.height, width: b.width}
	nextBoard.Init()

	for i := 0; i < b.height; i++ {
		for j := 0; j < b.width; j++ {
			liveNeighbours := b.CountLiveNeighbours(i, j)

			if liveNeighbours == 2 || liveNeighbours == 3 {
				nextBoard.cells[i][j].alive = true
			} else {
				nextBoard.cells[i][j].alive = false
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
				termbox.SetCell(j, i, 'O', termbox.ColorGreen, termbox.ColorDefault)
			}
		}
	}

	termbox.Flush()
}

func main() {
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

	board := Board{height: termHeight, width: termWidth}
	board.Init()
	//board.Randomize()
	//board.RandomizeArea(0, termHeight/2, 0, termWidth/2)
	board.RandomizeArea(termHeight/4, 3*termHeight/4, termWidth/4, 3*termWidth/4)
	board.Draw()

	delay := 1000 // Milliseconds

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
					delay += 100
				case termbox.KeyArrowDown:
					if delay > 100 {
						delay -= 100
					} else if delay > 10 {
						delay -= 10
					}
				}
			}

		default:
			time.Sleep(time.Duration(delay) * time.Millisecond)
			board.Next()
			board.Draw()
		}
	}
}
