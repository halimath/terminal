package main

import (
	"github.com/halimath/termx"
	"github.com/halimath/termx/csi"
	"github.com/halimath/termx/keypress"
)

func main() {
	t := termx.New()

	if err := t.EnterRawMode(); err != nil {
		panic(err)
	}
	defer t.ExitRawMode()

	t.WriteString(csi.EnterCaMode)
	defer t.WriteString(csi.ExitCaMode)

	t.Print(csi.SetWindowTitle("termx example app"))

	t.WriteString(csi.CursorHide)
	defer t.WriteString(csi.CursorShow)

	t.WriteString(csi.ClearScreen)
	t.Print(csi.SetCursorPosition(1, 1))

	t.WriteString("Press any key to see its internal representation; press C-x to display Bg color info; press C-c to quit")

	t.Print(csi.MoveCursorBackward(200))
	t.Print(csi.MoveCursorDown(1))

	// Enable cursor application mode
	t.WriteString("\x1b[?1h")

	for {
		evt, raw, err := t.ReadKeyPress()

		t.Print(csi.MoveCursorBackward(200))
		t.WriteString(csi.ClearLine)

		t.Printf("%s %#v %v", evt, raw, err)

		if evt == keypress.Ctrl('c') || evt == keypress.Char('q') {
			break
		}

		if evt == keypress.Ctrl('x') {
			r, g, b, err := csi.GetBackgroundColor(t)
			if err != nil {
				panic(err)
			}

			t.Print(csi.MoveCursorBackward(200))
			t.Print(csi.MoveCursorDown(1))
			t.WriteString(csi.ClearLine)
			t.Printf("(%d, %d, %d) %v", r, g, b, err)
			t.Print(csi.MoveCursorUp(1))
		}

		if evt == keypress.Ctrl('v') {
			x, y, err := csi.GetCursorPosition(t)

			t.Print(csi.MoveCursorBackward(200))
			t.Print(csi.MoveCursorDown(1))
			t.WriteString(csi.ClearLine)
			t.Printf("(%d,%d) %v", x, y, err)
			t.Print(csi.MoveCursorUp(1))
		}
	}

	t.Print(csi.MoveCursorUp(1))
	t.Print(csi.MoveCursorBackward(200))

}
