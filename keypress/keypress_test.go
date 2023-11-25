package keypress

import (
	"testing"

	"github.com/halimath/expect"
	"github.com/halimath/expect/is"
)

func TestDecode(t *testing.T) {
	type testCase struct {
		in   []byte
		want KeyPress
		err  error
	}

	tests := []testCase{
		{nil, nil, ErrInvalidKeyPressBytes},

		{[]byte("a"), Char('a'), nil},
		{[]byte(" "), Char(' '), nil},
		{[]byte{0x0}, Ctrl(' '), nil},
		{[]byte{0x1}, Ctrl('a'), nil},
		{[]byte{0x1a}, Ctrl('z'), nil},
		{[]byte{0x9}, Tab, nil},
		{[]byte{0xd}, Return, nil},
		{[]byte{0x1b}, Escape, nil},
		{[]byte{0x7f}, Backspace, nil},

		{[]byte("ö"), Char('ö'), nil},

		{[]byte("世"), Char('世'), nil},
		{[]byte{0x1b, 0x5b, 0x46}, End, nil},
		{[]byte{0x1b, 0x5b, 0x48}, Home, nil},
	}

	for _, test := range tests {
		got, err := Decode(test.in)
		expect.That(t,
			is.DeepEqualTo(got, test.want),
			is.EqualTo(err, test.err),
		)
	}
}
