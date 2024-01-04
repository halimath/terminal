// Package terminal provides a plattform agnostic interface to interact with a terminal.
package terminal

import (
	"errors"
	"fmt"
	"os"

	"github.com/halimath/terminal/input"
	"github.com/halimath/terminal/rawmode"
	"github.com/halimath/terminal/sgr"
)

var (
	// ErrRawMode is a sentinel error value returned from New or NewWithFD in case switching the terminal to
	// raw mode failed.
	ErrRawMode = errors.New("failed to activate raw mode")
)

// TruecolorSupported returns whether the environment this process runs in supports truecolor.
// This function checks the environment variable COLORTERM to be set to "truecolor".
func IsTruecolorSupported() bool {
	return os.Getenv("COLORTERM") == "truecolor"
}

// Terminal implements both read and write access to the terminal. By default, it runs on STDIN/STDOUT but can
// be configured to work with other file descriptors as well.
type Terminal struct {
	r, w        *os.File
	inputReader *input.Reader

	rawModeRestoreState *rawmode.State
}

// New creates a new Terminal using os.Stdin for input and os.Stdout for output.
func New() *Terminal {
	return NewWithFile(os.Stdin, os.Stdout)
}

// NewWithFile creates and initializes a new terminal using r and w as reader and write. Both r and w may
// point to the same os.File.
func NewWithFile(r, w *os.File) *Terminal {
	t := Terminal{
		r:           r,
		w:           w,
		inputReader: &input.Reader{Reader: r},
	}

	return &t
}

// IsTerminal returns true if t is connected to a TTY device.
func (t *Terminal) IsTerminal() bool {
	return rawmode.IsTerminal(t.w.Fd())
}

// EnterRawMode activates the terminal raw mode. In raw mode, key presses are sent down to the FD directly
// and no line buffering happens (as in canonical mode).
//
// Make sure to pair a call to EnterRawMode with a call to ExitRawMode to prevent leaving the terminal in
// a broken state.
func (t *Terminal) EnterRawMode() (err error) {
	t.rawModeRestoreState, err = rawmode.Enter(t.r.Fd())
	return
}

// ExitRawMode exits raw mode if it was previously entered. It is safe to invoke this method even if
// EnterRawMode hasn't been called as well as invoking this method multiple times.
func (t *Terminal) ExitRawMode() error {
	if t.rawModeRestoreState == nil {
		return nil
	}

	return rawmode.Restore(t.r.Fd(), t.rawModeRestoreState)
}

// Write writes the bytes in buf to the terminal and returns the number of bytes written and any error.
// Note that an error may occur after buf has been partially written. In that case, the returned number of
// bytes represents the number of bytes sucessfully written.
func (t *Terminal) Write(buf []byte) (int, error) {
	if !t.IsTerminal() {
		buf = sgr.Remove(buf)
	}

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

// ReadInputEvent reads a single input event from r. It returns the parsed event (or nil) as well as the actual
// bytes read. If an error occurs during reading both event and buffer are nil. If reading was successful
// but parsing the read bytes produced an error, the read bytes are returned for client code to handle them
// manually but event is nil. In any case, the returned error is non nil.
func (t *Terminal) ReadInputEvent() (input.Event, []byte, error) {
	return t.inputReader.ReadInputEvent()
}

// Size returns the size of the terminal.
func (t *Terminal) Size() (w, h int, err error) {
	w, h, err = rawmode.Size(t.w.Fd())
	return
}

// WriteString writes s to t. This method makes *Terminal satisfy io.StringWriter.
func (t *Terminal) WriteString(s string) (n int, err error) {
	return t.w.WriteString(s)
}

// Print is a convenient shortcut to calling
//
//	fmt.Fprint(t, arg)
func (t *Terminal) Print(arg ...any) (int, error) {
	return fmt.Fprint(t, arg...)
}

// Print is a convenient shortcut to calling
//
//	fmt.Fprint(t, arg)
func (t *Terminal) Println(arg ...any) (int, error) {
	return fmt.Fprintln(t, arg...)
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
