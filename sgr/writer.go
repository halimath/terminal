package sgr

import (
	"fmt"
	"io"
)

// Writer defines an interface for types that wrap a io.Writer but support
// disabling formatted output.
type Writer interface {
	io.Writer

	SuppressSGR() bool
}

// writer implements Writer with a suppress toggle.
type writer struct {
	io.Writer
	suppressSGR bool
}

func (w *writer) SuppressSGR() bool { return w.suppressSGR }

// Suppress decorates w to suppress SGR formatted output.
func Suppress(w io.Writer) Writer {
	return &writer{Writer: w, suppressSGR: true}
}

// Print prints all a formatted with SGR to w. If w satisfies Writer and
// SuppressSGR returns true, this function works identical to fmt.Print.
func Print(w io.Writer, sgr SGR, a ...any) (int, error) {
	wr, ok := w.(Writer)

	if ok && wr.SuppressSGR() {
		return fmt.Fprint(w, a...)
	}

	return io.WriteString(w, Format(sgr, fmt.Sprint(a...)))
}

// Println prints all a formatted with SGR to w followed by a newline. If w
// satisfies Writer and SuppressSGR returns true, this function works identical
// to fmt.Print.
func Println(w io.Writer, sgr SGR, a ...any) (int, error) {
	wr, ok := w.(Writer)

	if ok && wr.SuppressSGR() {
		return fmt.Fprintln(w, a...)
	}

	return io.WriteString(w, Format(sgr, fmt.Sprintln(a...)))

}

// Printf prints all args applied to format formatted with SGR to w. If w
// satisfies Writer and SuppressSGR returns true, this function works identical
// to fmt.Print.
func Printf(w io.Writer, sgr SGR, format string, args ...any) (int, error) {
	wr, ok := w.(Writer)

	if ok && wr.SuppressSGR() {
		return fmt.Fprintf(w, format, args...)
	}

	return io.WriteString(w, Format(sgr, fmt.Sprintf(format, args...)))
}
