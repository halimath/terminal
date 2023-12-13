package csi

import (
	"bytes"
	"testing"

	"github.com/halimath/expect"
	"github.com/halimath/expect/is"
)

func TestMoveCursorUp(t *testing.T) {
	expect.That(t, is.EqualTo(MoveCursorUp(2), "\x1b[2A"))
}

func TestMoveCursorDown(t *testing.T) {
	expect.That(t, is.EqualTo(MoveCursorDown(2), "\x1b[2B"))
}

func TestMoveCursorForward(t *testing.T) {
	expect.That(t, is.EqualTo(MoveCursorForward(2), "\x1b[2C"))
}

func TestMoveCursorBackward(t *testing.T) {
	expect.That(t, is.EqualTo(MoveCursorBackward(2), "\x1b[2D"))
}

func TestSetCursorPosition(t *testing.T) {
	expect.That(t, is.EqualTo(SetCursorPosition(2, 3), "\x1b[3;2H"))
}

func TestGetCursorPosition(t *testing.T) {
	t.Run("validResponse", func(t *testing.T) {
		var rw rw
		rw.r.WriteString("\x1b[2;3R")

		x, y, err := GetCursorPosition(&rw)
		expect.That(t,
			is.NoError(err),
			is.EqualTo(3, x),
			is.EqualTo(2, y),
		)
	})
	t.Run("invalidResponse", func(t *testing.T) {
		var rw rw
		rw.r.WriteString("caboom")

		_, _, err := GetCursorPosition(&rw)
		expect.That(t,
			is.Error(err, ErrInvalidTerminalResponse),
		)
	})

	t.Run("invalidResponse2", func(t *testing.T) {
		var rw rw
		rw.r.WriteString("\x1b[a;bR")

		_, _, err := GetCursorPosition(&rw)
		expect.That(t,
			is.Error(err, ErrInvalidTerminalResponse),
		)
	})
}

type rw struct {
	r bytes.Buffer
	w bytes.Buffer
}

func (rw *rw) Read(p []byte) (n int, err error) {
	return rw.r.Read(p)
}

func (rw *rw) Write(p []byte) (n int, err error) {
	return rw.w.Write(p)
}
