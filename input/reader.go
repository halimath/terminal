package input

import (
	"bytes"
	"errors"
	"io"

	"github.com/halimath/terminal/csi"
)

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
