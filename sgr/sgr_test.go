package sgr

import (
	"testing"

	"github.com/halimath/expect"
	"github.com/halimath/expect/is"
)

func TestEscape(t *testing.T) {
	expect.That(t, is.EqualTo(Escape(FgRed), "\x1B[31m"))
}

func TestJoin(t *testing.T) {
	type testCase struct {
		in   []SGR
		want string
	}

	tests := []testCase{
		{in: []SGR{}, want: ""},
		{in: []SGR{FgRed}, want: "31"},
		{in: []SGR{FgRed, BgGreen}, want: "31;42"},
	}

	for _, test := range tests {
		got := Join(test.in...)
		expect.That(t, is.EqualTo(got, SGR(test.want)))
	}
}

func TestFgRGB(t *testing.T) {
	expect.That(t, is.EqualTo(FgRGB(0, 2, 5), "38;5;33"))
}

func TestBgRGB(t *testing.T) {
	expect.That(t, is.EqualTo(BgRGB(0, 2, 5), "48;5;33"))
}

func TestAssertValidRGB(t *testing.T) {
	defer func() {
		expect.That(t, is.EqualTo(recover(), "invalid ANSI color: 9;0;0"))
	}()

	assertValidRGB(9, 0, 0)
}

func TestFgTrueColor(t *testing.T) {
	expect.That(t, is.EqualTo(FgTrueColor(0, 128, 59), "38;2;0;128;59"))
}

func TestBgTrueColor(t *testing.T) {
	expect.That(t, is.EqualTo(BgTrueColor(0, 128, 59), "48;2;0;128;59"))
}

func TestFormat(t *testing.T) {
	expect.That(t, is.EqualTo(Format(FgRed, "hello, world"), "\x1B[31mhello, world\x1B[0m"))
}

func TestFormatf(t *testing.T) {
	expect.That(t, is.EqualTo(Formatf(FgRed, "hello, %s", "world"), "\x1B[31mhello, world\x1B[0m"))
}
