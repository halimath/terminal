package input

import (
	"testing"

	"github.com/halimath/expect"
	"github.com/halimath/expect/is"
)

func TestKeyPress_String(t *testing.T) {
	tests := map[KeyPress]string{
		Char(' '):      "<Space>",
		Char('a'):      "a",
		Ctrl(' '):      "C-<Space>",
		Ctrl('a'):      "C-a",
		Alt('a'):       "M-a",
		FunctionKey(1): "<F1>",
		Return:         "<Ret>",
		Backspace:      "<Backspace>",
		Tab:            "<Tab>",
		Delete:         "<Del>",
		CursorUp:       "<Up>",
		CursorDown:     "<Down>",
		CursorLeft:     "<Left>",
		CursorRight:    "<Right>",
		Escape:         "<Esc>",
		Home:           "<Home>",
		End:            "<End>",
		PageUp:         "<PgUp>",
		PageDown:       "<PgDn>",
	}

	for in, want := range tests {
		expect.That(t, is.EqualTo(in.String(), want))
	}

}
