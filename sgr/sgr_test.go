package sgr

import (
	"fmt"
	"testing"

	"github.com/halimath/expect"
	"github.com/halimath/expect/is"
)

func TestEscape(t *testing.T) {
	expect.That(t, is.EqualTo(FgRed.Escape(), "\x1B[31m"))
}

func TestJoin(t *testing.T) {
	type testCase struct {
		in   []SGR
		want string
	}

	base := Bold

	tests := []testCase{
		{in: []SGR{}, want: "1"},
		{in: []SGR{FgRed}, want: "1;31"},
		{in: []SGR{FgRed, BgGreen}, want: "1;31;42"},
	}

	for _, test := range tests {
		got := base.Join(test.in...)
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
		expect.That(t, is.EqualTo(fmt.Sprintf("%s", recover()), "invalid ANSI color: 9;0;0"))
	}()

	assertValidRGB(9, 0, 0)
}

func TestFgTrueColor(t *testing.T) {
	expect.That(t, is.EqualTo(FgTrueColor(0, 128, 59), "38;2;0;128;59"))
}

func TestBgTrueColor(t *testing.T) {
	expect.That(t, is.EqualTo(BgTrueColor(0, 128, 59), "48;2;0;128;59"))
}

func TestApply(t *testing.T) {
	expect.That(t, is.EqualTo(FgRed.Apply("hello, world"), "\x1B[31mhello, world\x1B[0m"))
}

func TestApplyf(t *testing.T) {
	expect.That(t, is.EqualTo(FgRed.Applyf("hello, %s", "world"), "\x1B[31mhello, world\x1B[0m"))
}

func TestRemove(t *testing.T) {
	tests := map[string]string{
		"foobar":          "foobar",
		Bold.Apply("foo"): "foo",
	}

	for in, want := range tests {
		expect.WithMessage(t, "input %q", in).
			That(is.EqualTo(string(Remove([]byte(in))), want))
	}
}

func BenchmarkRemove_withSGRs(b *testing.B) {
	s := []byte(Bold.Apply("Hello") + Faint.Apply("world"))

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		Remove(s)
	}
}

func BenchmarkRemove_withoutSGRs(b *testing.B) {
	s := []byte("hello world")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		Remove(s)
	}
}
