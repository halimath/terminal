// This example program demonstrates how to output and query colors and screen attributes.
package main

import (
	"fmt"

	"github.com/halimath/terminal"
	"github.com/halimath/terminal/sgr"
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
	t := terminal.New()

	if t.IsTerminal() {
		w, h, err := t.Size()
		if err != nil {
			panic(err)
		}

		fmt.Fprintf(t, "Screen size:\t\t%dx%d", w, h)
		fmt.Fprintln(t)
	}

	t.WriteString("Background Colors:\t")
	for _, c := range bgColors {
		fmt.Fprint(t, c.Apply("   "))
	}

	fmt.Fprintln(t)
	t.WriteString("Foreground Colors:\t")
	for _, c := range fgColors {
		fmt.Fprint(t, c.Apply("aBc"))
	}

	fmt.Fprintln(t)
	t.WriteString("ANSI Color:\t\t")
	fmt.Fprint(t, sgr.BgRGB(5, 0, 0).Join(sgr.FgRGB(0, 3, 3)).Apply("This should be printed green on red"))

	fmt.Fprintln(t)
	fmt.Fprintf(t, "TrueColor [%v]:\t", terminal.IsTruecolorSupported())
	fmt.Fprint(t, sgr.BgTrueColor(120, 0, 0).Join(sgr.FgTrueColor(0, 120, 120)).Apply("This should be printed green on red"))
	fmt.Fprintln(t)
}
