package termx

import (
	"fmt"
	"io"
	"strconv"
)

const (
	// Commands for general cursor management
	CursorSave    = "\x1B[s"    // Saves the current cursor
	CursorRestore = "\x1B[u"    // Restores a previously stored cursor
	CursorHide    = "\x1B[?25l" // Hide cursor
	CursorShow    = "\x1B[?25h" // Show cursor

	// Commands to set cursor style
	CursorBlinkingBlock     = "\x1B[\x31 q" // Change the cursor style to blinking block
	CursorSteadyBlock       = "\x1B[\x32 q" // Change the cursor style to steady block
	CursorBlinkingUnderline = "\x1B[\x33 q" // Change the cursor style to blinking underline
	CursorSteadyUnderline   = "\x1B[\x34 q" // Change the cursor style to steady underline
	CursorBlinkingBar       = "\x1B[\x35 q" // Change the cursor style to blinking bar
	CursorSteadyBar         = "\x1B[\x36 q" // Change the cursor style to steady bar
)

func MoveCursorUp(w io.Writer, n int) (int, error) {
	return fmt.Fprintf(w, "%s%dA", CSI, n)
}

func MoveCursorDown(w io.Writer, n int) (int, error) {
	return fmt.Fprintf(w, "%s%dB", CSI, n)
}

func MoveCursorForward(w io.Writer, n int) (int, error) {
	return fmt.Fprintf(w, "%s%dC", CSI, n)
}

func MoveCursorBackward(w io.Writer, n int) (int, error) {
	return fmt.Fprintf(w, "%s%dD", CSI, n)
}

// SetCursorPosition formats a CSI to position the cursor at (x,y). It writes that CSI to w and returns the
// number of bytes writen and any error.
//
// According to ANSI terminal specs both coordinates are 1 based. This function adheres to that spec.
func SetCursorPosition(w io.Writer, x, y int) (int, error) {
	return fmt.Fprintf(w, "%s%d;%dH", CSI, y, x)
}

type position struct {
	x, y int
}

const getCursorPositionQuery = "\x1B[6n"

// GetCursorPosition queries the current cursor position from t and returns the x and y coordinates. Note that
// according to ANSI specs coordinates are 1 based, so the upper left corner is (1, 1).
func GetCursorPosition(rw io.ReadWriter) (x, y int, err error) {
	p, err := execQuery(rw, getCursorPositionQuery, 64, func(res []byte) (p position, err error) {
		if len(res) < 6 || res[0] != '\x1B' || res[1] != '[' || res[len(res)-1] != 'R' {
			err = fmt.Errorf("invalid terminal response for get cursor: %v", res)
			return
		}

		s := string(res[2 : len(res)-1])
		var split int

		for i, r := range s {
			if r == ':' || r == ';' {
				split = i
				break
			}
		}

		if split == 0 {
			err = fmt.Errorf("invalid terminal response for get cursor: %v", res)
			return
		}

		p.y, err = strconv.Atoi(string(s[:split]))
		if err != nil {
			return
		}

		p.x, err = strconv.Atoi(string(s[split+1:]))
		return
	})
	return p.x, p.y, err
}
