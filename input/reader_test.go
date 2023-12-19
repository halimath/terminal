package input

import (
	"bytes"
	"testing"

	"github.com/halimath/expect"
	"github.com/halimath/expect/is"
)

func TestReader_ReadInputEvent(t *testing.T) {
	t.Run("regular_char", func(t *testing.T) {
		var buf bytes.Buffer
		r := &Reader{Reader: &buf}
		buf.WriteRune('a')

		got, _, err := r.ReadInputEvent()
		var want Event = Char('a')
		expect.That(t,
			is.NoError(err),
			is.DeepEqualTo(got, want),
		)
	})

	t.Run("special_char", func(t *testing.T) {
		var buf bytes.Buffer
		r := &Reader{Reader: &buf}
		buf.WriteString("\x1bO\x48")

		got, _, err := r.ReadInputEvent()
		var want Event = Home
		expect.That(t,
			is.NoError(err),
			is.DeepEqualTo(got, want),
		)
	})

	t.Run("two_special_chars", func(t *testing.T) {
		var buf bytes.Buffer
		r := &Reader{Reader: &buf}
		buf.WriteString("\x1bO\x48\x1bO\x48")

		got, _, err := r.ReadInputEvent()
		var want Event = Home
		expect.That(t,
			is.NoError(err),
			is.DeepEqualTo(got, want),
		)
	})
}
