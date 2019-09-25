package render

import "github.com/nsf/termbox-go"

type Renderer interface {
	Clear()
	Flush()
}

type TermboxRenderer struct{}

func (r *TermboxRenderer) Clear() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
}

func (r *TermboxRenderer) Flush() {
	termbox.Flush()
}

func (r *TermboxRenderer) SetCell(col, row int, ch rune) {
	termbox.SetCell(col, row, ch, termbox.ColorDefault, termbox.ColorDefault)
}
