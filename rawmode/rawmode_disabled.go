//go:build !(windows || aix || linux || solaris || zos || darwin || dragonfly || freebsd || netbsd || openbsd)

package rawmode

import "errors"

var errNotImplemented = errors.New("rawmode not implemented for current OS")

type state struct {
}

func isTerminal(fd uintptr) bool {
	return false
}

func enter(fd uintptr) (*State, error) {
	return nil, errNotImplemented
}

func restore(fd uintptr, state *State) error {
	return errNotImplemented
}

func size(fd uintptr) (width, height int, err error) {
	return 0, 0, errNotImplemented
}
