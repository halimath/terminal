package input

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/halimath/termx/csi"
)

var ErrInvalidInputBytes = errors.New("invalid input byte sequence")

const (
	keyCodeTab         = 0x09 // <Tab>
	keyCodeReturn      = 0x0d // <Ret>
	keyCodeEscape      = 0x1b // <Esc>
	keyCodeOpenBracket = 0x5b // [
)

func Decode(b []byte) (Event, error) {
	if len(b) == 0 {
		return nil, ErrInvalidInputBytes
	}

	switch len(b) {
	case 1:
		return decodeSingleByteKeyPress(b[0])
	case 2:
		// Any two byte sequence is a unicode character
		return decodeUnicodeRune(b)
	case 3:
		return decodeThreeByteKeyPress(b)
	case 4:
		return decodeFourByteKeyPress(b)

	default:
		return decodeMouseEvent(b)
	}
}

func decodeUnicodeRune(b []byte) (KeyPress, error) {
	r, l := utf8.DecodeRune(b)
	if l != len(b) {
		return nil, fmt.Errorf("%w: invalid unicode character: %#v", ErrInvalidInputBytes, b)
	}
	return Char(r), nil

}

// decodeMouseEvent decodes an X10 encoded mouse event according to [xterm].
//
// [xterm]: https://invisible-island.net/xterm/ctlseqs/ctlseqs.html#h3-X10-compatibility-mode
func decodeMouseEvent(b []byte) (Event, error) {
	if bytes.HasPrefix(b, x10MouseEventPrefix) && len(b) == 6 {
		return decodeX10MouseEvent(b)
	}

	return decodeSGRMouseEvent(b)
}

var sgrMouseEventPrefix = []byte(csi.CSI + "<")

// decodeSGRMouseEvent decodes an SGR encoded mouse event. According to [xterm] an SGR encoded event is
// encoded as follows
//
//	CSI < Btn ; Px ; Py [M|m]
//
// The final character encoded if the button was pressed (M) or released (m).
// [xterm]: https://invisible-island.net/xterm/ctlseqs/ctlseqs.html#h3-Extended-coordinates
func decodeSGRMouseEvent(b []byte) (Event, error) {
	if !bytes.HasPrefix(b, sgrMouseEventPrefix) {
		return nil, fmt.Errorf("%w: invalid mouse event prefix: %s", ErrInvalidInputBytes, string(b[1:]))
	}

	if b[len(b)-1] != 'm' && b[len(b)-1] != 'M' {
		return nil, fmt.Errorf("%w: invalid mouse event suffix: %s", ErrInvalidInputBytes, string(b[1:]))
	}

	release := b[len(b)-1] == 'm'

	parts := strings.Split(string(b[3:len(b)-1]), ";")
	if len(parts) != 3 {
		return nil, fmt.Errorf("%w: invalid number of mouse event arguments: %s", ErrInvalidInputBytes, string(b[1:]))
	}

	flags, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, fmt.Errorf("%w: invalid mouse event btn: %s", ErrInvalidInputBytes, string(b[1:]))
	}
	x, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, fmt.Errorf("%w: invalid mouse event x pos: %s", ErrInvalidInputBytes, string(b[1:]))
	}
	y, err := strconv.Atoi(parts[2])
	if err != nil {
		return nil, fmt.Errorf("%w: invalid mouse event y pos: %s", ErrInvalidInputBytes, string(b[1:]))
	}

	return MouseEvent{
		Button:  determineMouseBtn(flags),
		X:       x,
		Y:       y,
		Release: release,
	}, nil
}

const (
	x10MouseRelease = 3

	//   4=Shift,
	//   8=Meta, and
	//   16=Control.
)

func determineMouseBtn(flags int) int {
	if flags&x10MouseRelease == 1 {
		return 2
	} else if flags&x10MouseRelease == 2 {
		return 3
	}

	return 1
}

var x10MouseEventPrefix = []byte(csi.CSI + "M")

func decodeX10MouseEvent(b []byte) (Event, error) {
	const x10MouseByteOffset = 32

	flags := int(b[3]) - x10MouseByteOffset
	x := int(b[4]) - x10MouseByteOffset
	y := int(b[5]) - x10MouseByteOffset

	if flags < 0 {
		return nil, fmt.Errorf("%w: invalid X10 mouse input: invalid button: %v", ErrInvalidInputBytes, b)
	}

	return MouseEvent{
		Button:  determineMouseBtn(flags),
		X:       x,
		Y:       y,
		Release: flags&x10MouseRelease == 3}, nil
}

func decodeSingleByteKeyPress(b byte) (KeyPress, error) {
	switch b {
	case 0:
		return Ctrl(' '), nil
	case keyCodeTab:
		return Tab, nil
	case keyCodeReturn:
		return Return, nil
	case keyCodeEscape:
		return Escape, nil
	case 127:
		return Backspace, nil
	}

	// Bytes 1 - 26 represent CTRL-a .. CTRL-z
	if b < keyCodeEscape {
		return Ctrl(rune('a' + b - 1)), nil
	}

	// Otherwise its a plain character
	return Char(b), nil
}

func decodeThreeByteKeyPress(b []byte) (KeyPress, error) {
	if b[0] != keyCodeEscape {
		return decodeUnicodeRune(b)
	}

	// These sequences are used in normal mode. See
	// https://invisible-island.net/xterm/ctlseqs/ctlseqs.html#h2-Special-Keyboard-Keys
	if b[1] == keyCodeOpenBracket {
		switch b[2] {
		case 0x41:
			return CursorUp, nil
		case 0x42:
			return CursorDown, nil
		case 0x43:
			return CursorRight, nil
		case 0x44:
			return CursorLeft, nil
		case 0x46:
			return End, nil
		case 0x48:
			return Home, nil

			// F5       | CSI 1 5 ~
			// F6       | CSI 1 7 ~
			// F7       | CSI 1 8 ~
			// F8       | CSI 1 9 ~
			// F9       | CSI 2 0 ~
			// F10      | CSI 2 1 ~
			// F11      | CSI 2 3 ~
			// F12      | CSI 2 4 ~

		default:
			return nil, ErrInvalidInputBytes
		}
	}

	// These sequences are used in application mode. See
	// https://invisible-island.net/xterm/ctlseqs/ctlseqs.html#h2-Special-Keyboard-Keys
	if b[1] == 0x4f {
		switch b[2] {
		case 0x41:
			return CursorUp, nil
		case 0x42:
			return CursorDown, nil
		case 0x43:
			return CursorRight, nil
		case 0x44:
			return CursorLeft, nil
		case 0x46:
			return End, nil
		case 0x48:
			return Home, nil
		case 0x50, 0x51, 0x52, 0x53:
			return FunctionKey(b[2] - 79), nil
		default:
			return nil, ErrInvalidInputBytes
		}

	}

	return nil, ErrInvalidInputBytes
}

func decodeFourByteKeyPress(b []byte) (KeyPress, error) {
	if b[0] == keyCodeEscape && b[1] == keyCodeOpenBracket {
		if b[2] == 0x33 && b[3] == 0x7e {
			return Delete, nil
		}

		if b[2] == 0x35 && b[3] == 0x7e {
			return PageUp, nil
		}

		if b[2] == 0x36 && b[3] == 0x7e {
			return PageDown, nil
		}
	}

	return decodeUnicodeRune(b)
}

const readerBufSize = 256

// Reader is a wrapper around an io.Reader which supports reading multiple
// control sequences that go together.
type Reader struct {
	io.Reader
	lastReadBuffer []byte
}

// ReadInputEvent reads a single input event from r. It returns the parsed event (or nil) as well as the actual
// bytes read. If an error occurs during reading both event and buffer are nil. If reading was successful
// but parsing the read bytes produced an error, the read bytes are returned for client code to handle them
// manually but event is nil. In any case, the returned error is non nil.
func (r *Reader) ReadInputEvent() (Event, []byte, error) {
	// If there is something left in the buffer, consume it first.
	if len(r.lastReadBuffer) > 0 {
		return r.readFromLastReadBuffer()
	}

	// Read from underlying reader up to 256 bytes
	var singleReadBuf [readerBufSize]byte
	var buf bytes.Buffer
	for {
		l, err := r.Read(singleReadBuf[:])
		if err != nil {
			if errors.Is(err, io.EOF) && buf.Len() > 0 {
				break
			}

			return nil, buf.Bytes(), err
		}
		buf.Write(singleReadBuf[:l])

		// If we haven't read up to the limit, we're done and continue with
		// decoding. Otherwise read again to pick up any remaining bytes.
		if l < readerBufSize {
			break
		}
	}

	// Copy everything to the last read buffer...
	r.lastReadBuffer = buf.Bytes()
	// ... and read from there
	return r.readFromLastReadBuffer()
}

func (r *Reader) readFromLastReadBuffer() (Event, []byte, error) {
	idx := r.findSecondEventOffsetLastReadBuffer()
	var buf []byte

	if idx == -1 {
		buf = r.lastReadBuffer
		r.lastReadBuffer = nil
	} else {
		buf = r.lastReadBuffer[0:idx]
		r.lastReadBuffer = r.lastReadBuffer[idx:]
	}

	evt, err := Decode(buf)
	return evt, buf, err
}

func (r *Reader) findSecondEventOffsetLastReadBuffer() int {
	idx := bytes.Index(r.lastReadBuffer[1:], []byte(csi.ESC))
	if idx == -1 {
		return -1
	}
	return idx + 1
}
