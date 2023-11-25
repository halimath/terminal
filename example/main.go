package main

import (
	"fmt"

	"github.com/halimath/termx"
	"github.com/halimath/termx/keypress"
	"github.com/halimath/termx/sgr"
)

var fgColors = []sgr.SGR{
	sgr.FgBlack,
	sgr.FgRed,
	sgr.FgGreen,
	sgr.FgYellow,
	sgr.FgBlue,
	sgr.FgMagenta,
	sgr.FgCyan,
	sgr.FgWhite,
	sgr.FgLightBlack,
	sgr.FgLightRed,
	sgr.FgLightGreen,
	sgr.FgLightYellow,
	sgr.FgLightBlue,
	sgr.FgLightMagenta,
	sgr.FgLightCyan,
	sgr.FgLightWhite,
}

var bgColors = []sgr.SGR{
	sgr.BgBlack,
	sgr.BgRed,
	sgr.BgGreen,
	sgr.BgYellow,
	sgr.BgBlue,
	sgr.BgMagenta,
	sgr.BgCyan,
	sgr.BgWhite,
	sgr.BgLightBlack,
	sgr.BgLightRed,
	sgr.BgLightGreen,
	sgr.BgLightYellow,
	sgr.BgLightBlue,
	sgr.BgLightMagenta,
	sgr.BgLightCyan,
	sgr.BgLightWhite,
}

func main() {
	t := termx.New()
	defer t.Close()

	if err := t.EnableRawMode(true); err != nil {
		panic(err)
	}

	termx.SetWindowTitle(t, "termx example app")

	t.WriteString(termx.CursorHide)
	defer t.WriteString(termx.CursorShow)

	w, h, err := t.Size()
	if err != nil {
		panic(err)
	}

	t.WriteString(termx.ClearScreen)
	termx.SetCursorPosition(t, 1, 1)

	fmt.Fprintf(t, "Screen size:\t\t%dx%d", w, h)

	termx.MoveCursorBackward(t, w)
	termx.MoveCursorDown(t, 1)

	t.WriteString("Background Colors:\t")
	for _, c := range bgColors {
		sgr.Print(t, c, "   ")
	}

	termx.MoveCursorBackward(t, w)
	termx.MoveCursorDown(t, 1)

	t.WriteString("Foreground Colors:\t")
	for _, c := range fgColors {
		sgr.Print(t, c, "aBc")
	}

	termx.MoveCursorBackward(t, w)
	termx.MoveCursorDown(t, 1)

	t.WriteString("ANSI Color:\t\t")
	sgr.Print(t, sgr.Join(sgr.BgRGB(5, 0, 0), sgr.FgRGB(0, 3, 3)), "This should be printed green on red")

	termx.MoveCursorBackward(t, w)
	termx.MoveCursorDown(t, 1)

	fmt.Fprintf(t, "TrueColor [%v]:\t", t.TruecolorSupported())
	sgr.Print(t, sgr.Join(sgr.BgTrueColor(120, 0, 0), sgr.FgTrueColor(0, 120, 120)), "This should be printed green on red")

	termx.MoveCursorBackward(t, w)
	termx.MoveCursorDown(t, 1)

	t.WriteString("Press any key to see its internal representation; press C-x to display Bg color info; press C-c to quit")

	termx.MoveCursorBackward(t, w)
	termx.MoveCursorDown(t, 1)

	// Enable cursor application mode
	t.WriteString("\x1b[?1h")

	for {
		evt, raw, err := t.ReadKeyPress()

		termx.MoveCursorBackward(t, w)
		t.WriteString(termx.ClearLine)

		t.Printf("%s %#v %v", evt, raw, err)

		if evt == keypress.Ctrl('c') || evt == keypress.Char('q') {
			break
		}

		if evt == keypress.Ctrl('x') {
			r, g, b, err := termx.GetBackgroundColor(t)
			if err != nil {
				panic(err)
			}

			termx.MoveCursorBackward(t, w)
			termx.MoveCursorDown(t, 1)
			t.WriteString(termx.ClearLine)
			t.Printf("(%d, %d, %d) %v", r, g, b, err)
			termx.MoveCursorUp(t, 1)
		}

		if evt == keypress.Ctrl('v') {
			x, y, err := termx.GetCursorPosition(t)

			termx.MoveCursorBackward(t, w)
			termx.MoveCursorDown(t, 1)
			t.WriteString(termx.ClearLine)
			t.Printf("(%d,%d) %v", x, y, err)
			termx.MoveCursorUp(t, 1)
		}
	}

	termx.MoveCursorUp(t, 1)
	termx.MoveCursorBackward(t, w)

}
