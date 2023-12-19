package input

import "fmt"

type MouseEvent struct {
	Button  int
	X, Y    int
	Release bool
}

func (m MouseEvent) evt() {}
func (m MouseEvent) String() string {
	action := "press"
	if m.Release {
		action = "release"
	}

	return fmt.Sprintf("<%s %d at (%d,%d)>", action, m.Button, m.X, m.Y)
}
