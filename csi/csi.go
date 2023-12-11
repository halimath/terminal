package csi

import (
	"bytes"
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

	// Commands to enter/exit CA mode, a.k.a. as CUP mode. See the terminfo documentation on smcup.
	EnterCaMode = CSI + "?1049h"
	ExitCaMode  = CSI + "?1049l"

	// Clear commands - all these are CSI sequences to be written to the terminal
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
			err = fmt.Errorf("invalid terminal response for get background color: %v", string(res[1:]))
			return
		}

		s := string(res[len(OSC) : len(res)-len(StringTerminator)])
		parts := strings.Split(s, ";")

		if len(parts) != 2 || parts[0] != "11" || !strings.HasPrefix(parts[1], rgbPrefix) {
			err = fmt.Errorf("invalid terminal response for get background color: %v", string(res[1:]))
			return
		}

		components := strings.Split(strings.TrimPrefix(parts[1], rgbPrefix), "/")
		if len(components) != 3 {
			err = fmt.Errorf("invalid terminal response for get background color: %v", string(res[1:]))
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
