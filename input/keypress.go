package input

import (
	"fmt"
)

type Event interface {
	evt()
}

// KeyPress defines an interface for all valid key presses. All KeyPress types define a String method. The
// returned string is safe for comparison and doing so is the recommended way to compare key presses.
type KeyPress interface {
	Event

	String() string

	keyPress()
}

// Char is a KeyPress that's a regular key.
type Char rune

func (Char) evt()      {}
func (Char) keyPress() {}
func (c Char) String() string {
	if c == ' ' {
		return "<Space>"
	}

	return fmt.Sprintf("%c", c)
}

// Ctrl is a KeyPress with a regular key combined with the control key.
type Ctrl rune

func (Ctrl) evt()      {}
func (Ctrl) keyPress() {}
func (c Ctrl) String() string {
	if c == ' ' {
		return "C-<Space>"
	}

	return fmt.Sprintf("C-%c", c)
}

// Alt is a KeyPress with a regular key combined with the alt or modifier key.
type Alt rune

func (Alt) evt()      {}
func (Alt) keyPress() {}

func (a Alt) String() string {
	return fmt.Sprintf("M-%c", a)
}

// FunctionKey is a KeyPress of a function key (or F-key). The integer number
// carries the key pressed. Note that only F1 to F10 are supported.
type FunctionKey uint8

func (FunctionKey) evt()      {}
func (FunctionKey) keyPress() {}
func (f FunctionKey) String() string {
	return fmt.Sprintf("<F%d>", f)
}

// SpecialKey implements special keys like cursor, delete, backspace...
// It's an enumerated type and only constants below are well-defined.
type SpecialKey int

func (SpecialKey) evt()      {}
func (SpecialKey) keyPress() {}
func (s SpecialKey) String() string {
	switch s {
	case Return:
		return "<Ret>"
	case Backspace:
		return "<Backspace>"
	case Tab:
		return "<Tab>"
	case Delete:
		return "<Del>"
	case CursorUp:
		return "<Up>"
	case CursorDown:
		return "<Down>"
	case CursorLeft:
		return "<Left>"
	case CursorRight:
		return "<Right>"
	case Escape:
		return "<Esc>"
	case Home:
		return "<Home>"
	case End:
		return "<End>"
	case PageUp:
		return "<PgUp>"
	case PageDown:
		return "<PgDn>"
	default:
		return "<?>"
	}
}

const (
	Return SpecialKey = iota + 1
	Backspace
	Tab
	Delete
	CursorUp
	CursorDown
	CursorLeft
	CursorRight
	Escape
	Home
	End
	PageUp
	PageDown
)
