// This example program demonstrates how to output and query colors and screen attributes.
package main

import (
	"fmt"

	"github.com/halimath/termx"
	"github.com/halimath/termx/csi"
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

	if t.IsTerminal() {
		w, h, err := t.Size()
		if err != nil {
			panic(err)
		}

		t.WriteString(csi.ClearScreen)
		t.Print(csi.SetCursorPosition(1, 1))

		fmt.Fprintf(t, "Screen size:\t\t%dx%d", w, h)
	}

	t.Print(csi.MoveCursorBackward(200))
	t.Print(csi.MoveCursorDown(1))

	t.WriteString("Background Colors:\t")
	for _, c := range bgColors {
		sgr.Print(t, c, "   ")
	}

	t.Print(csi.MoveCursorBackward(200))
	t.Print(csi.MoveCursorDown(1))

	t.WriteString("Foreground Colors:\t")
	for _, c := range fgColors {
		sgr.Print(t, c, "aBc")
	}

	t.Print(csi.MoveCursorBackward(200))
	t.Print(csi.MoveCursorDown(1))

	t.WriteString("ANSI Color:\t\t")
	sgr.Print(t, sgr.Join(sgr.BgRGB(5, 0, 0), sgr.FgRGB(0, 3, 3)), "This should be printed green on red")

	t.Print(csi.MoveCursorBackward(200))
	t.Print(csi.MoveCursorDown(1))

	fmt.Fprintf(t, "TrueColor [%v]:\t", termx.IsTruecolorSupported())
	sgr.Print(t, sgr.Join(sgr.BgTrueColor(120, 0, 0), sgr.FgTrueColor(0, 120, 120)), "This should be printed green on red")
}
