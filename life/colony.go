package life

import (
	"github.com/nsf/termbox-go"
)

// A Colony represents a type of population, and can feature distinctive
// attributes.
type Colony struct {
	symbol rune
	color  termbox.Attribute
}

// The available Colonies.
var Colonies = []Colony{
	{
		symbol: 'O',
		color:  termbox.ColorDefault,
	},
	{
		symbol: '-',
		color:  termbox.ColorRed,
	},
	{
		symbol: '+',
		color:  termbox.ColorGreen,
	},
	{
		symbol: 'X',
		color:  termbox.ColorBlue,
	},
	{
		symbol: '$',
		color:  termbox.ColorYellow,
	},
	{
		symbol: '*',
		color:  termbox.ColorMagenta,
	},
	{
		symbol: '%',
		color:  termbox.ColorCyan,
	},
	{
		symbol: '8',
		color:  termbox.ColorWhite,
	},
}
