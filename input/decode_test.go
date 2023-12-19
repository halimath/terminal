package input

import (
	"testing"

	"github.com/halimath/expect"
	"github.com/halimath/expect/is"
)

func TestDecode(t *testing.T) {
	type testCase struct {
		in   []byte
		want Event
		err  error
	}

	tests := []testCase{
		// Invalid input
		{nil, nil, ErrInvalidInputBytes},
		{[]byte("\x1b[M\x1f\x1f\x1f"), nil, ErrInvalidInputBytes},
		{[]byte{0x1b, 0x5b, 0x99}, nil, ErrInvalidInputBytes},

		// Single byte keys
		{[]byte("a"), Char('a'), nil},
		{[]byte(" "), Char(' '), nil},
		{[]byte{0x0}, Ctrl(' '), nil},
		{[]byte{0x1}, Ctrl('a'), nil},
		{[]byte{0x1a}, Ctrl('z'), nil},
		{[]byte{0x9}, Tab, nil},
		{[]byte{0xd}, Return, nil},
		{[]byte{0x1b}, Escape, nil},
		{[]byte{0x7f}, Backspace, nil},

		// UTF8 encoded runes
		{[]byte("ö"), Char('ö'), nil}, // Two bytes
		{[]byte("世"), Char('世'), nil}, // Three bytes
		{[]byte("𐍈"), Char('𐍈'), nil}, // Four bytes

		// Multi byte special keys in normal mode
		{[]byte{0x1b, 0x5b, 0x41}, CursorUp, nil},
		{[]byte{0x1b, 0x5b, 0x42}, CursorDown, nil},
		{[]byte{0x1b, 0x5b, 0x43}, CursorRight, nil},
		{[]byte{0x1b, 0x5b, 0x44}, CursorLeft, nil},
		{[]byte{0x1b, 0x5b, 0x46}, End, nil},
		{[]byte{0x1b, 0x5b, 0x48}, Home, nil},

		// Multi byte special keys in application mode
		{[]byte{0x1b, 0x4f, 0x41}, CursorUp, nil},
		{[]byte{0x1b, 0x4f, 0x42}, CursorDown, nil},
		{[]byte{0x1b, 0x4f, 0x43}, CursorRight, nil},
		{[]byte{0x1b, 0x4f, 0x44}, CursorLeft, nil},
		{[]byte{0x1b, 0x4f, 0x46}, End, nil},
		{[]byte{0x1b, 0x4f, 0x48}, Home, nil},
		{[]byte{0x1b, 0x4f, 0x50}, FunctionKey(1), nil},
		{[]byte{0x1b, 0x4f, 0x53}, FunctionKey(4), nil},
		{[]byte{0x1b, 0x5b, 0x33, 0x7e}, Delete, nil},
		{[]byte{0x1b, 0x5b, 0x35, 0x7e}, PageUp, nil},
		{[]byte{0x1b, 0x5b, 0x36, 0x7e}, PageDown, nil},

		// X10 mouse events
		{[]byte("\x1b[M!!!"), MouseEvent{Button: 2, X: 1, Y: 1}, nil},
		{[]byte("\x1b[M\"!!"), MouseEvent{Button: 3, X: 1, Y: 1}, nil},
		{[]byte("\x1b[M#!!"), MouseEvent{Button: 1, X: 1, Y: 1, Release: true}, nil},

		// SGR encoded mouse events
		{[]byte("\x1b[<0;1;1M"), MouseEvent{Button: 1, X: 1, Y: 1}, nil},
		{[]byte("\x1b[<1;1;1m"), MouseEvent{Button: 2, X: 1, Y: 1, Release: true}, nil},
	}

	for _, test := range tests {
		got, err := Decode(test.in)
		expect.WithMessage(t, "pattern %v", test.in).
			That(
				is.DeepEqualTo(got, test.want),
				is.Error(err, test.err),
			)
	}
}
