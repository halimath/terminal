// Package rawmode provides a minimal set of functions to enter and exit raw input mode. Entering raw mode
// requires a system call which is OS specific. This package collects all these syscall functions.
//
// The implementation has been largely inspired by golang.org/x/term but has been simplified as well as
// adopted for the specific needs and direction of this module.
package rawmode

// State represents an opaque value describing the OS specific state a terminal was in when entering raw mode.
// It can only be used to be restored which means leaving raw mode.
type State struct {
	state
}

// Enter enters raw input mode for the input file descriptor fd. This is usually os.Stdin.Fd() but may be
// selected differently for specific applications. Enter performs the required system calls and returns an
// opaque State which describes the terminal state before entering raw mode. Restoring that state later means
// to exit raw mode and return to canonical mode.
// In case the operation fails a non-nil error is returned and State is undefined.
func Enter(fd uintptr) (*State, error) {
	return enter(fd)
}

// Restore restores the terminal state on the input channel pointed to by fd with state. It returns any error
// that might happen during the system call.
func Restore(fd uintptr, state *State) error {
	return restore(fd, state)
}

// Size returns the terminal's dimension in character slots. fd should be the file descriptor referencing the
// input channel of the terminal.
func Size(fd uintptr) (width, height int, err error) {
	return size(fd)
}

// IsTerminal returns true, iff fd references an input channel connected to a terminal.
func IsTerminal(fd uintptr) bool {
	return isTerminal(fd)
}
