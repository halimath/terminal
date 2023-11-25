package keypress

import (
	"errors"
	"fmt"
	"unicode/utf8"
)

// KeyPress defines an interface for all valid key presses. All KeyPress types define a String method. The
// returned string is safe for comparison and doing so is the recommended way to compare key presses.
type KeyPress interface {
	String() string

	keyPress()
}

// Char is a KeyPress that's a regular key.
type Char rune

func (Char) keyPress() {}
func (c Char) String() string {
	if c == ' ' {
		return "<Space>"
	}

	return fmt.Sprintf("%c", c)
}

// Ctrl is a KeyPress with a regular key combined with the control key.
type Ctrl rune

func (Ctrl) keyPress() {}
func (c Ctrl) String() string {
	if c == ' ' {
		return "C-<Space>"
	}

	return fmt.Sprintf("C-%c", c)
}

// Alt is a KeyPress with a regular key combined with the alt or modifier key.
type Alt rune

func (Alt) keyPress() {}

func (a Alt) String() string {
	return fmt.Sprintf("M-%c", a)
}

// FunctionKey is a KeyPress of a function key (or F-key). The integer number
// carries the key pressed. Note that only F1 to F10 are supported.
type FunctionKey uint8

func (FunctionKey) keyPress() {}
func (f FunctionKey) String() string {
	return fmt.Sprintf("<F%d>", f)
}

// SpecialKey implements special keys like cursor, delete, backspace...
// It's an enumerated type and only constants below are well-defined.
type SpecialKey int

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

var ErrInvalidKeyPressBytes = errors.New("invalid key-press byte sequence")

const (
	keyCodeTab         = 9
	keyCodeReturn      = 13
	keyCodeEscape      = 27
	keyCodeOpenBracket = '['
)

func Decode(b []byte) (KeyPress, error) {
	if len(b) == 0 {
		return nil, ErrInvalidKeyPressBytes
	}

	switch len(b) {
	case 1:
		return decodeSingleByteKeyPress(b[0])
	case 2:
		// Any two byte sequence is a unicode character
		r, _ := utf8.DecodeRune(b)
		return Char(r), nil
	case 3:
		return decodeThreeByteKeyPress(b)
	case 4:
		return decodeFourByteKeyPress(b)
	}

	return nil, ErrInvalidKeyPressBytes
}

func decodeSingleByteKeyPress(b byte) (KeyPress, error) {
	switch b {
	case 0:
		return Ctrl(' '), nil
	case keyCodeTab:
		return Tab, nil
	case keyCodeReturn:
		return Return, nil
	case keyCodeEscape:
		return Escape, nil
	case 127:
		return Backspace, nil
	}

	// Bytes 1 - 26 represent CTRL-a .. CTRL-z
	if b < keyCodeEscape {
		return Ctrl(rune('a' + b - 1)), nil
	}

	// Otherwise its a plain character
	return Char(b), nil
}

func decodeThreeByteKeyPress(b []byte) (KeyPress, error) {
	if b[0] != keyCodeEscape {
		r, _ := utf8.DecodeRune(b)
		return Char(r), nil
	}

	if b[1] == keyCodeOpenBracket {
		switch b[2] {
		case 0x48:
			return Home, nil
		case 0x46:
			return End, nil
		case 65:
			return CursorUp, nil
		case 66:
			return CursorDown, nil
		case 67:
			return CursorRight, nil
		case 68:
			return CursorLeft, nil
		default:
			return nil, ErrInvalidKeyPressBytes
		}
	}

	if b[1] == 0x4f {
		switch b[2] {
		case 0x41:
			return CursorUp, nil
		case 0x42:
			return CursorDown, nil
		case 0x43:
			return CursorRight, nil
		case 0x44:
			return CursorLeft, nil
		case 0x46:
			return End, nil
		case 0x48:
			return Home, nil

		default:
			return nil, ErrInvalidKeyPressBytes
		}
	} else if b[1] == 79 && b[2] >= 80 && b[2] <= 89 {
		return FunctionKey(b[2] - 79), nil
	}

	return nil, ErrInvalidKeyPressBytes
}

func decodeFourByteKeyPress(b []byte) (KeyPress, error) {
	if b[0] == keyCodeEscape && b[1] == keyCodeOpenBracket {
		if b[2] == 51 && b[3] == 126 {
			return Delete, nil
		}

		if b[2] == 0x35 && b[3] == 0x7e {
			return PageUp, nil // Non Ca Mode
		}

		if b[2] == 0x36 && b[3] == 0x7e {
			return PageDown, nil // Non Ca Mode
		}
	}

	r, _ := utf8.DecodeRune(b)
	return Char(r), nil
}
