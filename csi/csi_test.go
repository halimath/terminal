package csi

import (
	"testing"

	"github.com/halimath/expect"
	"github.com/halimath/expect/is"
)

func TestSetWindowTitle(t *testing.T) {
	expect.That(t, is.EqualTo(SetWindowTitle("hello, world"), "\x1b]2;hello, world\x1b\\"))
}

func TestGetBackgroundColor(t *testing.T) {
	t.Run("validResponse", func(t *testing.T) {
		var rw rw
		rw.r.WriteString("\x1b]11;rgb:aaa/bbb/ccc\x1b\\")

		r, g, b, err := GetBackgroundColor(&rw)
		expect.That(t,
			is.NoError(err),
			is.EqualTo(2730, r),
			is.EqualTo(3003, g),
			is.EqualTo(3276, b),
		)
	})
	t.Run("invalidResponse", func(t *testing.T) {
		var rw rw
		rw.r.WriteString("caboom")

		_, _, _, err := GetBackgroundColor(&rw)
		expect.That(t,
			is.Error(err, ErrInvalidTerminalResponse),
		)
	})
}
