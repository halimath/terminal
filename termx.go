// Package termx provides extensions for the golang.org/x/term package providing convenience functions
// and types.
package termx

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/halimath/termx/keypress"
	"golang.org/x/term"
)

var (
	// ErrRawMode is a sentinel error value returned from New or NewWithFD in case switching the terminal to
	// raw mode failed.
	ErrRawMode = errors.New("failed to activate raw mode")
)

// Terminal implements both read and write access to the terminal. By default, it runs on STDIN/STDOUT but can
// be configured to work with other file descriptors as well.
type Terminal struct {
	r, w *os.File

	// Whether this terminal support true colors
	truecolorSupported bool

	// the terminal state found when starting raw mode
	restoreState *term.State
}

// New creates a new Terminal using os.Stdin for input and os.Stdout for output.
func New() *Terminal {
	t := NewWithFile(os.Stdin, os.Stdout)
	t.truecolorSupported = os.Getenv("COLORTERM") == "truecolor"
	return t
}

// NewWithFile creates and initializes a new terminal using r and w as reader and write. Both r and w may
// point to the same os.File.
func NewWithFile(r, w *os.File) *Terminal {
	t := Terminal{
		r: r,
		w: w,
	}

	return &t
}

// TruecolorSupported returns whether t supports truecolor. Note that this method will only return true, iff
// the terminal is bound to STDIN/STDOUT (has been created with New) and the environment variable COLORTERM is
// set to "truecolor".
//
// For all other situations this method simply returns false. Nevertheless, the terminal may support truecolor
// anyhow.
func (t *Terminal) TruecolorSupported() bool {
	return t.truecolorSupported
}

// EnableRawMode activates the terminal raw mode. In raw mode, key presses are sent down to the FD directly
// and no line buffering happens (as in canonical mode).
func (t *Terminal) EnableRawMode(enableCaMode bool) error {
	var err error
	t.restoreState, err = term.MakeRaw(int(t.r.Fd()))
	if err != nil {
		return err
	}

	if enableCaMode {
		t.WriteString(EnterCaMode)
	}

	return nil
}

// DisableRawMode disables a previously activated raw mode and restores the state the terminal was in before
// activating the raw mode (usually canonical mode).
func (t *Terminal) DisableRawMode() error {
	t.WriteString(ExitCaMode)

	if t.restoreState == nil {
		return nil
	}

	return term.Restore(int(t.r.Fd()), t.restoreState)
}

// Close closes the terminal and must be called before the application terminates to restore terminal state.
// If not called the application may leave the terminal in a destroyed and unusable state.
func (t *Terminal) Close() error {
	if t.restoreState == nil {
		return nil
	}

	return t.DisableRawMode()
}

// Write writes the bytes in buf to the terminal and returns the number of bytes written and any error.
// Note that an error may occur after buf has been partially written. In that case, the returned number of
// bytes represents the number of bytes sucessfully written.
func (t *Terminal) Write(buf []byte) (int, error) {
	n, err := t.w.Write(buf)

	if err == nil && n < len(buf) {
		err = fmt.Errorf("partial write: %d < %d", n, len(buf))
	}

	if err != nil {
		return n, err
	}

	return n, nil
}

// Read reads bytes from the terminal until buf is filled. It returns the number of bytes read (which may be
// less then len(buf)) as well as any error that occured.
func (t *Terminal) Read(buf []byte) (int, error) {
	return t.r.Read(buf)
}

// ReadKeyPress reads a single keypress from the t. It returns the
// parsed keypress (or nil) as well as the actual bytes read. If an
// error occurs the parsed keypress is nil. If reading was successful
// but parsing the read bytes produced an error, the bytes are returned
// for client code to handle them manually.
func (t *Terminal) ReadKeyPress() (keypress.KeyPress, []byte, error) {
	var buf [4]byte

	n, err := t.Read(buf[:])
	if err != nil {
		return nil, nil, err
	}

	k, err := keypress.Decode(buf[:n])
	if err != nil {
		return nil, buf[:n], err
	}

	return k, buf[:n], nil
}

// Size returns the size of the terminal.
func (t *Terminal) Size() (w, h int, err error) {
	w, h, err = term.GetSize(int(t.w.Fd()))
	return
}

// Printf is a convenience method adding fmt.Printf support for a *Terminal. Basicall, invoking
//
//	t.Printf(format, args...)
//
// is equivalent to
//
//	fmt.Fprintf(t, format, args...)
func (t *Terminal) Printf(format string, args ...any) (int, error) {
	return fmt.Fprintf(t, format, args...)
}

// WriteString writes s to t. This method makes *Terminal satisfy io.StringWriter.
func (t *Terminal) WriteString(s string) (n int, err error) {
	return t.w.WriteString(s)
}

const (
	// Core definitions for partial control sequences to be used in other definitions
	ESC              = "\x1b"   // The Escape character
	CSI              = "\x1b["  // Control Sequence Introducer (0x9b)
	StringTerminator = "\x1b\\" // String terminator sequence (0x9c)- used to terminate some sequences
	OSC              = "\x1b]"  // Operating System Command (0x9d)

	ResetTerminal = "\x1Bc" // Reset all terminal attributes to their default

	// CaMode Commands
	EnterCaMode = "\x1b[?1049h"
	ExitCaMode  = "\x1b[?1049l"

	// Clear commands - all these are CSI sequences to be written to the terminal
	ClearScreen       = "\x1B[2J" // Clear whole screen
	ClearLine         = "\x1B[2K" // Clear current line
	ClearAfterCursor  = "\x1B[J"  // Clear everything after the cursor
	ClearBeforeCursor = "\x1B[1J" // Clear everything before the cursor
	ClearUntilNewline = "\x1B[K"  // Clear from cursor to newline
)

// SetWindowTitle creates a control sequence to set the window's title to title and writes that sequence to
// w returning the number of bytes written and any error.
func SetWindowTitle(w io.Writer, title string) (int, error) {
	return fmt.Fprintf(w, "%s2;%s%s", OSC, title, StringTerminator)
}

var queryBackgroundColor = OSC + "11;?" + StringTerminator

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
