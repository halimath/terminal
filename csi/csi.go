// Package csi contains constants and functions to produce _control sequence introducer_ strings that instruct
// a terminal to perform special operations, i.e. moving the cursor, clearing the screen and so on.
// Static (parameter less) CSI are defined as constants, dynamic (parameterized) CSI are defined as functions
// returning strings. These sequences can be written to a terminal or any other buffer/file and be output.
// Note that the sequences contain special characters (most notably the escape character '\x1b') which can
// cause scrabled output when displayed by a device that does not interpret the sequences.
//
// In addition this package contains some query operations. These are defined as functions that receive an
// io.ReadWriter and return the query's response.
//
// All sequences defined in this package are based on the xterm definitions and may not work on certain
// devices. As almost all of the terminal emulations in use nowadays are compatible with xterm (to some
// extend, at least), this should not be considered a limitation. But be warned that some sequences may not
// work as intended when sent to a VT100 terminal.
package csi

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

const (
	// Core definitions for partial control sequences to be used in other definitions
	ESC              = "\x1b"     // The Escape character
	CSI              = ESC + "["  // Control Sequence Introducer (0x9b)
	StringTerminator = ESC + "\\" // String terminator sequence (0x9c)- used to terminate some sequences
	OSC              = ESC + "]"  // Operating System Command (0x9d)

	ResetTerminal = ESC + "c" // Reset all terminal attributes to their default

	// Commands to use alternate screen buffer. terminfo calls this CUP mode. See the terminfo
	// documentation on smcup and rmcup.
	//
	// In alternate screen buffer mode no scrolling is possible - the buffer has the exact size of the
	// screen. Thus, applications using alternate screen buffers will receive keypresses for scrolling
	// keys such as PgUp, PgDown, Home and End.
	//
	// It is recommended to switch to an alternate buffer at the very beginning of the application and return
	// to main buffer as one of the very last operations.
	UseAlternateScreenBuffer = CSI + "?1049h" // Switches to a new alternate screen buffer
	UseMainScreenBuffer      = CSI + "?1049l" // Switches to the main buffer

	// Sequences to enable/disable application mode. In application mode, the terminal sends different
	// control sequences.
	EnableApplicationMode  = CSI + "?1h"
	DisableApplicationMode = CSI + "?1l"

	// Commands to clear certain areas of the screen
	ClearScreen       = CSI + "2J" // Clear whole screen
	ClearLine         = CSI + "2K" // Clear current line
	ClearAfterCursor  = CSI + "J"  // Clear everything after the cursor
	ClearBeforeCursor = CSI + "1J" // Clear everything before the cursor
	ClearUntilNewline = CSI + "K"  // Clear from cursor to newline
)

// SetWindowTitle creates a control sequence to set the window's title to title.
func SetWindowTitle(title string) string {
	return fmt.Sprintf("%s2;%s%s", OSC, title, StringTerminator)
}

const queryBackgroundColor = OSC + "11;?" + StringTerminator

const rgbPrefix = "rgb:"

// GetBackgroundColor retrieves the background color of the terminal and returns it as r,g,b values each
// representing a single color component in 16bit resolution.
func GetBackgroundColor(rw io.ReadWriter) (r, g, b uint16, err error) {
	c, err := execQuery(rw, queryBackgroundColor, 128, func(res []byte) (rgb [3]uint16, err error) {
		if !bytes.HasPrefix(res, []byte(OSC)) || !(bytes.HasSuffix(res, []byte(StringTerminator)) || bytes.HasSuffix(res, []byte{'\a'})) {
			err = fmt.Errorf("%w: get background color: %v", ErrInvalidTerminalResponse, string(res[1:]))
			return
		}

		s := string(res[len(OSC) : len(res)-len(StringTerminator)])
		parts := strings.Split(s, ";")

		if len(parts) != 2 || parts[0] != "11" || !strings.HasPrefix(parts[1], rgbPrefix) {
			err = fmt.Errorf("%w: get background color: %v", ErrInvalidTerminalResponse, string(res[1:]))
			return
		}

		components := strings.Split(strings.TrimPrefix(parts[1], rgbPrefix), "/")
		if len(components) != 3 {
			err = fmt.Errorf("%w: get background color: %v", ErrInvalidTerminalResponse, string(res[1:]))
			return
		}

		var v uint64
		for i := 0; i < 3; i++ {
			v, err = strconv.ParseUint(components[i], 16, 16)
			if err != nil {
				return
			}
			rgb[i] = uint16(v)
		}

		return
	})

	return c[0], c[1], c[2], err
}

// ErrInvalidTerminalResponse is a sentinel error value returned from queries issued to the terminal.
var ErrInvalidTerminalResponse = errors.New("invalid terminal response")

// queryHandler is used to define functions that handle terminal query responses.
type queryHandler[T any] func([]byte) (T, error)

// execQuery writes query to t and reads up to responseLimit bytes. The actual bytes read are passed to
// handler to produce a response.
func execQuery[T any](rw io.ReadWriter, query string, responseLimit int, handler queryHandler[T]) (result T, err error) {
	_, err = rw.Write([]byte(query))
	if err != nil {
		return
	}

	buf := make([]byte, responseLimit)

	var n int
	n, err = rw.Read(buf)
	if err != nil {
		return
	}

	return handler(buf[:n])
}
