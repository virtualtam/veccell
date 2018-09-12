package automaton

import (
	"github.com/nsf/termbox-go"
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
