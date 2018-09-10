package automaton

import (
	"time"

	"github.com/nsf/termbox-go"
)

// The Controller handles Termbox events, user input and render queues.
type Controller struct {
	breakQueue chan bool
	drawQueue  chan bool
	eventQueue chan termbox.Event

	automaton Automaton
	delay     *int
}

// NewController creates and initializes a Controller.
func NewController(automaton Automaton, delay *int) Controller {
	c := Controller{automaton: automaton, delay: delay}

	c.breakQueue = make(chan bool)
	c.drawQueue = make(chan bool)
	c.eventQueue = make(chan termbox.Event)

	go func(delay *int) {
		for {
			time.Sleep(time.Duration(*delay) * time.Millisecond)
			c.drawQueue <- true
		}
	}(c.delay)

	go func() {
		for {
			c.eventQueue <- termbox.PollEvent()
		}
	}()

	return c
}

// HandleUserInput handles user events such as key presses.
func (c *Controller) HandleUserInput(ev termbox.Event) {
	switch ev.Type {
	case termbox.EventKey:
		switch ev.Key {
		case termbox.KeyEsc:
			c.Break()

		case termbox.KeyCtrlC:
			c.Break()

		case termbox.KeyArrowUp:
			switch {
			case *c.delay < 10:
				*c.delay++
			case *c.delay < 100:
				*c.delay += 10
			default:
				*c.delay += 100
			}

		case termbox.KeyArrowDown:
			switch {
			case *c.delay > 100:
				*c.delay -= 100
			case *c.delay > 10:
				*c.delay -= 10
			case *c.delay > 2:
				*c.delay--
			}
		}

		switch ev.Ch {
		case 'q':
			c.Break()
		case 'r':
			c.automaton.Randomize()
		}
	}
}

// Break triggers the interruption of the main render loop.
func (c *Controller) Break() {
	go func() {
		// Run asynchronously not to block the main goroutine
		c.breakQueue <- true
	}()
}

// Loop cycles through the automaton's iterations, and renders its state.
func (c *Controller) Loop() {
mainloop:
	for {
		select {
		case <-c.drawQueue:
			c.automaton.Next()
			c.automaton.Draw()

		case <-c.breakQueue:
			break mainloop

		case ev := <-c.eventQueue:
			c.HandleUserInput(ev)
		}
	}
}
