package main

import (
	"flag"

	"github.com/halimath/terminal"
	"github.com/halimath/terminal/csi"
	"github.com/halimath/terminal/input"
	"github.com/halimath/terminal/sgr"
)

func main() {
	useAlternateScreenBuffer := flag.Bool("alt-buffer", false, "Use alternative buffer")
	useApplicationMode := flag.Bool("app-mode", false, "Use application mode")
	enableMouse := flag.Bool("mouse", false, "Enable mouse tracking")
	flag.Parse()

	t := terminal.New()

	if *useAlternateScreenBuffer {
		t.Print(csi.UseAlternateScreenBuffer)
		defer t.Print(csi.UseMainScreenBuffer)
	}

	if *useApplicationMode {
		t.WriteString(csi.EnableApplicationMode)
		defer t.WriteString(csi.DisableApplicationMode)
	}

	if err := t.EnterRawMode(); err != nil {
		panic(err)
	}
	defer t.ExitRawMode()

	if *enableMouse {
		t.Print(csi.EnableMouseTracking, csi.EnableMouseSGREncoding, csi.EnableMouseButtonEvent)
		defer t.Print(csi.DisableMouseTracking, csi.DisableMouseSGREncoding, csi.DisableMouseButtonEvent)
	}

	t.Print(csi.SetWindowTitle("termx input example app"))

	t.WriteString(csi.CursorHide)
	defer t.WriteString(csi.CursorShow)

	t.WriteString(csi.ClearScreen)
	t.Print(csi.SetCursorPosition(1, 1))

	sgr.Print(t, sgr.Bold, "github.com/halimath/termx input example application")
	t.Print(csi.MoveCursorBackward(200))
	t.Print(csi.MoveCursorDown(1))

	t.Printf("Application Mode: %s; Alternative Buffer: %s, Mouse Tracking: %s",
		sgr.Formatf(sgr.Bold, "%v", *useApplicationMode),
		sgr.Formatf(sgr.Bold, "%v", *useAlternateScreenBuffer),
		sgr.Formatf(sgr.Bold, "%v", *enableMouse),
	)
	t.Print(csi.MoveCursorBackward(200))
	t.Print(csi.MoveCursorDown(1))

	sgr.Print(t, sgr.Faint, "Press any key to see its internal representation; press C-x to display Bg color info; press C-c to quit")
	t.Print(csi.MoveCursorBackward(200))
	t.Print(csi.MoveCursorDown(1))

	for {
		evt, raw, err := t.ReadInputEvent()

		t.Print(csi.MoveCursorBackward(200))
		t.WriteString(csi.ClearLine)

		t.Printf("%s %#v %v", evt, raw, err)

		if evt == input.Ctrl('c') || evt == input.Char('q') {
			break
		}

		if evt == input.Ctrl('x') {
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

		if evt == input.Ctrl('v') {
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
