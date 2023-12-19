//go:build windows

package rawmode

import (
	"golang.org/x/sys/windows"
)

type state struct {
	mode uint32
}

func isTerminal(fd uintptr) bool {
	var st uint32
	err := windows.GetConsoleMode(windows.Handle(fd), &st)
	return err == nil
}

func enter(fd uintptr) (*State, error) {
	var consoleMode uint32

	if err := windows.GetConsoleMode(windows.Handle(fd), &consoleMode); err != nil {
		return nil, err
	}

	oldState := State{state{consoleMode}}

	consoleMode &^= windows.ENABLE_ECHO_INPUT | windows.ENABLE_PROCESSED_INPUT | windows.ENABLE_LINE_INPUT
	consoleMode |= windows.ENABLE_PROCESSED_OUTPUT | windows.ENABLE_VIRTUAL_TERMINAL_INPUT

	if err := windows.SetConsoleMode(windows.Handle(fd), consoleMode); err != nil {
		return nil, err
	}

	return &oldState, nil
}

func restore(fd uintptr, state *State) error {
	return windows.SetConsoleMode(windows.Handle(fd), state.mode)
}

func size(fd uintptr) (width, height int, err error) {
	var info windows.ConsoleScreenBufferInfo
	if err = windows.GetConsoleScreenBufferInfo(windows.Handle(fd), &info); err != nil {
		return
	}

	width = int(info.Window.Right - info.Window.Left + 1)
	height = int(info.Window.Bottom - info.Window.Top + 1)

	return
}
