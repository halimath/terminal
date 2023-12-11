package sgr

import (
	"strings"
	"testing"

	"github.com/halimath/expect"
	"github.com/halimath/expect/is"
)

func TestPackage(t *testing.T) {
	var sb strings.Builder

	Print(&sb, FgRed, "hello, world\n")
	Println(&sb, FgRed, "hello, world")
	Printf(&sb, FgRed, "hello, %s\n", "world")

	expect.That(t,
		is.EqualTo(sb.String(), "\x1B[31mhello, world\n\x1B[0m\x1B[31mhello, world\n\x1B[0m\x1B[31mhello, world\n\x1B[0m"),
	)
}

func TestPackage_suppress(t *testing.T) {
	var sb strings.Builder

	w := Suppress(&sb)

	Print(w, FgRed, "hello, world\n")
	Println(w, FgRed, "hello, world")
	Printf(w, FgRed, "hello, %s\n", "world")

	expect.That(t,
		is.EqualTo(sb.String(), "hello, world\nhello, world\nhello, world\n"),
	)
}
