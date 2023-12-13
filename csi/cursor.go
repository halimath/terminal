package csi

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
)

const (
	// Commands for general cursor management
	CursorSave    = CSI + "s"    // Saves the current cursor
	CursorRestore = CSI + "u"    // Restores a previously stored cursor
	CursorHide    = CSI + "?25l" // Hide cursor
	CursorShow    = CSI + "?25h" // Show cursor

	// Commands to set cursor style
	CursorBlinkingBlock     = CSI + "\x31 q" // Change the cursor style to blinking block
	CursorSteadyBlock       = CSI + "\x32 q" // Change the cursor style to steady block
	CursorBlinkingUnderline = CSI + "\x33 q" // Change the cursor style to blinking underline
	CursorSteadyUnderline   = CSI + "\x34 q" // Change the cursor style to steady underline
	CursorBlinkingBar       = CSI + "\x35 q" // Change the cursor style to blinking bar
	CursorSteadyBar         = CSI + "\x36 q" // Change the cursor style to steady bar
)

// MoveCursorUp formats a CSI to move the cursor up by n rows.
func MoveCursorUp(n int) string {
	return fmt.Sprintf("%s%dA", CSI, n)
}

// MoveCursorDown formats a CSI to move the cursor down by n rows.
func MoveCursorDown(n int) string {
	return fmt.Sprintf("%s%dB", CSI, n)
}

// MoveCursorForward formats a CSI to move the cursor forward by n rows.
func MoveCursorForward(n int) string {
	return fmt.Sprintf("%s%dC", CSI, n)
}

// MoveCursorBackward formats a CSI to move the cursor backwards by n rows.
func MoveCursorBackward(n int) string {
	return fmt.Sprintf("%s%dD", CSI, n)
}

// SetCursorPosition formats a CSI to position the cursor at (x,y).
//
// According to ANSI terminal specs both coordinates are 1 based. This function adheres to that spec.
func SetCursorPosition(x, y int) string {
	return fmt.Sprintf("%s%d;%dH", CSI, y, x)
}

type position struct {
	x, y int
}

const getCursorPositionQuery = "\x1B[6n"

// GetCursorPosition queries the current cursor position from t and returns the x and y coordinates. Note that
// according to ANSI specs coordinates are 1 based, so the upper left corner is (1, 1).
func GetCursorPosition(rw io.ReadWriter) (x, y int, err error) {
	p, err := execQuery(rw, getCursorPositionQuery, 64, func(res []byte) (p position, err error) {
		if len(res) < 6 || !bytes.HasPrefix(res, []byte(CSI)) || res[len(res)-1] != 'R' {
			err = fmt.Errorf("%w: get cursor: %v", ErrInvalidTerminalResponse, err)
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
			err = fmt.Errorf("%w: get cursor: %v", ErrInvalidTerminalResponse, err)
			return
		}

		p.y, err = strconv.Atoi(string(s[:split]))
		if err != nil {
			err = fmt.Errorf("%w: get cursor: %v", ErrInvalidTerminalResponse, err)
			return
		}

		p.x, err = strconv.Atoi(string(s[split+1:]))
		return
	})
	return p.x, p.y, err
}
