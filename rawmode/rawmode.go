package rawmode

type State struct {
	state
}

func Enter(fd uintptr) (*State, error) {
	return enter(fd)
}

func Restore(fd uintptr, state *State) error {
	return restore(fd, state)
}

func Size(fd uintptr) (width, height int, err error) {
	return size(fd)
}

func IsTerminal(fd uintptr) bool {
	return isTerminal(fd)
}
