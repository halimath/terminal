// Package sgr contains definitions for _select graphic rendition_ instructions that, when applied to some
// string, instruct the terminal emulator to select different colors, fonts, decorations, ... for the display
// of the string.
//
// The package defines a type SGR which is just a string. SGRs are defined as constants (for static, parameter
// less) instructions or factory functions for parameterized instructions (such as RGB colors). Multiple
// instructions can be joined together to form a composite SGR. The Join function performs this task.
// When applied to a string (i.e. via the Format function) the resulting string can be written to a terminal
// or some buffer as a plain string. Keep in mind that the sequences contain special characters which can
// scramble the ouput when send to a device that does not interpret them.
package sgr

import (
	"fmt"
	"strings"

	"github.com/halimath/termx/csi"
)

// SGR defines a type for Select Graphic Rendition instructions.
type SGR string

const (
	sgrTerminator = 'm' // The final byte (defined as a string) to append to all SGRs
	sgrSeparator  = ';' // The parameter separator within SGRs - ';' seems to be most compatible
)

// Escape creates a CSI escape sequence activating all rendition instructions given in s.
func Escape(s SGR) string {
	return fmt.Sprintf("%s%s%c", csi.CSI, s, sgrTerminator)
}

// Join joins all SGRs in s together.
func Join(s ...SGR) SGR {
	if len(s) == 0 {
		return ""
	}

	if len(s) == 1 {
		return s[0]
	}

	var b strings.Builder
	for i, s := range s {
		if i > 0 {
			b.WriteByte(sgrSeparator)
		}
		b.WriteString(string(s))
	}

	return SGR(b.String())
}

const (
	// Basic rendition instructions
	ResetAll   SGR = "0" // reset all SGR effects to their default
	Bold       SGR = "1" // bold or increased intensity
	Faint      SGR = "2" // faint or decreased intensity
	Italic     SGR = "3" // Italic mode
	Underlined SGR = "4" // singly underlined
	Blink      SGR = "5" // slow blink
	Invert     SGR = "7" // Invert Fg/Bg colors

	// Rendition instructions for standard foregroud colors
	FgBlack   SGR = "30"
	FgRed     SGR = "31"
	FgGreen   SGR = "32"
	FgYellow  SGR = "33"
	FgBlue    SGR = "34"
	FgMagenta SGR = "35"
	FgCyan    SGR = "36"
	FgWhite   SGR = "37"

	// Rendition instructions for standard background colors
	BgBlack   SGR = "40"
	BgRed     SGR = "41"
	BgGreen   SGR = "42"
	BgYellow  SGR = "43"
	BgBlue    SGR = "44"
	BgMagenta SGR = "45"
	BgCyan    SGR = "46"
	BgWhite   SGR = "47"

	// Rendition instructions for light foregroud colors
	FgLightBlack   SGR = "90"
	FgLightRed     SGR = "91"
	FgLightGreen   SGR = "92"
	FgLightYellow  SGR = "93"
	FgLightBlue    SGR = "94"
	FgLightMagenta SGR = "95"
	FgLightCyan    SGR = "96"
	FgLightWhite   SGR = "97"

	// Rendition instructions for light background colors
	BgLightBlack   SGR = "100"
	BgLightRed     SGR = "101"
	BgLightGreen   SGR = "102"
	BgLightYellow  SGR = "103"
	BgLightBlue    SGR = "104"
	BgLightMagenta SGR = "105"
	BgLightCyan    SGR = "106"
	BgLightWhite   SGR = "107"
)

// FgRGB creates a SGR to set the foreground color to one of the 256 colors based on red, green and blue
// components. Note that for ANSI r, g, b must be >= 0 and <= 5. Any other value will cause a panic.
func FgRGB(r, g, b int) SGR {
	assertValidRGB(r, g, b)
	return SGR(fmt.Sprintf("38;5;%d", rgbColorValue(r, g, b)))
}

// BgRGB creates a SGR to set the background color to one of the 256 colors based on red, green and blue
// components. Note that for ANSI r, g, b must be >= 0 and <= 5. Any other value will cause a panic.
func BgRGB(r, g, b int) SGR {
	assertValidRGB(r, g, b)
	return SGR(fmt.Sprintf("48;5;%d", rgbColorValue(r, g, b)))
}

func assertValidRGB(r, g, b int) {
	if r < 0 || r > 5 || b < 0 || b > 5 || g < 0 || g > 5 {
		panic(fmt.Sprintf("invalid ANSI color: %d;%d;%d", r, g, b))
	}
}

func rgbColorValue(r, g, b int) int {
	return 16 + 36*r + 6*g + b
}

// FgTrueColor creates a SGR that sets the foreground color to the true color value given with r, g, b.
func FgTrueColor(r, g, b uint8) SGR {
	return SGR(fmt.Sprintf("38;2;%d;%d;%d", r, g, b))
}

// BgTrueColor creates a SGR that sets the background color to the true color value given with r, g, b.
func BgTrueColor(r, g, b uint8) SGR {
	return SGR(fmt.Sprintf("48;2;%d;%d;%d", r, g, b))
}

// Format applies sgr to s and returns the resulting string.
func Format(sgr SGR, s string) string {
	var buf strings.Builder
	buf.Grow(len(s) + len(sgr) + len(ResetAll) + 6)

	buf.WriteString(Escape(sgr))
	buf.WriteString(s)
	buf.WriteString(Escape(ResetAll))
	return buf.String()
}

// Formatf works like Format with fmt.Sprintf applied to format and args before.
func Formatf(sgr SGR, format string, args ...any) string {
	return Format(sgr, fmt.Sprintf(format, args...))
}
